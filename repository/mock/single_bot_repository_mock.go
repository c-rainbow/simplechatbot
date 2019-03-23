// Code generated by MockGen. DO NOT EDIT.
// Source: repository\single_bot_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	models "github.com/c-rainbow/simplechatbot/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSingleBotRepositoryT is a mock of SingleBotRepositoryT interface
type MockSingleBotRepositoryT struct {
	ctrl     *gomock.Controller
	recorder *MockSingleBotRepositoryTMockRecorder
}

// MockSingleBotRepositoryTMockRecorder is the mock recorder for MockSingleBotRepositoryT
type MockSingleBotRepositoryTMockRecorder struct {
	mock *MockSingleBotRepositoryT
}

// NewMockSingleBotRepositoryT creates a new mock instance
func NewMockSingleBotRepositoryT(ctrl *gomock.Controller) *MockSingleBotRepositoryT {
	mock := &MockSingleBotRepositoryT{ctrl: ctrl}
	mock.recorder = &MockSingleBotRepositoryTMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSingleBotRepositoryT) EXPECT() *MockSingleBotRepositoryTMockRecorder {
	return m.recorder
}

// GetBotInfo mocks base method
func (m *MockSingleBotRepositoryT) GetBotInfo() *models.Bot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBotInfo")
	ret0, _ := ret[0].(*models.Bot)
	return ret0
}

// GetBotInfo indicates an expected call of GetBotInfo
func (mr *MockSingleBotRepositoryTMockRecorder) GetBotInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBotInfo", reflect.TypeOf((*MockSingleBotRepositoryT)(nil).GetBotInfo))
}

// GetAllChannels mocks base method
func (m *MockSingleBotRepositoryT) GetAllChannels() []*models.Channel {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChannels")
	ret0, _ := ret[0].([]*models.Channel)
	return ret0
}

// GetAllChannels indicates an expected call of GetAllChannels
func (mr *MockSingleBotRepositoryTMockRecorder) GetAllChannels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChannels", reflect.TypeOf((*MockSingleBotRepositoryT)(nil).GetAllChannels))
}

// GetCommandByChannelAndName mocks base method
func (m *MockSingleBotRepositoryT) GetCommandByChannelAndName(channel, commandName string) *models.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommandByChannelAndName", channel, commandName)
	ret0, _ := ret[0].(*models.Command)
	return ret0
}

// GetCommandByChannelAndName indicates an expected call of GetCommandByChannelAndName
func (mr *MockSingleBotRepositoryTMockRecorder) GetCommandByChannelAndName(channel, commandName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommandByChannelAndName", reflect.TypeOf((*MockSingleBotRepositoryT)(nil).GetCommandByChannelAndName), channel, commandName)
}

// AddCommand mocks base method
func (m *MockSingleBotRepositoryT) AddCommand(channel string, command *models.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCommand", channel, command)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCommand indicates an expected call of AddCommand
func (mr *MockSingleBotRepositoryTMockRecorder) AddCommand(channel, command interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCommand", reflect.TypeOf((*MockSingleBotRepositoryT)(nil).AddCommand), channel, command)
}

// EditCommand mocks base method
func (m *MockSingleBotRepositoryT) EditCommand(channel string, command *models.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditCommand", channel, command)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditCommand indicates an expected call of EditCommand
func (mr *MockSingleBotRepositoryTMockRecorder) EditCommand(channel, command interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCommand", reflect.TypeOf((*MockSingleBotRepositoryT)(nil).EditCommand), channel, command)
}

// DeleteCommand mocks base method
func (m *MockSingleBotRepositoryT) DeleteCommand(channel string, command *models.Command) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCommand", channel, command)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommand indicates an expected call of DeleteCommand
func (mr *MockSingleBotRepositoryTMockRecorder) DeleteCommand(channel, command interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommand", reflect.TypeOf((*MockSingleBotRepositoryT)(nil).DeleteCommand), channel, command)
}

// ListCommands mocks base method
func (m *MockSingleBotRepositoryT) ListCommands(channel string) ([]*models.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCommands", channel)
	ret0, _ := ret[0].([]*models.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCommands indicates an expected call of ListCommands
func (mr *MockSingleBotRepositoryTMockRecorder) ListCommands(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCommands", reflect.TypeOf((*MockSingleBotRepositoryT)(nil).ListCommands), channel)
}