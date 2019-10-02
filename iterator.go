package goenv

import (
	"strings"
)

type optsIterator struct {
	opts    string
	current string
}

func newOptsIterator(opts string) *optsIterator {
	return &optsIterator{opts: opts}
}

func (it *optsIterator) Next() bool {
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

func (it *optsIterator) Name() string {
	if idx := strings.IndexByte(it.current, '='); idx >= 0 {
		return it.current[:idx]
	}
	return it.current
}

func (it *optsIterator) Value() string {
	if idx := strings.IndexByte(it.current, '='); idx >= 0 {
		return it.current[idx+1:]
	}
	return ""
}
