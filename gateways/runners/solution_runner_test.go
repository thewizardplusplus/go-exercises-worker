package runners

import (
	"encoding/json"
	"reflect"
	"regexp"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testrunner "github.com/thewizardplusplus/go-code-runner/test-runner"
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
		{
			name: "success",
			fields: fields{
				AllowedImports: nil,
				RunningTimeout: 10 * time.Second,
				Logger:         new(MockLogger),
			},
			args: args{
				solution: entities.Solution{
					ID: 23,
					Code: `
						package main

						func main() {
							var x, y int
							if _, err := fmt.Scan(&x, &y); err != nil {
								log.Fatal(err)
							}

							fmt.Println(x + y)
						}
					`,
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedSolution: func(test *testing.T, receivedSolution entities.Solution) {
				wantedSolution := entities.Solution{
					ID:        23,
					IsCorrect: true,
					Result:    json.RawMessage("{}"),
				}
				assert.Equal(test, wantedSolution, receivedSolution)
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with code compiling",
			fields: fields{
				AllowedImports: nil,
				RunningTimeout: 10 * time.Second,
				Logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"[error] unable to compile solution #%d: %s",
							uint(23),
							mock.MatchedBy(func(err error) bool { return err != nil }),
						).
						Return()

					return logger
				}(),
			},
			args: args{
				solution: entities.Solution{
					ID: 23,
					Code: `
						package main

						func main() {
							var x, y int
						}
					`,
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedSolution: func(test *testing.T, receivedSolution entities.Solution) {
				assert.Equal(test, uint(23), receivedSolution.ID)
				assert.False(test, receivedSolution.IsCorrect)
				if assert.IsType(test, ErrFailedCompiling{}, receivedSolution.Result) {
					result := receivedSolution.Result.(ErrFailedCompiling)
					assert.NotEmpty(test, result.ErrMessage)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with import checking",
			fields: fields{
				AllowedImports: mapset.NewSet("log"),
				RunningTimeout: 10 * time.Second,
				Logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"[error] unable to compile solution #%d: %s",
							uint(23),
							mock.MatchedBy(func(err error) bool {
								return err.Error() == `failed import checking: disallowed import "fmt"`
							}),
						).
						Return()

					return logger
				}(),
			},
			args: args{
				solution: entities.Solution{
					ID: 23,
					Code: `
						package main

						func main() {
							var x, y int
							if _, err := fmt.Scan(&x, &y); err != nil {
								log.Fatal(err)
							}

							fmt.Println(x + y)
						}
					`,
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedSolution: func(test *testing.T, receivedSolution entities.Solution) {
				assert.Equal(test, uint(23), receivedSolution.ID)
				assert.False(test, receivedSolution.IsCorrect)
				if assert.IsType(test, ErrFailedCompiling{}, receivedSolution.Result) {
					result := receivedSolution.Result.(ErrFailedCompiling)
					assert.NotEmpty(test, result.ErrMessage)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with code running",
			fields: fields{
				AllowedImports: nil,
				RunningTimeout: 10 * time.Second,
				Logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							"[error] unable to run solution #%d: %s",
							uint(23),
							mock.MatchedBy(func(err error) bool { return err != nil }),
						).
						Return()

					return logger
				}(),
			},
			args: args{
				solution: entities.Solution{
					ID: 23,
					Code: `
						package main

						func main() {
							panic("test")
						}
					`,
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedSolution: func(test *testing.T, receivedSolution entities.Solution) {
				assert.Equal(test, uint(23), receivedSolution.ID)
				assert.False(test, receivedSolution.IsCorrect)
				if assert.IsType(
					test,
					testrunner.ErrFailedRunning{},
					receivedSolution.Result,
				) {
					wantedTestCase := testrunner.TestCase{
						Input:          "5 12",
						ExpectedOutput: "17\n",
					}

					result := receivedSolution.Result.(testrunner.ErrFailedRunning)
					assert.Equal(test, wantedTestCase, result.TestCase)
					assert.NotEmpty(test, result.ErrMessage)
				}
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with the running timeout",
			fields: fields{
				AllowedImports: nil,
				RunningTimeout: 100 * time.Millisecond,
				Logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							mock.MatchedBy(func(message string) bool {
								const pattern = `^\[error\] unable to (compile|run) solution #%d: %s$`
								return regexp.MustCompile(pattern).MatchString(message)
							}),
							uint(23),
							mock.MatchedBy(func(err error) bool { return err != nil }),
						).
						Return()

					return logger
				}(),
			},
			args: args{
				solution: entities.Solution{
					ID: 23,
					Code: `
						package main

						func main() {
							for {
							}
						}
					`,
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedSolution: func(test *testing.T, receivedSolution entities.Solution) {
				assert.Equal(test, uint(23), receivedSolution.ID)
				assert.False(test, receivedSolution.IsCorrect)

				resultType := reflect.TypeOf(receivedSolution.Result)
				isErrFailedCompiling := resultType == reflect.TypeOf(ErrFailedCompiling{})
				isErrFailedRunning :=
					resultType == reflect.TypeOf(testrunner.ErrFailedRunning{})
				assert.True(test, isErrFailedCompiling || isErrFailedRunning)
			},
			wantedErr: assert.NoError,
		},
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
