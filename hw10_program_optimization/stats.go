package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, errors.New("unimplemented reader")
	}

	var user User
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		err := user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			return nil, err
		}

		if strings.Contains(user.Email, domain) {
			email := strings.Split(user.Email, "@")
			result[strings.ToLower(email[1])]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
