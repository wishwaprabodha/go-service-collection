package rmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqMessage struct {
	delivery amqp.Delivery
}

// GetRabbitMQDelivery Get full message
func (msg *RabbitmqMessage) GetRabbitMQDelivery() amqp.Delivery {
	return msg.delivery
}

// GetID get message ID
func (msg *RabbitmqMessage) GetID() string {
	return msg.delivery.MessageId
}

func (msg *RabbitmqMessage) Body() []byte {
	return msg.delivery.Body
}

func (msg *RabbitmqMessage) Headers() map[string]interface{} {
	return msg.delivery.Headers
}

func (msg *RabbitmqMessage) Ack(flag bool) error {
	if flag {
		err := msg.delivery.Ack(false)
		if err != nil {
			return err
		}
	} else {
		err := msg.delivery.Reject(false)
		if err != nil {
			return err
		}
	}
	return nil
}
