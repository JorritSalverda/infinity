// Code generated by MockGen. DO NOT EDIT.
// Source: random_string_generator.go

// Package lib is a generated GoMock package.
package lib

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRandomStringGenerator is a mock of RandomStringGenerator interface.
type MockRandomStringGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockRandomStringGeneratorMockRecorder
}

// MockRandomStringGeneratorMockRecorder is the mock recorder for MockRandomStringGenerator.
type MockRandomStringGeneratorMockRecorder struct {
	mock *MockRandomStringGenerator
}

// NewMockRandomStringGenerator creates a new mock instance.
func NewMockRandomStringGenerator(ctrl *gomock.Controller) *MockRandomStringGenerator {
	mock := &MockRandomStringGenerator{ctrl: ctrl}
	mock.recorder = &MockRandomStringGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRandomStringGenerator) EXPECT() *MockRandomStringGeneratorMockRecorder {
	return m.recorder
}

// GenerateRandomString mocks base method.
func (m *MockRandomStringGenerator) GenerateRandomString(n int) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRandomString", n)
	ret0, _ := ret[0].(string)
	return ret0
}

// GenerateRandomString indicates an expected call of GenerateRandomString.
func (mr *MockRandomStringGeneratorMockRecorder) GenerateRandomString(n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRandomString", reflect.TypeOf((*MockRandomStringGenerator)(nil).GenerateRandomString), n)
}
