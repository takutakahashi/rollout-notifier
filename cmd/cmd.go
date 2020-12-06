package main

import (
	"log"
	"os"

	"github.com/takutakahashi/rollout-notifier/pkg/daemon"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "rollout notifier"
	app.Usage = "notifier rollout on kubernetes"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "namespace",
			Usage: "k8s namespace",
		},
		cli.StringFlag{
			Name:  "webhook",
			Usage: "webhook url",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "notification type",
		},
	}
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) error {
	namespace := c.String("namespace")
	if namespace == "" {
		cli.ShowAppHelp(c)
		return nil
	}
	webhook := c.String("webhook")
	if webhook == "" {
		cli.ShowAppHelp(c)
		return nil
	}
	notifyType := c.String("type")
	if notifyType == "" {
		cli.ShowAppHelp(c)
		return nil
	}
	config, err := daemon.NewConfig(namespace, webhook, notifyType)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	d, err := daemon.NewDaemon(config)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	d.Start()
	return nil
}
