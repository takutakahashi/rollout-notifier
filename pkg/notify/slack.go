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
	return n.notifyError("test", "comment")
}

func (n SlackNotify) Finish(target, comment string) error {
	return n.notifySuccess(fmt.Sprintf("%s Rollout finished successflly! :+1:", target), comment)
}

func (n SlackNotify) Failed(target, comment string) error {
	return n.notifyError(fmt.Sprintf("%s Rollout failed. :cry:", target), comment)
}

func (n SlackNotify) Start(target, comment string) error {
	return n.notifySuccess(fmt.Sprintf("%s Rollout started! :rocket:", target), comment)
}

func (n SlackNotify) notifySuccess(message, comment string) error {
	return n.notify(message, comment, "#009900")
}

func (n SlackNotify) notifyError(message, comment string) error {
	return n.notify(message, comment, "#c62828")
}

func (n SlackNotify) notify(message, comment, color string) error {
	var attachment slack.Attachment
	if comment != "" {
		attachment = slack.Attachment{
			Title: message,
			Text:  comment,
			Color: color,
		}
	} else {
		attachment = slack.Attachment{
			Title: message,
			Color: color,
		}

	}

	_, _, err := n.c.PostMessage(
		n.config.ChannelID,

		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	return err
}
