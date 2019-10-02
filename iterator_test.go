package goenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	itr := newOptsIterator("required,default=ENV_VALUE")

	assert.True(t, itr.Next())
	assert.Equal(t, "required", itr.Name())
	assert.Equal(t, "", itr.Value())

	assert.True(t, itr.Next())
	assert.Equal(t, "default", itr.Name())
	assert.Equal(t, "ENV_VALUE", itr.Value())

	assert.False(t, itr.Next())
}
