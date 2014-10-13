package kiicli

import (
	"io/ioutil"
	"log"
)

var logger = log.New(ioutil.Discard, "", log.LstdFlags)

func Logger() *log.Logger {
	return logger
}
