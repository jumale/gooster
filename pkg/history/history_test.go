package history

import (
	"github.com/jumale/gooster/pkg/filesys/fstub"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func TestHistory(t *testing.T) {
	assert := require.New(t)
	fsProps := fstub.Config{
		WorkDir: "/wd",
		HomeDir: "/hd",
	}
	newFs := func() *fstub.Stub {
		return fstub.New(fsProps)
	}

	t.Run("Constructor", func(t *testing.T) {
		t.Run("should create a new clean manager", func(t *testing.T) {
			mng, err := NewManager(Config{FileSys: newFs()})
			assert.NoError(err)
			assert.Equal("", mng.filePath)
			assert.Len(mng.stack, 0)
			assert.Len(mng.set, 0)
		})

		t.Run("should create a manager and load history file", func(t *testing.T) {
			fs := newFs()
			fs.Root().Add("foo/bar.txt", fstub.NewFile("foo", "bar"))

			mng, err := NewManager(Config{HistoryFile: "foo/bar.txt", FileSys: fs})
			assert.NoError(err)

			assert.Equal([]string{"foo", "bar"}, mng.stack)
			assert.Contains(mng.set, "foo")
			assert.Contains(mng.set, "bar")
			assert.True(fs.Get("foo/bar.txt").Closed)
		})

		t.Run("should load history file with ~", func(t *testing.T) {
			homeFilePath := path.Join(fsProps.HomeDir, "foo/bar.txt")
			fs := newFs()
			fs.Root().Add(homeFilePath, fstub.NewFile("foo", "bar"))

			mng, err := NewManager(Config{HistoryFile: "~/foo/bar.txt", FileSys: fs})
			assert.NoError(err)

			assert.Equal([]string{"foo", "bar"}, mng.stack)
			assert.Contains(mng.set, "foo")
			assert.Contains(mng.set, "bar")
			assert.True(fs.Get(homeFilePath).Closed)
		})

		t.Run("should return error if not possible to load the history file", func(t *testing.T) {
			mng, err := NewManager(Config{HistoryFile: "foo/bar.txt", FileSys: newFs()})
			assert.Error(err)
			assert.Nil(mng)
		})
	})

	t.Run("Add", func(t *testing.T) {
		file := "history.txt"

		t.Run("should add new commands at the end of history", func(t *testing.T) {
			fs := newFs()
			fs.Root().Add(file, fstub.NewFile("foo"))
			mng, err := NewManager(Config{HistoryFile: file, FileSys: fs})
			assert.NoError(err)

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
			fs := newFs()
			fs.Root().Add(file, fstub.NewFile("foo"))

			mng, err := NewManager(Config{HistoryFile: file, FileSys: fs})
			assert.NoError(err)
			assert.Equal([]string{"foo"}, fs.Get(file).ContentLines())

			mng.Add("bar")
			assert.Equal([]string{"foo", "bar"}, fs.Get(file).ContentLines())
		})
	})

	create := func(lines ...string) *Manager {
		fs := newFs()
		fs.Root().Add("history.txt", fstub.NewFile(lines...))

		mng, err := NewManager(Config{HistoryFile: "history.txt", FileSys: fs})
		assert.NoError(err)

		return mng
	}

	t.Run("Next", func(t *testing.T) {
		t.Run("should return empty val for empty history", func(t *testing.T) {
			assert.Empty(create().Next())
		})

		t.Run("should return empty val for inactive history (index is not set)", func(t *testing.T) {
			mng := create("foo", "bar", "baz")
			assert.Empty(mng.Next())
		})

		t.Run("should return next value after the index", func(t *testing.T) {
			mng := create("foo", "bar", "baz")
			mng.index = 1
			assert.Equal("baz", mng.Next())
			mng.index = 0
			assert.Equal("bar", mng.Next())
		})

		t.Run("should return empty value and reset index after reaching the end", func(t *testing.T) {
			mng := create("foo", "bar", "baz")
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
			mng := create("foo", "bar", "baz")
			assert.Equal("baz", mng.Prev())
		})

		t.Run("should return previous value before the index", func(t *testing.T) {
			mng := create("foo", "bar", "baz")
			mng.index = 1
			assert.Equal("foo", mng.Prev())
		})

		t.Run("should jump to the end when there is not prev", func(t *testing.T) {
			mng := create("foo", "bar", "baz")
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
