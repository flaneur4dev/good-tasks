package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

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
	dec := json.NewDecoder(r)

	for dec.More() {
		if err := dec.Decode(&user); err != nil {
			return nil, err
		}

		if strings.Contains(user.Email, domain) {
			email := strings.Split(user.Email, "@")
			result[strings.ToLower(email[1])]++
		}
	}

	return result, nil
}
