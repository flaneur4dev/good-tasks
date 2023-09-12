package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	value struct {
		s string
		i int64
	}
	validator      func(value) error
	validatorsWrap map[string]func(limit string, kind reflect.Kind) (validator, error)
	constraints    map[string]validator
)

var (
	ErrType       = errors.New("unsupported type")
	ErrNotValidIn = errors.New("invalid \"in\"")
)

var vw = validatorsWrap{
	"len": func(limit string, _ reflect.Kind) (validator, error) {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}

		return func(v value) error {
			if len(v.s) != l {
				return errors.New("invalid \"len\"")
			}
			return nil
		}, nil
	},

	"regexp": func(limit string, _ reflect.Kind) (validator, error) {
		re, err := regexp.Compile(limit)
		if err != nil {
			return nil, err
		}

		return func(v value) error {
			if !re.MatchString(v.s) {
				return errors.New("invalid \"regexp\"")
			}
			return nil
		}, nil
	},

	"min": func(limit string, _ reflect.Kind) (validator, error) {
		min, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return nil, err
		}

		return func(v value) error {
			if v.i < min {
				return errors.New("invalid \"min\"")
			}
			return nil
		}, nil
	},

	"max": func(limit string, _ reflect.Kind) (validator, error) {
		max, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return nil, err
		}

		return func(v value) error {
			if v.i > max {
				return errors.New("invalid \"max\"")
			}
			return nil
		}, nil
	},

	"in": func(limit string, k reflect.Kind) (validator, error) {
		switch k { //nolint:exhaustive
		case reflect.String:
			m := map[string]struct{}{}
			for _, s := range strings.Split(limit, ",") {
				m[s] = struct{}{}
			}

			return func(v value) error {
				if _, ok := m[v.s]; ok {
					return nil
				}
				return ErrNotValidIn
			}, nil

		case reflect.Int:
			m := map[int64]struct{}{}
			for _, s := range strings.Split(limit, ",") {
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return nil, err
				}
				m[n] = struct{}{}
			}

			return func(v value) error {
				if _, ok := m[v.i]; ok {
					return nil
				}
				return ErrNotValidIn
			}, nil

		default:
			return nil, ErrType
		}
	},
}

func parseKind(rv reflect.Value) (reflect.Kind, bool, error) {
	kind := rv.Kind()
	switch kind { //nolint:exhaustive
	case reflect.String, reflect.Int:
		return kind, true, nil

	case reflect.Slice:
		if rv.Len() == 0 {
			return kind, false, nil
		}

		switch rv.Index(0).Kind() { //nolint:exhaustive
		case reflect.String, reflect.Int:
			return kind, true, nil
		default:
			return kind, false, fmt.Errorf("[%w] %s", ErrType, rv.Type())
		}
	default:
		return kind, false, fmt.Errorf("[%w] %s", ErrType, rv.Type())
	}
}

func parseTag(str string, k reflect.Kind) (constraints, error) {
	c := constraints{}

	for _, s := range strings.Split(str, "|") {
		kv := strings.Split(s, ":")
		name, limit := kv[0], kv[1]

		w, ok := vw[name]
		if !ok {
			return nil, errors.New("unexpected validator")
		}

		v, err := w(limit, k)
		if err != nil {
			return nil, err
		}
		c[name] = v
	}

	return c, nil
}

func validateValue(verrs ValidationErrors, name string, v value, cs constraints) ValidationErrors {
	for _, vd := range cs {
		if err := vd(v); err != nil {
			verrs = append(verrs, ValidationError{name, err})
		}
	}
	return verrs
}

func validateSlice(verrs ValidationErrors, name string, rv reflect.Value, cs constraints) ValidationErrors {
	switch rv.Index(0).Kind() { //nolint:exhaustive
	case reflect.String:
		for i := 0; i < rv.Len(); i++ {
			n := fmt.Sprintf("%s[%d]", name, i)
			verrs = validateValue(verrs, n, value{s: rv.Index(i).String()}, cs)
		}

	case reflect.Int:
		for i := 0; i < rv.Len(); i++ {
			n := fmt.Sprintf("%s[%d]", name, i)
			verrs = validateValue(verrs, n, value{i: rv.Index(i).Int()}, cs)
		}
	}

	return verrs
}
