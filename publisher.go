package rmq

func (client *RabbitMQClient) Publish(exchangeName string, key string, payload interface{}, kind string) error {
	return client.createChannel(0).declareExchange(exchangeName, kind).publish(exchangeName, key, payload)
}
