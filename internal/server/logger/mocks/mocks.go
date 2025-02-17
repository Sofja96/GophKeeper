// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/server/logger/logger.go

// Package mock_logging is a generated GoMock package.
package mock_logging

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	logrus "github.com/sirupsen/logrus"
)

// MockILogger is a mock of ILogger interface.
type MockILogger struct {
	ctrl     *gomock.Controller
	recorder *MockILoggerMockRecorder
}

// MockILoggerMockRecorder is the mock recorder for MockILogger.
type MockILoggerMockRecorder struct {
	mock *MockILogger
}

// NewMockILogger creates a new mock instance.
func NewMockILogger(ctrl *gomock.Controller) *MockILogger {
	mock := &MockILogger{ctrl: ctrl}
	mock.recorder = &MockILoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockILogger) EXPECT() *MockILoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockILogger) Debug(format string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockILoggerMockRecorder) Debug(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockILogger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockILogger) Error(format string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockILoggerMockRecorder) Error(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockILogger)(nil).Error), varargs...)
}

// Info mocks base method.
func (m *MockILogger) Info(format string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockILoggerMockRecorder) Info(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockILogger)(nil).Info), varargs...)
}

// Log mocks base method.
func (m *MockILogger) Log() *logrus.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*logrus.Logger)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockILoggerMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockILogger)(nil).Log))
}
