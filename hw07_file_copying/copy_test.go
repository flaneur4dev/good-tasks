package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

const (
	in  = "testdata/input.txt"
	out = "test_out.txt"
)

func hashData(path string) (string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	data, err := io.ReadAll(fd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func TestCopy(t *testing.T) {
	cases := [...]string{
		"out_offset0_limit0",
		"out_offset0_limit10",
		"out_offset0_limit1000",
		"out_offset0_limit10000",
		"out_offset100_limit1000",
		"out_offset6000_limit1000",
	}
	m := make(map[string]string, len(cases))

	for _, c := range cases {
		h, err := hashData("testdata/" + c + ".txt")
		if err != nil {
			panic(err)
		}
		m[c] = h
	}

	tests := [...]struct {
		name   string
		offset int64
		limit  int64
	}{
		{cases[0], 0, 0},
		{cases[1], 0, 10},
		{cases[2], 0, 1000},
		{cases[3], 0, 10000},
		{cases[4], 100, 1000},
		{cases[5], 6000, 1000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := Copy(in, out, tc.offset, tc.limit); err != nil {
				t.Fatalf("copy error: case %s got %v", tc.name, err)
			}

			got, err := hashData(out)
			if err != nil {
				t.Fatalf("hash error: case %s got %v", tc.name, err)
			}
			expected := m[tc.name]
			if got != expected {
				t.Fatalf("error: got %s expected %s", got, expected)
			}
		})
	}

	if err := os.Remove(out); err != nil {
		t.Log(err)
	}
}

func TestCopyWithErrors(t *testing.T) {
	t.Run("files not exists", func(t *testing.T) {
		if err := Copy("fake_in.txt", "fake_out.txt", 0, 0); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("error: got %v but want %v", err, os.ErrNotExist)
		}
	})

	t.Run("offset > file size", func(t *testing.T) {
		if err := Copy(in, out, 6800, 100); !errors.Is(err, ErrOffsetExceedsFileSize) {
			t.Fatalf("error: got %v but want %v", err, ErrOffsetExceedsFileSize)
		}
	})

	t.Run("unsupported file", func(t *testing.T) {
		if err := Copy("testdata", out, 0, 10); !errors.Is(err, ErrUnsupportedFile) {
			t.Fatalf("error: got %v but want %v", err, ErrUnsupportedFile)
		}
	})
}
