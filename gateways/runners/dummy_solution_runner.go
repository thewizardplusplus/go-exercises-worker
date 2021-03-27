package runners

import (
	"time"

	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

// DummySolutionRunner ...
type DummySolutionRunner struct {
	Logger log.Logger
}

// RunSolution ...
func (runner DummySolutionRunner) RunSolution(
	solution entities.Solution,
) (entities.Solution, error) {
	type result struct{ Timestamp time.Time }

	runner.Logger.Logf("[debug] solution has been received: %+v", solution)

	updatedSolution := entities.Solution{
		ID:        solution.ID,
		IsCorrect: true,
		Result:    result{Timestamp: time.Now()},
	}
	return updatedSolution, nil
}
