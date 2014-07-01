// +build !windows

package main

import (
	"os"
)

// HomeDir returns home directory of current user.
func HomeDir() (string, error) {
	return os.Getenv("HOME"), nil
}
