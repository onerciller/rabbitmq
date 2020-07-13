package rmq

import "github.com/streadway/amqp"

type Consumer struct {
	QueueName           string
	DeadLetterQueueName string
	ExchangeName        string
	Handler             func(value []byte) error
	Connected           bool
	PrefetchCount       int
	ExchangeKind        string
	RoutingKey          string
	ApplicationName     string
}

func (client *RabbitMQClient) RegisterConsumer(consumer Consumer) {
	client.consumers = append(client.consumers, consumer)
}

func (client *RabbitMQClient) declareExchangeAndBindQueue(consumer Consumer) *RabbitMQClient {
	client.logger.Debugln("=======DeclareExchangeAndBindQueue========")
	client.prefectCount(consumer.PrefetchCount).
		declareExchange(consumer.ExchangeName, consumer.ExchangeKind).
		declareQueue(consumer.QueueName).
		declareDeadLetterQueue(consumer.DeadLetterQueueName).
		bindQueue(consumer.QueueName, consumer.ExchangeName, consumer.RoutingKey)
	client.logger.Debugln("=======DeclareExchangeAndBindQueue========")

	return client
}

func (client *RabbitMQClient) RunConsumers() {

	client.logger.Debugln(ConnectionState, client.connectionState)

	client.createChannel(0)

	for _, consumer := range client.consumers {
		deliveries := client.declareExchangeAndBindQueue(consumer).
			consume(consumer.QueueName, generateChannelName(consumer.ApplicationName))

		for i := 0; i < consumer.PrefetchCount; i++ {
			go client.sendDeliveryToHandler(consumer, deliveries)
		}
	}
}

func (client *RabbitMQClient) sendDeliveryToHandler(consumer Consumer, deliveries <-chan amqp.Delivery) {
	for delivery := range deliveries {
		err := consumer.Handler(delivery.Body)
		client.logger.Debugln(Payload, string(delivery.Body), Exchange, consumer.ExchangeName, Queue, consumer.QueueName)
		if err != nil {
			client.logger.Warn("Failed to handler a consumer on:", consumer.ExchangeName, err)
			_ = delivery.Reject(false)
		} else {
			_ = delivery.Ack(false)
		}
	}
}
