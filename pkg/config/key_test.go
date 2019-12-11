package config

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestKey(t *testing.T) {
	assert := require.New(t)

	t.Run("NewKey", func(t *testing.T) {

		testCases := []struct {
			origin   tcell.Key
			expected Key
		}{
			{tcell.KeyUp, Key{tcell.KeyUp, 0, 0}},
			{tcell.KeyF9, Key{tcell.KeyF9, 0, 0}},
			{tcell.KeyCtrlD, Key{tcell.KeyCtrlD, 'D', tcell.ModCtrl}},
			{tcell.KeyCtrlSpace, Key{tcell.KeyCtrlSpace, ' ', tcell.ModCtrl}},
			{tcell.KeyCtrlBackslash, Key{tcell.KeyCtrlBackslash, '\\', tcell.ModCtrl}},
		}

		for _, testCase := range testCases {
			k := NewKey(testCase.origin)
			assert.Equal(testCase.expected, k)
		}
	})

	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			val Key
			str string
		}{
			{Key{tcell.KeyUp, 0, 0}, "Up"},
			{Key{tcell.KeyUp, 0, tcell.ModAlt}, "Alt-Up"},
			{Key{tcell.KeyCtrlD, 'D', 0}, "Ctrl-D"},
			{Key{tcell.KeyCtrlSpace, ' ', tcell.ModAlt}, "Alt-Ctrl-Space"},
		}

		for _, testCase := range testCases {
			assert.Equal(testCase.str, testCase.val.String())
		}
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		decode := func(shortcut string) (Key, error) {
			var target Key
			err := (&target).UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, shortcut)))
			return target, err
		}

		t.Run("should correctly decode all predefined keys", func(t *testing.T) {
			// for every predefined key
			for keyVal, keyName := range tcell.KeyNames {
				// escape trailing backslash to not break JSON value
				encoded := keyName
				if strings.HasSuffix(encoded, "\\") {
					encoded = strings.Replace(encoded, "\\", "\\\\", 1)
				}

				actual, err := decode(encoded)
				assert.NoError(err, "Failed parsing key %s", encoded)
				assert.EqualValues(tcellKeyToExpectedVal(keyName, keyVal), actual, "Failed parsing key %s", encoded)
			}
		})

		t.Run("should decode a simple rune", func(t *testing.T) {
			actual, err := decode("a")
			assert.NoError(err)
			assert.Equal(Key{Type: tcell.KeyRune, Rune: 'a'}, actual)
		})

		t.Run("should decode a complex key", func(t *testing.T) {
			actual, err := decode("Shift-Meta-Alt-Ctrl-Space")
			assert.NoError(err)
			assert.Equal(Key{Type: tcell.KeyCtrlSpace, Mod: tcell.ModShift | tcell.ModMeta | tcell.ModAlt | tcell.ModCtrl}, actual)
		})

		t.Run("should fail if invalid json value", func(t *testing.T) {
			_, err := decode(`a"b`)
			assert.Error(err)
			assert.Contains(err.Error(), "invalid character")
		})

		t.Run("should fail if value does not match the shortcut pattern", func(t *testing.T) {
			_, err := decode(`some random value`)
			assert.Error(err)
			kpe := err.(KeyParseError)
			assert.Contains(kpe.Reason, "unexpected format")
		})

		t.Run("should fail if value key is not either predefined key nor valid rune", func(t *testing.T) {
			_, err := decode(`Alt-foo`)
			assert.Error(err)
			kpe := err.(KeyParseError)
			assert.Contains(kpe.Reason, "predefined key")
			assert.Contains(kpe.Reason, "valid rune")
		})

		t.Run("should fail if modifier is unknown", func(t *testing.T) {
			_, err := decode(`Alt-Foo-A`)
			assert.Error(err)
			ime := err.(InvalidModError)
			assert.Equal("Foo", ime.Mod)
		})

		t.Run("should fail if value does not support Ctrl modifier", func(t *testing.T) {
			_, err := decode(`Ctrl-Up`)
			assert.Error(err)
			ime := err.(InvalidModError)
			assert.Equal("Ctrl", ime.Mod)
		})

		t.Run("should not allow to apply Shift for runes (which produces another runes)", func(t *testing.T) {
			_, err := decode(`Shift-r`)
			assert.Error(err)
			ime := err.(InvalidModError)
			assert.Equal("Shift", ime.Mod)
		})

		t.Run("should fail if Shift modifier does not affect the final result", func(t *testing.T) {
			_, err := decode(`Shift-Ctrl-]`)
			assert.Error(err)
			ime := err.(InvalidModError)
			assert.Equal("Shift", ime.Mod)
		})
	})
}

func tcellKeyToExpectedVal(name string, val tcell.Key) Key {
	var mod tcell.ModMask
	var char rune
	if strings.HasPrefix(name, "Ctrl-") {
		mod |= tcell.ModCtrl
		char = rune(val)
	}
	return Key{Type: val, Rune: char, Mod: mod}
}
