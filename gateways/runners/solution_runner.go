package runners

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/pkg/errors"
	coderunner "github.com/thewizardplusplus/go-code-runner"
	systemutils "github.com/thewizardplusplus/go-code-runner/system-utils"
	testrunner "github.com/thewizardplusplus/go-code-runner/test-runner"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

// SolutionRunner ...
type SolutionRunner struct {
	AllowedImports mapset.Set
	RunningTimeout time.Duration
	Logger         log.Logger
}

// RunSolution ...
func (runner SolutionRunner) RunSolution(
	solution entities.Solution,
) (entities.Solution, error) {
	pathToCode, err := systemutils.SaveTemporaryText(solution.Code, ".go")
	if err != nil {
		return entities.Solution{}, errors.Wrap(err, "unable to save the solution")
	}
	defer os.RemoveAll(filepath.Dir(pathToCode)) // nolint: errcheck

	ctx, cancel := context.WithTimeout(context.Background(), runner.RunningTimeout)
	defer cancel()

	updatedSolution := entities.Solution{ID: solution.ID}
	pathToExecutable, err :=
		coderunner.CompileCode(ctx, pathToCode, runner.AllowedImports)
	if err != nil {
		runner.Logger.
			Log(errors.Wrapf(err, "[error] unable to compile solution #%d", solution.ID))
		updatedSolution.Result = ErrFailedCompiling{ErrMessage: err.Error()}

		return updatedSolution, nil
	}

	if err := testrunner.RunTestCases(
		ctx,
		solution.Task.TestCases,
		func(ctx context.Context, input string) (output string, err error) {
			return systemutils.RunCommand(ctx, input, pathToExecutable)
		},
	); err != nil {
		runner.Logger.
			Log(errors.Wrapf(err, "[error] unable to run solution #%d", solution.ID))
		// the error is already wrapped in the testrunner.RunCode() function
		updatedSolution.Result = err

		return updatedSolution, nil
	}

	updatedSolution.IsCorrect = true
	updatedSolution.Result = json.RawMessage("{}") // empty JSON object

	return updatedSolution, nil
}
