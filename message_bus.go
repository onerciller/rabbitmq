package rmq

type MessageBus interface {
	Publish(exchangeName string, key string, payload interface{}, kind string) error
	RegisterConsumer(consumer Consumer)
	RunConsumers()
}

type messageBus struct {
	messageBus MessageBus
}

func New(mbus MessageBus) *messageBus {
	return &messageBus{messageBus: mbus}
}

func (c messageBus) Publish(exchangeName string, key string, payload interface{}, kind string) error {
	return c.messageBus.Publish(exchangeName, key, payload, kind)
}

func (c messageBus) RegisterConsumer(consumer Consumer) {
	 c.messageBus.RegisterConsumer(consumer)
}

func (c messageBus) RunConsumers() {
	c.messageBus.RunConsumers()
}





