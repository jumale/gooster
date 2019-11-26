package command

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComplete(t *testing.T) {
	assert := require.New(t)
	complete := ApplyCompletion

	t.Run("should complete simple command", func(t *testing.T) {
		assert.Equal("cd /foo/bar", complete("cd /f", "/foo/bar"))
	})

	t.Run("should complete value without any prefix", func(t *testing.T) {
		assert.Equal("/foo/bar", complete("/f", "/foo/bar"))
	})

	t.Run("should recognise escaped spaces", func(t *testing.T) {
		assert.Equal(`cd /foo\ bar/baz`, complete(`cd /foo\ ba`, `/foo\ bar/baz`))
	})
}
