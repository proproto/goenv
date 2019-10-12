package goenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/proproto/camelcase"
	"github.com/proproto/goenv/internal/options"
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

// Binder struct
type Binder struct {
	// TagKey specifies the tag key while calling reflect.Field.Tag.Lookup(string)
	TagKey string
}

// DefaultBinder uses `env` as TagKey
var DefaultBinder = &Binder{TagKey: "env"}

// Bind binds environment variables to dst
func (b *Binder) Bind(dst interface{}) error {
	t := reflect.TypeOf(dst)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("goenv: dst must be a pointer to struct: %T", dst)
	}

	structElem := t.Elem()
	value := reflect.ValueOf(dst).Elem()

	err := bindErrors{}

	for i, n := 0, structElem.NumField(); i < n; i++ {
		f := structElem.Field(i)
		tag, ok := f.Tag.Lookup(b.TagKey)
		if ok && tag != "" {
			if handleErr := handleTagPresent(tag, i, f, value); handleErr != nil {
				err.Append(handleErr.Error())
			}
		} else {
			if handleErr := handleTagNotPresent(i, f, value); handleErr != nil {
				err.Append(handleErr.Error())
			}
		}
	}

	if len(err.msg) == 0 {
		return nil
	}
	return err
}

func handleTagPresent(tag string, index int, sf reflect.StructField, value reflect.Value) error {
	if tag == "-" {
		return nil
	}

	setting := buildBindSetting(tag)

	if setting.required {
		envValue, ok := os.LookupEnv(setting.envKey)
		if !ok {
			if setting.hasDefaultValue {
				if err := setValue(value.Field(index), setting.defaultValue); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("goenv: %s not set", setting.envKey)
			}
		} else {
			setting.envValueRaw = envValue
			if err := setValue(value.Field(index), setting.envValueRaw); err != nil {
				return err
			}
		}
	} else {
		setting.envValueRaw = os.Getenv(setting.envKey)
		if setting.envValueRaw == "" {
			if setting.hasDefaultValue {
				if err := setValue(value.Field(index), setting.defaultValue); err != nil {
					return err
				}
			}
		} else {
			if err := setValue(value.Field(index), setting.envValueRaw); err != nil {
				return err
			}
		}
	}

	return nil
}

func handleTagNotPresent(index int, sf reflect.StructField, value reflect.Value) error {
	envKey := camelcase.ToMacroCase(sf.Name)
	if ev := os.Getenv(envKey); ev != "" {
		return setValue(value.Field(index), ev)
	}

	return nil
}

type bindSetting struct {
	envKey          string
	envValueRaw     string
	required        bool
	hasDefaultValue bool
	defaultValue    string
}

func buildBindSetting(v string) *bindSetting {
	envKey, opts := parseTag(v)
	setting := bindSetting{envKey: envKey}

	for opts.Next() {
		switch opts.Name() {
		case "required":
			setting.required = true
		case "default":
			setting.hasDefaultValue = true
			setting.defaultValue = opts.Value()
		default:
			panic("goenv: unknown method: " + opts.Name())
		}
	}
	return &setting
}

func parseTag(tag string) (string, *options.Iterator) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], options.NewIterator(tag[idx+1:])
	}
	return tag, options.NewIterator("")
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

func setValue(v reflect.Value, stringValue string) error {
	switch v.Type() {
	case typeString:
		v.SetString(stringValue)
	case typeBool:
		b, err := strconv.ParseBool(stringValue)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case typeInt, typeInt8, typeInt16, typeInt32, typeInt64:
		i, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case typeUint, typeUint8, typeUint16, typeUint32, typeUint64:
		i, err := strconv.ParseUint(stringValue, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case typeDuration:
		d, err := time.ParseDuration(stringValue)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(d))
	default:
		return fmt.Errorf("goenv: unsupported bind type: %s", v.Type())
	}

	return nil
}

// MustBind binds environment variables to dst by DefaultBinder
// if fails, it would panic
func MustBind(dst interface{}) {
	if err := Bind(dst); err != nil {
		panic(err)
	}
}

// Bind binds environment variables to dst by DefaultBinder
func Bind(dst interface{}) error {
	return DefaultBinder.Bind(dst)
}
