package logger

import (
	"errors"
	"fmt"
)

type LogLevel int

const (
	LevelError LogLevel = 1 << iota
	LevelInfo           = 1 << iota
	LevelTrace          = 1 << iota

	LogLevelError = LevelError
	LogLevelInfo  = LogLevelError | LevelInfo
	LogLevelTrace = LogLevelInfo | LevelTrace

	LogLevelAll = LogLevelTrace
)

type DefaultLogger struct {
	level LogLevel
}

func New() *DefaultLogger {
	return &DefaultLogger{
		level: LogLevelError,
	}
}

func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l DefaultLogger) Log(a ...any) {
	fmt.Print(a...)
}

func (l DefaultLogger) Logf(format string, a ...any) {
	fmt.Printf(format, a...)
}

func (l DefaultLogger) Trace(a ...any) {
	if l.shouldLog(LevelTrace) {
		fmt.Print(a...)
	}
}

func (l DefaultLogger) Tracef(format string, a ...any) {
	if l.shouldLog(LevelTrace) {
		fmt.Printf(format, a...)
	}
}

func (l DefaultLogger) Error(a ...any) error {
	if l.shouldLog(LevelError) {
		return fmt.Errorf(``, a)
	}

	return errors.New(fmt.Sprint(a))
}

func (l DefaultLogger) Errorf(format string, a ...any) error {
	if l.shouldLog(LevelError) {
		return fmt.Errorf(format, a...)
	}

	return errors.New(fmt.Sprintf(format, a))
}

func (l DefaultLogger) shouldLog(level LogLevel) bool {
	return (l.level & level) > 0
}
