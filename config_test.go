package main

import (
	"testing"
)

func TestPickup(t *testing.T) {
	if pickup("foo", "bar", "baz") != "foo" {
		t.Fatalf("should return 1st %#v", "foo")
	}
	if pickup("", "bar", "baz") != "bar" {
		t.Fatalf("should return 2nd %#v", "bar")
	}
	if pickup("", "", "baz") != "baz" {
		t.Fatalf("should return 3rd %#v", "baz")
	}
}
