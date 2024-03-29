// Code generated by MockGen. DO NOT EDIT.
// Source: internal/users/internal/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=internal/users/internal/interfaces.go -destination=internal/mocks/users/internal/interfaces.go
//
// Package mock_internal is a generated GoMock package.
package users

import (
	models "our-little-chatik/internal/models"
	models0 "our-little-chatik/internal/users/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepo) CreateUser(user models.User) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepoMockRecorder) CreateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepo)(nil).CreateUser), user)
}

// DeactivateUser mocks base method.
func (m *MockUserRepo) DeactivateUser(user models.User) models.StatusCode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeactivateUser", user)
	ret0, _ := ret[0].(models.StatusCode)
	return ret0
}

// DeactivateUser indicates an expected call of DeactivateUser.
func (mr *MockUserRepoMockRecorder) DeactivateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeactivateUser", reflect.TypeOf((*MockUserRepo)(nil).DeactivateUser), user)
}

// FindUsers mocks base method.
func (m *MockUserRepo) FindUsers(nickname string) ([]models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUsers", nickname)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// FindUsers indicates an expected call of FindUsers.
func (mr *MockUserRepoMockRecorder) FindUsers(nickname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUsers", reflect.TypeOf((*MockUserRepo)(nil).FindUsers), nickname)
}

// GetUserForItsID mocks base method.
func (m *MockUserRepo) GetUserForItsID(user models.User) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserForItsID", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// GetUserForItsID indicates an expected call of GetUserForItsID.
func (mr *MockUserRepoMockRecorder) GetUserForItsID(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserForItsID", reflect.TypeOf((*MockUserRepo)(nil).GetUserForItsID), user)
}

// GetUserForItsNickname mocks base method.
func (m *MockUserRepo) GetUserForItsNickname(user models.User) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserForItsNickname", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// GetUserForItsNickname indicates an expected call of GetUserForItsNickname.
func (mr *MockUserRepoMockRecorder) GetUserForItsNickname(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserForItsNickname", reflect.TypeOf((*MockUserRepo)(nil).GetUserForItsNickname), user)
}

// UpdateUser mocks base method.
func (m *MockUserRepo) UpdateUser(user models.User) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepoMockRecorder) UpdateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepo)(nil).UpdateUser), user)
}

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// DeactivateUser mocks base method.
func (m *MockUserUsecase) DeactivateUser(user models.User) models.StatusCode {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeactivateUser", user)
	ret0, _ := ret[0].(models.StatusCode)
	return ret0
}

// DeactivateUser indicates an expected call of DeactivateUser.
func (mr *MockUserUsecaseMockRecorder) DeactivateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeactivateUser", reflect.TypeOf((*MockUserUsecase)(nil).DeactivateUser), user)
}

// FindUsers mocks base method.
func (m *MockUserUsecase) FindUsers(nickname string) ([]models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUsers", nickname)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// FindUsers indicates an expected call of FindUsers.
func (mr *MockUserUsecaseMockRecorder) FindUsers(nickname any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUsers", reflect.TypeOf((*MockUserUsecase)(nil).FindUsers), nickname)
}

// GetUser mocks base method.
func (m *MockUserUsecase) GetUser(request models0.GetUserRequest) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", request)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserUsecaseMockRecorder) GetUser(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserUsecase)(nil).GetUser), request)
}

// Login mocks base method.
func (m *MockUserUsecase) Login(request models0.LoginRequest) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", request)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserUsecaseMockRecorder) Login(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserUsecase)(nil).Login), request)
}

// SignUp mocks base method.
func (m *MockUserUsecase) SignUp(request models0.SignUpPersonRequest) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", request)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserUsecaseMockRecorder) SignUp(request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserUsecase)(nil).SignUp), request)
}

// UpdateUser mocks base method.
func (m *MockUserUsecase) UpdateUser(userToUpdate models.User, request models0.UpdateUserRequest) (models.User, models.StatusCode) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", userToUpdate, request)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.StatusCode)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserUsecaseMockRecorder) UpdateUser(userToUpdate, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserUsecase)(nil).UpdateUser), userToUpdate, request)
}
