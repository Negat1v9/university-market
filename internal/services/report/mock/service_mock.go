// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/report/type.go
//
// Generated by this command:
//
//	mockgen -source=./internal/services/report/type.go -destination=./internal/services/report/mock/service_mock.go -package=report_service_mock
//

// Package report_service_mock is a generated GoMock package.
package report_service_mock

import (
	context "context"
	reflect "reflect"

	reportmodel "github.com/Negat1v9/work-marketplace/model/report"
	gomock "go.uber.org/mock/gomock"
)

// MockReportService is a mock of ReportService interface.
type MockReportService struct {
	ctrl     *gomock.Controller
	recorder *MockReportServiceMockRecorder
	isgomock struct{}
}

// MockReportServiceMockRecorder is the mock recorder for MockReportService.
type MockReportServiceMockRecorder struct {
	mock *MockReportService
}

// NewMockReportService creates a new mock instance.
func NewMockReportService(ctrl *gomock.Controller) *MockReportService {
	mock := &MockReportService{ctrl: ctrl}
	mock.recorder = &MockReportServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReportService) EXPECT() *MockReportServiceMockRecorder {
	return m.recorder
}

// CreateReportOnUser mocks base method.
func (m *MockReportService) CreateReportOnUser(ctx context.Context, workerID string, report *reportmodel.NewReportReq) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReportOnUser", ctx, workerID, report)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReportOnUser indicates an expected call of CreateReportOnUser.
func (mr *MockReportServiceMockRecorder) CreateReportOnUser(ctx, workerID, report any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReportOnUser", reflect.TypeOf((*MockReportService)(nil).CreateReportOnUser), ctx, workerID, report)
}

// CreateReportOnWorker mocks base method.
func (m *MockReportService) CreateReportOnWorker(ctx context.Context, userID string, report *reportmodel.NewReportReq) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReportOnWorker", ctx, userID, report)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReportOnWorker indicates an expected call of CreateReportOnWorker.
func (mr *MockReportServiceMockRecorder) CreateReportOnWorker(ctx, userID, report any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReportOnWorker", reflect.TypeOf((*MockReportService)(nil).CreateReportOnWorker), ctx, userID, report)
}
