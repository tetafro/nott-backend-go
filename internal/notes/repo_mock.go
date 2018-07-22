// Code generated by MockGen. DO NOT EDIT.
// Source: internal/notes/repo.go

// Package notes is a generated GoMock package.
package notes

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepo is a mock of Repo interface
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockRepo) Get(arg0 filter) ([]Note, error) {
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockRepoMockRecorder) Get(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepo)(nil).Get), arg0)
}

// Create mocks base method
func (m *MockRepo) Create(arg0 Note) (Note, error) {
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockRepoMockRecorder) Create(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepo)(nil).Create), arg0)
}

// Update mocks base method
func (m *MockRepo) Update(arg0 Note) (Note, error) {
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockRepoMockRecorder) Update(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepo)(nil).Update), arg0)
}

// Delete mocks base method
func (m *MockRepo) Delete(arg0 Note) error {
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepoMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepo)(nil).Delete), arg0)
}
