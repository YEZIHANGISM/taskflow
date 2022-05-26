package taskflow

import "context"

// A Task is an interface that defines task-related functions
type Task interface {
	// Execute execute current task
	Execute(context.Context, *TaskParam) error
	// Rollback rollback current task
	Rollback(context.Context, *TaskParam) error
	// NeedRollback mark whether the current task needed to be rolled back
	NeedRollback() bool
	// GetName get task name
	GetName() string
}
