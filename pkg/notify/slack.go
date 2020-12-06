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
	return sn, nil
}

func (n SlackNotify) Finish(target string) error {
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "a",
					Value: "no",
				},
			},
		*/
	}

	_, _, err := n.c.PostMessage(
		n.config.ChannelID,
		slack.MsgOptionText("Some text", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil
	}
	return nil
}

func (n SlackNotify) Start(target string) error {

	return nil
}
