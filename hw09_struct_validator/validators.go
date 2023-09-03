package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type (
	constraints      map[string]string
	stringValidators map[string]func(value, limit string) error
	intValidators    map[string]func(value int64, limit string) error
)

var (
	ErrNotValid  = errors.New("validation error")
	ErrValidator = errors.New("unexpected validator")
	ErrType      = errors.New("unsupported type")
)

var svs = stringValidators{
	"len": func(v, l string) error {
		vl, err := strconv.Atoi(l)
		if err != nil {
			return err
		}
		if len(v) != vl {
			return fmt.Errorf("[%w] invalid len", ErrNotValid)
		}
		return nil
	},

	"regexp": func(v, l string) error {
		re, err := regexp.Compile(l)
		if err != nil {
			return err
		}
		if !re.MatchString(v) {
			return fmt.Errorf("[%w] invalid regexp", ErrNotValid)
		}
		return nil
	},

	"in": func(v, l string) error {
		for _, s := range strings.Split(l, ",") {
			if v == s {
				return nil
			}
		}
		return fmt.Errorf("[%w] invalid in", ErrNotValid)
	},
}

var ivs = intValidators{
	"min": func(v int64, l string) error {
		min, err := strconv.ParseInt(l, 10, 64)
		if err != nil {
			return err
		}
		if v < min {
			return fmt.Errorf("[%w] invalid min", ErrNotValid)
		}
		return nil
	},

	"max": func(v int64, l string) error {
		max, err := strconv.ParseInt(l, 10, 64)
		if err != nil {
			return err
		}
		if v > max {
			return fmt.Errorf("[%w] invalid max", ErrNotValid)
		}
		return nil
	},

	"in": func(v int64, l string) error {
		for _, s := range strings.Split(l, ",") {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			if v == n {
				return nil
			}
		}
		return fmt.Errorf("[%w] invalid in", ErrNotValid)
	},
}

func validateString(verrs ValidationErrors, name string, value string, cs constraints) (ValidationErrors, error) {
	for k, c := range cs {
		vf, ok := svs[k]
		if !ok {
			return nil, ErrValidator
		}

		if err := vf(value, c); err != nil {
			if errors.Is(err, ErrNotValid) {
				verrs = append(verrs, ValidationError{name, err})
			} else {
				return nil, err
			}
		}
	}
	return verrs, nil
}

func validateInt(verrs ValidationErrors, name string, value int64, cs constraints) (ValidationErrors, error) {
	for k, c := range cs {
		vf, ok := ivs[k]
		if !ok {
			return nil, ErrValidator
		}

		if err := vf(value, c); err != nil {
			if errors.Is(err, ErrNotValid) {
				verrs = append(verrs, ValidationError{name, err})
			} else {
				return nil, err
			}
		}
	}
	return verrs, nil
}

func parse(str string) constraints {
	c := constraints{}
	for _, s := range strings.Split(str, "|") {
		kv := strings.Split(s, ":")
		c[kv[0]] = kv[1]
	}
	return c
}
