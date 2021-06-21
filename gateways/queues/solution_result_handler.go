package queues

import (
	"reflect"

	"github.com/pkg/errors"
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

// MessageType ...
func (handler SolutionHandler) MessageType() reflect.Type {
	return reflect.TypeOf(entities.Solution{})
}

// HandleMessage ...
func (handler SolutionHandler) HandleMessage(message interface{}) error {
	solution, err :=
		handler.dependencies.SolutionRunner.RunSolution(message.(entities.Solution))
	if err != nil {
		return errors.Wrap(err, "unable to run the solution")
	}

	err = handler.client.PublishMessage(SolutionResultQueueName, "", solution)
	if err != nil {
		return errors.Wrap(err, "unable to publish the solution")
	}

	return nil
}
