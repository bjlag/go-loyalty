// Code generated by MockGen. DO NOT EDIT.
// Source: logger.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	logger "github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	gomock "github.com/golang/mock/gomock"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Debug", msg)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), msg)
}

// Debugf mocks base method.
func (m *MockLogger) Debugf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockLoggerMockRecorder) Debugf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*MockLogger)(nil).Debugf), varargs...)
}

// Error mocks base method.
func (m *MockLogger) Error(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", msg)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), msg)
}

// Errorf mocks base method.
func (m *MockLogger) Errorf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockLoggerMockRecorder) Errorf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockLogger)(nil).Errorf), varargs...)
}

// Info mocks base method.
func (m *MockLogger) Info(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", msg)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), msg)
}

// Infof mocks base method.
func (m *MockLogger) Infof(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockLoggerMockRecorder) Infof(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockLogger)(nil).Infof), varargs...)
}

// Warning mocks base method.
func (m *MockLogger) Warning(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Warning", msg)
}

// Warning indicates an expected call of Warning.
func (mr *MockLoggerMockRecorder) Warning(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warning", reflect.TypeOf((*MockLogger)(nil).Warning), msg)
}

// Warningf mocks base method.
func (m *MockLogger) Warningf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warningf", varargs...)
}

// Warningf indicates an expected call of Warningf.
func (mr *MockLoggerMockRecorder) Warningf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warningf", reflect.TypeOf((*MockLogger)(nil).Warningf), varargs...)
}

// WithError mocks base method.
func (m *MockLogger) WithError(err error) logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithError", err)
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// WithError indicates an expected call of WithError.
func (mr *MockLoggerMockRecorder) WithError(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithError", reflect.TypeOf((*MockLogger)(nil).WithError), err)
}

// WithField mocks base method.
func (m *MockLogger) WithField(key string, value interface{}) logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithField", key, value)
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// WithField indicates an expected call of WithField.
func (mr *MockLoggerMockRecorder) WithField(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithField", reflect.TypeOf((*MockLogger)(nil).WithField), key, value)
}
