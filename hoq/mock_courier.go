// Code generated by MockGen. DO NOT EDIT.
// Source: HOQ/hoq (interfaces: Courier)

// Package mock is a generated GoMock package.
package hoq

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCourier is a mock of Courier interface
type MockCourier struct {
	ctrl     *gomock.Controller
	recorder *MockCourierMockRecorder
}

// MockCourierMockRecorder is the mock recorder for MockCourier
type MockCourierMockRecorder struct {
	mock *MockCourier
}

// NewMockCourier creates a new mock instance
func NewMockCourier(ctrl *gomock.Controller) *MockCourier {
	mock := &MockCourier{ctrl: ctrl}
	mock.recorder = &MockCourierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCourier) EXPECT() *MockCourierMockRecorder {
	return m.recorder
}

// RoundTrip mocks base method
func (m *MockCourier) RoundTrip(arg0 *Request) (*Response, *RemoteInfo, error) {
	ret := m.ctrl.Call(m, "RoundTrip", arg0)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(*RemoteInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RoundTrip indicates an expected call of RoundTrip
func (mr *MockCourierMockRecorder) RoundTrip(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoundTrip", reflect.TypeOf((*MockCourier)(nil).RoundTrip), arg0)
}
