// Code generated by MockGen. DO NOT EDIT.
// Source: docker_runner.go

// Package lib is a generated GoMock package.
package lib

import (
	context "context"
	log "log"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockDockerRunner is a mock of DockerRunner interface.
type MockDockerRunner struct {
	ctrl     *gomock.Controller
	recorder *MockDockerRunnerMockRecorder
}

// MockDockerRunnerMockRecorder is the mock recorder for MockDockerRunner.
type MockDockerRunnerMockRecorder struct {
	mock *MockDockerRunner
}

// NewMockDockerRunner creates a new mock instance.
func NewMockDockerRunner(ctrl *gomock.Controller) *MockDockerRunner {
	mock := &MockDockerRunner{ctrl: ctrl}
	mock.recorder = &MockDockerRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDockerRunner) EXPECT() *MockDockerRunnerMockRecorder {
	return m.recorder
}

// ContainerGetExitCode mocks base method.
func (m *MockDockerRunner) ContainerGetExitCode(ctx context.Context, logger *log.Logger, containerID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerGetExitCode", ctx, logger, containerID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerGetExitCode indicates an expected call of ContainerGetExitCode.
func (mr *MockDockerRunnerMockRecorder) ContainerGetExitCode(ctx, logger, containerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerGetExitCode", reflect.TypeOf((*MockDockerRunner)(nil).ContainerGetExitCode), ctx, logger, containerID)
}

// ContainerLogs mocks base method.
func (m *MockDockerRunner) ContainerLogs(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string, start time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerLogs", ctx, logger, stage, containerID, start)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerLogs indicates an expected call of ContainerLogs.
func (mr *MockDockerRunnerMockRecorder) ContainerLogs(ctx, logger, stage, containerID, start interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerLogs", reflect.TypeOf((*MockDockerRunner)(nil).ContainerLogs), ctx, logger, stage, containerID, start)
}

// ContainerPull mocks base method.
func (m *MockDockerRunner) ContainerPull(ctx context.Context, logger *log.Logger, stage ManifestStage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerPull", ctx, logger, stage)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerPull indicates an expected call of ContainerPull.
func (mr *MockDockerRunnerMockRecorder) ContainerPull(ctx, logger, stage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerPull", reflect.TypeOf((*MockDockerRunner)(nil).ContainerPull), ctx, logger, stage)
}

// ContainerRemove mocks base method.
func (m *MockDockerRunner) ContainerRemove(ctx context.Context, logger *log.Logger, containerID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerRemove", ctx, logger, containerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerRemove indicates an expected call of ContainerRemove.
func (mr *MockDockerRunnerMockRecorder) ContainerRemove(ctx, logger, containerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerRemove", reflect.TypeOf((*MockDockerRunner)(nil).ContainerRemove), ctx, logger, containerID)
}

// ContainerStart mocks base method.
func (m *MockDockerRunner) ContainerStart(ctx context.Context, logger *log.Logger, stage ManifestStage, needsNetwork bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerStart", ctx, logger, stage, needsNetwork)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerStart indicates an expected call of ContainerStart.
func (mr *MockDockerRunnerMockRecorder) ContainerStart(ctx, logger, stage, needsNetwork interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerStart", reflect.TypeOf((*MockDockerRunner)(nil).ContainerStart), ctx, logger, stage, needsNetwork)
}

// ContainerStop mocks base method.
func (m *MockDockerRunner) ContainerStop(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerStop", ctx, logger, stage, containerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerStop indicates an expected call of ContainerStop.
func (mr *MockDockerRunnerMockRecorder) ContainerStop(ctx, logger, stage, containerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerStop", reflect.TypeOf((*MockDockerRunner)(nil).ContainerStop), ctx, logger, stage, containerID)
}

// NeedsNetwork mocks base method.
func (m *MockDockerRunner) NeedsNetwork(stages []*ManifestStage) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NeedsNetwork", stages)
	ret0, _ := ret[0].(bool)
	return ret0
}

// NeedsNetwork indicates an expected call of NeedsNetwork.
func (mr *MockDockerRunnerMockRecorder) NeedsNetwork(stages interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NeedsNetwork", reflect.TypeOf((*MockDockerRunner)(nil).NeedsNetwork), stages)
}

// NetworkCreate mocks base method.
func (m *MockDockerRunner) NetworkCreate(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkCreate", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// NetworkCreate indicates an expected call of NetworkCreate.
func (mr *MockDockerRunnerMockRecorder) NetworkCreate(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkCreate", reflect.TypeOf((*MockDockerRunner)(nil).NetworkCreate), ctx)
}

// NetworkRemove mocks base method.
func (m *MockDockerRunner) NetworkRemove(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkRemove", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// NetworkRemove indicates an expected call of NetworkRemove.
func (mr *MockDockerRunnerMockRecorder) NetworkRemove(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkRemove", reflect.TypeOf((*MockDockerRunner)(nil).NetworkRemove), ctx)
}

// StopRunningContainers mocks base method.
func (m *MockDockerRunner) StopRunningContainers(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopRunningContainers", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopRunningContainers indicates an expected call of StopRunningContainers.
func (mr *MockDockerRunnerMockRecorder) StopRunningContainers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopRunningContainers", reflect.TypeOf((*MockDockerRunner)(nil).StopRunningContainers), ctx)
}
