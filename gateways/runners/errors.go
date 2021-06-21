package runners

import (
	"fmt"
)

// ErrFailedCompiling ...
type ErrFailedCompiling struct {
	ErrMessage string
}

// Error ...
func (err ErrFailedCompiling) Error() string {
	return fmt.Sprintf("failed compiling: %s", err.ErrMessage)
}
