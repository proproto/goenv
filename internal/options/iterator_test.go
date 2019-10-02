package options_test

import (
	"testing"

	"github.com/proproto/goenv/internal/options"
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	itr := options.NewIterator("required,default=ENV_VALUE")

	assert.True(t, itr.Next())
	assert.Equal(t, "required", itr.Name())
	assert.Equal(t, "", itr.Value())

	assert.True(t, itr.Next())
	assert.Equal(t, "default", itr.Name())
	assert.Equal(t, "ENV_VALUE", itr.Value())

	assert.False(t, itr.Next())
}

func TestIterator_EmptyOption(t *testing.T) {
	itr := options.NewIterator("")
	assert.False(t, itr.Next())
}
