package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/caarlos0/env"
	"github.com/go-log/log/print"
	"github.com/thewizardplusplus/go-exercises-worker/gateways/queues"
	"github.com/thewizardplusplus/go-exercises-worker/gateways/runners"
)

type options struct {
	AllowedImportConfig string `env:"ALLOWED_IMPORT_CONFIG" envDefault:"./configs/allowed_imports.json"` // nolint: lll
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

	var allowedImports []string
	allowedImportsAsJSON, err := ioutil.ReadFile(options.AllowedImportConfig)
	if err != nil {
		logger.Fatalf("[error] unable to read the allowed import config: %v", err)
	}
	if err := json.Unmarshal(allowedImportsAsJSON, &allowedImports); err != nil {
		logger.Fatalf(
			"[error] unable to unmarshal the allowed import config: %v",
			err,
		)
	}

	messageBrokerClient, err := queues.NewClient(
		options.MessageBroker.Address,
		options.SolutionConsumer.BufferSize,
	)
	if err != nil {
		logger.Fatalf("[error] unable to create the message broker client: %v", err)
	}
	defer func() {
		if err := messageBrokerClient.Close(); err != nil {
			logger.Fatalf("[error] unable to close the message broker client: %v", err)
		}
	}()

	solutionConsumer, err := queues.NewSolutionConsumer(
		messageBrokerClient,
		queues.NewSolutionHandler(
			messageBrokerClient,
			queues.SolutionHandlerDependencies{
				SolutionRunner: runners.SolutionRunner{
					AllowedImports: allowedImports,
					Logger:         print.New(logger),
				},
				Logger: print.New(logger),
			},
		),
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
