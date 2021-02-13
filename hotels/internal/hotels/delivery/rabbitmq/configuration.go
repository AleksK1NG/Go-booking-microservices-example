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

	HotelsExchange = "hotels"

	UpdateImageQueue       = "update_hotel_image"
	UpdateImageBindingKey  = "update_hotel_image_key"
	UpdateImageWorkers     = 5
	UpdateImageConsumerTag = "update_hotel_image_consumer"
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_hotels_incoming_messages_total",
		Help: "The total number of incoming RabbitMQ messages",
	})
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_hotels_success_messages_total",
		Help: "The total number of success incoming success RabbitMQ messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_hotels_error_messages_total",
		Help: "The total number of error incoming success RabbitMQ messages",
	})
)

// Initialize consumers
func (c *hotelsConsumer) Initialize() error {
	if err := c.Dial(); err != nil {
		return errors.Wrap(err, "Consumer Dial")
	}

	updateImageChan, err := c.CreateExchangeAndQueue(HotelsExchange, UpdateImageQueue, UpdateImageBindingKey)
	if err != nil {
		return errors.Wrap(err, "CreateExchangeAndQueue")
	}

	c.channels = append(c.channels, updateImageChan)

	return nil
}

// CloseChannels close active channels
func (c *hotelsConsumer) CloseChannels() {
	for _, channel := range c.channels {
		go func(ch *amqp.Channel) {
			if err := ch.Close(); err != nil {
				c.logger.Errorf("CloseChannels ch.Close error: %v", err)
			}
		}(channel)
	}
}
