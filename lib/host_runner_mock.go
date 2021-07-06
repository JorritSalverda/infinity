// Code generated by MockGen. DO NOT EDIT.
// Source: host_runner.go

// Package lib is a generated GoMock package.
package lib

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	log "log"
	reflect "reflect"
)

// MockHostRunner is a mock of HostRunner interface
type MockHostRunner struct {
	ctrl     *gomock.Controller
	recorder *MockHostRunnerMockRecorder
}

// MockHostRunnerMockRecorder is the mock recorder for MockHostRunner
type MockHostRunnerMockRecorder struct {
	mock *MockHostRunner
}

// NewMockHostRunner creates a new mock instance
func NewMockHostRunner(ctrl *gomock.Controller) *MockHostRunner {
	mock := &MockHostRunner{ctrl: ctrl}
	mock.recorder = &MockHostRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHostRunner) EXPECT() *MockHostRunnerMockRecorder {
	return m.recorder
}

// RunStage mocks base method
func (m *MockHostRunner) RunStage(ctx context.Context, logger *log.Logger, stage ManifestStage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunStage", ctx, logger, stage)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunStage indicates an expected call of RunStage
func (mr *MockHostRunnerMockRecorder) RunStage(ctx, logger, stage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunStage", reflect.TypeOf((*MockHostRunner)(nil).RunStage), ctx, logger, stage)
}
