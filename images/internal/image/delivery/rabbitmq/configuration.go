package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/streadway/amqp"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false

	ImagesExchange = "images"

	ResizeQueueName   = "resize_queue"
	ResizeConsumerTag = "resize_consumer"
	ResizeWorkers     = 10
	ResizeBindingKey  = "resize_image_key"

	CreateQueueName   = "create_queue"
	CreateConsumerTag = "create_consumer"
	CreateWorkers     = 5
	CreateBindingKey  = "create_image_key"

	UploadHotelImageQueue       = "upload_hotel_image_queue"
	UploadHotelImageConsumerTag = "upload_hotel_image_consumer_tag"
	UploadHotelImageWorkers     = 10
	UploadHotelImageBindingKey  = "upload_hotel_image_binding_key"
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_images_incoming_messages_total",
		Help: "The total number of incoming RabbitMQ messages",
	})
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_images_success_messages_total",
		Help: "The total number of success incoming success RabbitMQ messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_images_error_messages_total",
		Help: "The total number of error incoming success RabbitMQ messages",
	})
)

// Initialize consumers
func (c *ImageConsumer) Initialize() error {
	if err := c.Dial(); err != nil {
		return errors.Wrap(err, "Consumer Dial")
	}

	updateImageChan, err := c.CreateExchangeAndQueue(ImagesExchange, UploadHotelImageQueue, UploadHotelImageBindingKey)
	if err != nil {
		return errors.Wrap(err, "CreateExchangeAndQueue")
	}
	c.channels = append(c.channels, updateImageChan)

	resizeChan, err := c.CreateExchangeAndQueue(ImagesExchange, ResizeQueueName, ResizeBindingKey)
	if err != nil {
		return errors.Wrap(err, "CreateExchangeAndQueue")
	}
	c.channels = append(c.channels, resizeChan)

	createImgChan, err := c.CreateExchangeAndQueue(ImagesExchange, CreateQueueName, CreateBindingKey)
	if err != nil {
		return errors.Wrap(err, "CreateExchangeAndQueue")
	}
	c.channels = append(c.channels, createImgChan)

	return nil
}

// CloseChannels close active channels
func (c *ImageConsumer) CloseChannels() {
	for _, channel := range c.channels {
		go func(ch *amqp.Channel) {
			if err := ch.Close(); err != nil {
				c.logger.Errorf("CloseChannels ch.Close error: %v", err)
			}
		}(channel)
	}
}
