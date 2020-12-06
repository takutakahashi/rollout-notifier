package rollout

import (
	"context"

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
	dp := m.cs.AppsV1().Deployments(m.c.Namespace)
	dpl, err := dp.List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return []string{}, err
	}
	for _, dp := range dpl.Items {
		s := dp.Status.DeepCopy()
		if s.UpdatedReplicas != s.Replicas && s.UnavailableReplicas != 0 {
			result = append(result, dp.Name)
		}
	}
	return result, nil
}

func (m Manager) StartWatch(ctx context.Context, name string) {

}
