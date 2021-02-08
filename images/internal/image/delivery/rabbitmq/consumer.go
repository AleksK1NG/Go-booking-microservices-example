package rabbitmq

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/rabbitmq"
)

type Consumer struct {
	Worker         func(ctx context.Context, wg *sync.WaitGroup, messages <-chan amqp.Delivery)
	WorkerPoolSize int
	QueueName      string
	ConsumerTag    string
}

type ImageConsumer struct {
	amqpConn  *amqp.Connection
	logger    logger.Logger
	cfg       *config.Config
	imageUC   image.UseCase
	consumers []*Consumer
	channels  []*amqp.Channel
}

func NewImageConsumer(logger logger.Logger, cfg *config.Config, imageUC image.UseCase) *ImageConsumer {
	return &ImageConsumer{logger: logger, cfg: cfg, imageUC: imageUC}
}

func (c *ImageConsumer) Dial() error {
	conn, err := rabbitmq.NewRabbitMQConn(c.cfg)
	if err != nil {
		return err
	}
	c.amqpConn = conn
	return nil
}

// Consume messages
func (c *ImageConsumer) CreateExchangeAndQueue(exchangeName, queueName, bindingKey string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	c.logger.Infof("Declaring exchange: %s", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
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

	c.logger.Infof("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,
		bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}

func (c *ImageConsumer) startConsume(
	ctx context.Context,
	worker func(ctx context.Context, wg *sync.WaitGroup, messages <-chan amqp.Delivery),
	workerPoolSize int,
	queueName string,
	consumerTag string,
) error {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return errors.Wrap(err, "c.amqpConn.Channel")
	}

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "ch.Consume")
	}

	wg := &sync.WaitGroup{}

	wg.Add(workerPoolSize)
	for i := 0; i < workerPoolSize; i++ {
		go worker(ctx, wg, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	c.logger.Errorf("ch.NotifyClose: %v", chanErr)

	wg.Wait()

	return chanErr
}

func (c *ImageConsumer) AddConsumer(consumer *Consumer) {
	c.consumers = append(c.consumers, consumer)
}

func (c *ImageConsumer) run(ctx context.Context, cancel context.CancelFunc) {
	for _, cs := range c.consumers {
		go func(consumer *Consumer) {
			if err := c.startConsume(
				ctx,
				consumer.Worker,
				consumer.WorkerPoolSize,
				consumer.QueueName,
				consumer.ConsumerTag,
			); err != nil {
				c.logger.Errorf("StartResizeConsumer: %v", err)
				cancel()
			}
		}(cs)
	}
}

func (c *ImageConsumer) RunConsumers(ctx context.Context, cancel context.CancelFunc) {
	c.AddConsumer(&Consumer{
		Worker:         c.resizeWorker,
		WorkerPoolSize: ResizeWorkers,
		QueueName:      ResizeQueueName,
		ConsumerTag:    ResizeConsumerTag,
	})
	c.AddConsumer(&Consumer{
		Worker:         c.createImageWorker,
		WorkerPoolSize: CreateWorkers,
		QueueName:      CreateQueueName,
		ConsumerTag:    CreateConsumerTag,
	})
	c.run(ctx, cancel)
}
