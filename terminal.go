// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import (
	"io"
	"os"
	"os/exec"
)

// stdOut is a writer to store the value of last
// stdout write.
type stdOut struct {
	str string
}

func (s *stdOut) Write(p []byte) (n int, err error) {
	s.str = string(p)
	return len(p), nil
}

// readKey reads one key silently and returns the
// details.
func readKey() (string, error) {
	so := &stdOut{}
	err := runCommand(so, "bash", "-c", "read -n1 -s key; "+
		"case \"$key\" in $'\\0a') key='enter' ;; $'\\177') key='backspace' ;; $'\\e') read -n2 -s key ;; esac; "+
		"echo -n $key")
	if err != nil {
		return "", err
	}

	if len(so.str) == 2 {
		switch so.str {
		case "[D":
			return "left-arrow", nil
		case "[C":
			return "right-arrow", nil
		}
	}
	return so.str, nil
}

// LineUp moves the cursor one line up.
func lineUp() error {
	return runCommand(nil, "tput", "cuu1")
}

// runCommand runs the command and the stdout is redirected
// to the given stdout (or OS's stdout if given is nil) and
// the OS's stderr.
func runCommand(stdout io.Writer, command string, args ...string) (exitCode error) {
	cmd := exec.Command(command, args...)
	if stdout == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
