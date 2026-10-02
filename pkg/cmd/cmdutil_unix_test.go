//go:build !windows

package cmd

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenSocketPairPagerSupportsArguments(t *testing.T) {
	t.Setenv("PAGER", "cat -n")

	pagerInput, pid, err := openSocketPairPager("stream test")
	require.NoError(t, err)
	defer pagerInput.Close()

	_, err = pagerInput.WriteString("Hello world\n")
	require.NoError(t, err)
	require.NoError(t, pagerInput.Close())

	var wstatus syscall.WaitStatus
	_, err = syscall.Wait4(pid, &wstatus, 0, nil)
	require.NoError(t, err)
	assert.True(t, wstatus.Exited())
	assert.Equal(t, 0, wstatus.ExitStatus())
}

func TestOpenSocketPairPagerRejectsMissingPager(t *testing.T) {
	t.Setenv("PAGER", "/path/to/missing-pager --flag")

	_, _, err := openSocketPairPager("stream test")
	require.Error(t, err)
}
