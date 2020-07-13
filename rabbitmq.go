package rmq

import (
	"github.com/streadway/amqp"
)

type connectionState string

const (
	ConnectionState = "CONNECTION_STATE:"
	Payload         = "PAYLOAD:"
	Exchange        = "EXCHANGE:"
	Queue           = "QUEUE:"
)

const (
	connected connectionState = "CONNECTED"
	closed    connectionState = "CLOSED"
	onRetry   connectionState = "ON_RETRY"
)

const (
	ExchangeDirect             = "direct"
	ExchangeFanout             = "fanout"
	ExchangeTopic              = "topic"
	ExchangeHeaders            = "headers"
	applicationJsonContentType = "application/json"
)

type RabbitMQConfig struct {
	Host     string
	VHost    string
	Username string
	Password string
	Port     int
	Debug    bool
}

type RabbitMQClient struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	logger          *Logger
	consumers       []Consumer
	connectionState connectionState
}

func NewRabbitMQClient(config RabbitMQConfig) *RabbitMQClient {

	logger := NewLogger("RMQ: ", config.Debug)

	conn, state, err := connect(config, logger)

	if err != nil {
		logger.Panic("Could not establish connection to rabbitMQ")
	}

	return &RabbitMQClient{
		conn: conn,
		connectionState: state,
		logger:          logger,
	}
}

