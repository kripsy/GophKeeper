// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1 (interfaces: GophKeeperService_BlockStoreServer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	GophKeeper "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	metadata "google.golang.org/grpc/metadata"
)

// MockGophKeeperService_BlockStoreServer is a mock of GophKeeperService_BlockStoreServer interface.
type MockGophKeeperService_BlockStoreServer struct {
	ctrl     *gomock.Controller
	recorder *MockGophKeeperService_BlockStoreServerMockRecorder
}

// MockGophKeeperService_BlockStoreServerMockRecorder is the mock recorder for MockGophKeeperService_BlockStoreServer.
type MockGophKeeperService_BlockStoreServerMockRecorder struct {
	mock *MockGophKeeperService_BlockStoreServer
}

// NewMockGophKeeperService_BlockStoreServer creates a new mock instance.
func NewMockGophKeeperService_BlockStoreServer(ctrl *gomock.Controller) *MockGophKeeperService_BlockStoreServer {
	mock := &MockGophKeeperService_BlockStoreServer{ctrl: ctrl}
	mock.recorder = &MockGophKeeperService_BlockStoreServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGophKeeperService_BlockStoreServer) EXPECT() *MockGophKeeperService_BlockStoreServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).Context))
}

// Recv mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) Recv() (*GophKeeper.BlockStoreRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*GophKeeper.BlockStoreRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) Send(arg0 *GophKeeper.BlockStoreResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockGophKeeperService_BlockStoreServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockGophKeeperService_BlockStoreServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockGophKeeperService_BlockStoreServer)(nil).SetTrailer), arg0)
}
