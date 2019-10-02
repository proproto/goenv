package options

import (
	"strings"
)

// Iterator struct
type Iterator struct {
	opts    string
	current string
}

// NewIterator creates a new iterator
func NewIterator(opts string) *Iterator {
	return &Iterator{opts: opts}
}

// Next returns true if iterator has next value
func (it *Iterator) Next() bool {
	if (it.opts == "") || (it.current == it.opts) {
		return false
	}

	i := strings.IndexByte(it.opts, ',')
	if i >= 0 {
		it.current, it.opts = it.opts[:i], it.opts[i+1:]
	} else {
		it.current = it.opts
	}

	return true
}

// Name returns current option name
func (it *Iterator) Name() string {
	if idx := strings.IndexByte(it.current, '='); idx >= 0 {
		return it.current[:idx]
	}
	return it.current
}

// Value returns current option value
func (it *Iterator) Value() string {
	if idx := strings.IndexByte(it.current, '='); idx >= 0 {
		return it.current[idx+1:]
	}
	return ""
}
