package completion

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompletion(t *testing.T) {
	assert := require.New(t)

	t.Run("ApplyTo", func(t *testing.T) {
		t.Run("should return as is, if completion is not selected", func(t *testing.T) {
			cmp := Completion{}
			assert.Equal("cd /foo", cmp.ApplyTo("cd /foo"))
		})

		t.Run("should complete dir with adding slash", func(t *testing.T) {
			cmp := Completion{Type: TypeDir, Selected: "/foo/bar"}
			assert.Equal("cd /foo/bar/", cmp.ApplyTo("cd /f"))
		})

		t.Run("should complete value with adding corresponding suffix", func(t *testing.T) {
			testCases := []struct {
				completeType Type
				expected     string
			}{
				{TypeDir, "foo/"},
				{TypeCommand, "foo "},
				{TypeArg, "foo "},
				{TypeVar, "foo "},
				{TypeFile, "foo "},
				{TypeCustom, "foo "},
			}
			for _, testCase := range testCases {
				cmp := Completion{Type: testCase.completeType, Selected: "foo"}
				assert.Equal(testCase.expected, cmp.ApplyTo("f"))
			}
		})

		t.Run("should complete only latest argument", func(t *testing.T) {
			cmp := Completion{Selected: "foo"}
			assert.Equal("echo foo ", cmp.ApplyTo("echo fo"))
		})

		t.Run("should recognise escaped spaces in the last argument", func(t *testing.T) {
			cmp := Completion{Selected: `/foo\ bar/baz`}
			assert.Equal(`cd /foo\ bar/baz `, cmp.ApplyTo(`cd /foo\ ba`))
		})
	})
}
