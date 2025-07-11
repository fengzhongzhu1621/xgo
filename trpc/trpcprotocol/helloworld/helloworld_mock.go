// Code generated by MockGen. DO NOT EDIT.
// Source: helloworld.trpc.go

// Package helloworld is a generated GoMock package.
package helloworld

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	client "trpc.group/trpc-go/trpc-go/client"
)

// MockGreeterService is a mock of GreeterService interface.
type MockGreeterService struct {
	ctrl     *gomock.Controller
	recorder *MockGreeterServiceMockRecorder
}

// MockGreeterServiceMockRecorder is the mock recorder for MockGreeterService.
type MockGreeterServiceMockRecorder struct {
	mock *MockGreeterService
}

// NewMockGreeterService creates a new mock instance.
func NewMockGreeterService(ctrl *gomock.Controller) *MockGreeterService {
	mock := &MockGreeterService{ctrl: ctrl}
	mock.recorder = &MockGreeterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGreeterService) EXPECT() *MockGreeterServiceMockRecorder {
	return m.recorder
}

// SayHello mocks base method.
func (m *MockGreeterService) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SayHello", ctx, req)
	ret0, _ := ret[0].(*HelloReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SayHello indicates an expected call of SayHello.
func (mr *MockGreeterServiceMockRecorder) SayHello(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockGreeterService)(nil).SayHello), ctx, req)
}

// MockGreeterClientProxy is a mock of GreeterClientProxy interface.
type MockGreeterClientProxy struct {
	ctrl     *gomock.Controller
	recorder *MockGreeterClientProxyMockRecorder
}

// MockGreeterClientProxyMockRecorder is the mock recorder for MockGreeterClientProxy.
type MockGreeterClientProxyMockRecorder struct {
	mock *MockGreeterClientProxy
}

// NewMockGreeterClientProxy creates a new mock instance.
func NewMockGreeterClientProxy(ctrl *gomock.Controller) *MockGreeterClientProxy {
	mock := &MockGreeterClientProxy{ctrl: ctrl}
	mock.recorder = &MockGreeterClientProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGreeterClientProxy) EXPECT() *MockGreeterClientProxyMockRecorder {
	return m.recorder
}

// SayHello mocks base method.
func (m *MockGreeterClientProxy) SayHello(ctx context.Context, req *HelloRequest, opts ...client.Option) (*HelloReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, req}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SayHello", varargs...)
	ret0, _ := ret[0].(*HelloReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SayHello indicates an expected call of SayHello.
func (mr *MockGreeterClientProxyMockRecorder) SayHello(ctx, req interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, req}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SayHello", reflect.TypeOf((*MockGreeterClientProxy)(nil).SayHello), varargs...)
}
