// Package taskflow implements a tasks manage package.
// It defines a type TaskFlow that provides the following features:
// Build a TaskFlow object; Execute tasks by order;
// Cancel task; Pause task; Check current task name, status,
// process etc.; metrics with each task.
package taskflow

import (
	"context"
	"taskflow/task"
	"time"
)

// Each TaskFlow object provide a series of methods to manage tasks
type TaskFlow struct {
	tasks *[]task.Task    // a slice of task
	err   error           // an error occurred during task execution
	ctx   context.Context // context message
}

// Run execute tasks, return err if taskflow failed.
// The panic from inside the TaskFlow(set flag true) will be hidden,
// panic from build-in will continue to throw out.
func (f *TaskFlow) Run() error {
	return nil
}

// Cancel cancel taskflow, current task will be executed.
// Rollback all tasks that have been executed.
func (f *TaskFlow) Cancel() error {
	return nil
}

// Pause pause taskflow for secends, current task will be executed.
func (f *TaskFlow) Pause(secends int64) error {
	return nil
}

// PauseUntil pause until t, current task will be executed.
func (f *TaskFlow) PauseUntil(t time.Time) error {
	return nil
}

// New Init a TaskFlow
func New(tasks *[]task.Task, err error, ctx context.Context) *TaskFlow {
	return &TaskFlow{tasks: tasks, err: err, ctx: ctx}
}
