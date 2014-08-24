package main

import "testing"

func Test_converLogFormat(t *testing.T) {
	f := "${time} [${level}]"
	expected := "{{.time}} [{{.level}}]"
	k := convertLogFormat(f)
	if k != expected {
		t.Errorf("expected %v, but %v", expected, k)
	}
}
