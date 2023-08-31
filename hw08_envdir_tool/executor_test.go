package main

import (
	"bytes"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	os.Setenv("HELLO", "SHOULD_REPLACE")
	os.Setenv("FOO", "SHOULD_REPLACE")
	os.Setenv("UNSET", "SHOULD_REMOVE")
	os.Setenv("ADDED", "from original env")
	os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

	env := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": EnvValue{"", true},
	}

	var b bytes.Buffer

	expected := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`

	c := run([]string{"testdata/echo.sh", "arg1=1", "arg2=2"}, env, &b)
	if c != 0 {
		t.Fatalf("error: wrong exit code: want 0 but got %d", c)
	}

	output := b.String()
	if output != expected {
		t.Fatalf("error: invalid output: got:\n%s\nbut want:\n%s\n", output, expected)
	}
}

func TestRunCmdWithErrors(t *testing.T) {
	env := make(Environment)
	var b bytes.Buffer

	if c := run([]string{"some-not-exists-command", "arg=1", "arg=2"}, env, &b); c != 1 {
		t.Fatalf("error: wrong exit code: want 1 but got %d", c)
	}
}
