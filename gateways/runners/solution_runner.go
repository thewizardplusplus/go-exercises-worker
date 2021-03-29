package runners

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-log/log"
	"github.com/pkg/errors"
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
	Logger         log.Logger
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
		runner.Logger.
			Log(errors.Wrapf(err, "[error] unable to compile solution #%d", solution.ID))
		updatedSolution.Result = ErrFailedCompiling{ErrMessage: err.Error()}

		return updatedSolution, nil
	}

	if err := testrunner.RunCode(
		pathToExecutable,
		solution.Task.TestCases,
	); err != nil {
		runner.Logger.
			Log(errors.Wrapf(err, "[error] unable to run solution #%d", solution.ID))
		// the error is already wrapped in the testrunner.RunCode() function
		updatedSolution.Result = err

		return updatedSolution, nil
	}

	updatedSolution.IsCorrect = true
	return updatedSolution, nil
}
