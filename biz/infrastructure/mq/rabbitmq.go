package mq

import (
	"time"

	"github.com/buyandship/bns-golib/config"
	"github.com/buyandship/bns-golib/rabbitmq"
	ServiceConfig "github.com/buyandship/supply-service/biz/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Cli *rabbitmq.Client
)

func Init() {
	cli, err := rabbitmq.NewInsecureTLSClient(config.GlobalAppConfig.RabbitMQ)
	if err != nil {
		hlog.Fatalf("failed to create rabbitmq client: %v", err)
		return
	}

	if err := cli.DeclareQueue(
		ServiceConfig.MessageDelayedQueue,
		ServiceConfig.MessageDelayedExchange,
		ServiceConfig.MessageDelayedRoutingKey,
		amqp.Table{},
		amqp.Table{
			"x-single-active-consumer": true,
		},
	); err != nil {
		hlog.Fatalf("failed to declare queue: %v", err)
		return
	}

	args := amqp.Table{
		"x-message-ttl":             int32(time.Hour / time.Millisecond),
		"x-dead-letter-exchange":    ServiceConfig.MessageDelayedExchange,
		"x-dead-letter-routing-key": ServiceConfig.MessageDelayedRoutingKey,
	}

	if err := cli.DeclareQueue(
		ServiceConfig.OneHourBufferMessageQueue,
		ServiceConfig.OneHourBufferMessageExchange,
		ServiceConfig.OneHourBufferMessageRoutingKey,
		amqp.Table{},
		args,
	); err != nil {
		hlog.Fatalf("failed to declare queue: %v", err)
		return
	}

	Cli = cli
}

type Message struct {
	Topic      string `json:"topic"`
	RoutingKey string `json:"routing_key"`
	Msg        string `json:"msg"`
}

func SendMessage(msg Message) error {
	return Cli.PublishMessage(msg.Topic, msg.RoutingKey, msg.Msg)
}
