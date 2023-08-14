package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	bufSize = 1024
	pb      = "**********"
	sp      = "          "
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Println("src file: invalid close")
		}
	}()

	srcStat, err := src.Stat()
	if err != nil {
		return err
	}

	srcMode := srcStat.Mode()
	if !srcMode.IsRegular() {
		return ErrUnsupportedFile
	}

	srcSize := srcStat.Size()
	if offset > srcSize {
		return ErrOffsetExceedsFileSize
	}

	dist, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if err := dist.Close(); err != nil {
			fmt.Println("dist file: invalid close")
		}
	}()

	buf := make([]byte, bufSize)
	sum := 0
	isCopied := false
	copySize := srcSize - offset

	if 0 < limit && limit < copySize {
		copySize -= limit
	}
	if offset > 0 {
		if _, err = src.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	for !isCopied {
		nr, err := src.Read(buf)
		switch {
		case errors.Is(err, io.EOF):
			isCopied = true
		case err != nil:
			return err
		}

		sum += nr
		if 0 < limit && limit <= int64(sum) {
			nr -= sum - int(limit)
			sum = int(copySize)
			isCopied = true
		}
		p := sum * 100 / int(copySize)
		i := p / 10

		_, err = dist.Write(buf[:nr])
		if err != nil {
			return err
		}

		fmt.Printf("\r0%% %s%s %d%% ", pb[:i], sp[i:], p)
		// time.Sleep(time.Millisecond * 100) // for testing progressbar
	}

	return nil
}
