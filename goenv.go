package goenv

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type ParseErrors struct {
	msg []string
}

func (err *ParseErrors) Append(msg string) {
	err.msg = append(err.msg, msg)
}

func (err ParseErrors) Error() string {
	return strings.Join(err.msg, ", ")
}

func (err ParseErrors) String() string {
	return err.Error()
}

func Parse(dst interface{}) error {
	ele := reflect.TypeOf(dst).Elem()
	value := reflect.ValueOf(dst).Elem()

	err := ParseErrors{}

	for i, n := 0, ele.NumField(); i < n; i++ {
		f := ele.Field(i)
		v, ok := f.Tag.Lookup("env")
		if ok {
			values := strings.Split(v, ",")
			if len(values) == 0 {
				err.Append("goenv: env tag has empty value")
				return err
			}

			envKey := values[0]
			isRequired := contains(values, "required")
			if isRequired {
				envValue, ok := os.LookupEnv(envKey)
				if !ok {
					if v, ok := containsValue(values, "default"); ok {
						value.Field(i).SetString(v)
					} else {
						err.Append(fmt.Sprintf("goenv: %s not set", envKey))
					}
				}
				value.Field(i).SetString(envValue)
			} else {
				envValue := os.Getenv(envKey)
				if v, ok := containsValue(values, "default"); ok {
					value.Field(i).SetString(v)
				} else {
					value.Field(i).SetString(envValue)
				}
			}
		}
	}

	if len(err.msg) == 0 {
		return nil
	}
	return err
}

func MustParse(i interface{}) {
	if err := Parse(i); err != nil {
		panic(err)
	}
}
