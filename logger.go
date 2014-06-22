package main

import "log"

var logger Logger = &FakeLogger{}

type Logger interface {
	Printf(format string, args ...interface{})
}

type FakeLogger struct {
}

func (self *FakeLogger) Printf(format string, args ...interface{}) {
	// NOP
}

type _Logger struct {
}

func (self *_Logger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
