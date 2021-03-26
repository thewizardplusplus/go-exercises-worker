package runners

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

// DummySolutionRunner ...
type DummySolutionRunner struct{}

// RunSolution ...
func (DummySolutionRunner) RunSolution(
	solution entities.Solution,
) (entities.Solution, error) {
	type result struct{ Timestamp time.Time }

	resultAsJSON, err := json.Marshal(result{Timestamp: time.Now()})
	if err != nil {
		return entities.Solution{},
			errors.Wrap(err, "unable to marshal the solution result")
	}

	updatedSolution := entities.Solution{
		ID:        solution.ID,
		IsCorrect: true,
		Result:    string(resultAsJSON),
	}
	return updatedSolution, nil
}
