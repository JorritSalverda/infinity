// Code generated by MockGen. DO NOT EDIT.
// Source: scaffolder.go

// Package lib is a generated GoMock package.
package lib

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockScaffolder is a mock of Scaffolder interface.
type MockScaffolder struct {
	ctrl     *gomock.Controller
	recorder *MockScaffolderMockRecorder
}

// MockScaffolderMockRecorder is the mock recorder for MockScaffolder.
type MockScaffolderMockRecorder struct {
	mock *MockScaffolder
}

// NewMockScaffolder creates a new mock instance.
func NewMockScaffolder(ctrl *gomock.Controller) *MockScaffolder {
	mock := &MockScaffolder{ctrl: ctrl}
	mock.recorder = &MockScaffolderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScaffolder) EXPECT() *MockScaffolderMockRecorder {
	return m.recorder
}

// Scaffold mocks base method.
func (m *MockScaffolder) Scaffold(ctx context.Context, applicationType ApplicationType, language Language, applicationName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scaffold", ctx, applicationType, language, applicationName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scaffold indicates an expected call of Scaffold.
func (mr *MockScaffolderMockRecorder) Scaffold(ctx, applicationType, language, applicationName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scaffold", reflect.TypeOf((*MockScaffolder)(nil).Scaffold), ctx, applicationType, language, applicationName)
}
