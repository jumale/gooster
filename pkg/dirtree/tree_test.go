package dirtree

import (
	"github.com/jumale/gooster/pkg/filesys/fstub"
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestDirTree(t *testing.T) {
	assert := _assert.New(t)
	fsProps := &fstub.Props{
		WorkDir: "/wd",
		HomeDir: "/hd",
	}

	t.Run("Create a tree at the target dir", func(t *testing.T) {
		fs := fstub.New(fsProps).FromSchema(
			fstub.Dir("target",
				"baz.txt",
				".bashrc",
				fstub.Dir("foo",
					"mad.txt",
					fstub.Dir("sad"),
					"cat",
				),
				"bar",
			),
		)
		tree := newTree(fs, Config{})

		err := tree.Refresh("/wd/target")
		assert.NoError(err)

		t.Run("should render dir items sorted alphabetically", func(t *testing.T) {
			items := tree.Root().TreeNode.GetChildren()
			assert.Equal(".bashrc", items[0].GetText())
			assert.Equal("bar", items[1].GetText())
			assert.Equal("baz.txt", items[2].GetText())
			assert.Equal("foo", items[3].GetText())

			assert.Equal("/wd/target", tree.Path())
		})

		t.Run("should find node by abs path", func(t *testing.T) {
			fooNode := tree.Find("/wd/target/foo")
			assert.NotNil(fooNode)
			assert.Equal("/wd/target/foo", fooNode.Path)
			assert.Equal("foo", fooNode.Info.Name())

			t.Run("and expand the node", func(t *testing.T) {
				assert.Empty(fooNode.TreeNode.GetChildren())

				tree.ExpandNode(fooNode.TreeNode)
				children := fooNode.TreeNode.GetChildren()
				assert.Len(children, 3)
				assert.Equal("cat", children[0].GetText())
				assert.Equal("mad.txt", children[1].GetText())
				assert.Equal("sad", children[2].GetText())

				t.Run("and find nested node", func(t *testing.T) {
					currNode := tree.Find("foo/mad.txt")
					assert.NotNil(currNode)
					assert.Equal("/wd/target/foo/mad.txt", currNode.Path)
					assert.Equal("mad.txt", currNode.Info.Name())

					prevNode := tree.Find("foo/mad.txt", FindPrev)
					assert.NotNil(prevNode)
					assert.Equal("/wd/target/foo/cat", prevNode.Path)
					assert.Equal("cat", prevNode.Info.Name())

					nextNode := tree.Find("foo/mad.txt", FindNext)
					assert.NotNil(nextNode)
					assert.Equal("/wd/target/foo/sad", nextNode.Path)
					assert.Equal("sad", nextNode.Info.Name())
				})
			})
		})

		t.Run("should find node by relative path", func(t *testing.T) {
			node := tree.Find("foo")
			assert.NotNil(node)
			assert.Equal("/wd/target/foo", node.Path)
		})

		t.Run("should find nil for unknown path", func(t *testing.T) {
			node := tree.Find("/some/wrong/path")
			assert.Nil(node)
		})

		t.Run("should find next node", func(t *testing.T) {
			node := tree.Find("/wd/target/bar", FindNext)
			assert.NotNil(node)
			assert.Equal("/wd/target/baz.txt", node.Path)

			t.Run("or return nil when there is no next node", func(t *testing.T) {
				node := tree.Find("/wd/target/foo", FindNext)
				assert.Nil(node)
			})
		})

		t.Run("should find prev node", func(t *testing.T) {
			node := tree.Find("/wd/target/bar", FindPrev)
			assert.NotNil(node)
			assert.Equal("/wd/target/.bashrc", node.Path)

			t.Run("or return nil when there is no prev node", func(t *testing.T) {
				node := tree.Find("/wd/target/.bashrc", FindPrev)
				assert.Nil(node)
			})
		})
	})

}
