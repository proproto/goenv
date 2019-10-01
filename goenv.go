package goenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
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
	t := reflect.TypeOf(dst)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("goenv: dst must be a pointer: %T", dst)
	}

	structElem := t.Elem()
	value := reflect.ValueOf(dst).Elem()

	err := ParseErrors{}

	for i, n := 0, structElem.NumField(); i < n; i++ {
		f := structElem.Field(i)
		v, ok := f.Tag.Lookup("env")
		if ok {
			values := strings.Split(v, ",")
			if len(values) == 0 || values[0] == "" {
				panic(fmt.Sprintf("goenv: field %s has empty env tag", f.Name))
			}
			setting := buildParseSetting(values)

			if setting.required {
				envValue, ok := os.LookupEnv(setting.envKey)
				if !ok {
					if setting.hasDefaultValue {
						value.Field(i).SetString(setting.defaultValue)
					} else {
						err.Append(fmt.Sprintf("goenv: %s not set", setting.envKey))
					}
				} else {
					setting.envValueRaw = envValue
					setValue(value.Field(i), setting, &err)
				}
			} else {
				setting.envValueRaw = os.Getenv(setting.envKey)
				if setting.envValueRaw == "" {
					if setting.hasDefaultValue {
						value.Field(i).SetString(setting.defaultValue)
					}
				} else {
					setValue(value.Field(i), setting, &err)
				}
			}
		}
	}

	if len(err.msg) == 0 {
		return nil
	}
	return err
}

type parseSetting struct {
	envKey          string
	envValueRaw     string
	required        bool
	hasDefaultValue bool
	defaultValue    string
}

func buildParseSetting(values []string) *parseSetting {
	setting := parseSetting{
		envKey: values[0],
	}

	for _, value := range values[1:] {
		if value == "required" {
			setting.required = true
		} else if strings.HasPrefix(value, "default=") {
			setting.hasDefaultValue = true
			setting.defaultValue = strings.TrimPrefix(value, "default=")
		} else {
			panic("goenv: unknown method: " + value)
		}
	}

	return &setting
}

func setValue(v reflect.Value, setting *parseSetting, err *ParseErrors) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(setting.envValueRaw)
	case reflect.Bool:
		b, e := strconv.ParseBool(setting.envValueRaw)
		if e != nil {
			err.Append(e.Error())
		} else {
			v.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		integer, e := strconv.ParseInt(setting.envValueRaw, 10, 64)
		if e != nil {
			err.Append(e.Error())
		} else {
			v.SetInt(integer)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		integer, e := strconv.ParseUint(setting.envValueRaw, 10, 64)
		if e != nil {
			err.Append(e.Error())
		} else {
			v.SetUint(integer)
		}
	default:
		panic("cannot handle kind: " + v.Kind().String())
	}

}

func MustParse(i interface{}) {
	if err := Parse(i); err != nil {
		panic(err)
	}
}
