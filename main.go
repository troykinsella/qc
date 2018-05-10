package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
)

const (
	AppName = "qc"

	optConsume     = "c"
	optConsumeLong = optConsume + ", consume"
	optContentType = "content-type"
	optDurable     = "d"
	optDurableLong = optDurable + ", durable"
	optExType      = "t"
	optExTypeLong  = optExType + ", exchange-type"
	optUrl         = "u"
	optUrlLong     = optUrl + ", url"

	optDefaultContentType = "plain/text"
	optDefaultExType      = "direct"
	optDefaultUrl         = "amqp://guest:guest@localhost:5672/"
)

var (
	AppVersion = "0.0.0-dev.0"
)

func requireString(c *cli.Context, opt string) (string, error) {
	val := c.String(opt)
	if val == "" {
		return "", fmt.Errorf("option required: %s", opt)
	}
	return val, nil
}

func newQC(c *cli.Context) (*QC, error) {

	url, err := requireString(c, optUrl)
	if err != nil {
		return nil, err
	}

	exType, err := requireString(c, optExType)
	if err != nil {
		return nil, err
	}

	consume := c.Bool(optConsume)

	var contentType string
	if !consume {
		contentType, err = requireString(c, optContentType)
		if err != nil {
			return nil, err
		}
	}

	durable := c.Bool(optDurable)

	args := c.Args()
	if len(args) < 2 {
		return nil, errors.New("must supply arguments: exchange and at least one routing key")
	}

	exName := args[0]
	routingKeys := args[1:]

	qc := New(
		url,
		exName,
		exType,
		routingKeys,
		consume,
		contentType,
		durable,
	)

	return qc, nil
}

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = AppVersion
	app.Usage = "\"queue cat\" - publish and subscribe to AMQP 0.9.1 message queues, like RabbitMQ"
	app.UsageText = AppName + " [options] exchange routing-key [routing-key...]"
	app.Author = "Troy Kinsella (troy.kinsella@startmail.com)"
	app.Action = func(c *cli.Context) error {
		qc, err := newQC(c)
		if err != nil {
			return err
		}

		err = qc.Connect()
		if err != nil {
			return err
		}

		err = qc.Run()
		return err
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  optConsumeLong,
			Usage: "consume queue messages, otherwise publish",
		},
		cli.StringFlag{
			Name:  optContentType,
			Value: optDefaultContentType,
			Usage: "publication `CONTENT_TYPE`",
		},
		cli.BoolFlag{
			Name:  optDurableLong,
			Usage: "if declaring an exchange and/or a queue, make it durable",
		},
		cli.StringFlag{
			Name:  optExTypeLong,
			Value: optDefaultExType,
			Usage: "exchange `TYPE` [direct, fanout, topic]",
		},
		cli.StringFlag{
			Name:  optUrlLong,
			Value: optDefaultUrl,
			Usage: "`URL` of the message queue broker",
		},
	}

	return app
}

func main() {
	app := newCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
