package entities

// Solution ...
type Solution struct {
	ID        uint
	Task      Task
	Code      string
	IsCorrect bool
	Result    string
}

// Task ...
type Task struct {
	TestCases string
}
