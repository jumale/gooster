package fstub

import (
	_assert "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuilder(t *testing.T) {
	assert := _assert.New(t)

	t.Run("Add", func(t *testing.T) {
		t.Run("should add single file", func(t *testing.T) {
			stub := New(nil)
			builder := &builder{stub: stub}

			builder.Add("foo.txt", NewFile())

			assertFiles(t, stub.files, expectedFiles{
				"foo.txt": file{name: "foo.txt"},
			})
		})

		t.Run("should add file with nested path", func(t *testing.T) {
			stub := New(nil)
			builder := &builder{stub: stub}

			builder.Add("foo/bar/baz.txt", NewFile())

			assertFiles(t, stub.files, expectedFiles{
				"foo":             dir{name: "foo"},
				"foo/bar":         dir{name: "bar"},
				"foo/bar/baz.txt": file{name: "baz.txt"},
			})
		})

		t.Run("should add nested file to a sub-dir", func(t *testing.T) {
			stub := New(nil)
			builder := &builder{stub: stub, path: "foo"}

			builder.Add("bar/baz.txt", NewFile())

			assertFiles(t, stub.files, expectedFiles{
				"foo":             dir{name: "foo"},
				"foo/bar":         dir{name: "bar"},
				"foo/bar/baz.txt": file{name: "baz.txt"},
			})
		})

		t.Run("should add file with abs path", func(t *testing.T) {
			stub := New(nil)
			builder := &builder{stub: stub}

			builder.Add("/foo.txt", NewFile())

			assertFiles(t, stub.files, expectedFiles{
				"/foo.txt": file{name: "foo.txt"},
			})
		})

		t.Run("should add file with nested abs path", func(t *testing.T) {
			stub := New(nil)
			builder := &builder{stub: stub}

			builder.Add("/foo/bar.txt", NewFile())

			assertFiles(t, stub.files, expectedFiles{
				"/foo":         dir{name: "foo"},
				"/foo/bar.txt": file{name: "bar.txt"},
			})
		})

		t.Run("should ignore context path when abs path is provided", func(t *testing.T) {
			stub := New(nil)
			files := stub.files
			builder := &builder{stub: stub, path: "some/context/path"}

			builder.Add("/foo/bar.txt", NewFile())

			assert.Len(files, 2)

			assert.Equal("foo", files["/foo"].Info.Name)
			assert.True(files["/foo"].Info.IsDir)

			assert.Equal("bar.txt", files["/foo/bar.txt"].Info.Name)
			assert.False(files["/foo/bar.txt"].Info.IsDir)
		})
	})

	t.Run("AddDir", func(t *testing.T) {
		t.Run("should add dir and return a new builder with sub-dir", func(t *testing.T) {
			stub := New(nil)
			files := stub.files
			builder := &builder{stub: stub, path: "foo"}

			builder.
				AddDir("bar").
				AddDir("baz").
				Add("new.txt", NewFile())

			assert.Len(files, 4)

			assert.Equal("foo", files["foo"].Info.Name)
			assert.True(files["foo"].Info.IsDir)

			assert.Equal("bar", files["foo/bar"].Info.Name)
			assert.True(files["foo/bar"].Info.IsDir)

			assert.Equal("baz", files["foo/bar/baz"].Info.Name)
			assert.True(files["foo/bar/baz"].Info.IsDir)

			assert.Equal("new.txt", files["foo/bar/baz/new.txt"].Info.Name)
			assert.False(files["foo/bar/baz/new.txt"].Info.IsDir)
		})

		t.Run("should add dir with abs path", func(t *testing.T) {
			stub := New(nil)
			files := stub.files
			builder := &builder{stub: stub}

			builder.AddDir("/foo")

			assert.Len(files, 1)

			assert.Equal("foo", files["/foo"].Info.Name)
			assert.True(files["/foo"].Info.IsDir)
		})

		t.Run("should add dir with nested abs path", func(t *testing.T) {
			stub := New(nil)
			files := stub.files
			builder := &builder{stub: stub}

			builder.AddDir("/foo/bar")

			assert.Len(files, 2)

			assert.Equal("foo", files["/foo"].Info.Name)
			assert.True(files["/foo"].Info.IsDir)
		})
	})
}

func TestStub(t *testing.T) {
	assert := _assert.New(t)

	t.Run("ReadDir", func(t *testing.T) {
		t.Run("should return all items in directory", func(t *testing.T) {
			stub := New(nil)

			stub.Root().
				AddDir("bar/foo").
				Add("barFoo.txt", NewFile()).
				Add("subBarFoo/suBarFoo.txt", NewFile())

			stub.Root().
				AddDir("baz/foo").
				Add("bazFoo.txt", NewFile()).
				Add("subBazFoo/suBazFoo.txt", NewFile())

			files, err := stub.ReadDir("baz/foo")

			assert.NoError(err)
			assert.Len(files, 2)
			assert.Equal("bazFoo.txt", files[0].Name())
			assert.Equal("subBazFoo", files[1].Name())
		})

		t.Run("should find by both local and global paths", func(t *testing.T) {
			stub := New(&Props{WorkDir: "/wd"})

			stub.Root().
				AddDir("foo/bar").
				Add("baz.txt", NewFile())

			files, err := stub.ReadDir("foo/bar")
			assert.NoError(err)
			assert.Len(files, 1)
			assert.Equal("baz.txt", files[0].Name())

			files, err = stub.ReadDir("/wd/foo/bar")
			assert.NoError(err)
			assert.Len(files, 1)
			assert.Equal("baz.txt", files[0].Name())
		})
	})

	t.Run("Open", func(t *testing.T) {
		t.Run("should open file if exists", func(t *testing.T) {
			stub := New(nil)
			stub.Root().Add("foo/bar.txt", NewFile())

			f, err := stub.Open("foo/bar.txt")
			assert.NoError(err)
			assert.Equal("bar.txt", f.(*FileStub).Info.Name)
			assert.Equal(os.O_RDONLY, f.(*FileStub).Info.Flag)
		})

		t.Run("should return error if file does not exist", func(t *testing.T) {
			stub := New(nil)

			f, err := stub.Open("foo.txt")
			assert.Error(err)
			assert.Nil(f)
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("should open file if exists", func(t *testing.T) {
			stub := New(nil)
			stub.Root().Add("foo/bar.txt", NewFile("lorem ipsum"))

			f, err := stub.Create("foo/bar.txt")
			actual := f.(*FileStub)
			assert.NoError(err)
			assert.Equal("bar.txt", actual.Info.Name)
			assert.Equal(os.O_RDWR|os.O_CREATE|os.O_TRUNC, actual.Info.Flag)
			assert.Equal(os.FileMode(0), actual.Info.Mode)
			assert.Equal("lorem ipsum", actual.ContentString())
		})

		t.Run("should create a new file if does not exist", func(t *testing.T) {
			stub := New(nil)

			f, err := stub.Create("foo/bar.txt")
			actual := f.(*FileStub)
			assert.NoError(err)
			assert.Equal("bar.txt", actual.Info.Name)
			assert.Equal(os.O_RDWR|os.O_CREATE|os.O_TRUNC, actual.Info.Flag)
			assert.Equal(os.FileMode(0666), actual.Info.Mode)
			assert.Empty(actual.Content())
		})
	})
}

func assertFiles(t *testing.T, actual map[filePath]*FileStub, expected expectedFiles) {
	assert := _assert.New(t)
	assert.Len(actual, len(expected))

	for pth, item := range expected {
		info := actual[pth].Info
		switch v := item.(type) {
		case file:
			assert.Equal(v.name, info.Name)
			assert.False(info.IsDir)

		case dir:
			assert.Equal(v.name, info.Name)
			assert.True(info.IsDir)

		default:
			assert.Fail("non supported file expectation")
		}
	}
}

type expectedFiles map[filePath]interface{}

type file struct {
	name string
}
type dir struct {
	name string
}
