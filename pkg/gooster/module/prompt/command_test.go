package prompt

import (
	"github.com/jumale/gooster/pkg/filesys/fstub"
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestDetectWorkDirPath(t *testing.T) {
	assert := _assert.New(t)
	fsProps := &fstub.Props{
		WorkDir: "/work/dir",
		HomeDir: "/home/dir",
	}
	fs := fstub.New(fsProps)
	fs.Root().Add("/work/dir/foo/bar", fstub.NewDir())
	fs.Root().Add("/home/dir/foo/bar", fstub.NewDir())
	fs.Root().Add("/foo/bar/baz", fstub.NewDir())
	fs.Root().Add("/work/dir/some/file", fstub.NewFile())

	t.Run("should detect absolute path from cd command", func(t *testing.T) {
		assert.Equal("/foo/bar/baz", detectWorkDirPath(fs, "cd /foo/bar/baz"))
	})

	t.Run("should detect relative path from cd command", func(t *testing.T) {
		assert.Equal("/work/dir/foo/bar", detectWorkDirPath(fs, "cd foo/bar"))
	})

	t.Run("should detect home path from cd command", func(t *testing.T) {
		assert.Equal("/home/dir/foo/bar", detectWorkDirPath(fs, "cd ~/foo/bar"))
	})

	t.Run("should detect home dir from cd command", func(t *testing.T) {
		assert.Equal("/home/dir", detectWorkDirPath(fs, "cd ~"))
	})

	t.Run("should detect raw absolute path", func(t *testing.T) {
		assert.Equal("/foo/bar/baz", detectWorkDirPath(fs, "/foo/bar/baz"))
	})

	t.Run("should detect raw relative path", func(t *testing.T) {
		assert.Equal("/work/dir/foo/bar", detectWorkDirPath(fs, "./foo/bar"))
	})

	t.Run("should return empty for some random commands", func(t *testing.T) {
		assert.Empty(detectWorkDirPath(fs, "ls"))
		assert.Empty(detectWorkDirPath(fs, "cp ./foo/bar"))
		assert.Empty(detectWorkDirPath(fs, "rm /foo/bar"))
	})

	t.Run("should return empty if it looks like a path, but does not start with '/', './', or '../'", func(t *testing.T) {
		assert.Empty(detectWorkDirPath(fs, "foo/bar/baz"))
	})

	t.Run("should return empty if the target is not a directory", func(t *testing.T) {
		assert.Empty(detectWorkDirPath(fs, "./some/file"))
	})
}
