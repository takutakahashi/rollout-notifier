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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch

func (r *DeploymentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("deployment", req.NamespacedName)
	n, err := notify.NewNotify("noop", "/")
	if err != nil {
		return ctrl.Result{}, err
	}
	var d *appsv1.Deployment
	// your logic here
	n.Start(fmt.Sprintf("%s/%s", d.Namespace, d.Name))
	go func() {
		sec := 0
		err := r.Get(ctx, req.NamespacedName, d)
		dd := d.DeepCopy()
		if err != nil {
			r.Log.Error(err, "failed to get deployment")
			return
		}
		for {
			if rollout.Finished(dd) {
				n.Finish(fmt.Sprintf("%s/%s", dd.Namespace, dd.Name))
			}
			time.Sleep(10 * time.Second)
			sec += 10
			if sec > 600 {
				n.Failed(fmt.Sprintf("%s/%s", dd.Namespace, dd.Name))
				return
			}
		}
	}()

	return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
