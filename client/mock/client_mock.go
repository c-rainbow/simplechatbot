// Code generated by MockGen. DO NOT EDIT.
// Source: client/client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	go_twitch_irc "github.com/gempir/go-twitch-irc"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTwitchClientT is a mock of TwitchClientT interface
type MockTwitchClientT struct {
	ctrl     *gomock.Controller
	recorder *MockTwitchClientTMockRecorder
}

// MockTwitchClientTMockRecorder is the mock recorder for MockTwitchClientT
type MockTwitchClientTMockRecorder struct {
	mock *MockTwitchClientT
}

// NewMockTwitchClientT creates a new mock instance
func NewMockTwitchClientT(ctrl *gomock.Controller) *MockTwitchClientT {
	mock := &MockTwitchClientT{ctrl: ctrl}
	mock.recorder = &MockTwitchClientTMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTwitchClientT) EXPECT() *MockTwitchClientTMockRecorder {
	return m.recorder
}

// Connect mocks base method
func (m *MockTwitchClientT) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (mr *MockTwitchClientTMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockTwitchClientT)(nil).Connect))
}

// Disconnect mocks base method
func (m *MockTwitchClientT) Disconnect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect
func (mr *MockTwitchClientTMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockTwitchClientT)(nil).Disconnect))
}

// Join mocks base method
func (m *MockTwitchClientT) Join(channel string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Join", channel)
}

// Join indicates an expected call of Join
func (mr *MockTwitchClientTMockRecorder) Join(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Join", reflect.TypeOf((*MockTwitchClientT)(nil).Join), channel)
}

// Depart mocks base method
func (m *MockTwitchClientT) Depart(channel string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Depart", channel)
}

// Depart indicates an expected call of Depart
func (mr *MockTwitchClientTMockRecorder) Depart(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Depart", reflect.TypeOf((*MockTwitchClientT)(nil).Depart), channel)
}

// Say mocks base method
func (m *MockTwitchClientT) Say(channel, text string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Say", channel, text)
}

// Say indicates an expected call of Say
func (mr *MockTwitchClientTMockRecorder) Say(channel, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Say", reflect.TypeOf((*MockTwitchClientT)(nil).Say), channel, text)
}

// OnConnect mocks base method
func (m *MockTwitchClientT) OnConnect(callback func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnConnect", callback)
}

// OnConnect indicates an expected call of OnConnect
func (mr *MockTwitchClientTMockRecorder) OnConnect(callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnConnect", reflect.TypeOf((*MockTwitchClientT)(nil).OnConnect), callback)
}

// OnPrivateMessage mocks base method
func (m *MockTwitchClientT) OnPrivateMessage(callback func(go_twitch_irc.PrivateMessage)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPrivateMessage", callback)
}

// OnPrivateMessage indicates an expected call of OnPrivateMessage
func (mr *MockTwitchClientTMockRecorder) OnPrivateMessage(callback interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPrivateMessage", reflect.TypeOf((*MockTwitchClientT)(nil).OnPrivateMessage), callback)
}
