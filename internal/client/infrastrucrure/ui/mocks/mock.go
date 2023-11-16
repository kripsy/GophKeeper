// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/client/infrastrucrure/ui/interface.go

// Package mock_ui is a generated GoMock package.
package mock_ui

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	filemanager "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	models "github.com/kripsy/GophKeeper/internal/models"
)

// MockUserInterface is a mock of UserInterface interface.
type MockUserInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserInterfaceMockRecorder
}

// MockUserInterfaceMockRecorder is the mock recorder for MockUserInterface.
type MockUserInterfaceMockRecorder struct {
	mock *MockUserInterface
}

// NewMockUserInterface creates a new mock instance.
func NewMockUserInterface(ctrl *gomock.Controller) *MockUserInterface {
	mock := &MockUserInterface{ctrl: ctrl}
	mock.recorder = &MockUserInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserInterface) EXPECT() *MockUserInterfaceMockRecorder {
	return m.recorder
}

// AddBasicAuth mocks base method.
func (m *MockUserInterface) AddBasicAuth() (filemanager.BasicAuth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBasicAuth")
	ret0, _ := ret[0].(filemanager.BasicAuth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBasicAuth indicates an expected call of AddBasicAuth.
func (mr *MockUserInterfaceMockRecorder) AddBasicAuth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBasicAuth", reflect.TypeOf((*MockUserInterface)(nil).AddBasicAuth))
}

// AddCard mocks base method.
func (m *MockUserInterface) AddCard() (filemanager.CardData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard")
	ret0, _ := ret[0].(filemanager.CardData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCard indicates an expected call of AddCard.
func (mr *MockUserInterfaceMockRecorder) AddCard() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockUserInterface)(nil).AddCard))
}

// AddMetaInfo mocks base method.
func (m *MockUserInterface) AddMetaInfo() (models.DataInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetaInfo")
	ret0, _ := ret[0].(models.DataInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMetaInfo indicates an expected call of AddMetaInfo.
func (mr *MockUserInterfaceMockRecorder) AddMetaInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetaInfo", reflect.TypeOf((*MockUserInterface)(nil).AddMetaInfo))
}

// AddNote mocks base method.
func (m *MockUserInterface) AddNote() (filemanager.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNote")
	ret0, _ := ret[0].(filemanager.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNote indicates an expected call of AddNote.
func (mr *MockUserInterfaceMockRecorder) AddNote() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNote", reflect.TypeOf((*MockUserInterface)(nil).AddNote))
}

// ChooseSecretType mocks base method.
func (m *MockUserInterface) ChooseSecretType() (int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChooseSecretType")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ChooseSecretType indicates an expected call of ChooseSecretType.
func (mr *MockUserInterfaceMockRecorder) ChooseSecretType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChooseSecretType", reflect.TypeOf((*MockUserInterface)(nil).ChooseSecretType))
}

// Clear mocks base method.
func (m *MockUserInterface) Clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear")
}

// Clear indicates an expected call of Clear.
func (mr *MockUserInterfaceMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockUserInterface)(nil).Clear))
}

// DeleteSecret mocks base method.
func (m *MockUserInterface) DeleteSecret(metaData models.MetaData) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", metaData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockUserInterfaceMockRecorder) DeleteSecret(metaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockUserInterface)(nil).DeleteSecret), metaData)
}

// Exit mocks base method.
func (m *MockUserInterface) Exit() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Exit")
}

// Exit indicates an expected call of Exit.
func (mr *MockUserInterfaceMockRecorder) Exit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exit", reflect.TypeOf((*MockUserInterface)(nil).Exit))
}

// GetFilePath mocks base method.
func (m *MockUserInterface) GetFilePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFilePath indicates an expected call of GetFilePath.
func (mr *MockUserInterfaceMockRecorder) GetFilePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilePath", reflect.TypeOf((*MockUserInterface)(nil).GetFilePath))
}

// GetRepeatedPassword mocks base method.
func (m *MockUserInterface) GetRepeatedPassword() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepeatedPassword")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepeatedPassword indicates an expected call of GetRepeatedPassword.
func (mr *MockUserInterfaceMockRecorder) GetRepeatedPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepeatedPassword", reflect.TypeOf((*MockUserInterface)(nil).GetRepeatedPassword))
}

// GetSecret mocks base method.
func (m *MockUserInterface) GetSecret(metaData models.MetaData) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", metaData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockUserInterfaceMockRecorder) GetSecret(metaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockUserInterface)(nil).GetSecret), metaData)
}

// GetUser mocks base method.
func (m *MockUserInterface) GetUser() (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser")
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserInterfaceMockRecorder) GetUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserInterface)(nil).GetUser))
}

// IsLocalStorage mocks base method.
func (m *MockUserInterface) IsSyncStorage() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSyncStorage")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsLocalStorage indicates an expected call of IsLocalStorage.
func (mr *MockUserInterfaceMockRecorder) IsLocalStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSyncStorage", reflect.TypeOf((*MockUserInterface)(nil).IsSyncStorage))
}

// Menu mocks base method.
func (m *MockUserInterface) Menu(isSyncStorage bool) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Menu", isSyncStorage)
	ret0, _ := ret[0].(int)
	return ret0
}

// Menu indicates an expected call of Menu.
func (mr *MockUserInterfaceMockRecorder) Menu(isSyncStorage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Menu", reflect.TypeOf((*MockUserInterface)(nil).Menu), isSyncStorage)
}

// PrintErr mocks base method.
func (m *MockUserInterface) PrintErr(err string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PrintErr", err)
}

// PrintErr indicates an expected call of PrintErr.
func (mr *MockUserInterfaceMockRecorder) PrintErr(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintErr", reflect.TypeOf((*MockUserInterface)(nil).PrintErr), err)
}

// Sync mocks base method.
func (m *MockUserInterface) Sync(stop <-chan struct{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Sync", stop)
}

// Sync indicates an expected call of Sync.
func (mr *MockUserInterfaceMockRecorder) Sync(stop interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockUserInterface)(nil).Sync), stop)
}

// TryAgain mocks base method.
func (m *MockUserInterface) TryAgain() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryAgain")
	ret0, _ := ret[0].(bool)
	return ret0
}

// TryAgain indicates an expected call of TryAgain.
func (mr *MockUserInterfaceMockRecorder) TryAgain() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryAgain", reflect.TypeOf((*MockUserInterface)(nil).TryAgain))
}

// UpdateSecret mocks base method.
func (m *MockUserInterface) UpdateSecret(metaData models.MetaData) (string, int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecret", metaData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// UpdateSecret indicates an expected call of UpdateSecret.
func (mr *MockUserInterfaceMockRecorder) UpdateSecret(metaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecret", reflect.TypeOf((*MockUserInterface)(nil).UpdateSecret), metaData)
}

// UploadFileTo mocks base method.
func (m *MockUserInterface) UploadFileTo(cfgDir string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFileTo", cfgDir)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// UploadFileTo indicates an expected call of UploadFileTo.
func (mr *MockUserInterfaceMockRecorder) UploadFileTo(cfgDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFileTo", reflect.TypeOf((*MockUserInterface)(nil).UploadFileTo), cfgDir)
}

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// GetRepeatedPassword mocks base method.
func (m *MockAuth) GetRepeatedPassword() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepeatedPassword")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepeatedPassword indicates an expected call of GetRepeatedPassword.
func (mr *MockAuthMockRecorder) GetRepeatedPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepeatedPassword", reflect.TypeOf((*MockAuth)(nil).GetRepeatedPassword))
}

// GetUser mocks base method.
func (m *MockAuth) GetUser() (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser")
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAuthMockRecorder) GetUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuth)(nil).GetUser))
}

// IsLocalStorage mocks base method.
func (m *MockAuth) IsSyncStorage() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSyncStorage")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsLocalStorage indicates an expected call of IsLocalStorage.
func (mr *MockAuthMockRecorder) IsLocalStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSyncStorage", reflect.TypeOf((*MockAuth)(nil).IsSyncStorage))
}

// TryAgain mocks base method.
func (m *MockAuth) TryAgain() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryAgain")
	ret0, _ := ret[0].(bool)
	return ret0
}

// TryAgain indicates an expected call of TryAgain.
func (mr *MockAuthMockRecorder) TryAgain() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryAgain", reflect.TypeOf((*MockAuth)(nil).TryAgain))
}

// MockSecretCreator is a mock of SecretCreator interface.
type MockSecretCreator struct {
	ctrl     *gomock.Controller
	recorder *MockSecretCreatorMockRecorder
}

// MockSecretCreatorMockRecorder is the mock recorder for MockSecretCreator.
type MockSecretCreatorMockRecorder struct {
	mock *MockSecretCreator
}

// NewMockSecretCreator creates a new mock instance.
func NewMockSecretCreator(ctrl *gomock.Controller) *MockSecretCreator {
	mock := &MockSecretCreator{ctrl: ctrl}
	mock.recorder = &MockSecretCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretCreator) EXPECT() *MockSecretCreatorMockRecorder {
	return m.recorder
}

// AddBasicAuth mocks base method.
func (m *MockSecretCreator) AddBasicAuth() (filemanager.BasicAuth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBasicAuth")
	ret0, _ := ret[0].(filemanager.BasicAuth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBasicAuth indicates an expected call of AddBasicAuth.
func (mr *MockSecretCreatorMockRecorder) AddBasicAuth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBasicAuth", reflect.TypeOf((*MockSecretCreator)(nil).AddBasicAuth))
}

// AddCard mocks base method.
func (m *MockSecretCreator) AddCard() (filemanager.CardData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard")
	ret0, _ := ret[0].(filemanager.CardData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCard indicates an expected call of AddCard.
func (mr *MockSecretCreatorMockRecorder) AddCard() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockSecretCreator)(nil).AddCard))
}

// AddMetaInfo mocks base method.
func (m *MockSecretCreator) AddMetaInfo() (models.DataInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetaInfo")
	ret0, _ := ret[0].(models.DataInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMetaInfo indicates an expected call of AddMetaInfo.
func (mr *MockSecretCreatorMockRecorder) AddMetaInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetaInfo", reflect.TypeOf((*MockSecretCreator)(nil).AddMetaInfo))
}

// AddNote mocks base method.
func (m *MockSecretCreator) AddNote() (filemanager.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNote")
	ret0, _ := ret[0].(filemanager.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNote indicates an expected call of AddNote.
func (mr *MockSecretCreatorMockRecorder) AddNote() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNote", reflect.TypeOf((*MockSecretCreator)(nil).AddNote))
}

// ChooseSecretType mocks base method.
func (m *MockSecretCreator) ChooseSecretType() (int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChooseSecretType")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ChooseSecretType indicates an expected call of ChooseSecretType.
func (mr *MockSecretCreatorMockRecorder) ChooseSecretType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChooseSecretType", reflect.TypeOf((*MockSecretCreator)(nil).ChooseSecretType))
}

// MockFileDirector is a mock of FileDirector interface.
type MockFileDirector struct {
	ctrl     *gomock.Controller
	recorder *MockFileDirectorMockRecorder
}

// MockFileDirectorMockRecorder is the mock recorder for MockFileDirector.
type MockFileDirectorMockRecorder struct {
	mock *MockFileDirector
}

// NewMockFileDirector creates a new mock instance.
func NewMockFileDirector(ctrl *gomock.Controller) *MockFileDirector {
	mock := &MockFileDirector{ctrl: ctrl}
	mock.recorder = &MockFileDirectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileDirector) EXPECT() *MockFileDirectorMockRecorder {
	return m.recorder
}

// GetFilePath mocks base method.
func (m *MockFileDirector) GetFilePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFilePath indicates an expected call of GetFilePath.
func (mr *MockFileDirectorMockRecorder) GetFilePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilePath", reflect.TypeOf((*MockFileDirector)(nil).GetFilePath))
}

// UploadFileTo mocks base method.
func (m *MockFileDirector) UploadFileTo(cfgDir string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFileTo", cfgDir)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// UploadFileTo indicates an expected call of UploadFileTo.
func (mr *MockFileDirectorMockRecorder) UploadFileTo(cfgDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFileTo", reflect.TypeOf((*MockFileDirector)(nil).UploadFileTo), cfgDir)
}

// MockSecretManager is a mock of SecretManager interface.
type MockSecretManager struct {
	ctrl     *gomock.Controller
	recorder *MockSecretManagerMockRecorder
}

// MockSecretManagerMockRecorder is the mock recorder for MockSecretManager.
type MockSecretManagerMockRecorder struct {
	mock *MockSecretManager
}

// NewMockSecretManager creates a new mock instance.
func NewMockSecretManager(ctrl *gomock.Controller) *MockSecretManager {
	mock := &MockSecretManager{ctrl: ctrl}
	mock.recorder = &MockSecretManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretManager) EXPECT() *MockSecretManagerMockRecorder {
	return m.recorder
}

// AddBasicAuth mocks base method.
func (m *MockSecretManager) AddBasicAuth() (filemanager.BasicAuth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBasicAuth")
	ret0, _ := ret[0].(filemanager.BasicAuth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBasicAuth indicates an expected call of AddBasicAuth.
func (mr *MockSecretManagerMockRecorder) AddBasicAuth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBasicAuth", reflect.TypeOf((*MockSecretManager)(nil).AddBasicAuth))
}

// AddCard mocks base method.
func (m *MockSecretManager) AddCard() (filemanager.CardData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard")
	ret0, _ := ret[0].(filemanager.CardData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCard indicates an expected call of AddCard.
func (mr *MockSecretManagerMockRecorder) AddCard() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockSecretManager)(nil).AddCard))
}

// AddMetaInfo mocks base method.
func (m *MockSecretManager) AddMetaInfo() (models.DataInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetaInfo")
	ret0, _ := ret[0].(models.DataInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMetaInfo indicates an expected call of AddMetaInfo.
func (mr *MockSecretManagerMockRecorder) AddMetaInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetaInfo", reflect.TypeOf((*MockSecretManager)(nil).AddMetaInfo))
}

// AddNote mocks base method.
func (m *MockSecretManager) AddNote() (filemanager.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNote")
	ret0, _ := ret[0].(filemanager.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNote indicates an expected call of AddNote.
func (mr *MockSecretManagerMockRecorder) AddNote() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNote", reflect.TypeOf((*MockSecretManager)(nil).AddNote))
}

// ChooseSecretType mocks base method.
func (m *MockSecretManager) ChooseSecretType() (int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChooseSecretType")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ChooseSecretType indicates an expected call of ChooseSecretType.
func (mr *MockSecretManagerMockRecorder) ChooseSecretType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChooseSecretType", reflect.TypeOf((*MockSecretManager)(nil).ChooseSecretType))
}

// DeleteSecret mocks base method.
func (m *MockSecretManager) DeleteSecret(metaData models.MetaData) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", metaData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretManagerMockRecorder) DeleteSecret(metaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecretManager)(nil).DeleteSecret), metaData)
}

// GetSecret mocks base method.
func (m *MockSecretManager) GetSecret(metaData models.MetaData) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", metaData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretManagerMockRecorder) GetSecret(metaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretManager)(nil).GetSecret), metaData)
}

// UpdateSecret mocks base method.
func (m *MockSecretManager) UpdateSecret(metaData models.MetaData) (string, int, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecret", metaData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// UpdateSecret indicates an expected call of UpdateSecret.
func (mr *MockSecretManagerMockRecorder) UpdateSecret(metaData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecret", reflect.TypeOf((*MockSecretManager)(nil).UpdateSecret), metaData)
}
