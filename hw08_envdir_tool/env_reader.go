package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		k, v, r, err := process(dir, f.Name())
		if err != nil {
			fmt.Printf("processing %s failed: %v\n", f.Name(), err)
			continue
		}

		env[k] = EnvValue{v, r}
	}

	return env, nil
}

func process(dir, file string) (k string, v string, r bool, e error) {
	fd, err := os.Open(filepath.Join(dir, file))
	if err != nil {
		e = err
		return
	}
	defer func() {
		if err := fd.Close(); err != nil {
			fmt.Println("invalid close")
		}
	}()

	fStat, err := fd.Stat()
	if err != nil {
		e = err
		return
	}

	if fStat.Size() == 0 {
		r = true
	}

	k = strings.ReplaceAll(file, "=", "")

	sc := bufio.NewScanner(fd)
	if sc.Scan() {
		v = strings.TrimRight(strings.ReplaceAll(sc.Text(), "\x00", "\n"), " \t")
	}

	return
}
