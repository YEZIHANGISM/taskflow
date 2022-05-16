// Package task define the basic unit of TaskFlow
package task

// A Task is an interface that defines task-related functions
type Task interface {
	// Execute execute current task
	Execute() error
	// Rollback rollback current task
	Rollback() error

	// GetName get task name
	GetName() string
}
