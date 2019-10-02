package goenv_test

import (
	"fmt"
	"os"
	"time"

	"github.com/proproto/goenv"
)

func ExampleBind() {
	type config struct {
		Host string `env:"HOST,required"`
	}

	os.Clearenv()
	os.Setenv("HOST", "0.0.0.0")

	c := config{}
	err := goenv.Bind(&c)

	fmt.Println("err:", err)
	fmt.Println("Host:", c.Host)
	// OUTPUT: err: <nil>
	// Host: 0.0.0.0
}

func ExampleBind_error() {
	type config struct {
		Host string `env:"HOST,required"`
	}

	os.Clearenv()

	c := config{}
	err := goenv.Bind(&c)

	fmt.Println(err)
	// OUTPUT: goenv: HOST not set
}

func ExampleMustBind() {
	type MySQLConfig struct {
		Host     string        `env:"MYSQL_HOST,default=localhost:3306"`
		User     string        `env:"MYSQL_USER,default=root"`
		Password string        `env:"MYSQL_PASSWORD"`
		Database string        `env:"MYSQL_DATABASE,required"`
		Timeout  time.Duration `env:"MYSQL_TIMEOUT,default=10s"`
		TLS      bool          `env:"MYSQL_TLS_ENABLED"`
		MaxConns int           `env:"MYSQL_MAX_CONNS"`
	}

	config := MySQLConfig{}

	os.Clearenv()
	os.Setenv("MYSQL_PASSWORD", "db_password")
	os.Setenv("MYSQL_DATABASE", "db_name")
	os.Setenv("MYSQL_TLS_ENABLED", "true")
	os.Setenv("MYSQL_MAX_CONNS", "32")

	goenv.MustBind(&config)

	fmt.Printf("Host: %s\n", config.Host)
	fmt.Printf("User: %s\n", config.User)
	fmt.Printf("Password: %s\n", config.Password)
	fmt.Printf("Database: %s\n", config.Database)
	fmt.Printf("Timeout: %v\n", config.Timeout.String())
	fmt.Printf("TLS: %v\n", config.TLS)
	fmt.Printf("MaxConns: %v\n", config.MaxConns)
	// Output:
	// Host: localhost:3306
	// User: root
	// Password: db_password
	// Database: db_name
	// Timeout: 10s
	// TLS: true
	// MaxConns: 32
}
