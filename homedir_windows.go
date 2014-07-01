// +build windows

package main

import (
	"os/user"
)

// HomeDir returns home directory of current user.
func HomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}
