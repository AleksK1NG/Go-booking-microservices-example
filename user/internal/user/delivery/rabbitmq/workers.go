package rabbitmq

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
)

func (c *UserConsumer) imagesWorker(ctx context.Context, wg *sync.WaitGroup, messages <-chan amqp.Delivery) {
	defer wg.Done()

	for delivery := range messages {
		span, ctx := opentracing.StartSpanFromContext(ctx, "ImageConsumer.resizeWorker")

		c.logger.Infof("processDeliveries deliveryTag% v", delivery.DeliveryTag)

		incomingMessages.Inc()

		err := c.userUC.UpdateUploadedAvatar(ctx, delivery)
		if err != nil {
			if err := delivery.Reject(false); err != nil {
				c.logger.Errorf("Err delivery.Reject: %v", err)
			}
			c.logger.Errorf("Failed to process delivery: %v", err)
			errorMessages.Inc()
		} else {
			err = delivery.Ack(false)
			if err != nil {
				c.logger.Errorf("Failed to acknowledge delivery: %v", err)
				errorMessages.Inc()
				continue
			}
			successMessages.Inc()
		}
		span.Finish()
	}

	c.logger.Info("Deliveries channel closed")
}
