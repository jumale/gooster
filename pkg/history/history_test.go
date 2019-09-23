package history

import (
	"github.com/jumale/gooster/pkg/stub"
	_assert "github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestHistory(t *testing.T) {
	assert := _assert.New(t)

	t.Run("Constructor", func(t *testing.T) {
		t.Run("should create a new constructor", func(t *testing.T) {
			mng := NewManager("/foo/bar")
			assert.Equal("/foo/bar", mng.filePath)
		})

		t.Run("should format path with home dir", func(t *testing.T) {
			home, _ := os.UserHomeDir()
			mng := NewManager("~/foo/bar")
			assert.Equal(home+"/foo/bar", mng.filePath)
		})
	})

	var historyFile *stub.FileStub
	historyPath := "/foo/bar"

	create := func(lines ...string) *Manager {
		m := NewManager(historyPath)
		historyFile = stub.NewFileStub(lines)
		m.openReadFile = func(path string) (closer io.ReadCloser, e error) {
			assert.Equal(historyPath, path)
			return historyFile, nil
		}
		m.openWriteFile = func(path string) (closer io.WriteCloser, e error) {
			assert.Equal(historyPath, path)
			return historyFile, nil
		}
		return m
	}

	t.Run("Load", func(t *testing.T) {
		t.Run("should load history lines and close file", func(t *testing.T) {
			mng := create("foo", "bar", "baz")
			assert.False(historyFile.Closed)
			mng.Load()

			assert.Equal([]string{"foo", "bar", "baz"}, mng.stack)
			assert.Contains(mng.set, "foo")
			assert.Contains(mng.set, "bar")
			assert.Contains(mng.set, "baz")
			assert.True(historyFile.Closed)
		})
	})

	t.Run("Add", func(t *testing.T) {
		t.Run("should add new commands at the end of history", func(t *testing.T) {
			mng := create("foo").Load()

			assert.Equal([]string{"foo"}, mng.stack)
			assert.Contains(mng.set, "foo")

			mng.Add("bar")
			assert.Equal([]string{"foo", "bar"}, mng.stack)
			assert.Contains(mng.set, "foo")
			assert.Contains(mng.set, "bar")

			mng.Add("baz")
			assert.Equal([]string{"foo", "bar", "baz"}, mng.stack)
			assert.Contains(mng.set, "foo")
			assert.Contains(mng.set, "bar")
			assert.Contains(mng.set, "baz")
		})

		t.Run("should write new commands to the file", func(t *testing.T) {
			mng := create()
			mng.Add("foo")
			assert.Equal([]string{"foo"}, historyFile.Writes())
			mng.Add("bar")
			assert.Equal([]string{"foo", "bar"}, historyFile.Writes())
		})
	})

	t.Run("Next", func(t *testing.T) {
		t.Run("should return empty val for empty history", func(t *testing.T) {
			assert.Empty(create().Next())
		})

		t.Run("should return empty val for inactive history (index is not set)", func(t *testing.T) {
			mng := create("foo", "bar", "baz")
			assert.Empty(mng.Next())
		})

		t.Run("should return next value after the index", func(t *testing.T) {
			mng := create("foo", "bar", "baz").Load()
			mng.index = 1
			assert.Equal("baz", mng.Next())
			mng.index = 0
			assert.Equal("bar", mng.Next())
		})

		t.Run("should return empty value and reset index after reaching the end", func(t *testing.T) {
			mng := create("foo", "bar", "baz").Load()
			mng.index = 1
			assert.Equal("baz", mng.Next())
			assert.Equal("", mng.Next())
			assert.Equal(-1, mng.index)
			assert.Equal("", mng.Next()) // just in case, one more time
		})
	})

	t.Run("Prev", func(t *testing.T) {
		t.Run("should return empty val for empty history", func(t *testing.T) {
			assert.Empty(create().Prev())
		})

		t.Run("should return latest val for inactive history (index is not set)", func(t *testing.T) {
			mng := create("foo", "bar", "baz").Load()
			assert.Equal("baz", mng.Prev())
		})

		t.Run("should return previous value before the index", func(t *testing.T) {
			mng := create("foo", "bar", "baz").Load()
			mng.index = 1
			assert.Equal("foo", mng.Prev())
		})

		t.Run("should jump to the end when there is not prev", func(t *testing.T) {
			mng := create("foo", "bar", "baz").Load()
			mng.index = 0
			assert.Equal("baz", mng.Prev())
		})
	})

	t.Run("Reset", func(t *testing.T) {
		t.Run("should reset index", func(t *testing.T) {
			mng := create()
			assert.Equal(-1, mng.index)

			mng.index = 5
			mng.Reset()
			assert.Equal(-1, mng.index)
		})
	})
}
