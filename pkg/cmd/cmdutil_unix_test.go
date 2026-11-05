//go:build !windows

package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPagerHasExitedWhenFailedOutputReturns(t *testing.T) {
	streamers := map[string]func(string, func(*os.File) error) error{
		"socket pair": streamOutputOSSpecific,
		"pipe":        streamToPagerWithPipe,
	}

	for name, stream := range streamers {
		t.Run(name, func(t *testing.T) {
			// Like a pager that waits for q, this one outlives the end of its input.
			script := "#!/bin/sh\ncat >/dev/null\nsleep 0.1\ntouch \"$0.exited\"\n"
			pager := filepath.Join(t.TempDir(), "pager")
			require.NoError(t, os.WriteFile(pager, []byte(script), 0o755))
			t.Setenv("PAGER", pager)
			errOutput := errors.New("the stream failed")

			err := stream("label", func(*os.File) error { return errOutput })

			require.ErrorIs(t, err, errOutput)
			require.FileExists(t, pager+".exited")
		})
	}
}
