package entities

import (
	testrunner "github.com/thewizardplusplus/go-code-runner/test-runner"
)

// Solution ...
type Solution struct {
	ID        uint
	Task      Task
	Code      string
	IsCorrect bool
	Result    interface{}
}

// Task ...
type Task struct {
	TestCases []testrunner.TestCase
}
