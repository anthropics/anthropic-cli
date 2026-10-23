package cmd

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type budgetTestIterator struct {
	items []string
	index int
	calls int
	err   error
	done  bool
}

func (it *budgetTestIterator) Next() bool {
	it.calls++
	if it.index == len(it.items) {
		it.done = true
		return false
	}
	it.index++
	return true
}

func (it *budgetTestIterator) Current() string {
	return it.items[it.index-1]
}

func (it *budgetTestIterator) Err() error {
	if it.done {
		return it.err
	}
	return nil
}

func TestItemLimitIterator(t *testing.T) {
	t.Parallel()
	items := make([]string, 50)
	for i := range items {
		items[i] = strconv.Itoa(i)
	}
	for _, tc := range []struct {
		name      string
		limit     int64
		wantItems int
		wantCalls int
	}{
		{"zero", 0, 0, 0},
		{"less than explorer preload", 2, 2, 2},
		{"at explorer preload", 20, 20, 20},
		{"beyond explorer preload", 25, 25, 25},
		{"exact source length", 50, 50, 50},
		{"more than source length", 60, 50, 51},
		{"unlimited", -1, 50, 51},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			source := &budgetTestIterator{items: items}
			iter := &itemLimitIterator[string]{Iterator: source, remaining: tc.limit}
			got := make([]string, 0)
			for iter.Next() {
				got = append(got, iter.Current())
			}
			require.Equal(t, items[:tc.wantItems], got)
			require.Equal(t, tc.wantCalls, source.calls)
			require.NoError(t, iter.Err())
			if tc.limit >= 0 && tc.limit <= int64(len(items)) {
				require.False(t, iter.Next())
				require.Equal(t, tc.wantCalls, source.calls, "re-polling a capped explorer must not fetch another item")
			}
		})
	}
}

func TestItemLimitIteratorPreservesErrors(t *testing.T) {
	t.Parallel()
	for _, limit := range []int64{-1, 5} {
		sourceErr := errors.New("next page failed")
		source := &budgetTestIterator{items: []string{"first", "second"}, err: sourceErr}
		iter := &itemLimitIterator[string]{Iterator: source, remaining: limit}
		for iter.Next() {
		}
		require.ErrorIs(t, iter.Err(), sourceErr)
	}
}
