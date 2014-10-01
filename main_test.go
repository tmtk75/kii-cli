package main

import (
	"testing"

	"github.com/tmtk75/cli"
)

func TestFlatten(t *testing.T) {
	a_cmds := []cli.Command{
		{Name: "add"},
		{Name: "sub"},
	}
	b_cmds := []cli.Command{
		{Name: "time"},
		{Name: "div"},
	}
	a := [][]cli.Command{a_cmds, b_cmds}

	actual := Flatten(a)

	if len(actual) != 4 {
		t.Errorf("len is not 4: %d", len(actual))
	}
	if actual[0].Name != "add" {
		t.Errorf("[0] is not `add`: %s", actual[0])
	}
}
