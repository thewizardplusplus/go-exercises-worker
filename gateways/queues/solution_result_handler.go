package queues

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
	rabbitmqutils "github.com/thewizardplusplus/go-rabbitmq-utils"
)

// SolutionRunner ...
type SolutionRunner interface {
	RunSolution(solution entities.Solution) (entities.Solution, error)
}

// SolutionHandlerDependencies ...
type SolutionHandlerDependencies struct {
	SolutionRunner SolutionRunner
}

// SolutionHandler ...
type SolutionHandler struct {
	client       rabbitmqutils.Client
	dependencies SolutionHandlerDependencies
}

// NewSolutionHandler ...
func NewSolutionHandler(
	client rabbitmqutils.Client,
	dependencies SolutionHandlerDependencies,
) SolutionHandler {
	return SolutionHandler{client: client, dependencies: dependencies}
}

// HandleMessage ...
func (handler SolutionHandler) HandleMessage(message amqp.Delivery) error {
	var solution entities.Solution
	if err := json.Unmarshal(message.Body, &solution); err != nil {
		return errors.Wrap(err, "unable to unmarshal the solution")
	}

	solution, err := handler.dependencies.SolutionRunner.RunSolution(solution)
	if err != nil {
		return errors.Wrap(err, "unable to run the solution")
	}

	err = handler.client.PublishMessage(SolutionResultQueueName, "", solution)
	if err != nil {
		return errors.Wrap(err, "unable to publish the solution")
	}

	return nil
}
