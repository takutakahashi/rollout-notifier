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
	ConfigPath string
	NotifyType string
}

type Daemon struct {
	config Config
}

func NewConfig(namespace, configPath, notifyType string) (Config, error) {
	return Config{Namespace: namespace, ConfigPath: configPath, NotifyType: notifyType}, nil
}

func NewDaemon(c Config) (Daemon, error) {
	return Daemon{config: c}, nil
}

func (d Daemon) Start() error {
	n, err := notify.NewNotify(d.config.NotifyType, d.config.ConfigPath)
	if err != nil {
		return err
	}
	kconf := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(kconf)
	manager := rollout.NewNamager(clientset, d.config.Namespace)
	log.Info("starting daemon...")
	tracing := map[string]context.Context{}
	for {
		time.Sleep(15 * time.Second)
		targets, err := manager.GetTargets()
		if err != nil {
			log.Error(err)
			continue
		}
		for t := range tracing {
			finished, err := manager.Finished(t)
			if err != nil {
				log.Error(err)
				continue
			}
			if finished {
				delete(tracing, t)
				log.Infof("notify finish. %s", t)
				err = n.Finish(t)
				if err != nil {
					log.Error(err)
				}
			}
		}
		for _, t := range targets {
			beforeContext, alreadyStarted := tracing[t]
			if alreadyStarted && beforeContext != nil {
				continue
			}
			ctx := context.TODO()
			tracing[t] = ctx
			log.Infof("notify start. %s", t)
			err = n.Start(t)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
