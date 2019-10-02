package goenv

import "strings"

type optsIterator struct {
	opts       string
	current    string
	begin, end int
}

func NewOptsIterator(opts string) *optsIterator {
	return &optsIterator{
		opts: opts,
	}
}

func (it *optsIterator) Next() bool {
	return false
}

func (it *optsIterator) Name() string {
	idx := strings.IndexByte(it.current, '=')
	if idx == -1 {
		return it.current
	}
	return it.current[:idx+1]
}

func (it *optsIterator) Value() string {
	idx := strings.IndexByte(it.current, '=')
	if idx == -1 {
		return ""
	}
	return it.current[idx+1:]
}
