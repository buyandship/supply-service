package yahoo

import (
	"github.com/buyandship/bns-golib/config"
	"github.com/buyandship/bns-golib/rabbitmq"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	ServiceConfig "github.com/buyandship/supply-service/biz/common/config"
)

func MQWatcher() {
	cli, err := rabbitmq.NewInsecureTLSClient(config.GlobalAppConfig.RabbitMQ)
	if err != nil {
		hlog.Fatalf("failed to create rabbitmq client: %v", err)
		return
	}

	ch, err := cli.Channel()
	if err != nil {
		hlog.Fatalf("failed to create channel: %v", err)
		return
	}

	msgs, err := ch.Consume(
		ServiceConfig.MessageDelayedQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		hlog.Fatalf("failed to consume messages: %v", err)
		return
	}

	go func() {
		for msg := range msgs {
			hlog.Infof("received message: %s", string(msg.Body))
			if err := msg.Ack(false); err != nil {
				hlog.Errorf("failed to ack message: %v", err)
			}
		}
	}()
}
