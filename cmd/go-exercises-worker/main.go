package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-log/log/print"
	"github.com/thewizardplusplus/go-exercises-worker/gateways/config"
	"github.com/thewizardplusplus/go-exercises-worker/gateways/queues"
	"github.com/thewizardplusplus/go-exercises-worker/gateways/runners"
	rabbitmqutils "github.com/thewizardplusplus/go-rabbitmq-utils"
)

type options struct {
	AllowedImportConfig string        `env:"ALLOWED_IMPORT_CONFIG" envDefault:"./configs/allowed_imports.json"` // nolint: lll
	RunningTimeout      time.Duration `env:"RUNNING_TIMEOUT" envDefault:"10s"`
	MessageBroker       struct {
		Address    string `env:"MESSAGE_BROKER_ADDRESS" envDefault:"amqp://rabbitmq:rabbitmq@localhost:5672"` // nolint: lll
		BufferSize int    `env:"MESSAGE_BROKER_BUFFER_SIZE" envDefault:"1000"`
	}
	SolutionConsumer struct {
		Concurrency int `env:"SOLUTION_CONSUMER_CONCURRENCY" envDefault:"1000"`
	}
}

const (
	solutionQueueName       = "solution_queue"
	solutionResultQueueName = "solution_result_queue"
)

func main() {
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	var options options
	if err := env.Parse(&options); err != nil {
		logger.Fatalf("[error] unable to parse the options: %v", err)
	}

	allowedImports, err := config.LoadAllowedImports(options.AllowedImportConfig)
	if err != nil {
		logger.Fatalf("[error] unable to load the allowed imports: %v", err)
	}

	messageBrokerClient, err := rabbitmqutils.NewClient(
		options.MessageBroker.Address,
		rabbitmqutils.WithMaximalQueueSize(options.MessageBroker.BufferSize),
		rabbitmqutils.WithQueues([]string{
			solutionQueueName,
			solutionResultQueueName,
		}),
	)
	if err != nil {
		logger.Fatalf("[error] unable to create the message broker client: %v", err)
	}
	defer func() {
		if err := messageBrokerClient.Close(); err != nil {
			logger.Fatalf("[error] unable to close the message broker client: %v", err)
		}
	}()

	solutionConsumer, err := rabbitmqutils.NewMessageConsumer(
		messageBrokerClient,
		solutionQueueName,
		rabbitmqutils.Acknowledger{
			MessageHandling: rabbitmqutils.TwiceMessageHandling,
			MessageHandler: rabbitmqutils.JSONMessageHandler{
				MessageHandler: queues.SolutionHandler{
					SolutionResultQueueName: solutionResultQueueName,
					SolutionRunner: runners.SolutionRunner{
						AllowedImports: allowedImports,
						RunningTimeout: options.RunningTimeout,
						Logger:         print.New(logger),
					},
					MessagePublisher: messageBrokerClient,
				},
			},
			Logger: print.New(logger),
		},
	)
	if err != nil {
		logger.Fatalf("[error] unable to create the solution consumer: %v", err)
	}
	go solutionConsumer.StartConcurrently(options.SolutionConsumer.Concurrency)
	defer func() {
		if err := solutionConsumer.Stop(); err != nil {
			logger.Fatalf("[error] unable to stop the solution consumer: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
