package publisher

import (
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/rabbitmq"
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
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_incoming_rabbitmq_messages_total",
		Help: "The total number of incoming RabbitMQ messages",
	})
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_success_incoming_rabbitmq_messages_total",
		Help: "The total number of success incoming success RabbitMQ messages",
	})
	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_error_incoming_rabbitmq_message_total",
		Help: "The total number of error incoming success RabbitMQ messages",
	})
)

type ImagePublisher struct {
	amqpConn *amqp.Connection
	cfg      *config.Config
	logger   logger.Logger
}

func NewImagePublisher(cfg *config.Config, logger logger.Logger) (*ImagePublisher, error) {
	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		return nil, err
	}
	// amqpChan, err := mqConn.Channel()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "p.amqpConn.Channel")
	// }
	return &ImagePublisher{cfg: cfg, logger: logger, amqpConn: amqpConn}, nil
}

func (p *ImagePublisher) CreateExchangeAndQueue(exchange, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	amqpChan, err := p.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "p.amqpConn.Channel")
	}

	p.logger.Infof("Declaring exchange: %s", exchange)
	if err := amqpChan.ExchangeDeclare(
		exchange,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	); err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := amqpChan.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	p.logger.Infof("Declared queue, binding it to exchange: Queue: %v, messageCount: %v, "+
		"consumerCount: %v, exchange: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchange,
		bindingKey,
	)

	err = amqpChan.QueueBind(
		queue.Name,
		bindingKey,
		exchange,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	p.logger.Infof(
		"Queue bound to exchange, starting to consume from queue, consumerTag: %v",
		consumerTag,
	)
	return amqpChan, nil
}

// Publish message
func (p *ImagePublisher) Publish(body []byte, contentType string) error {
	amqpChan, err := p.amqpConn.Channel()
	if err != nil {
		return errors.Wrap(err, "p.amqpConn.Channel")
	}
	defer amqpChan.Close()

	p.logger.Infof("Publishing message Exchange: %s, RoutingKey: %s", p.cfg.RabbitMQ.Exchange, p.cfg.RabbitMQ.RoutingKey)

	if err := amqpChan.Publish(
		p.cfg.RabbitMQ.Exchange,
		p.cfg.RabbitMQ.RoutingKey,
		publishMandatory,
		publishImmediate,
		amqp.Publishing{
			ContentType:  contentType,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.NewV4().String(),
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		return errors.Wrap(err, "ch.Publish")
	}

	successMessages.Inc()
	return nil
}
