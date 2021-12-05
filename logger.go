package main

import "log"

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

var _ Logger = log.Default()

type dummyLogger struct{}

func (l dummyLogger) Print(v ...interface{})                 {}
func (l dummyLogger) Printf(format string, v ...interface{}) {}
