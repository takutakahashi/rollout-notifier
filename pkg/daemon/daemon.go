package daemon

import (
	"context"
	"time"

	"github.com/takutakahashi/rollout-notifier/pkg/rollout"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Config struct {
	Namespace  string
	Webhook    string
	NotifyType string
}

type Daemon struct {
	config Config
}

func NewConfig(namespace, webhook, notifyType string) (Config, error) {
	return Config{Namespace: namespace, Webhook: webhook, NotifyType: notifyType}, nil
}

func NewDaemon(c Config) (Daemon, error) {
	return Daemon{config: c}, nil
}

func (d Daemon) Start() {
	kconf := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(kconf)
	manager := rollout.NewNamager(clientset, d.config.Namespace)
	tracing := map[string]context.Context{}
	_ = tracing
	for {
		targets, err := manager.GetTargets()
		if err != nil {
			continue
		}
		for _, t := range targets {
			beforeContext, alreadyStarted := tracing[t]
			if alreadyStarted && beforeContext != nil {
				continue
			}
			ctx, _ := context.WithTimeout(context.Background(), 80*time.Minute)
			tracing[t] = ctx
			go manager.StartWatch(ctx, t)
			go func() {
				select {
				case <-ctx.Done():
					delete(tracing, t)
				}
			}()

		}
		time.Sleep(10 * time.Second)
	}

}
