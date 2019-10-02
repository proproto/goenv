package goenv

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse_Argument(t *testing.T) {
	t.Run("non-pointer", func(t *testing.T) {
		assert.EqualError(t, Parse(struct{}{}), "goenv: dst must be a pointer to struct: struct {}")
	})

	t.Run("non-pointer-to-struct", func(t *testing.T) {
		assert.EqualError(t, Parse(""), "goenv: dst must be a pointer to struct: string")
	})
}

func TestEnvTag(t *testing.T) {
	type emptyEnv struct {
		Field string `env:""`
	}

	assert.PanicsWithValue(t, "goenv: field Field has empty env tag", func() { Parse(&emptyEnv{}) })
}

func TestEnvUnknownMethod(t *testing.T) {
	type unknown struct {
		Field string `env:"ENV_KEY,unknown"`
	}

	assert.PanicsWithValue(t, "goenv: unknown method: unknown", func() { Parse(&unknown{}) })
}

func TestBindDuration(t *testing.T) {
	type T struct {
		Timeout time.Duration `env:"HTTP_TIMEOUT,default=10s"`
	}
}

func TestParseMySQLConfig(t *testing.T) {
	type MySQLConfig struct {
		Host     string `env:"MYSQL_HOST,default=localhost:3306"`
		User     string `env:"MYSQL_USER,default=root"`
		Password string `env:"MYSQL_PASSWORD"`
		Database string `env:"MYSQL_DATABASE,required"`
		MaxConns int    `env:"MYSQL_MAX_CONNS"`
		Ping     bool   `env:"MYSQL_AUTO_PING"`
	}

	cases := map[string]struct {
		SetupFunc func()
		Config    MySQLConfig
		Error     error
	}{
		"Example1": {
			SetupFunc: func() {
				os.Setenv("MYSQL_PASSWORD", "rootpassword")
				os.Setenv("MYSQL_DATABASE", "db")
			},
			Config: MySQLConfig{
				Host:     "localhost:3306",
				User:     "root",
				Password: "rootpassword",
				Database: "db",
			},
		},
		"Example2": {
			SetupFunc: func() {
				os.Setenv("MYSQL_HOST", "http://10.42.8.63")
				os.Setenv("MYSQL_USER", "admin")
				os.Setenv("MYSQL_DATABASE", "service")
				os.Setenv("MYSQL_MAX_CONNS", "100")
				os.Setenv("MYSQL_AUTO_PING", "true")
			},
			Config: MySQLConfig{
				Host:     "http://10.42.8.63",
				User:     "admin",
				Database: "service",
				MaxConns: 100,
				Ping:     true,
			},
		},
		"Required": {
			Error: errors.New("goenv: MYSQL_DATABASE not set"),
		},
	}

	for testname, testcase := range cases {
		t.Run(testname, func(t *testing.T) {
			os.Clearenv()
			if testcase.SetupFunc != nil {
				testcase.SetupFunc()
			}

			config := MySQLConfig{}
			err := Parse(&config)

			if testcase.Error != nil {
				assert.EqualError(t, testcase.Error, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testcase.Config, config)
				t.Logf("%s: %#v", testname, config)
			}
		})
	}
}
