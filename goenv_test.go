package goenv

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBind_Argument(t *testing.T) {
	t.Run("non-pointer", func(t *testing.T) {
		assert.EqualError(t, Bind(struct{}{}), "goenv: dst must be a pointer to struct: struct {}")
	})

	t.Run("non-pointer-to-struct", func(t *testing.T) {
		assert.EqualError(t, Bind(""), "goenv: dst must be a pointer to struct: string")
	})
}

func TestEnvTag(t *testing.T) {
	type emptyEnv struct {
		Field string `env:""`
	}

	assert.PanicsWithValue(t, "goenv: field Field has empty env tag", func() { Bind(&emptyEnv{}) })
}

func TestEnvUnknownMethod(t *testing.T) {
	type unknown struct {
		Field string `env:"ENV_KEY,unknown"`
	}

	assert.PanicsWithValue(t, "goenv: unknown method: unknown", func() { Bind(&unknown{}) })
}

func TestBindDuration(t *testing.T) {
	type ServerConfig struct {
		Timeout time.Duration `env:"HTTP_TIMEOUT,default=10s"`
	}

	t.Run("default", func(t *testing.T) {
		config := &ServerConfig{}
		err := Bind(config)
		assert.NoError(t, err)
		assert.Equal(t, 10*time.Second, config.Timeout)
	})

	t.Run("explicit", func(t *testing.T) {
		os.Setenv("HTTP_TIMEOUT", "30s")
		config := &ServerConfig{}
		err := Bind(config)
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, config.Timeout)
	})
}
