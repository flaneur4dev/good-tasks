package main

import "testing"

func TestReadDir(t *testing.T) {
	expected := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": EnvValue{"", true},
	}

	env, err := ReadDir("testdata/env")
	if err != nil {
		t.Fatalf("error: ReadDir got %v", err)
	}

	envLen, expLen := len(env), len(expected)
	if envLen != expLen {
		t.Fatalf("error: wrong number of variables: want %d but got %d", expLen, envLen)
	}

	for k := range env {
		if env[k] != expected[k] {
			t.Fatalf("error: wrong value for %s", k)
			break
		}
	}
}

func TestReadDirWithErrors(t *testing.T) {
	if _, err := ReadDir("some-not-exists-folder/env"); err == nil {
		t.Fatalf("error: want error but got %v", err)
	}
}
