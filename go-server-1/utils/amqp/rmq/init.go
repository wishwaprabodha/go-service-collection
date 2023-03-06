package rmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/service/handler"
)

type QueueConfig struct {
	rmqConfig   *Config
	userHandler handler.UserHandler
}

func (q *QueueConfig) InitializeConnection(url string) (*amqp.Connection, error) {
	rmqConnection, err := q.rmqConfig.Connection(url)
	if err != nil {
		return nil, err
	}
	return rmqConnection, nil
}
