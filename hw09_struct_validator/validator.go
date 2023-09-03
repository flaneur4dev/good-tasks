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
		s = append(s, fmt.Sprintf("%s: %s", e.Field, e.Err))
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

		cs := parse(tv)
		fv := rv.Field(i)
		var err error

		switch fv.Kind() {
		case reflect.String:
			errs, err = validateString(errs, f.Name, fv.String(), cs)
			if err != nil {
				return err
			}

		case reflect.Int:
			errs, err = validateInt(errs, f.Name, fv.Int(), cs)
			if err != nil {
				return err
			}

		case reflect.Slice:
			sLen := fv.Len()
			if sLen == 0 {
				break
			}

			switch fv.Index(0).Kind() {
			case reflect.String:
				for i := 0; i < sLen; i++ {
					fn := fmt.Sprintf("%s[%d]", f.Name, i)
					errs, err = validateString(errs, fn, fv.Index(i).String(), cs)
					if err != nil {
						return err
					}
				}
			case reflect.Int:
				for i := 0; i < sLen; i++ {
					fn := fmt.Sprintf("%s[%d]", f.Name, i)
					errs, err = validateInt(errs, fn, fv.Index(i).Int(), cs)
					if err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("[%w] %s", ErrType, fv.Type())
			}

		default:
			return fmt.Errorf("[%w] %s", ErrType, fv.Type())
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
