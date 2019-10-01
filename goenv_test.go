package goenv

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MySQLConfig struct {
	Host     string `env:"MYSQL_HOST,default=localhost:3306"`
	User     string `env:"MYSQL_USER,default=root"`
	Password string `env:"MYSQL_PASSWORD"`
	Database string `env:"MYSQL_DATABASE,required"`
}

func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		c.User,
		c.Password,
		c.Host,
		c.Database,
	)
}

func TestParse(t *testing.T) {
	cases := map[string]struct {
		SetupFunc func()
		DSN       string
		Error     error
	}{
		"Example": {
			SetupFunc: func() {
				os.Setenv("MYSQL_PASSWORD", "rootpassword")
				os.Setenv("MYSQL_DATABASE", "db")
			},
			DSN: "root:rootpassword@tcp(localhost:3306)/db",
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
				assert.Equal(t, testcase.DSN, config.DSN())
				t.Logf("%s: %#v", testname, config)
			}
		})
	}
}
