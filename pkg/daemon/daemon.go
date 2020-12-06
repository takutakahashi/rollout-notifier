package daemon

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/takutakahashi/rollout-notifier/pkg/notify"
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
	n := notify.NewNotify(d.config.NotifyType, d.config.Webhook)
	kconf := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(kconf)
	manager := rollout.NewNamager(clientset, d.config.Namespace)
	log.Info("starting daemon...")
	tracing := map[string]context.Context{}
	for {
		targets, err := manager.GetTargets()
		log.Info(targets)
		if err != nil {
			continue
		}
		for t := range tracing {
			finished, err := manager.Finished(t)
			if err != nil {
				continue
			}
			if finished {
				delete(tracing, t)
				n.Finish(t)
			}
		}
		for _, t := range targets {
			beforeContext, alreadyStarted := tracing[t]
			if alreadyStarted && beforeContext != nil {
				continue
			}
			ctx, _ := context.WithTimeout(context.Background(), 80*time.Minute)
			tracing[t] = ctx
			n.Start(t)

		}
		time.Sleep(1 * time.Second)
	}

}
