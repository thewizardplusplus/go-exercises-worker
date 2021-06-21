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
func (runner SolutionRunner) RunSolution(solution entities.Solution) (
	entities.Solution,
	error,
) {
	pathToCode, err := systemutils.SaveTemporaryText(solution.Code, ".go")
	if err != nil {
		return entities.Solution{}, errors.Wrap(err, "unable to save the solution")
	}
	defer os.RemoveAll(filepath.Dir(pathToCode)) // nolint: errcheck

	ctx, cancel := context.WithTimeout(context.Background(), runner.RunningTimeout)
	defer cancel()

	pathToExecutable, err :=
		coderunner.CompileCode(ctx, pathToCode, runner.AllowedImports)
	if err != nil {
		runner.Logger.
			Logf("[error] unable to compile solution #%d: %s", solution.ID, err)

		updatedSolution := entities.Solution{
			ID:     solution.ID,
			Result: ErrFailedCompiling{ErrMessage: err.Error()},
		}
		return updatedSolution, nil
	}

	err = testrunner.RunTestCases(
		ctx,
		solution.Task.TestCases,
		func(ctx context.Context, input string) (output string, err error) {
			return systemutils.RunCommand(ctx, input, pathToExecutable)
		},
	)
	if err != nil {
		runner.Logger.Logf("[error] unable to run solution #%d: %s", solution.ID, err)

		updatedSolution := entities.Solution{
			ID:     solution.ID,
			Result: err, // error has already been wrapped in the testrunner package
		}
		return updatedSolution, nil
	}

	updatedSolution := entities.Solution{
		ID:        solution.ID,
		IsCorrect: true,
		Result:    json.RawMessage("{}"), // empty JSON object
	}
	return updatedSolution, nil
}
