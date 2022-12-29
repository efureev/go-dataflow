package dataflow

import (
	"errors"
	"fmt"
)

type Logger interface {
	Log(...any)
	Logf(string, ...any)
	Trace(...any)
	Tracef(string, ...any)
	Error(...any) error
	Errorf(string, ...any) error
}

type ProxyLogger struct {
	client Logger
}

func (l ProxyLogger) hasClientLogger() bool {
	return l.client != nil
}

func (l ProxyLogger) Log(a ...any) {
	if l.hasClientLogger() {
		l.client.Log(a...)
	}
}

func (l ProxyLogger) Logf(format string, a ...any) {
	if l.hasClientLogger() {
		l.client.Logf(format, a...)
	}
}

func (l ProxyLogger) Trace(a ...any) {
	if l.hasClientLogger() {
		l.client.Trace(a...)
	}
}

func (l ProxyLogger) Tracef(format string, a ...any) {
	if l.hasClientLogger() {
		l.client.Tracef(format, a...)
	}
}

func (l ProxyLogger) Error(a ...any) error {
	if l.hasClientLogger() {
		return l.client.Error(a...)
	}

	return errors.New(fmt.Sprint(a...))
}

func (l ProxyLogger) Errorf(format string, a ...any) error {
	if l.hasClientLogger() {
		return l.client.Error(a...)
	}

	return errors.New(fmt.Sprintf(format, a...))
}
