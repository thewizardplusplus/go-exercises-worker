package runners

import (
	"time"

	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

// DummySolutionRunner ...
type DummySolutionRunner struct{}

// RunSolution ...
func (DummySolutionRunner) RunSolution(
	solution entities.Solution,
) (entities.Solution, error) {
	type result struct{ Timestamp time.Time }

	updatedSolution := entities.Solution{
		ID:        solution.ID,
		IsCorrect: true,
		Result:    result{Timestamp: time.Now()},
	}
	return updatedSolution, nil
}
