package ext

import (
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/filesys/fstub"
	"github.com/stretchr/testify/require"
	"testing"
)

type Input []*dirtree.Node
type Expected []string

func TestSortExtension(t *testing.T) {
	assert := require.New(t)

	filesAndDirs := Input{
		file("ccc.txt"),
		dir("bbb"),
		file("aaa.txt"),
		dir("ddd"),
	}

	t.Run("should by default sort by name ASC", func(t *testing.T) {
		assert.EqualValues(
			Expected{
				"aaa.txt",
				"bbb",
				"ccc.txt",
				"ddd",
			},
			sortByExtension(0, filesAndDirs),
		)
	})

	t.Run("should sort by name DESC", func(t *testing.T) {
		assert.EqualValues(
			Expected{
				"ddd",
				"ccc.txt",
				"bbb",
				"aaa.txt",
			},
			sortByExtension(SortDesc, filesAndDirs),
		)
	})

	t.Run("should sort dirs & files ACS", func(t *testing.T) {
		assert.EqualValues(
			Expected{
				"bbb",
				"ddd",
				"aaa.txt",
				"ccc.txt",
			},
			Actual(SortTree{SortTreeConfig{Mode: SortByType}}.sort(Input{
				file("ccc.txt"),
				dir("bbb"),
				file("aaa.txt"),
				dir("ddd"),
			})),
		)
	})

	t.Run("should sort dirs & files DESC", func(t *testing.T) {
		assert.EqualValues(
			Expected{
				"ccc.txt",
				"aaa.txt",
				"ddd",
				"bbb",
			},
			sortByExtension(SortByType|SortDesc, Input{
				file("ccc.txt"),
				dir("bbb"),
				file("aaa.txt"),
				dir("ddd"),
			}),
		)
	})

	fileExtensions := []*dirtree.Node{
		file("ddd.mmm"),
		file("aaa.aaa"),
		file("bbb.mmm"),
		file("ttt.aaa"),
	}

	t.Run("should sort files by extension ASC", func(t *testing.T) {
		assert.EqualValues(
			Expected{
				"aaa.aaa",
				"ttt.aaa",
				"bbb.mmm",
				"ddd.mmm",
			},
			sortByExtension(SortByType, fileExtensions),
		)
	})

	t.Run("should sort files by extension DESC", func(t *testing.T) {
		assert.EqualValues(
			Expected{
				"ddd.mmm",
				"bbb.mmm",
				"ttt.aaa",
				"aaa.aaa",
			},
			sortByExtension(SortByType|SortDesc, fileExtensions),
		)
	})

	dotDirsAndFiles := []*dirtree.Node{
		dir("foo"),
		dir(".bar"),
		file("baz.txt"),
		file(".bat"),
	}

	t.Run("dot-files and dirs should be first", func(t *testing.T) {
		assert.EqualValues(
			Expected{
				".bar",
				"foo",
				".bat",
				"baz.txt",
			},
			sortByExtension(SortByType, dotDirsAndFiles),
		)
	})
}

func dir(name string) *dirtree.Node {
	return &dirtree.Node{
		Info: fstub.FileInfo{NAME: name, DIR: true},
	}
}

func file(name string) *dirtree.Node {
	return &dirtree.Node{
		Info: fstub.FileInfo{NAME: name, DIR: false},
	}
}

func sortByExtension(mode SortMode, nodes []*dirtree.Node) []string {
	input := append(nodes[:0:0], nodes...)
	return Actual(SortTree{SortTreeConfig{Mode: mode}}.sort(input))
}

func Actual(nodes []*dirtree.Node) (names []string) {
	for _, item := range nodes {
		names = append(names, item.Info.Name())
	}
	return names
}
