package notify

import "github.com/labstack/gommon/log"

type NotifyType string

var notifyTypeNoop NotifyType = "noop"
var notifyTypeSlack NotifyType = "slack"

type NoopNotify struct {
	nType   NotifyType
	webhook string
}

type SlackNotify struct {
	webhook string
}

type Notifier interface {
	Start(target string) error
	Finish(target string) error
}

func NewNotify(t, w string) Notifier {
	switch t {
	case "slack":
		return SlackNotify{webhook: w}
	default:
		return NoopNotify{}
	}
}

func (n NoopNotify) Finish(target string) error {
	log.Info("finish", target)
	return nil
}

func (n NoopNotify) Start(target string) error {
	log.Info("start", target)
	return nil
}

func (n SlackNotify) Finish(target string) error {

	return nil
}

func (n SlackNotify) Start(target string) error {

	return nil
}
