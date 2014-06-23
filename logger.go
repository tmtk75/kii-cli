package main

import (
	"fmt"
	"log"
	"os"
)

var logger Logger = &FakeLogger{}

type Logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type FakeLogger struct {
}

func (self *FakeLogger) Printf(format string, args ...interface{}) {
	// NOP
}

func (self *FakeLogger) Errorf(format string, args ...interface{}) {
	// NOP
}

type _Logger struct {
}

func (self *_Logger) Printf(format string, args ...interface{}) {
	//_, fn, line, _ := runtime.Caller(1)
	//f := fmt.Sprintf("%s:%d -- %s", filepath.Base(fn), line, format)
	log.Printf(format, args...)
}

func (self *_Logger) Errorf(format string, args ...interface{}) {
	//_, fn, line, _ := runtime.Caller(1)
	//f := fmt.Sprintf("%s:%d -- %s", filepath.Base(fn), line, format)
	fmt.Fprintf(os.Stderr, format, args...)
}
