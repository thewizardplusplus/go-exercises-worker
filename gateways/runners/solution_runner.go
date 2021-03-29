package runners

import (
	"os"
	"path/filepath"

	coderunner "github.com/thewizardplusplus/go-code-runner"
	testrunner "github.com/thewizardplusplus/go-code-runner/test-runner"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

// SolutionRunner ...
type SolutionRunner struct {
	AllowedImports []string
}

// RunSolution ...
func (runner SolutionRunner) RunSolution(
	solution entities.Solution,
) (entities.Solution, error) {
	pathToCode, err := coderunner.SaveTemporaryCode(solution.Code)
	if err != nil {
		return entities.Solution{}, errors.Wrap(err, "unable to save the solution")
	}
	defer os.RemoveAll(filepath.Dir(pathToCode)) // nolint: errcheck

	updatedSolution := entities.Solution{ID: solution.ID}
	pathToExecutable, err :=
		coderunner.CompileCode(pathToCode, runner.AllowedImports)
	if err != nil {
		updatedSolution.Result = testrunner.ErrFailedRunning{ErrMessage: err.Error()}
		return updatedSolution, nil
	}

	runningErr := testrunner.RunCode(pathToExecutable, solution.Task.TestCases)
	updatedSolution.IsCorrect = runningErr == nil
	updatedSolution.Result = runningErr

	return updatedSolution, nil
}
