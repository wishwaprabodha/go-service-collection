package handler

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/utils"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/utils/amqp"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/utils/amqp/rmq"
	"log"
)

type UserHandler struct {
	Config *rmq.Config
}
type Notification struct {
	EventId          string `json:"event_id" structs:"event_id" mapstructure:"event_id"`
	NotificationId   string `json:"notification_id" structs:"notification_id" mapstructure:"notification_id"`
	StartTimestampMs int64  `json:"start_timestamp_ms" structs:"start_timestamp_ms" mapstructure:"start_timestamp_ms"`
	UserId           string `json:"user_id" structs:"user_id" mapstructure:"user_id"`
	Email            string `json:"email" structs:"email" mapstructure:"email"`
	Mobile           string `json:"mobile" structs:"mobile" mapstructure:"mobile"`
}

func (u *UserHandler) HandleEvent(ctx context.Context, message amqp.Message) error {
	m := make(map[string]interface{})
	decoder := utils.JSONDecoder(string(message.Body()))
	decoder.Decode(&m)
	//	log2.Info(ctx)
	notification := Notification{}
	mapstructure.Decode(m, &notification)
	eventId := message.GetID()
	log.Println("event received: ", eventId)
	messageBody := message.Body()
	fmt.Println(messageBody)
	return nil
}

func (u *UserHandler) HandleDeadMessage(ctx context.Context, message amqp.Message, err error) {
}

func (u *UserHandler) Retry(ctx context.Context, message amqp.Message, err error) error {
	err = message.Ack(false)
	if err != nil {
		return err
	}
	return nil
}
