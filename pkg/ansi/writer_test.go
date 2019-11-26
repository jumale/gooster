package ansi

import (
	"bytes"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestAsciiToTviewColors(t *testing.T) {
	assert := require.New(t)

	write := createWriteTester(t, WriterConfig{
		DefaultFg: tcell.ColorDefault,
		DefaultBg: tcell.ColorDefault,
	})

	t.Run("convert reset all", func(t *testing.T) {
		assert.Equal(`foo[-:-:-]`, write("foo\033[0m"))
	})
	t.Run("override reset colors", func(t *testing.T) {
		write := createWriteTester(t, WriterConfig{
			DefaultFg: tcell.ColorRed,
			DefaultBg: tcell.ColorBlue,
		})
		assert.Equal(`foo[red:blue:-]`, write("foo\033[0m"))
	})
	t.Run("convert reset fg color", func(t *testing.T) {
		assert.Equal(`foo[-]`, write("foo\033[39m"))
	})
	t.Run("convert reset bg color", func(t *testing.T) {
		assert.Equal(`foo[:-]`, write("foo\033[49m"))
	})
	t.Run("convert reset flags", func(t *testing.T) {
		assert.Equal(`foo[::-]`, write("foo\033[21m"))
		assert.Equal(`foo[::-]`, write("foo\033[22m"))
		assert.Equal(`foo[::-]`, write("foo\033[24m"))
		assert.Equal(`foo[::-]`, write("foo\033[25m"))
		assert.Equal(`foo[::-]`, write("foo\033[27m"))
	})

	t.Run("convert fg color", func(t *testing.T) {
		assert.Equal(`[red]foo`, write("\033[31mfoo"))
	})
	t.Run("convert bg color", func(t *testing.T) {
		assert.Equal(`[:blue]foo`, write("\033[44mfoo"))
	})
	t.Run("convert flag", func(t *testing.T) {
		assert.Equal(`[::b]foo`, write("\033[1mfoo"))
	})
	t.Run("apply custom color map", func(t *testing.T) {
		write := createWriteTester(t, WriterConfig{
			ColorMap: map[ColorId]tcell.Color{
				31: tcell.ColorYellow,
			},
		})
		assert.Equal(`[yellow]foo`, write("\033[31mfoo"))
	})

	t.Run("convert many", func(t *testing.T) {
		assert.Equal(`foo [white:red:br]bar`, write("foo \033[1;41;97;7mbar"))
	})

	t.Run("handle heading and intermediate dividers", func(t *testing.T) {
		assert.Equal(`foo [white:red]bar`, write("foo \033[;41;;97mbar"))
	})

	t.Run("ignore tag if it has tailing divider", func(t *testing.T) {
		assert.Equal("foo \033[41;97;mbar", write("foo \033[41;97;mbar"))
	})

	t.Run("ignore tag if it's empty", func(t *testing.T) {
		assert.Equal("foo \033[mbar", write("foo \033[mbar"))
	})

	t.Run("ignore tag if there are unsupported chars", func(t *testing.T) {
		assert.Equal("foo \033[1;41[97mbar", write("foo \033[1;41[97mbar"))
	})

	t.Run("don't be greedy, stop at the first closing case", func(t *testing.T) {
		assert.Equal("foo [:red:b]97mbar", write("foo \033[1;41m97mbar"))
	})
}

func createWriteTester(t *testing.T, cfg WriterConfig) func(string) string {
	assert := require.New(t)

	return func(v string) string {
		data := []byte(v)
		buf := bytes.NewBuffer([]byte{})
		writer := NewWriter(buf, cfg)
		_, err := writer.Write(data)

		assert.NoError(err)

		return buf.String()
	}
}

var line = "\033[1;4;107;31mLorem\033[24m ipsum\033[21m dolor\033[49m sit amet\033[39m, consectetur adipiscing elit.\033[0m\n"

// This is the benchmark for our custom parser,
// it should be faster or equal than the tview implementation.
func BenchmarkWriter(b *testing.B) {
	text := []byte(strings.Repeat(line, 1000))
	target := benchBuf{}
	writer := NewWriter(target, WriterConfig{})

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write(text)
	}
}

// This is the benchmark for the tview version of parser.
func BenchmarkTviewWriter(b *testing.B) {
	text := []byte(strings.Repeat(line, 1000))
	target := benchBuf{}
	writer := tview.ANSIWriter(target)

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write(text)
	}
}

type benchBuf struct {
}

func (b benchBuf) Write(p []byte) (n int, err error) {
	return len(p), nil
}
