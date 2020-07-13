package rmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)
func getConnectionUrl(config RabbitMQConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.VHost)
}

func connect(config RabbitMQConfig, logger *Logger) (*amqp.Connection, connectionState, error) {
	var conn *amqp.Connection
	var err error
	var state connectionState

	conn, err = amqp.Dial(getConnectionUrl(config))

	if err == nil {
		state = connected
	}

	go func() {
		for {
			select {
			case <-conn.NotifyClose(make(chan *amqp.Error)):
				{
					var newErr error
					var newConn *amqp.Connection
					for {
						newConn, newErr = amqp.Dial(getConnectionUrl(config))
						if newErr != nil {
							state = onRetry
							time.Sleep(1 * time.Second)
							logger.Debugln(ConnectionState, state)
						} else {
							state = connected
							conn = newConn
							logger.Debugln(ConnectionState, state)
							return
						}
					}
				}
			}
		}
	}()

	return conn, state, err
}

func (client *RabbitMQClient) Close() error {
	err := client.conn.Close()
	if err == nil {
		client.connectionState = closed
	}
	return err
}

