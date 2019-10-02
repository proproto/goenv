package goenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type bindErrors struct {
	msg []string
}

func (err *bindErrors) Append(msg string) {
	err.msg = append(err.msg, msg)
}

func (err bindErrors) Error() string {
	return strings.Join(err.msg, ", ")
}

func (err bindErrors) String() string {
	return err.Error()
}

// MustBind panics if Bind fails
func MustBind(i interface{}) {
	if err := Bind(i); err != nil {
		panic(err)
	}
}

// Bind binds environment variables to dst
func Bind(dst interface{}) error {
	t := reflect.TypeOf(dst)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("goenv: dst must be a pointer to struct: %T", dst)
	}

	structElem := t.Elem()
	value := reflect.ValueOf(dst).Elem()

	err := bindErrors{}

	for i, n := 0, structElem.NumField(); i < n; i++ {
		f := structElem.Field(i)
		v, ok := f.Tag.Lookup("env")
		if ok {
			it := NewOptsIterator(v)
			for it.Next() {
				fmt.Printf("name: %s, value: %s\n", it.Name(), it.Value())
			}

			values := strings.Split(v, ",")
			if len(values) == 0 || values[0] == "" {
				panic(fmt.Sprintf("goenv: field %s has empty env tag", f.Name))
			}
			setting := buildBindSetting(values)

			if setting.required {
				envValue, ok := os.LookupEnv(setting.envKey)
				if !ok {
					if setting.hasDefaultValue {
						setValue(value.Field(i), setting.defaultValue, &err)
					} else {
						err.Append(fmt.Sprintf("goenv: %s not set", setting.envKey))
					}
				} else {
					setting.envValueRaw = envValue
					setValue(value.Field(i), setting.envValueRaw, &err)
				}
			} else {
				setting.envValueRaw = os.Getenv(setting.envKey)
				if setting.envValueRaw == "" {
					if setting.hasDefaultValue {
						setValue(value.Field(i), setting.defaultValue, &err)
					}
				} else {
					setValue(value.Field(i), setting.envValueRaw, &err)
				}
			}
		}
	}

	if len(err.msg) == 0 {
		return nil
	}
	return err
}

type bindSetting struct {
	envKey          string
	envValueRaw     string
	required        bool
	hasDefaultValue bool
	defaultValue    string
}

func buildBindSetting(values []string) *bindSetting {
	setting := bindSetting{
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

var (
	typeBool     = reflect.TypeOf((*bool)(nil)).Elem()
	typeDuration = reflect.TypeOf((*time.Duration)(nil)).Elem()
	typeInt      = reflect.TypeOf((*int)(nil)).Elem()
	typeInt8     = reflect.TypeOf((*int8)(nil)).Elem()
	typeInt16    = reflect.TypeOf((*int16)(nil)).Elem()
	typeInt32    = reflect.TypeOf((*int32)(nil)).Elem()
	typeInt64    = reflect.TypeOf((*int64)(nil)).Elem()
	typeString   = reflect.TypeOf((*string)(nil)).Elem()
	typeUint     = reflect.TypeOf((*uint)(nil)).Elem()
	typeUint8    = reflect.TypeOf((*uint8)(nil)).Elem()
	typeUint16   = reflect.TypeOf((*uint16)(nil)).Elem()
	typeUint32   = reflect.TypeOf((*uint32)(nil)).Elem()
	typeUint64   = reflect.TypeOf((*uint64)(nil)).Elem()
)

func setValue(v reflect.Value, stringValue string, err *bindErrors) {
	switch v.Type() {
	case typeString:
		v.SetString(stringValue)
	case typeBool:
		b, e := strconv.ParseBool(stringValue)
		if e != nil {
			err.Append(e.Error())
		} else {
			v.SetBool(b)
		}
	case typeInt, typeInt8, typeInt16, typeInt32, typeInt64:
		integer, e := strconv.ParseInt(stringValue, 10, 64)
		if e != nil {
			err.Append(e.Error())
		} else {
			v.SetInt(integer)
		}
	case typeUint, typeUint8, typeUint16, typeUint32, typeUint64:
		integer, e := strconv.ParseUint(stringValue, 10, 64)
		if e != nil {
			err.Append(e.Error())
		} else {
			v.SetUint(integer)
		}
	case typeDuration:
		d, e := time.ParseDuration(stringValue)
		if e != nil {
			err.Append(e.Error())
		} else {
			v.Set(reflect.ValueOf(d))
		}
	default:
		panic("goenv: unsupported bind type: " + v.Type().String())
	}
}
