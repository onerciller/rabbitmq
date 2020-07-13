package rmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

func (client *RabbitMQClient) createChannel(prefetchCount int) *RabbitMQClient {
	channel, err := client.conn.Channel()

	if err != nil {
		_ = channel.Close()
		client.logger.Panic("Channel could not created. Terminating...", err)
	}

	err = channel.Qos(prefetchCount, 0, false)

	if err != nil {
		client.logger.Panic("PrefetchCount could not defined. Terminating...", err)
	}

	client.channel = channel

	return client
}

func (client *RabbitMQClient) consume(queueName string, channelName string) <-chan amqp.Delivery {
	delivery, err := client.channel.Consume(queueName, channelName, false, false, false, false, nil)
	if err != nil {
		client.logger.Panic("An error occurred while consuming. Terminating...", err)
	}
	return delivery
}

func (client *RabbitMQClient) publish(exchangeName string, key string, payload interface{}) error {

	bytes, err := json.Marshal(payload)

	if err != nil {
		client.logger.Warn("Failed to publish payload, serialization error", err)
	}

	err = client.channel.Publish(exchangeName, key, false, false, amqp.Publishing{
		Body:        bytes,
		ContentType: applicationJsonContentType,
	})

	client.logger.Debugln("publish: ", payload)

	if err != nil {
		client.logger.Warn("Failed to publish payload", "exchange: ", exchangeName, "key: ", key, "payload: ", payload, "error: ", err)
	}

	return err
}

func (client *RabbitMQClient) prefectCount(prefetchCount int) *RabbitMQClient {
	err := client.channel.Qos(prefetchCount, 0, false)
	if err != nil {
		client.logger.Panic("PrefetchCount could not defined. Terminating...", err)
	}
	client.logger.Debugln("prefetchCount:", prefetchCount)
	return client
}

func (client *RabbitMQClient) declareExchange(name string, kind string) *RabbitMQClient {
	err := client.channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
	if err != nil {
		client.logger.Panic("Exchange could not declared. Terminating....", err)
	}

	client.logger.Debugln("declareExchange:", name, "kind:", kind)

	return client
}

func (client *RabbitMQClient) declareQueue(name string) *RabbitMQClient {
	deadLetterArgs := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": name + ".dead-letter",
	}
	_, err := client.channel.QueueDeclare(name, true, false, false, false, deadLetterArgs)

	if err != nil {
		client.logger.Panic("Queue could not declared. Terminating...", err)
	}

	client.logger.Debugln("declareQueue:", name)

	return client
}

func (client *RabbitMQClient) declareDeadLetterQueue(name string) *RabbitMQClient {
	_, err := client.channel.QueueDeclare(name+".deadLetter", true, false, false, false, nil)

	if err != nil {
		client.logger.Panic("DeadLetter queue could not declared. Terminating...", err)
	}

	return client
}

func (client *RabbitMQClient) bindQueue(queueName string, exchangeName string, routingKey string) *RabbitMQClient {
	err := client.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		client.logger.Panic("DeadLetter queue could not declared. Terminating...", err)
	}

	client.logger.Debugln("declareDeadLetterQueue:", queueName)

	return client
}

