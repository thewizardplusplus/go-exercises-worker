package queues

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	rabbitmqutils "github.com/thewizardplusplus/go-rabbitmq-utils"
)

// MessageHandler ...
type MessageHandler interface {
	HandleMessage(message amqp.Delivery)
}

// SolutionConsumer ...
type SolutionConsumer struct {
	client            rabbitmqutils.Client
	messages          <-chan amqp.Delivery
	stoppingCtx       context.Context
	stoppingCtxCancel context.CancelFunc
	messageHandler    MessageHandler
}

// NewSolutionConsumer ...
func NewSolutionConsumer(
	client rabbitmqutils.Client,
	messageHandler MessageHandler,
) (SolutionConsumer, error) {
	messages, err := client.ConsumeMessages(SolutionQueueName)
	if err != nil {
		return SolutionConsumer{},
			errors.Wrap(err, "unable to start the message consumption")
	}

	stoppingCtx, stoppingCtxCancel := context.WithCancel(context.Background())
	consumer := SolutionConsumer{
		client:            client,
		messages:          messages,
		stoppingCtx:       stoppingCtx,
		stoppingCtxCancel: stoppingCtxCancel,
		messageHandler:    messageHandler,
	}

	return consumer, nil
}

// Start ...
func (consumer SolutionConsumer) Start() {
	for message := range consumer.messages {
		consumer.messageHandler.HandleMessage(message)
	}
}

// StartConcurrently ...
func (consumer SolutionConsumer) StartConcurrently(concurrency int) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer waitGroup.Done()

			consumer.Start()
		}()
	}

	waitGroup.Wait()
	consumer.stoppingCtxCancel()
}

// Stop ...
func (consumer SolutionConsumer) Stop() error {
	if err := consumer.client.CancelConsuming(SolutionQueueName); err != nil {
		return errors.Wrap(err, "unable to cancel the message consumption")
	}

	<-consumer.stoppingCtx.Done()
	return nil
}
