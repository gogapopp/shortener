// Code generated by MockGen. DO NOT EDIT.
// Source: userurls.go

// Package mock_userurls is a generated GoMock package.
package mock_userurls

import (
	reflect "reflect"

	models "github.com/gogapopp/shortener/internal/app/lib/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserURLsGetter is a mock of UserURLsGetter interface.
type MockUserURLsGetter struct {
	ctrl     *gomock.Controller
	recorder *MockUserURLsGetterMockRecorder
}

// MockUserURLsGetterMockRecorder is the mock recorder for MockUserURLsGetter.
type MockUserURLsGetterMockRecorder struct {
	mock *MockUserURLsGetter
}

// NewMockUserURLsGetter creates a new mock instance.
func NewMockUserURLsGetter(ctrl *gomock.Controller) *MockUserURLsGetter {
	mock := &MockUserURLsGetter{ctrl: ctrl}
	mock.recorder = &MockUserURLsGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserURLsGetter) EXPECT() *MockUserURLsGetterMockRecorder {
	return m.recorder
}

// GetUserURLs mocks base method.
func (m *MockUserURLsGetter) GetUserURLs(userID string) ([]models.UserURLs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", userID)
	ret0, _ := ret[0].([]models.UserURLs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockUserURLsGetterMockRecorder) GetUserURLs(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockUserURLsGetter)(nil).GetUserURLs), userID)
}