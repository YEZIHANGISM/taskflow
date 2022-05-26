package taskflow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTask struct {
	mock.Mock
	Pause int
}

func (t *MockTask) Execute(ctx context.Context, params *TaskParam) error {
	time.Sleep(time.Second * time.Duration(t.Pause))
	args := t.Called(ctx, params)
	return args.Error(0)
}

func (t *MockTask) Rollback(ctx context.Context, params *TaskParam) error {
	args := t.Called(ctx, params)
	return args.Error(0)
}

func (t *MockTask) GetName() string {
	args := t.Called()
	return args.String(0)
}

func (t *MockTask) NeedRollback() bool {
	args := t.Called()
	return args.Bool(0)
}

func TestExecuteSuccess(t *testing.T) {
	task1 := &MockTask{}
	task2 := &MockTask{}
	task3 := &MockTask{}
	params := &TaskParam{}

	ctx := context.Background()

	task1.On("Execute", ctx, params).Return(nil)
	task2.On("Execute", ctx, params).Return(nil)
	task3.On("Execute", ctx, params).Return(nil)

	tasks := &[]Task{task1, task2, task3}
	flow := New(tasks)
	flow.Run()
	assert.Nil(t, flow.err)
	assert.True(t, flow.Done())
}

func TestFirstExecuteFailed(t *testing.T) {
	task1 := &MockTask{}
	task2 := &MockTask{}
	task3 := &MockTask{}
	params := &TaskParam{}

	ctx := context.Background()
	err := errors.New("fail")

	task1.On("Execute", ctx, params).Return(err)
	task1.On("NeedRollback").Return(true)
	task1.On("GetName").Return("")
	task1.On("Rollback", ctx, params).Return(nil)

	tasks := &[]Task{task1, task2, task3}
	flow := New(tasks)
	flow.Run()
	assert.NotNil(t, flow.err)
	assert.True(t, flow.Done())
}

func TestLastExecuteFailed(t *testing.T) {
	task1 := &MockTask{}
	task2 := &MockTask{}
	task3 := &MockTask{}
	params := &TaskParam{}

	ctx := context.Background()
	err := errors.New("fail")

	task1.On("Execute", ctx, params).Return(nil)
	task1.On("NeedRollback").Return(true)
	task1.On("GetName").Return("")
	task1.On("Rollback", ctx, params).Return(nil)

	task2.On("Execute", ctx, params).Return(nil)
	task2.On("NeedRollback").Return(true)
	task2.On("GetName").Return("")
	task2.On("Rollback", ctx, params).Return(nil)

	task3.On("Execute", ctx, params).Return(err)
	task3.On("NeedRollback").Return(false)
	task3.On("GetName").Return("")

	tasks := &[]Task{task1, task2, task3}
	flow := New(tasks)
	flow.Run()
	assert.NotNil(t, flow.err)
	assert.True(t, flow.Done())
}

func TestCancelSecondTaskSuccess(t *testing.T) {
	task1 := &MockTask{}
	task2 := &MockTask{Pause: 2}
	task3 := &MockTask{}
	param := &TaskParam{}

	ctx := context.Background()
	err := errors.New("fail")

	task1.On("Execute", ctx, param).Return(nil)
	task1.On("NeedRollback").Return(true)
	task1.On("GetName").Return("")
	task1.On("Rollback", ctx, param).Return(nil)

	task2.On("Execute", ctx, param).Return(err)
	task2.On("NeedRollback").Return(true)
	task2.On("GetName").Return("")
	task2.On("Rollback", ctx, param).Return(nil)

	tasks := &[]Task{task1, task2, task3}
	flow := New(tasks)
	go flow.Run()
	time.Sleep(time.Second)
	flow.Cancel()
	for {
		if flow.Done() {
			assert.NotNil(t, flow.err)
			assert.True(t, flow.Done())
			break
		}
	}
}
