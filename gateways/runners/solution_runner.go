package runners

import (
	"fmt"
	"os"
	"path/filepath"

	coderunner "github.com/thewizardplusplus/go-code-runner"
	testrunner "github.com/thewizardplusplus/go-code-runner/test-runner"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

// ErrFailedCompiling ...
type ErrFailedCompiling struct {
	ErrMessage string
}

// Error ...
func (err ErrFailedCompiling) Error() string {
	return fmt.Sprintf("failed compiling: %s", err.ErrMessage)
}

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
		updatedSolution.Result = ErrFailedCompiling{ErrMessage: err.Error()}
		return updatedSolution, nil
	}

	runningErr := testrunner.RunCode(pathToExecutable, solution.Task.TestCases)
	updatedSolution.IsCorrect = runningErr == nil
	updatedSolution.Result = runningErr

	return updatedSolution, nil
}
