// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock_session is a generated GoMock package.
package mock_session

import (
	session "basicLoginRest/internal/session"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Refresh mocks base method.
func (m *MockManager) Refresh(ctx context.Context, oldSid string) (session.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, oldSid)
	ret0, _ := ret[0].(session.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockManagerMockRecorder) Refresh(ctx, oldSid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockManager)(nil).Refresh), ctx, oldSid)
}

// Start mocks base method.
func (m *MockManager) Start(ctx context.Context, sid string) (session.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx, sid)
	ret0, _ := ret[0].(session.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Start indicates an expected call of Start.
func (mr *MockManagerMockRecorder) Start(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockManager)(nil).Start), ctx, sid)
}