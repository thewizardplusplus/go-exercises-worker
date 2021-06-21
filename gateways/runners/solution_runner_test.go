package runners

import (
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

func TestSolutionRunner_RunSolution(test *testing.T) {
	type fields struct {
		AllowedImports mapset.Set
		RunningTimeout time.Duration
		Logger         log.Logger
	}
	type args struct {
		solution entities.Solution
	}

	for _, data := range []struct {
		name           string
		fields         fields
		args           args
		wantedSolution func(test *testing.T, solution entities.Solution)
		wantedErr      assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			runner := SolutionRunner{
				AllowedImports: data.fields.AllowedImports,
				RunningTimeout: data.fields.RunningTimeout,
				Logger:         data.fields.Logger,
			}
			receivedSolution, receivedErr := runner.RunSolution(data.args.solution)

			mock.AssertExpectationsForObjects(test, data.fields.Logger)
			data.wantedSolution(test, receivedSolution)
			data.wantedErr(test, receivedErr)
		})
	}
}
