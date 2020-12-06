package notify

import (
	"fmt"
	"io/ioutil"

	"github.com/labstack/gommon/log"
	"github.com/slack-go/slack"
	"gopkg.in/yaml.v2"
)

type SlackConfig struct {
	Token     string `yaml:"token"`
	ChannelID string `yaml:"channelID"`
}

type SlackNotify struct {
	config SlackConfig
	c      *slack.Client
}

func NewSlackNotify(configPath string) (SlackNotify, error) {
	sn := SlackNotify{}
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return sn, err
	}
	sn.config = SlackConfig{}
	err = yaml.Unmarshal(buf, &sn.config)
	log.Info(err)
	log.Info(sn)
	if err != nil {
		return sn, err
	}
	sn.c = slack.New(sn.config.Token)
	sn.Test()
	return sn, nil
}

func (n SlackNotify) Test() error {
	return n.notifyError("test")
}

func (n SlackNotify) Finish(target string) error {
	return n.notifySuccess(fmt.Sprintf("Rollout finished successflly! :+1: target: %s", target))
}

func (n SlackNotify) Start(target string) error {
	return n.notifySuccess(fmt.Sprintf("Rollout started! :rocket: target: %s", target))
}

func (n SlackNotify) notifySuccess(message string) error {
	return n.notify(message, "#009900")
}

func (n SlackNotify) notifyError(message string) error {
	return n.notify(message, "#c62828")
}

func (n SlackNotify) notify(message, color string) error {
	attachment := slack.Attachment{
		Text:  message,
		Color: color,
	}

	_, _, err := n.c.PostMessage(
		n.config.ChannelID,

		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	return err
}
