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

	// Initialize delayed queue
	if err := initDelayedQueue(cli); err != nil {
		hlog.Fatalf("failed to initialize delayed queue: %v", err)
		return
	}

	// Initialize retry queue with retry strategy
	if err := initRetryQueue(cli); err != nil {
		hlog.Fatalf("failed to initialize retry queue: %v", err)
		return
	}

	Cli = cli
}

func initDelayedQueue(cli *rabbitmq.Client) error {
	if err := cli.DeclareQueue(
		ServiceConfig.MessageDelayedQueue,
		ServiceConfig.MessageDelayedExchange,
		ServiceConfig.MessageDelayedRoutingKey,
		amqp.Table{},
		amqp.Table{
			"x-single-active-consumer": true,
		},
	); err != nil {
		hlog.Errorf("failed to declare queue: %v", err)
		return err
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
		hlog.Errorf("failed to declare queue: %v", err)
		return err
	}

	return nil
}

func initRetryQueue(cli *rabbitmq.Client) error {
	// Step 1: Declare the main processing queue
	// This queue will receive messages to process
	// Failed messages will be sent to retry queues via dead letter exchange
	if err := cli.DeclareQueue(
		ServiceConfig.RetryQueue,
		ServiceConfig.RetryExchange,
		ServiceConfig.RetryRoutingKey,
		amqp.Table{},
		amqp.Table{
			"x-single-active-consumer": true,
		},
	); err != nil {
		return err
	}

	// Step 2: Declare dead letter queue for messages that exceeded max retries
	if err := cli.DeclareQueue(
		ServiceConfig.RetryDLQ,
		ServiceConfig.RetryDLX,
		ServiceConfig.RetryDLRoutingKey,
		amqp.Table{},
		amqp.Table{},
	); err != nil {
		return err
	}

	// Step 3: Declare retry queues with exponential backoff delays
	// Retry delays: 1min, 5min, 15min, 1hour
	retryDelays := []struct {
		delay      time.Duration
		queue      string
		routingKey string
	}{
		{1 * time.Minute, "supply:v1:retry_queue_1min", "retry_routing_key_1min"},
		{5 * time.Minute, "supply:v1:retry_queue_5min", "retry_routing_key_5min"},
		{15 * time.Minute, "supply:v1:retry_queue_15min", "retry_routing_key_15min"},
		{1 * time.Hour, "supply:v1:retry_queue_1hour", "retry_routing_key_1hour"},
	}

	for _, retry := range retryDelays {
		// Each retry queue has TTL and routes back to main queue via dead letter exchange
		args := amqp.Table{
			"x-message-ttl":             int32(retry.delay / time.Millisecond),
			"x-dead-letter-exchange":    ServiceConfig.RetryExchange,
			"x-dead-letter-routing-key": ServiceConfig.RetryRoutingKey,
		}

		if err := cli.DeclareQueue(
			retry.queue,
			ServiceConfig.RetryExchange,
			retry.routingKey,
			amqp.Table{},
			args,
		); err != nil {
			return err
		}
	}

	return nil
}

type Message struct {
	Exchange   string          `json:"exchange"`
	RoutingKey string          `json:"routing_key"`
	Publishing amqp.Publishing `json:"publishing"`
}

func SendMessage(msg Message) error {
	return Cli.PublishMessage(msg.Exchange, msg.RoutingKey, msg.Publishing)
}
