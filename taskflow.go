// Package taskflow implements a tasks manage package.
// It defines a type TaskFlow that provides the following features:
// Build a TaskFlow object; Execute tasks by order;
// Cancel task; Pause task; Check current task name, status,
// process etc.; metrics with each task.
package taskflow

import (
	"context"
	"time"
)

var Flag struct{} // Raise flag when TaskFlow failed.

type TaskParam map[string]interface{}

// Each TaskFlow object provide a series of methods to manage tasks
type TaskFlow struct {
	tasks  *[]Task         // a slice of tasks
	params *TaskParam      // params for tasks
	err    error           // an error occurred during task execution
	ctx    context.Context // context message
	cancel chan struct{}   // flag of cancel
	done   bool            // flag of flow execution status
}

// Run execute tasks, return err if taskflow failed.
// The panic from inside the TaskFlow(set flag true) will be hidden,
// panic from build-in will continue to throw out.
func (f *TaskFlow) Run() {
	defer func() {
		switch p := recover(); p {
		case nil:
		case Flag:
		default:
			panic(p)
		}
	}()

	defer func() {
		f.done = true
	}()

	f.run()
}

// Cancel cancel taskflow before it completes, current task will be executed.
// Rollback all tasks that have been executed,
func (f *TaskFlow) Cancel() {
	f.cancel <- struct{}{}
}

// Pause pause taskflow for secends, current task will be executed.
func (f *TaskFlow) Pause(secends int64) error {
	return nil
}

// PauseUntil pause until t, current task will be executed.
func (f *TaskFlow) PauseUntil(t time.Time) error {
	return nil
}

// Done return done flag, no matter flow failed or success.
func (f *TaskFlow) Done() bool {
	return f.done
}

func (f *TaskFlow) run() {
	for _, task := range *f.tasks {
		// registe rollback function
		defer func(ft Task, fp *TaskParam) {
			if p := recover(); p != nil {
				if p != nil && ft.NeedRollback() {
					ft.Rollback(f.ctx, fp)
				}
				defer panic(p)
			}

		}(task, f.params)

		if err := task.Execute(f.ctx, f.params); err != nil {
			f.err = err
			panic(Flag)
		}

		// cancel sign
		select {
		case <-f.cancel:
			panic(Flag)
		default:
		}
	}
}

type Option func(*TaskFlow)

func WithParams(params *TaskParam) Option {
	return func(f *TaskFlow) {
		f.params = params
	}
}

func WithError(err error) Option {
	return func(f *TaskFlow) {
		f.err = err
	}
}

func WithContext(ctx context.Context) Option {
	return func(f *TaskFlow) {
		f.ctx = ctx
	}
}

// New Init a TaskFlow
func New(tasks *[]Task, opts ...Option) *TaskFlow {
	f := &TaskFlow{
		tasks:  tasks,
		params: &TaskParam{},
		err:    nil,
		ctx:    context.Background(),
		cancel: make(chan struct{}, 1),
		done:   false,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}
