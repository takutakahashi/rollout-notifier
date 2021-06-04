package notify

import "github.com/labstack/gommon/log"

type NoopNotify struct {
}

func (n NoopNotify) Init() error {
	return nil
}

func (n NoopNotify) Finish(target, comment string) error {
	log.Info("finish", target, comment)
	return nil
}

func (n NoopNotify) Start(target, comment string) error {
	log.Info("start", target, comment)
	return nil
}

func (n NoopNotify) Failed(target, comment string) error {
	log.Info("failed", target, comment)
	return nil
}
