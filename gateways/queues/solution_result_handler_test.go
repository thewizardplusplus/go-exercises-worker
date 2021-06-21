package queues

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-exercises-worker/entities"
)

func TestSolutionHandler_MessageType(test *testing.T) {
	var handler SolutionHandler
	receivedMessageType := handler.MessageType()

	wantedMessageType := reflect.TypeOf(entities.Solution{})
	assert.Equal(test, wantedMessageType, receivedMessageType)
}
