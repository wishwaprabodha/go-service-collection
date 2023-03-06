package amqp

import (
	"context"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, message Message) error
	Retry(ctx context.Context, message Message, err error) error
	HandleDeadMessage(ctx context.Context, message Message, err error)
}
