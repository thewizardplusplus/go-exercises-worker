package queues

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

//go:generate mockery --name=SolutionRunner --inpackage --case=underscore --testonly

// SolutionRunner ...
type SolutionRunner interface {
	RunSolution(solution entities.Solution) (entities.Solution, error)
}

//go:generate mockery --name=MessagePublisher --inpackage --case=underscore --testonly

// MessagePublisher ...
type MessagePublisher interface {
	PublishMessage(queue string, messageID string, messageData interface{}) error
}

// SolutionHandler ...
type SolutionHandler struct {
	SolutionResultQueueName string
	SolutionRunner          SolutionRunner
	MessagePublisher        MessagePublisher
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

	messageID := fmt.Sprintf("solution-%d-result", solution.ID)
	err = handler.MessagePublisher.
		PublishMessage(handler.SolutionResultQueueName, messageID, solution)
	if err != nil {
		return errors.Wrap(err, "unable to publish the solution result")
	}

	return nil
}
