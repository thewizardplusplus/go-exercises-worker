package queues

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		// TODO: Add test cases.
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
