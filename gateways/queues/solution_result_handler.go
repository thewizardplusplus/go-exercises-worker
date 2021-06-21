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

// SolutionHandler ...
type SolutionHandler struct {
	SolutionResultQueueName string
	SolutionRunner          SolutionRunner
	Client                  rabbitmqutils.Client
}

// MessageType ...
func (handler SolutionHandler) MessageType() reflect.Type {
	return reflect.TypeOf(entities.Solution{})
}

// HandleMessage ...
func (handler SolutionHandler) HandleMessage(message interface{}) error {
	solution, err := handler.SolutionRunner.
		RunSolution(message.(entities.Solution))
	if err != nil {
		return errors.Wrap(err, "unable to run the solution")
	}

	err = handler.Client.
		PublishMessage(handler.SolutionResultQueueName, "", solution)
	if err != nil {
		return errors.Wrap(err, "unable to publish the solution")
	}

	return nil
}
