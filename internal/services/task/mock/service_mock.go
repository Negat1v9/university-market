// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/task/type.go
//
// Generated by this command:
//
//	mockgen -source=./internal/services/task/type.go -destination=./internal/services/task/mock/service_mock.go -package=task_service_mock
//

// Package task_service_mock is a generated GoMock package.
package task_service_mock

import (
	context "context"
	url "net/url"
	reflect "reflect"

	taskmodel "github.com/Negat1v9/work-marketplace/model/taskModel"
	gomock "go.uber.org/mock/gomock"
)

// MockTaskService is a mock of TaskService interface.
type MockTaskService struct {
	ctrl     *gomock.Controller
	recorder *MockTaskServiceMockRecorder
	isgomock struct{}
}

// MockTaskServiceMockRecorder is the mock recorder for MockTaskService.
type MockTaskServiceMockRecorder struct {
	mock *MockTaskService
}

// NewMockTaskService creates a new mock instance.
func NewMockTaskService(ctrl *gomock.Controller) *MockTaskService {
	mock := &MockTaskService{ctrl: ctrl}
	mock.recorder = &MockTaskServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskService) EXPECT() *MockTaskServiceMockRecorder {
	return m.recorder
}

// AttachFiles mocks base method.
func (m *MockTaskService) AttachFiles(ctx context.Context, taskID, fileID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachFiles", ctx, taskID, fileID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AttachFiles indicates an expected call of AttachFiles.
func (mr *MockTaskServiceMockRecorder) AttachFiles(ctx, taskID, fileID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachFiles", reflect.TypeOf((*MockTaskService)(nil).AttachFiles), ctx, taskID, fileID)
}

// CompleteTask mocks base method.
func (m *MockTaskService) CompleteTask(ctx context.Context, taskID, userID string) (*taskmodel.InfoTaskRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteTask", ctx, taskID, userID)
	ret0, _ := ret[0].(*taskmodel.InfoTaskRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompleteTask indicates an expected call of CompleteTask.
func (mr *MockTaskServiceMockRecorder) CompleteTask(ctx, taskID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteTask", reflect.TypeOf((*MockTaskService)(nil).CompleteTask), ctx, taskID, userID)
}

// Create mocks base method.
func (m *MockTaskService) Create(ctx context.Context, userID string, meta *taskmodel.TaskMeta) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userID, meta)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTaskServiceMockRecorder) Create(ctx, userID, meta any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTaskService)(nil).Create), ctx, userID, meta)
}

// DeleteTask mocks base method.
func (m *MockTaskService) DeleteTask(ctx context.Context, taskID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, taskID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskServiceMockRecorder) DeleteTask(ctx, taskID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskService)(nil).DeleteTask), ctx, taskID, userID)
}

// FindOne mocks base method.
func (m *MockTaskService) FindOne(ctx context.Context, userID, taskID string) (*taskmodel.InfoTaskRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", ctx, userID, taskID)
	ret0, _ := ret[0].(*taskmodel.InfoTaskRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockTaskServiceMockRecorder) FindOne(ctx, userID, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockTaskService)(nil).FindOne), ctx, userID, taskID)
}

// FindUserTasks mocks base method.
func (m *MockTaskService) FindUserTasks(ctx context.Context, userID string, v url.Values) ([]taskmodel.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserTasks", ctx, userID, v)
	ret0, _ := ret[0].([]taskmodel.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserTasks indicates an expected call of FindUserTasks.
func (mr *MockTaskServiceMockRecorder) FindUserTasks(ctx, userID, v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserTasks", reflect.TypeOf((*MockTaskService)(nil).FindUserTasks), ctx, userID, v)
}

// PublishTask mocks base method.
func (m *MockTaskService) PublishTask(ctx context.Context, taskID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishTask", ctx, taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishTask indicates an expected call of PublishTask.
func (mr *MockTaskServiceMockRecorder) PublishTask(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishTask", reflect.TypeOf((*MockTaskService)(nil).PublishTask), ctx, taskID)
}

// RaiseTask mocks base method.
func (m *MockTaskService) RaiseTask(ctx context.Context, taskID, userID string) (*taskmodel.InfoTaskRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RaiseTask", ctx, taskID, userID)
	ret0, _ := ret[0].(*taskmodel.InfoTaskRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RaiseTask indicates an expected call of RaiseTask.
func (mr *MockTaskServiceMockRecorder) RaiseTask(ctx, taskID, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RaiseTask", reflect.TypeOf((*MockTaskService)(nil).RaiseTask), ctx, taskID, userID)
}

// SelectWorker mocks base method.
func (m *MockTaskService) SelectWorker(ctx context.Context, taskID, userID, workerID string) (*taskmodel.InfoTaskRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectWorker", ctx, taskID, userID, workerID)
	ret0, _ := ret[0].(*taskmodel.InfoTaskRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectWorker indicates an expected call of SelectWorker.
func (mr *MockTaskServiceMockRecorder) SelectWorker(ctx, taskID, userID, workerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectWorker", reflect.TypeOf((*MockTaskService)(nil).SelectWorker), ctx, taskID, userID, workerID)
}

// UpdateTaskMeta mocks base method.
func (m *MockTaskService) UpdateTaskMeta(ctx context.Context, taskID, userID string, data *taskmodel.UpdateTaskMeta) (*taskmodel.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskMeta", ctx, taskID, userID, data)
	ret0, _ := ret[0].(*taskmodel.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTaskMeta indicates an expected call of UpdateTaskMeta.
func (mr *MockTaskServiceMockRecorder) UpdateTaskMeta(ctx, taskID, userID, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskMeta", reflect.TypeOf((*MockTaskService)(nil).UpdateTaskMeta), ctx, taskID, userID, data)
}
