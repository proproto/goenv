package goenv

import (
	"strings"
)

type optsIterator struct {
	isBegin bool
	opts    string
	current string
	end     int
}

func newOptsIterator(opts string) *optsIterator {
	it := &optsIterator{
		opts: opts,
		end:  -1,
	}

	return it
}

func (it *optsIterator) Next() bool {
	if it.isBegin && (it.end == -1 || it.end+1 >= len(it.opts)) {
		return false
	}

	it.isBegin = true

	it.opts = it.opts[it.end+1:]
	if idx := strings.IndexByte(it.opts, ','); idx != -1 { // has multiple options
		it.current = it.opts[:idx]
		it.end = idx
	} else { // single option
		it.current = it.opts
		it.end = len(it.opts)
	}

	return true
}

func (it *optsIterator) Name() string {
	idx := strings.IndexByte(it.current, '=')
	if idx == -1 {
		return it.current
	}
	return it.current[:idx]
}

func (it *optsIterator) Value() string {
	idx := strings.IndexByte(it.current, '=')
	if idx == -1 {
		return ""
	}
	return it.current[idx+1:]
}
