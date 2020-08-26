// Code generated by MockGen. DO NOT EDIT.
// Source: internals/repositories/repositories.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockExampleRepositoriesContract is a mock of ExampleContract interface
type MockExampleRepositoriesContract struct {
	ctrl     *gomock.Controller
	recorder *MockExampleRepositoriesContractMockRecorder
}

// MockExampleRepositoriesContractMockRecorder is the mock recorder for MockExampleRepositoriesContract
type MockExampleRepositoriesContractMockRecorder struct {
	mock *MockExampleRepositoriesContract
}

// NewMockExampleRepositoriesContract creates a new mock instance
func NewMockExampleRepositoriesContract(ctrl *gomock.Controller) *MockExampleRepositoriesContract {
	mock := &MockExampleRepositoriesContract{ctrl: ctrl}
	mock.recorder = &MockExampleRepositoriesContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExampleRepositoriesContract) EXPECT() *MockExampleRepositoriesContractMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockExampleRepositoriesContract) Get() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockExampleRepositoriesContractMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockExampleRepositoriesContract)(nil).Get))
}
