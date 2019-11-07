package cmd

import (
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestBashCompletion(t *testing.T) {
	assert := _assert.New(t)
	parse := ParseCommands

	t.Run("should parse command with mixed quotes", func(t *testing.T) {
		input := `echo -e "foo \n\"bar\" 'baz'" 'some "other \' quoted" string' --cat="fat"`
		result, err := parse(input)
		assert.NoError(err)
		assert.Len(result, 1)
		assert.Equal(Definition{
			"echo",
			[]string{
				"-e",
				`foo \n\"bar\" 'baz'`,           // double-quoted argument, parent quotes removed
				`some "other \' quoted" string`, // single-quoted argument, parent quotes removed
				`--cat="fat"`,
			},
		}, result[0])
	})

	t.Run("should parse multiple chained commands", func(t *testing.T) {
		input := `echo "foo; bar" | grep "foo | bar"; tail`
		result, err := parse(input)
		assert.NoError(err)
		assert.Len(result, 3)

		assert.Equal(Definition{"echo", []string{"foo; bar"}}, result[0])
		assert.Equal(Definition{"grep", []string{"foo | bar"}}, result[1])
		assert.Equal(Definition{"tail", nil}, result[2])

		t.Run("even without spaces", func(t *testing.T) {
			input := `echo "foo; bar"|grep "foo | bar";tail`
			result, err := parse(input)
			assert.NoError(err)
			assert.Len(result, 3)

			assert.Equal(Definition{"echo", []string{"foo; bar"}}, result[0])
			assert.Equal(Definition{"grep", []string{"foo | bar"}}, result[1])
			assert.Equal(Definition{"tail", nil}, result[2])
		})
	})

	t.Run("should also parse non-finished commands, but also return an error", func(t *testing.T) {
		input := `echo "foo \"te`
		result, err := parse(input)
		assert.Equal(QuoteErr, err)
		assert.Len(result, 1)
		assert.Equal(Definition{"echo", []string{`"foo \"te`}}, result[0])
	})
}

/*
1450-1500 ns/op
*/
var results []Definition

func BenchmarkParseCommandRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		results, _ = ParseCommands(`echo "foo; bar" | grep "foo | bar"; tail`)
	}
}
