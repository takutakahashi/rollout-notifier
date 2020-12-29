package rollout

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	Namespace string
}
type Manager struct {
	cs *kubernetes.Clientset
	c  Config
}

func NewNamager(cs *kubernetes.Clientset, ns string) Manager {
	return Manager{cs: cs, c: Config{Namespace: ns}}

}

func (m Manager) GetTargets() ([]string, error) {
	result := []string{}
	dpc := m.cs.AppsV1().Deployments(m.c.Namespace)
	dpl, err := dpc.List(v1.ListOptions{})
	if err != nil {
		return []string{}, err
	}
	for _, dp := range dpl.Items {
		if !Finished(&dp) {
			result = append(result, dp.Name)
		}
	}
	return result, nil
}

func (m Manager) Finished(name string) (bool, error) {
	dpc := m.cs.AppsV1().Deployments(m.c.Namespace)
	dp, err := dpc.Get(name, v1.GetOptions{})
	if err != nil {
		return false, err
	}
	return Finished(dp), nil

}

func Finished(dp *appsv1.Deployment) bool {
	s := dp.Status.DeepCopy()
	replicas := *dp.Spec.DeepCopy().Replicas
	return s.ReadyReplicas == replicas && s.UpdatedReplicas == replicas

}
