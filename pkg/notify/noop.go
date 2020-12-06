package notify

import "github.com/labstack/gommon/log"

type NoopNotify struct {
}

func (n NoopNotify) Init() error {
	return nil
}

func (n NoopNotify) Finish(target string) error {
	log.Info("finish", target)
	return nil
}

func (n NoopNotify) Start(target string) error {
	log.Info("start", target)
	return nil
}
