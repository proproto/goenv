# goenv
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![GoDoc](https://godoc.org/github.com/proproto/goenv?status.svg)](https://godoc.org/github.com/proproto/goenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/proproto/goenv)](https://goreportcard.com/report/github.com/proproto/goenv)
[![License](https://img.shields.io/badge/License-MIT-red.svg)](https://github.com/proproto/goenv/blob/master/LICENSE)
[![CircleCI](https://circleci.com/gh/proproto/goenv.svg?style=svg)](https://circleci.com/gh/proproto/goenv)

The goenv is used to map environment variables to go struct.

## Installation
```
go get -u github.com/proproto/goenv
```

## Example
```go
type MySQLConfig struct {
	Host     string        `env:"MYSQL_HOST,default=localhost:3306"`
	User     string        `env:"MYSQL_USER,default=root"`
	Password string        `env:"MYSQL_PASSWORD"`
	Database string        `env:"MYSQL_DATABASE,required"`
	Timeout  time.Duration `env:"MYSQL_TIMEOUT,default=10s"`
	TLS      bool          `env:"MYSQL_TLS_ENABLED"`
	MaxConns int           `env:"MYSQL_MAX_CONNS"`
}

os.Clearenv()
os.Setenv("MYSQL_PASSWORD", "db_password")
os.Setenv("MYSQL_DATABASE", "db_name")
os.Setenv("MYSQL_TLS_ENABLED", "true")
os.Setenv("MYSQL_MAX_CONNS", "32")

config := MySQLConfig{}
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
```
