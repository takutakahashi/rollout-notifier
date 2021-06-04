/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/takutakahashi/rollout-notifier/pkg/notify"
	"github.com/takutakahashi/rollout-notifier/pkg/rollout"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var commentAnnotation = "rollout-notifier.io/comment"

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	Progressing map[types.NamespacedName]bool
	Notify      notify.Notifier
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch

func (r *DeploymentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("deployment", req.NamespacedName)
	if r.Notify == nil {
		n, err := notify.NewNotify("slack", "/etc/rollout-notifier/config.yaml")
		if err != nil {
			n, err = notify.NewNotify("noop", "/")
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		r.Notify = n
	}
	n := r.Notify
	if r.Progressing == nil {
		r.Progressing = map[types.NamespacedName]bool{}
	} else if r.Progressing[req.NamespacedName] {
		return ctrl.Result{}, nil
	}
	var d appsv1.Deployment
	err := r.Get(ctx, req.NamespacedName, &d)
	if err != nil {
		return ctrl.Result{}, err
	}
	if rollout.Finished(&d) || rollout.Timeout(&d) {
		return ctrl.Result{}, nil
	}
	r.Progressing[req.NamespacedName] = true
	target := fmt.Sprintf("%s/%s", d.Namespace, d.Name)
	comment := d.Annotations[commentAnnotation]
	n.Start(target, comment)
	go func() {
		for {
			err := r.Get(ctx, req.NamespacedName, &d)
			if err != nil {
				r.Log.Error(err, "failed to get deployment")
				return
			}
			dd := d.DeepCopy()
			if rollout.Finished(dd) {
				n.Finish(target, comment)
				delete(r.Progressing, req.NamespacedName)
				return
			} else if rollout.Timeout(dd) {
				n.Failed(target, comment)
				delete(r.Progressing, req.NamespacedName)
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()

	return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
