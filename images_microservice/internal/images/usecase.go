package images

import (
	"context"

	"github.com/streadway/amqp"
)

type UseCase interface {
	ResizeImage(ctx context.Context, delivery amqp.Delivery) error
	Create(ctx context.Context, delivery amqp.Delivery) error
}
