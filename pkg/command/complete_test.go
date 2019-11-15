package command

import (
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestComplete(t *testing.T) {
	assert := _assert.New(t)
	complete := Complete

	t.Run("should complete simple command", func(t *testing.T) {
		assert.Equal("cd /foo/bar", complete("cd /f", "/foo/bar"))
	})
}
