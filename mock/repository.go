// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	domain "github.com/matiasvarela/minesweeper-API/internal/core/domain"
	reflect "reflect"
)

// MockGameRepository is a mock of GameRepository interface
type MockGameRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGameRepositoryMockRecorder
}

// MockGameRepositoryMockRecorder is the mock recorder for MockGameRepository
type MockGameRepositoryMockRecorder struct {
	mock *MockGameRepository
}

// NewMockGameRepository creates a new mock instance
func NewMockGameRepository(ctrl *gomock.Controller) *MockGameRepository {
	mock := &MockGameRepository{ctrl: ctrl}
	mock.recorder = &MockGameRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGameRepository) EXPECT() *MockGameRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockGameRepository) Get(userID, gameID string) (*domain.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", userID, gameID)
	ret0, _ := ret[0].(*domain.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockGameRepositoryMockRecorder) Get(userID, gameID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGameRepository)(nil).Get), userID, gameID)
}

// GetAll mocks base method
func (m *MockGameRepository) GetAll(userID string) ([]domain.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", userID)
	ret0, _ := ret[0].([]domain.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockGameRepositoryMockRecorder) GetAll(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockGameRepository)(nil).GetAll), userID)
}

// Save mocks base method
func (m *MockGameRepository) Save(game domain.Game) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", game)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockGameRepositoryMockRecorder) Save(game interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockGameRepository)(nil).Save), game)
}
