package entities

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
	TestCases []TestCase
}

// TestCase ...
type TestCase struct {
	Input          string
	ExpectedOutput string
}
