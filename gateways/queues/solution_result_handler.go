package queues

import (
	"encoding/json"

	"github.com/go-log/log"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

// SolutionRunner ...
type SolutionRunner interface {
	RunSolution(solution entities.Solution) (entities.Solution, error)
}

// SolutionHandlerDependencies ...
type SolutionHandlerDependencies struct {
	SolutionRunner SolutionRunner
	Logger         log.Logger
}

// SolutionHandler ...
type SolutionHandler struct {
	client       Client
	dependencies SolutionHandlerDependencies
}

// NewSolutionHandler ...
func NewSolutionHandler(
	client Client,
	dependencies SolutionHandlerDependencies,
) SolutionHandler {
	return SolutionHandler{client: client, dependencies: dependencies}
}

// HandleMessage ...
func (handler SolutionHandler) HandleMessage(message amqp.Delivery) {
	solution, err := handler.performHandling(message)
	if err != nil {
		err = errors.Wrapf(err, "[error] unable to handle solution #%d", solution.ID)
		handler.dependencies.Logger.Log(err)

		message.Reject(true) // nolint: gosec, errcheck
		return
	}

	handler.dependencies.Logger.
		Logf("[info] solution #%d has been handled", solution.ID)
	message.Ack(false) // nolint: gosec, errcheck
}

func (handler SolutionHandler) performHandling(
	message amqp.Delivery,
) (entities.Solution, error) {
	var solution entities.Solution
	if err := json.Unmarshal(message.Body, &solution); err != nil {
		return solution, errors.Wrap(err, "unable to unmarshal the solution")
	}

	solution, err := handler.dependencies.SolutionRunner.RunSolution(solution)
	if err != nil {
		return solution, errors.Wrap(err, "unable to run the solution")
	}

	solutionAsJSON, err := json.Marshal(solution)
	if err != nil {
		return solution, errors.Wrap(err, "unable to marshal the solution")
	}

	if err := handler.client.channel.Publish(
		"",                      // exchange
		SolutionResultQueueName, // queue name
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        solutionAsJSON,
		},
	); err != nil {
		return solution, errors.Wrap(err, "unable to publish the solution")
	}

	return solution, nil
}
