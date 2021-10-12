// Code generated by MockGen. DO NOT EDIT.
// Source: command_runner.go

// Package lib is a generated GoMock package.
package lib

import (
	context "context"
	log "log"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCommandRunner is a mock of CommandRunner interface.
type MockCommandRunner struct {
	ctrl     *gomock.Controller
	recorder *MockCommandRunnerMockRecorder
}

// MockCommandRunnerMockRecorder is the mock recorder for MockCommandRunner.
type MockCommandRunnerMockRecorder struct {
	mock *MockCommandRunner
}

// NewMockCommandRunner creates a new mock instance.
func NewMockCommandRunner(ctrl *gomock.Controller) *MockCommandRunner {
	mock := &MockCommandRunner{ctrl: ctrl}
	mock.recorder = &MockCommandRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandRunner) EXPECT() *MockCommandRunnerMockRecorder {
	return m.recorder
}

// RunCommand mocks base method.
func (m *MockCommandRunner) RunCommand(ctx context.Context, logger *log.Logger, dir, command string, args []string, env ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, logger, dir, command, args}
	for _, a := range env {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunCommand", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunCommand indicates an expected call of RunCommand.
func (mr *MockCommandRunnerMockRecorder) RunCommand(ctx, logger, dir, command, args interface{}, env ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, logger, dir, command, args}, env...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCommand", reflect.TypeOf((*MockCommandRunner)(nil).RunCommand), varargs...)
}

// RunCommandWithOutput mocks base method.
func (m *MockCommandRunner) RunCommandWithOutput(ctx context.Context, logger *log.Logger, dir, command string, args []string, env ...string) ([]byte, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, logger, dir, command, args}
	for _, a := range env {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunCommandWithOutput", varargs...)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunCommandWithOutput indicates an expected call of RunCommandWithOutput.
func (mr *MockCommandRunnerMockRecorder) RunCommandWithOutput(ctx, logger, dir, command, args interface{}, env ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, logger, dir, command, args}, env...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCommandWithOutput", reflect.TypeOf((*MockCommandRunner)(nil).RunCommandWithOutput), varargs...)
}
