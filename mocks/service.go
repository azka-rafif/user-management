// Code generated by MockGen. DO NOT EDIT.
// Source: internals/services/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	dto "github.com/evermos/boilerplate-go/internals/dto"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockExampleServiceContract is a mock of ExampleContract interface
type MockExampleServiceContract struct {
	ctrl     *gomock.Controller
	recorder *MockExampleServiceContractMockRecorder
}

// MockExampleServiceContractMockRecorder is the mock recorder for MockExampleServiceContract
type MockExampleServiceContractMockRecorder struct {
	mock *MockExampleServiceContract
}

// NewMockExampleServiceContract creates a new mock instance
func NewMockExampleServiceContract(ctrl *gomock.Controller) *MockExampleServiceContract {
	mock := &MockExampleServiceContract{ctrl: ctrl}
	mock.recorder = &MockExampleServiceContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExampleServiceContract) EXPECT() *MockExampleServiceContractMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockExampleServiceContract) Get() (dto.Example, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(dto.Example)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockExampleServiceContractMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockExampleServiceContract)(nil).Get))
}
