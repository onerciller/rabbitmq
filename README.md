# rabbitmq

RabbitMQ wrapper over amqp exchanges and queues. In memory retries for consuming messages when an error occured. it provides multiple consumers in a single process and create goroutines and consume messages concurrently also retries to connect another node when RabbitMQ node is down or broken connection.

## Installation

``` 
$ go get https://github.com/onerciller/rabbitmq
```


## Usage 

### Consumer

```go 

// Configuration
rabbitMQClient := rabbitmq.NewRabbitMQClient(rabbitmq.RabbitMQConfig{
		Host:     "",
		Port:     15672
		Username: "user",
		Password: "pass",
		VHost:    "",
	})
  
  // New RabbitMQ Instance
	messageBus := rabbitmq.New(rabbitMQClient)
  
  // Register Consumer
	messageBus.RegisterConsumer(rabbitmq.Consumer{
		QueueName:           "example-queue",
		DeadLetterQueueName: "deadletter-queue-name",
		PrefetchCount:       4,
		ExchangeName:        "exchange-name",
		ExchangeKind:        rabbitmq.ExchangeFanout,
		RoutingKey:          "",
		Handler: func(value []byte) error {
			return nil
		},
	})
  
	go messageBus.RunConsumers()


```
