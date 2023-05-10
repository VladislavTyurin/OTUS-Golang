package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}

	command := cmd[0]
	args := make([]string, 0)
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	c := exec.Command(command, args...)
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	err := c.Run()
	if err != nil {
		// Return not 0 code
		return 1
	}
	return c.ProcessState.ExitCode()
}
