package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	s := make([]string, 0, len(v))
	for _, e := range v {
		s = append(s, fmt.Sprintf("%s: [validation error] %s", e.Field, e.Err.Error()))
	}
	return strings.Join(s, "\n")
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("[%w] %s", ErrType, rv.Type())
	}

	var errs ValidationErrors
	t := rv.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}

		tv, ok := f.Tag.Lookup("validate")
		if !ok {
			continue
		}

		fv := rv.Field(i)
		k, ok, err := parseKind(fv)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		cs, err := parseTag(tv, k)
		if err != nil {
			return err
		}

		switch k { //nolint:exhaustive
		case reflect.String:
			errs = validateValue(errs, f.Name, value{s: fv.String()}, cs)
		case reflect.Int:
			errs = validateValue(errs, f.Name, value{i: fv.Int()}, cs)
		case reflect.Slice:
			errs = validateSlice(errs, f.Name, fv, cs)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
