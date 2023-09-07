package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, errors.New("unimplemented reader")
	}

	re, err := regexp.Compile(`"Email":"[^@]*@[^@]*\.` + domain)
	if err != nil {
		return nil, err
	}

	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		email := strings.Split(string(re.Find(scanner.Bytes())), "@")
		if len(email) != 2 {
			continue
		}
		result[strings.ToLower(email[1])]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
