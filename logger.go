package main

import (
	"io/ioutil"
	"log"
)

var logger = log.New(ioutil.Discard, "", log.LstdFlags)
