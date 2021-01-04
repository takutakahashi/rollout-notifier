package notify

import (
	"fmt"
	"io/ioutil"

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
	if err != nil {
		return sn, err
	}
	sn.c = slack.New(sn.config.Token)
	return sn, nil
}

func (n SlackNotify) Test() error {
	return n.notifyError("test")
}

func (n SlackNotify) Finish(target string) error {
	return n.notifySuccess(fmt.Sprintf("%s Rollout finished successflly! :+1:", target))
}

func (n SlackNotify) Failed(target string) error {
	return n.notifyError(fmt.Sprintf("%s Rollout failed. :cry:", target))
}

func (n SlackNotify) Start(target string) error {
	return n.notifySuccess(fmt.Sprintf("%s Rollout started! :rocket:", target))
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
