// +build integration

package tests

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testrunner "github.com/thewizardplusplus/go-code-runner/test-runner"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
	rabbitmqutils "github.com/thewizardplusplus/go-rabbitmq-utils"
)

func TestIntegration(test *testing.T) {
	const solutionQueueName = "solution_queue"
	const solutionResultQueueName = "solution_result_queue"

	dsn, ok := os.LookupEnv("MESSAGE_BROKER_ADDRESS")
	if !ok {
		dsn = "amqp://rabbitmq:rabbitmq@localhost:5672"
	}

	for _, data := range []struct {
		name                  string
		solution              entities.Solution
		wantedReceivedMessage func(test *testing.T, receivedMessage amqp.Delivery)
		wantedSolutionResult  func(test *testing.T, solutionResult entities.Solution)
	}{
		{
			name: "success",
			solution: entities.Solution{
				ID: 23,
				Code: `
					package main

					func main() {
						var x, y int
						fmt.Scan(&x, &y)

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
			wantedReceivedMessage: func(test *testing.T, receivedMessage amqp.Delivery) {
				assert.Equal(test, "solution-23-result", receivedMessage.MessageId)
				assert.WithinDuration(
					test,
					time.Now(),
					receivedMessage.Timestamp,
					time.Hour,
				)
				assert.Equal(test, "application/json", receivedMessage.ContentType)
			},
			wantedSolutionResult: func(
				test *testing.T,
				solutionResult entities.Solution,
			) {
				assert.Equal(test, uint(23), solutionResult.ID)
				assert.True(test, solutionResult.IsCorrect)
				assert.Equal(test, map[string]interface{}{}, solutionResult.Result)
			},
		},
		{
			name: "error with code compiling",
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
			wantedReceivedMessage: func(test *testing.T, receivedMessage amqp.Delivery) {
				assert.Equal(test, "solution-23-result", receivedMessage.MessageId)
				assert.WithinDuration(
					test,
					time.Now(),
					receivedMessage.Timestamp,
					time.Hour,
				)
				assert.Equal(test, "application/json", receivedMessage.ContentType)
			},
			wantedSolutionResult: func(
				test *testing.T,
				solutionResult entities.Solution,
			) {
				assert.Equal(test, uint(23), solutionResult.ID)
				assert.False(test, solutionResult.IsCorrect)
				if assert.IsType(test, map[string]interface{}{}, solutionResult.Result) {
					data := solutionResult.Result.(map[string]interface{})
					if errMessage := data["ErrMessage"]; assert.IsType(test, "", errMessage) {
						assert.True(
							test,
							strings.HasPrefix(errMessage.(string), "unable to compile the code"),
						)
					}
				}
			},
		},
		{
			name: "error with code running",
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
			wantedReceivedMessage: func(test *testing.T, receivedMessage amqp.Delivery) {
				assert.Equal(test, "solution-23-result", receivedMessage.MessageId)
				assert.WithinDuration(
					test,
					time.Now(),
					receivedMessage.Timestamp,
					time.Hour,
				)
				assert.Equal(test, "application/json", receivedMessage.ContentType)
			},
			wantedSolutionResult: func(
				test *testing.T,
				solutionResult entities.Solution,
			) {
				assert.Equal(test, uint(23), solutionResult.ID)
				assert.False(test, solutionResult.IsCorrect)
				if assert.IsType(test, map[string]interface{}{}, solutionResult.Result) {
					data := solutionResult.Result.(map[string]interface{})
					assert.Equal(test, "5 12", data["Input"])
					assert.Equal(test, "17\n", data["ExpectedOutput"])
					if errMessage := data["ErrMessage"]; assert.IsType(test, "", errMessage) {
						assert.True(
							test,
							strings.HasPrefix(errMessage.(string), "unable to run the command"),
						)
					}
				}
			},
		},
		{
			name: "error with an unexpected output",
			solution: entities.Solution{
				ID: 23,
				Code: `
					package main

					func main() {
						var x, y int
						fmt.Scan(&x, &y)

						fmt.Println(x + y)
					}
				`,
				Task: entities.Task{
					TestCases: []testrunner.TestCase{
						{Input: "5 12", ExpectedOutput: "17\n"},
						{Input: "23 42", ExpectedOutput: "100\n"},
					},
				},
			},
			wantedReceivedMessage: func(test *testing.T, receivedMessage amqp.Delivery) {
				assert.Equal(test, "solution-23-result", receivedMessage.MessageId)
				assert.WithinDuration(
					test,
					time.Now(),
					receivedMessage.Timestamp,
					time.Hour,
				)
				assert.Equal(test, "application/json", receivedMessage.ContentType)
			},
			wantedSolutionResult: func(
				test *testing.T,
				solutionResult entities.Solution,
			) {
				assert.Equal(test, uint(23), solutionResult.ID)
				assert.False(test, solutionResult.IsCorrect)
				if assert.IsType(test, map[string]interface{}{}, solutionResult.Result) {
					data := solutionResult.Result.(map[string]interface{})
					assert.Equal(test, "23 42", data["Input"])
					assert.Equal(test, "100\n", data["ExpectedOutput"])
					assert.Equal(test, "65\n", data["ActualOutput"])
				}
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			// prepare the client
			client, err := rabbitmqutils.NewClient(
				dsn,
				rabbitmqutils.WithQueues([]string{
					solutionQueueName,
					solutionResultQueueName,
				}),
			)
			require.NoError(test, err)
			defer client.Close()

			// publish the solution
			err = client.PublishMessage(solutionQueueName, "", data.solution)
			require.NoError(test, err)

			// receive the message
			receivedMessage, err := client.GetMessage(solutionResultQueueName)
			require.NoError(test, err)
			defer receivedMessage.Ack(false /* multiple */)

			// clean the irrelevant message fields
			receivedMessage = amqp.Delivery{
				MessageId:   receivedMessage.MessageId,
				Timestamp:   receivedMessage.Timestamp,
				ContentType: receivedMessage.ContentType,
				Body:        receivedMessage.Body,
			}

			// unmarshal the solution result
			var solutionResult entities.Solution
			err = json.Unmarshal(receivedMessage.Body, &solutionResult)
			require.NoError(test, err)

			// check the results
			data.wantedReceivedMessage(test, receivedMessage)
			data.wantedSolutionResult(test, solutionResult)
		})
	}
}
