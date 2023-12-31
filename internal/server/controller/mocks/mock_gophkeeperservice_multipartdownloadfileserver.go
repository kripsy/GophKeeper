// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1 (interfaces: GophKeeperService_MultipartDownloadFileServer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	GophKeeper "github.com/kripsy/GophKeeper/gen/pkg/api/GophKeeper/v1"
	metadata "google.golang.org/grpc/metadata"
)

// MockGophKeeperService_MultipartDownloadFileServer is a mock of GophKeeperService_MultipartDownloadFileServer interface.
type MockGophKeeperService_MultipartDownloadFileServer struct {
	ctrl     *gomock.Controller
	recorder *MockGophKeeperService_MultipartDownloadFileServerMockRecorder
}

// MockGophKeeperService_MultipartDownloadFileServerMockRecorder is the mock recorder for MockGophKeeperService_MultipartDownloadFileServer.
type MockGophKeeperService_MultipartDownloadFileServerMockRecorder struct {
	mock *MockGophKeeperService_MultipartDownloadFileServer
}

// NewMockGophKeeperService_MultipartDownloadFileServer creates a new mock instance.
func NewMockGophKeeperService_MultipartDownloadFileServer(ctrl *gomock.Controller) *MockGophKeeperService_MultipartDownloadFileServer {
	mock := &MockGophKeeperService_MultipartDownloadFileServer{ctrl: ctrl}
	mock.recorder = &MockGophKeeperService_MultipartDownloadFileServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGophKeeperService_MultipartDownloadFileServer) EXPECT() *MockGophKeeperService_MultipartDownloadFileServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockGophKeeperService_MultipartDownloadFileServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockGophKeeperService_MultipartDownloadFileServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockGophKeeperService_MultipartDownloadFileServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m *MockGophKeeperService_MultipartDownloadFileServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockGophKeeperService_MultipartDownloadFileServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockGophKeeperService_MultipartDownloadFileServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockGophKeeperService_MultipartDownloadFileServer) Send(arg0 *GophKeeper.MultipartDownloadFileResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockGophKeeperService_MultipartDownloadFileServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockGophKeeperService_MultipartDownloadFileServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockGophKeeperService_MultipartDownloadFileServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockGophKeeperService_MultipartDownloadFileServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockGophKeeperService_MultipartDownloadFileServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockGophKeeperService_MultipartDownloadFileServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockGophKeeperService_MultipartDownloadFileServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockGophKeeperService_MultipartDownloadFileServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockGophKeeperService_MultipartDownloadFileServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockGophKeeperService_MultipartDownloadFileServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockGophKeeperService_MultipartDownloadFileServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockGophKeeperService_MultipartDownloadFileServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockGophKeeperService_MultipartDownloadFileServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockGophKeeperService_MultipartDownloadFileServer)(nil).SetTrailer), arg0)
}
