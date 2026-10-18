package cmd

import "github.com/anthropics/anthropic-cli/internal/jsonview"

type itemLimitIterator[T any] struct {
	jsonview.Iterator[T]
	remaining int64
}

func (it *itemLimitIterator[T]) Next() bool {
	if it.remaining == 0 || !it.Iterator.Next() {
		return false
	}
	if it.remaining > 0 {
		it.remaining--
	}
	return true
}
