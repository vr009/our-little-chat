// Code generated by MockGen. DO NOT EDIT.
// Source: internal/flusher/internal/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=internal/flusher/internal/interfaces.go -destination=internal/mocks/flusher/internal/interfaces.go
//
// Package mock_internal is a generated GoMock package.
package flusher

import (
	models "our-little-chatik/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockQueueRepo is a mock of QueueRepo interface.
type MockQueueRepo struct {
	ctrl     *gomock.Controller
	recorder *MockQueueRepoMockRecorder
}

// MockQueueRepoMockRecorder is the mock recorder for MockQueueRepo.
type MockQueueRepoMockRecorder struct {
	mock *MockQueueRepo
}

// NewMockQueueRepo creates a new mock instance.
func NewMockQueueRepo(ctrl *gomock.Controller) *MockQueueRepo {
	mock := &MockQueueRepo{ctrl: ctrl}
	mock.recorder = &MockQueueRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueueRepo) EXPECT() *MockQueueRepoMockRecorder {
	return m.recorder
}

// FetchAllMessages mocks base method.
func (m *MockQueueRepo) FetchAllMessages() ([]models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchAllMessages")
	ret0, _ := ret[0].([]models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchAllMessages indicates an expected call of FetchAllMessages.
func (mr *MockQueueRepoMockRecorder) FetchAllMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchAllMessages", reflect.TypeOf((*MockQueueRepo)(nil).FetchAllMessages))
}

// MockPersistantRepo is a mock of PersistantRepo interface.
type MockPersistantRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPersistantRepoMockRecorder
}

// MockPersistantRepoMockRecorder is the mock recorder for MockPersistantRepo.
type MockPersistantRepoMockRecorder struct {
	mock *MockPersistantRepo
}

// NewMockPersistantRepo creates a new mock instance.
func NewMockPersistantRepo(ctrl *gomock.Controller) *MockPersistantRepo {
	mock := &MockPersistantRepo{ctrl: ctrl}
	mock.recorder = &MockPersistantRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersistantRepo) EXPECT() *MockPersistantRepoMockRecorder {
	return m.recorder
}

// PersistAllMessages mocks base method.
func (m *MockPersistantRepo) PersistAllMessages(msgs []models.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PersistAllMessages", msgs)
	ret0, _ := ret[0].(error)
	return ret0
}

// PersistAllMessages indicates an expected call of PersistAllMessages.
func (mr *MockPersistantRepoMockRecorder) PersistAllMessages(msgs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PersistAllMessages", reflect.TypeOf((*MockPersistantRepo)(nil).PersistAllMessages), msgs)
}