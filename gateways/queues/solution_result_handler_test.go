package queues

import (
	"encoding/json"
	"reflect"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testrunner "github.com/thewizardplusplus/go-code-runner/test-runner"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

func TestSolutionHandler_MessageType(test *testing.T) {
	var handler SolutionHandler
	receivedMessageType := handler.MessageType()

	wantedMessageType := reflect.TypeOf(entities.Solution{})
	assert.Equal(test, wantedMessageType, receivedMessageType)
}

func TestSolutionHandler_HandleMessage(test *testing.T) {
	type fields struct {
		SolutionResultQueueName string
		SolutionRunner          SolutionRunner
		MessagePublisher        MessagePublisher
	}
	type args struct {
		message interface{}
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantedErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				SolutionResultQueueName: "one",
				SolutionRunner: func() SolutionRunner {
					solution := entities.Solution{
						ID: 23,
						Code: "package main\n" +
							"\n" +
							"func main() {\n" +
							"\tvar x, y int\n" +
							"\tif _, err := fmt.Scan(&x, &y); err != nil {\n" +
							"\t\tlog.Fatal(err)\n" +
							"\t}\n" +
							"\n" +
							"\tfmt.Println(x + y)\n" +
							"}\n",
						Task: entities.Task{
							TestCases: []testrunner.TestCase{
								{Input: "5 12", ExpectedOutput: "17\n"},
								{Input: "23 42", ExpectedOutput: "65\n"},
							},
						},
					}
					solutionResult := entities.Solution{
						ID:        23,
						IsCorrect: true,
						Result:    json.RawMessage("{}"),
					}

					solutionRunner := new(MockSolutionRunner)
					solutionRunner.On("RunSolution", solution).Return(solutionResult, nil)

					return solutionRunner
				}(),
				MessagePublisher: func() MessagePublisher {
					solutionResult := entities.Solution{
						ID:        23,
						IsCorrect: true,
						Result:    json.RawMessage("{}"),
					}

					messagePublisher := new(MockMessagePublisher)
					messagePublisher.
						On("PublishMessage", "one", "solution-23-result", solutionResult).
						Return(nil)

					return messagePublisher
				}(),
			},
			args: args{
				message: entities.Solution{
					ID: 23,
					Code: "package main\n" +
						"\n" +
						"func main() {\n" +
						"\tvar x, y int\n" +
						"\tif _, err := fmt.Scan(&x, &y); err != nil {\n" +
						"\t\tlog.Fatal(err)\n" +
						"\t}\n" +
						"\n" +
						"\tfmt.Println(x + y)\n" +
						"}\n",
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedErr: assert.NoError,
		},
		{
			name: "error with solution running",
			fields: fields{
				SolutionResultQueueName: "one",
				SolutionRunner: func() SolutionRunner {
					solution := entities.Solution{
						ID: 23,
						Code: "package main\n" +
							"\n" +
							"func main() {\n" +
							"\tvar x, y int\n" +
							"\tif _, err := fmt.Scan(&x, &y); err != nil {\n" +
							"\t\tlog.Fatal(err)\n" +
							"\t}\n" +
							"\n" +
							"\tfmt.Println(x + y)\n" +
							"}\n",
						Task: entities.Task{
							TestCases: []testrunner.TestCase{
								{Input: "5 12", ExpectedOutput: "17\n"},
								{Input: "23 42", ExpectedOutput: "65\n"},
							},
						},
					}

					solutionRunner := new(MockSolutionRunner)
					solutionRunner.
						On("RunSolution", solution).
						Return(entities.Solution{}, iotest.ErrTimeout)

					return solutionRunner
				}(),
				MessagePublisher: new(MockMessagePublisher),
			},
			args: args{
				message: entities.Solution{
					ID: 23,
					Code: "package main\n" +
						"\n" +
						"func main() {\n" +
						"\tvar x, y int\n" +
						"\tif _, err := fmt.Scan(&x, &y); err != nil {\n" +
						"\t\tlog.Fatal(err)\n" +
						"\t}\n" +
						"\n" +
						"\tfmt.Println(x + y)\n" +
						"}\n",
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedErr: assert.Error,
		},
		{
			name: "error with solution result publishing",
			fields: fields{
				SolutionResultQueueName: "one",
				SolutionRunner: func() SolutionRunner {
					solution := entities.Solution{
						ID: 23,
						Code: "package main\n" +
							"\n" +
							"func main() {\n" +
							"\tvar x, y int\n" +
							"\tif _, err := fmt.Scan(&x, &y); err != nil {\n" +
							"\t\tlog.Fatal(err)\n" +
							"\t}\n" +
							"\n" +
							"\tfmt.Println(x + y)\n" +
							"}\n",
						Task: entities.Task{
							TestCases: []testrunner.TestCase{
								{Input: "5 12", ExpectedOutput: "17\n"},
								{Input: "23 42", ExpectedOutput: "65\n"},
							},
						},
					}
					solutionResult := entities.Solution{
						ID:        23,
						IsCorrect: true,
						Result:    json.RawMessage("{}"),
					}

					solutionRunner := new(MockSolutionRunner)
					solutionRunner.On("RunSolution", solution).Return(solutionResult, nil)

					return solutionRunner
				}(),
				MessagePublisher: func() MessagePublisher {
					solutionResult := entities.Solution{
						ID:        23,
						IsCorrect: true,
						Result:    json.RawMessage("{}"),
					}

					messagePublisher := new(MockMessagePublisher)
					messagePublisher.
						On("PublishMessage", "one", "solution-23-result", solutionResult).
						Return(iotest.ErrTimeout)

					return messagePublisher
				}(),
			},
			args: args{
				message: entities.Solution{
					ID: 23,
					Code: "package main\n" +
						"\n" +
						"func main() {\n" +
						"\tvar x, y int\n" +
						"\tif _, err := fmt.Scan(&x, &y); err != nil {\n" +
						"\t\tlog.Fatal(err)\n" +
						"\t}\n" +
						"\n" +
						"\tfmt.Println(x + y)\n" +
						"}\n",
					Task: entities.Task{
						TestCases: []testrunner.TestCase{
							{Input: "5 12", ExpectedOutput: "17\n"},
							{Input: "23 42", ExpectedOutput: "65\n"},
						},
					},
				},
			},
			wantedErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := SolutionHandler{
				SolutionResultQueueName: data.fields.SolutionResultQueueName,
				SolutionRunner:          data.fields.SolutionRunner,
				MessagePublisher:        data.fields.MessagePublisher,
			}
			receivedErr := handler.HandleMessage(data.args.message)

			mock.AssertExpectationsForObjects(
				test,
				data.fields.SolutionRunner,
				data.fields.MessagePublisher,
			)
			data.wantedErr(test, receivedErr)
		})
	}
}
