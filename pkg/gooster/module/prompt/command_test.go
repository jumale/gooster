package prompt

import (
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestDetectWorkDirPath(t *testing.T) {
	assert := _assert.New(t)
	userHomeDir = func() (string, error) {
		return "/home/dir", nil
	}
	getWd = func() (string, error) {
		return "/work/dir", nil
	}

	t.Run("should detect absolute path from cd command", func(t *testing.T) {
		assert.Equal("/foo/bar/baz", detectWorkDirPath("cd /foo/bar/baz"))
	})

	t.Run("should detect relative path from cd command", func(t *testing.T) {
		assert.Equal("/work/dir/foo/bar", detectWorkDirPath("cd foo/bar"))
	})

	t.Run("should detect home path from cd command", func(t *testing.T) {
		assert.Equal("/home/dir/foo/bar", detectWorkDirPath("cd ~/foo/bar"))
	})

	t.Run("should detect home dir from cd command", func(t *testing.T) {
		assert.Equal("/home/dir", detectWorkDirPath("cd ~"))
	})

	t.Run("should detect raw absolute path", func(t *testing.T) {
		assert.Equal("/foo/bar/baz", detectWorkDirPath("/foo/bar/baz"))
	})

	t.Run("should detect raw relative path", func(t *testing.T) {
		assert.Equal("/work/dir/foo/bar", detectWorkDirPath("./foo/bar"))
	})

	t.Run("should return empty for some random commands", func(t *testing.T) {
		assert.Empty(detectWorkDirPath("ls"))
		assert.Empty(detectWorkDirPath("cp ./foo/bar"))
		assert.Empty(detectWorkDirPath("rm /foo/bar"))
	})

	t.Run("should return empty if it looks like a path, but does not start with '/', './', or '../'", func(t *testing.T) {
		assert.Empty(detectWorkDirPath("foo/bar/baz"))
	})
}
