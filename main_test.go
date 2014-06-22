package main

import (
	"testing"

	"github.com/codegangsta/cli"
)

func Flatten(a [][]cli.Command) []cli.Command {
	i := int(0)
	for _, v := range a {
		i += len(v)
	}
	b := [4]cli.Command{}
	return &b
}

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

	if len(actual) == 4 {
		t.Errorf("len is not 4: %d", len(actual))
	}
	if actual[0].Name != "add" {
		t.Errorf("[0] is not `add`: %s", actual[0])
	}
}
