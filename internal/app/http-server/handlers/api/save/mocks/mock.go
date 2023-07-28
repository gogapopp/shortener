// Code generated by MockGen. DO NOT EDIT.
// Source: save.go

// Package mock_save is a generated GoMock package.
package mock_save

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockURLSaver is a mock of URLSaver interface.
type MockURLSaver struct {
	ctrl     *gomock.Controller
	recorder *MockURLSaverMockRecorder
}

// MockURLSaverMockRecorder is the mock recorder for MockURLSaver.
type MockURLSaverMockRecorder struct {
	mock *MockURLSaver
}

// NewMockURLSaver creates a new mock instance.
func NewMockURLSaver(ctrl *gomock.Controller) *MockURLSaver {
	mock := &MockURLSaver{ctrl: ctrl}
	mock.recorder = &MockURLSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLSaver) EXPECT() *MockURLSaverMockRecorder {
	return m.recorder
}

// GetShortURL mocks base method.
func (m *MockURLSaver) GetShortURL(longURL string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURL", longURL)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetShortURL indicates an expected call of GetShortURL.
func (mr *MockURLSaverMockRecorder) GetShortURL(longURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURL", reflect.TypeOf((*MockURLSaver)(nil).GetShortURL), longURL)
}

// SaveURL mocks base method.
func (m *MockURLSaver) SaveURL(longURL, shortURL, correlationID, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURL", longURL, shortURL, correlationID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveURL indicates an expected call of SaveURL.
func (mr *MockURLSaverMockRecorder) SaveURL(longURL, shortURL, correlationID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURL", reflect.TypeOf((*MockURLSaver)(nil).SaveURL), longURL, shortURL, correlationID, userID)
}
