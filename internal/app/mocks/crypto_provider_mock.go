// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/apolsh/yapr-url-shortener/internal/app/crypto (interfaces: CryptographicProvider)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCryptographicProvider is a mock of CryptographicProvider interface.
type MockCryptographicProvider struct {
	ctrl     *gomock.Controller
	recorder *MockCryptographicProviderMockRecorder
}

// MockCryptographicProviderMockRecorder is the mock recorder for MockCryptographicProvider.
type MockCryptographicProviderMockRecorder struct {
	mock *MockCryptographicProvider
}

// NewMockCryptographicProvider creates a new mock instance.
func NewMockCryptographicProvider(ctrl *gomock.Controller) *MockCryptographicProvider {
	mock := &MockCryptographicProvider{ctrl: ctrl}
	mock.recorder = &MockCryptographicProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCryptographicProvider) EXPECT() *MockCryptographicProviderMockRecorder {
	return m.recorder
}

// Decrypt mocks base method.
func (m *MockCryptographicProvider) Decrypt(arg0 []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrypt", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrypt indicates an expected call of Decrypt.
func (mr *MockCryptographicProviderMockRecorder) Decrypt(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*MockCryptographicProvider)(nil).Decrypt), arg0)
}

// Encrypt mocks base method.
func (m *MockCryptographicProvider) Encrypt(arg0 []byte) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encrypt", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Encrypt indicates an expected call of Encrypt.
func (mr *MockCryptographicProviderMockRecorder) Encrypt(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*MockCryptographicProvider)(nil).Encrypt), arg0)
}