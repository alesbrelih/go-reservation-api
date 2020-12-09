package test_util

import (
	"io"
	"log"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/mock"
)

type HcLogMock struct {
	mock.Mock
}

func (h *HcLogMock) Log(level hclog.Level, msg string, args ...interface{}) {
	h.Called(level, msg, args)
	return
}

func (h *HcLogMock) Trace(msg string, args ...interface{}) {
	h.Called(msg, args)
	return
}

func (h *HcLogMock) Debug(msg string, args ...interface{}) {
	h.Called(msg, args)
	return
}

func (h *HcLogMock) Info(msg string, args ...interface{}) {
	h.Called(msg, args)
	return
}

func (h *HcLogMock) Warn(msg string, args ...interface{}) {
	h.Called(msg, args)
	return
}

func (h *HcLogMock) Error(msg string, args ...interface{}) {
	h.Called(msg, args)
	return
}

func (h *HcLogMock) IsTrace() bool {
	args := h.Called()
	return args.Get(0).(bool)
}

func (h *HcLogMock) IsDebug() bool {
	args := h.Called()
	return args.Get(0).(bool)
}

func (h *HcLogMock) IsInfo() bool {
	args := h.Called()
	return args.Get(0).(bool)
}

func (h *HcLogMock) IsWarn() bool {
	args := h.Called()
	return args.Get(0).(bool)
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (h *HcLogMock) IsError() bool {
	args := h.Called()
	return args.Get(0).(bool)
}

func (h *HcLogMock) ImpliedArgs() []interface{} {
	args := h.Called()
	return args.Get(0).([]interface{})
}

func (h *HcLogMock) With(args ...interface{}) hclog.Logger {
	h.Called(args)
	return nil
}

func (h *HcLogMock) Name() string {
	args := h.Called()
	return args.Get(0).(string)
}

func (h *HcLogMock) Named(name string) hclog.Logger {
	args := h.Called(name)
	return args.Get(0).(hclog.Logger)
}

func (h *HcLogMock) ResetNamed(name string) hclog.Logger {
	args := h.Called(name)
	return args.Get(0).(hclog.Logger)
}

func (h *HcLogMock) SetLevel(level hclog.Level) {
	h.Called(level)
	return
}

func (h *HcLogMock) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	args := h.Called(opts)
	return args.Get(0).(*log.Logger)
}

func (h *HcLogMock) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	args := h.Called(opts)
	return args.Get(0).(io.Writer)
}
