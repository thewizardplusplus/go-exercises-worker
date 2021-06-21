package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log/print"
	"github.com/thewizardplusplus/go-exercises-worker/gateways/queues"
	"github.com/thewizardplusplus/go-exercises-worker/gateways/runners"
	rabbitmqutils "github.com/thewizardplusplus/go-rabbitmq-utils"
)

type options struct {
	AllowedImportConfig string        `env:"ALLOWED_IMPORT_CONFIG" envDefault:"./configs/allowed_imports.json"` // nolint: lll
	RunningTimeout      time.Duration `env:"RUNNING_TIMEOUT" envDefault:"10s"`
	MessageBroker       struct {
		Address string `env:"MESSAGE_BROKER_ADDRESS" envDefault:"amqp://rabbitmq:rabbitmq@localhost:5672"` // nolint: lll
	}
	SolutionConsumer struct {
		BufferSize  int `env:"SOLUTION_CONSUMER_BUFFER_SIZE" envDefault:"1000"`
		Concurrency int `env:"SOLUTION_CONSUMER_CONCURRENCY" envDefault:"1000"`
	}
}

func main() {
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	var options options
	if err := env.Parse(&options); err != nil {
		logger.Fatalf("[error] unable to parse the options: %v", err)
	}

	allowedImportsAsJSON, err := ioutil.ReadFile(options.AllowedImportConfig)
	if err != nil {
		logger.Fatalf("[error] unable to read the allowed import config: %v", err)
	}
	var allowedImportAsSlice []string
	err = json.Unmarshal(allowedImportsAsJSON, &allowedImportAsSlice)
	if err != nil {
		logger.Fatalf(
			"[error] unable to unmarshal the allowed import config: %v",
			err,
		)
	}
	allowedImports := mapset.NewSet()
	for _, allowedImport := range allowedImportAsSlice {
		allowedImports.Add(allowedImport)
	}

	messageBrokerClient, err := rabbitmqutils.NewClient(
		options.MessageBroker.Address,
		rabbitmqutils.WithMaximalQueueSize(options.SolutionConsumer.BufferSize),
		rabbitmqutils.WithQueues([]string{
			queues.SolutionQueueName,
			queues.SolutionResultQueueName,
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
		queues.SolutionQueueName,
		rabbitmqutils.Acknowledger{
			MessageHandling: rabbitmqutils.TwiceMessageHandling,
			MessageHandler: rabbitmqutils.JSONMessageHandler{
				MessageHandler: queues.SolutionHandler{
					SolutionResultQueueName: queues.SolutionResultQueueName,
					SolutionRunner: runners.SolutionRunner{
						AllowedImports: allowedImports,
						RunningTimeout: options.RunningTimeout,
						Logger:         print.New(logger),
					},
					Client: messageBrokerClient,
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
