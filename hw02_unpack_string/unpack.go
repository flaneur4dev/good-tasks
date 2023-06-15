package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		sb    strings.Builder
		next  string
		slash bool
	)

	for _, r := range s {
		switch {
		case '0' <= r && r <= '9':
			if next == "" {
				return "", ErrInvalidString
			}

			if slash {
				sb.WriteString(next)
				next = string(r)
				slash = false
				break
			}

			rep, _ := strconv.Atoi(string(r)) // игнорируем err т.к. в этот блок попадаем с валидным r
			sb.WriteString(strings.Repeat(next, rep))
			next = ""
		case r == '\\':
			if slash {
				sb.WriteString(next)
				next = string(r)
				slash = false
				break
			}

			slash = true
		default:
			if slash {
				return "", ErrInvalidString
			}

			sb.WriteString(next)
			next = string(r)
		}
	}

	sb.WriteString(next)
	return sb.String(), nil
}
