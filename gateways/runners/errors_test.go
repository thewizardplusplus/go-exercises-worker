package runners

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrFailedCompiling_Error(test *testing.T) {
	err := ErrFailedCompiling{
		ErrMessage: "error",
	}

	const wantedErrMessage = "failed compiling: error"
	assert.EqualError(test, err, wantedErrMessage)
}
