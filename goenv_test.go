package goenv

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBind_Argument(t *testing.T) {
	t.Run("NonPointer", func(t *testing.T) {
		assert.EqualError(t, Bind(struct{}{}), "goenv: dst must be a pointer to struct: struct {}")
	})

	t.Run("NonPointerToStruct", func(t *testing.T) {
		assert.EqualError(t, Bind(""), "goenv: dst must be a pointer to struct: string")
	})
}

func TestBind_EmptyTag(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		type emptyEnv struct {
			CamelCaseField string `env:""`
		}

		os.Clearenv()
		os.Setenv("CAMEL_CASE_FIELD", "(｀･ω･´)")
		c := emptyEnv{}

		assert.NoError(t, Bind(&c))
		assert.Equal(t, "(｀･ω･´)", c.CamelCaseField)
	})

	t.Run("NotSet", func(t *testing.T) {
		type emptyEnv struct {
			Field string `env:""`
		}

		os.Clearenv()
		c := emptyEnv{}

		assert.NoError(t, Bind(&c))
		assert.Empty(t, c.Field)
	})
}

func TestBind_UnknownOption(t *testing.T) {
	type unknown struct {
		Field string `env:"ENV_KEY,unknown"`
	}

	assert.PanicsWithValue(t, "goenv: unknown method: unknown", func() { Bind(&unknown{}) })
}

func TestBind_String(t *testing.T) {
	t.Run("NoOption", func(t *testing.T) {
		type Config struct {
			Field string `env:"ENV_FIELD"`
		}
		cases := map[string]struct {
			SetupFunc func()
			Expect    Config
		}{
			"Implicitly": {
				SetupFunc: func() {},
			},
			"Explicitly": {
				SetupFunc: func() {
					os.Setenv("ENV_FIELD", "set")
				},
				Expect: Config{
					Field: "set",
				},
			},
		}

		for testname, testcase := range cases {
			t.Run(testname, func(t *testing.T) {
				os.Clearenv()
				testcase.SetupFunc()

				c := Config{}

				assert.NoError(t, Bind(&c))
				assert.Equal(t, testcase.Expect, c)
			})
		}
	})

	t.Run("WithDefault", func(t *testing.T) {
		type Config struct {
			Field string `env:"ENV_FIELD,default=value"`
		}

		cases := map[string]struct {
			SetupFunc func()
			Expect    Config
		}{
			"Implicitly": {
				SetupFunc: func() {},
				Expect: Config{
					Field: "value",
				},
			},
			"Explicitly": {
				SetupFunc: func() {
					os.Setenv("ENV_FIELD", "another")
				},
				Expect: Config{
					Field: "another",
				},
			},
		}

		for testname, testcase := range cases {
			t.Run(testname, func(t *testing.T) {
				os.Clearenv()
				testcase.SetupFunc()

				c := Config{}

				assert.NoError(t, Bind(&c))
				assert.Equal(t, testcase.Expect, c)
			})
		}
	})

	t.Run("WithRequired", func(t *testing.T) {
		type Config struct {
			Field string `env:"ENV_FIELD,required"`
		}

		cases := map[string]struct {
			SetupFunc func()
			Expect    Config
			Err       error
		}{
			"Implicitly": {
				SetupFunc: func() {},
				Err:       fmt.Errorf("goenv: %s not set", "ENV_FIELD"),
			},
			"Explicitly": {
				SetupFunc: func() {
					os.Setenv("ENV_FIELD", "ok")
				},
				Expect: Config{
					Field: "ok",
				},
			},
		}

		for testname, testcase := range cases {
			t.Run(testname, func(t *testing.T) {
				os.Clearenv()
				testcase.SetupFunc()

				c := Config{}
				err := Bind(&c)

				if testcase.Err == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, testcase.Err, err.Error())
				}
				assert.Equal(t, testcase.Expect, c)
			})
		}
	})
}

func TestBind_Int(t *testing.T) {
	t.Run("NoOption", func(t *testing.T) {
		type Config struct {
			Int int `env:"ENV_INT"`
		}

		t.Run("Implicitly", func(t *testing.T) {
			os.Clearenv()
			config := &Config{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 0, config.Int)
		})

		t.Run("Explicitly", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("ENV_INT", "1024")
			config := &Config{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 1024, config.Int)
		})

		t.Run("NaN", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("ENV_INT", "NaN")
			config := &Config{}
			err := Bind(config)
			_, expectErr := strconv.ParseInt("NaN", 10, 64)
			assert.EqualError(t, expectErr, err.Error())
		})
	})

	t.Run("WithDefault", func(t *testing.T) {
		type Config struct {
			Int int `env:"ENV_INT,default=256"`
		}

		t.Run("Implicitly", func(t *testing.T) {
			os.Clearenv()
			config := &Config{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 256, config.Int)
		})

		t.Run("Explicitly", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("ENV_INT", "1024")
			config := &Config{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 1024, config.Int)
		})
	})

	t.Run("WithRequired", func(t *testing.T) {
		type Config struct {
			Int int `env:"ENV_INT,required"`
		}

		t.Run("Implicitly", func(t *testing.T) {
			os.Clearenv()
			config := &Config{}
			err := Bind(config)
			assert.EqualError(t, fmt.Errorf("goenv: %s not set", "ENV_INT"), err.Error())
		})

		t.Run("Explicitly", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("ENV_INT", "1024")
			config := &Config{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 1024, config.Int)
		})
	})
}

func TestBind_Duration(t *testing.T) {
	t.Run("NoOption", func(t *testing.T) {
		type ServerConfig struct {
			Timeout time.Duration `env:"HTTP_TIMEOUT"`
		}

		t.Run("Implicitly", func(t *testing.T) {
			os.Clearenv()
			config := &ServerConfig{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, time.Duration(0), config.Timeout)
		})

		t.Run("Explicitly", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("HTTP_TIMEOUT", "30s")
			config := &ServerConfig{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 30*time.Second, config.Timeout)
		})

		t.Run("InvalidDuration", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("HTTP_TIMEOUT", "nat")
			config := &ServerConfig{}
			err := Bind(config)
			_, expectErr := time.ParseDuration("nat")
			assert.EqualError(t, expectErr, err.Error())
		})
	})

	t.Run("WithDefault", func(t *testing.T) {
		type ServerConfig struct {
			Timeout time.Duration `env:"HTTP_TIMEOUT,default=10s"`
		}

		t.Run("Implicitly", func(t *testing.T) {
			os.Clearenv()
			config := &ServerConfig{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 10*time.Second, config.Timeout)
		})

		t.Run("Explicitly", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("HTTP_TIMEOUT", "30s")
			config := &ServerConfig{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 30*time.Second, config.Timeout)
		})
	})

	t.Run("WithRequired", func(t *testing.T) {
		type ServerConfig struct {
			Timeout time.Duration `env:"HTTP_TIMEOUT,required"`
		}

		t.Run("Implicitly", func(t *testing.T) {
			os.Clearenv()
			config := &ServerConfig{}
			err := Bind(config)
			assert.EqualError(t, fmt.Errorf("goenv: %s not set", "HTTP_TIMEOUT"), err.Error())
		})

		t.Run("Explicitly", func(t *testing.T) {
			os.Clearenv()
			os.Setenv("HTTP_TIMEOUT", "30s")
			config := &ServerConfig{}
			err := Bind(config)
			assert.NoError(t, err)
			assert.Equal(t, 30*time.Second, config.Timeout)
		})
	})
}
