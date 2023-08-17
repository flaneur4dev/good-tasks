package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	return run(cmd, env, nil)
}

func run(cmd []string, env Environment, out io.Writer) int {
	if _, err := exec.LookPath(cmd[0]); err != nil {
		fmt.Printf("%v\n", err)
		return 1
	}

	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v.Value)
		}
	}

	com := exec.Command(cmd[0], cmd[1:]...)
	com.Stdin = os.Stdin
	com.Stderr = os.Stderr
	if out == nil {
		com.Stdout = os.Stdout
	} else {
		com.Stdout = out
	}

	if err := com.Run(); err != nil {
		fmt.Printf("command finished with error: %v\n", err)
		return 1
	}

	return 0
}
