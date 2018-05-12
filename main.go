package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

const (
	AppName = "qc"

	optConsume        = "c"
	optConsumeLong    = optConsume + ", consume"
	optContentType    = "content-type"
	optDurable        = "d"
	optDurableLong    = optDurable + ", durable"
	optExchange       = "e"
	optExchangeLong   = optExchange + ", exchange"
	optExType         = "t"
	optExTypeLong     = optExType + ", exchange-type"
	optRoutingKey     = "r"
	optRoutingKeyLong = optRoutingKey + ", routing-key"
	optUrl            = "u"
	optUrlLong        = optUrl + ", broker-url"

	optDefaultContentType = "plain/text"
	optDefaultExchange    = ""
	optDefaultExType      = "direct"
	optDefaultUrl         = "amqp://guest:guest@localhost:5672/"

	optEnvUrl = "QC_BROKER_URL"
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

	exchange := c.String(optExchange)

	exType, err := requireString(c, optExType)
	if err != nil {
		return nil, err
	}

	routingKeys := c.StringSlice(optRoutingKey)

	consume := c.Bool(optConsume)

	var contentType string
	if !consume {
		contentType, err = requireString(c, optContentType)
		if err != nil {
			return nil, err
		}
	}

	durable := c.Bool(optDurable)

	qc := New(
		url,
		exchange,
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
	app.UsageText = AppName + " [options]"
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
			Usage: "declare exchanges and queues as durable",
		},
		cli.StringFlag{
			Name:  optExchangeLong,
			Value: optDefaultExchange,
			Usage: "the name of the `EXCHANGE`",
		},
		cli.StringFlag{
			Name:  optExTypeLong,
			Value: optDefaultExType,
			Usage: "exchange `TYPE` [direct, fanout, headers, topic]",
		},
		cli.StringSliceFlag{
			Name:  optRoutingKeyLong,
			Usage: "bind consumer to, or publish to `ROUTING_KEY`",
		},
		cli.StringFlag{
			Name:   optUrlLong,
			Value:  optDefaultUrl,
			EnvVar: optEnvUrl,
			Usage:  "`URL` of the message queue broker",
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
