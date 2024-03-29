// Code generated by MockGen. DO NOT EDIT.
// Source: redirect.go

// Package mock_redirect is a generated GoMock package.
package mock_redirect

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockURLGetter is a mock of URLGetter interface.
type MockURLGetter struct {
	ctrl     *gomock.Controller
	recorder *MockURLGetterMockRecorder
}

// MockURLGetterMockRecorder is the mock recorder for MockURLGetter.
type MockURLGetterMockRecorder struct {
	mock *MockURLGetter
}

// NewMockURLGetter creates a new mock instance.
func NewMockURLGetter(ctrl *gomock.Controller) *MockURLGetter {
	mock := &MockURLGetter{ctrl: ctrl}
	mock.recorder = &MockURLGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLGetter) EXPECT() *MockURLGetterMockRecorder {
	return m.recorder
}

// GetURL mocks base method.
func (m *MockURLGetter) GetURL(shortURL, userID string) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", shortURL, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetURL indicates an expected call of GetURL.
func (mr *MockURLGetterMockRecorder) GetURL(shortURL, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockURLGetter)(nil).GetURL), shortURL, userID)
}
