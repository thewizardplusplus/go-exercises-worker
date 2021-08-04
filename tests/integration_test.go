// +build integration

package tests

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
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
		// TODO: Add test cases.
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
