package yahoo

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	ServiceConfig "github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrastructure/http"
	"github.com/buyandship/supply-service/biz/infrastructure/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Handle message to update the auction item price
func DelayedQueueConsumer() {
	// Use the shared connection from mq package
	if mq.Cli == nil {
		hlog.Fatalf("mq client is not initialized, call mq.Init() first")
		return
	}

	ch, err := mq.Cli.Channel()
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

func RetryQueueConsumer() {
	// Use the shared connection from mq package
	if mq.Cli == nil {
		hlog.Fatalf("mq client is not initialized, call mq.Init() first")
		return
	}

	ch, err := mq.Cli.Channel()
	if err != nil {
		hlog.Fatalf("failed to create channel: %v", err)
		return
	}

	msgs, err := ch.Consume(
		ServiceConfig.RetryQueue,
		"",
		false, // manual ack
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
			// Extract retry attempt from message headers (if present)
			// For first attempt, headers may be empty
			retryAttempt := getRetryAttemptFromHeaders(msg.Headers)

			hlog.Debugf("retry queue body: %s", string(msg.Body))

			orderNumber := getOrderNumberFromHeaders(msg.Headers)
			if orderNumber == "" {
				hlog.Errorf("order number not found in headers")
				if err := msg.Ack(false); err != nil {
					hlog.Errorf("failed to ack message: %v", err)
				}
				continue
			}
			// Process your message
			if err := http.GetNotifier().NotifyBiddingStatus(context.Background(), orderNumber, msg.Body); err != nil {

				retryAttempt++
				// Increment retry attempt
				if msg.Headers == nil {
					msg.Headers = make(amqp.Table)
				}
				msg.Headers["x-retry-count"] = retryAttempt

				publishing := amqp.Publishing{
					Body:    msg.Body,
					Headers: msg.Headers,
				}

				retryMessage := GetRetryMessage(retryAttempt, publishing)

				if err := mq.SendMessage(retryMessage); err != nil {
					hlog.Errorf("failed to send to retry queue: %v", err)
					// Optionally nack and requeue
					msg.Nack(false, true)
					continue
				}

				// Acknowledge the original message (it's now in retry queue)
				if err := msg.Ack(false); err != nil {
					hlog.Errorf("failed to ack message: %v", err)
				}
				continue
			}

			if err := msg.Ack(false); err != nil {
				hlog.Errorf("failed to ack message: %v", err)
			}
		}
	}()
}

func getRetryAttemptFromHeaders(headers amqp.Table) int {
	if headers == nil {
		return 0 // First attempt
	}

	retryCount, exists := headers["x-retry-count"]
	if !exists {
		return 0
	}

	// Try different numeric types that amqp.Table might use
	// RabbitMQ's amqp.Table can store numbers as int, int32, int64, or float64
	switch v := retryCount.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	case float32:
		return int(v)
	default:
		return 0
	}
}

func GetRetryMessage(attempt int, publishing amqp.Publishing) mq.Message {
	switch attempt {
	case 1:
		return mq.Message{
			Exchange:   ServiceConfig.RetryExchange,
			RoutingKey: "retry_routing_key_1min",
			Publishing: publishing,
		}
	case 2:
		return mq.Message{
			Exchange:   ServiceConfig.RetryExchange,
			RoutingKey: "retry_routing_key_5min",
			Publishing: publishing,
		}
	case 3:
		return mq.Message{
			Exchange:   ServiceConfig.RetryExchange,
			RoutingKey: "retry_routing_key_15min",
			Publishing: publishing,
		}
	case 4:
		return mq.Message{
			Exchange:   ServiceConfig.RetryExchange,
			RoutingKey: "retry_routing_key_1hour",
			Publishing: publishing,
		}
	default:
		// After max retries, send to dead letter queue
		return mq.Message{
			Exchange:   ServiceConfig.RetryDLX,
			RoutingKey: ServiceConfig.RetryDLRoutingKey,
			Publishing: publishing,
		}
	}
}

func getOrderNumberFromHeaders(headers amqp.Table) string {
	if headers == nil {
		return ""
	}
	orderNumber, exists := headers["x-order-number"]
	if !exists {
		return ""
	}
	return orderNumber.(string)
}
