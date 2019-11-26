package config

import (
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestColor(t *testing.T) {
	assert := require.New(t)

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("should decode predefined color names", func(t *testing.T) {
			var target Color
			err := (&target).UnmarshalJSON([]byte(`"red"`))
			assert.NoError(err)
			assert.Equal(tcell.ColorRed, target.Origin())
		})

		t.Run("should decode hex color", func(t *testing.T) {
			var target Color
			err := (&target).UnmarshalJSON([]byte(`"#00FF00"`))
			assert.NoError(err)
			assert.Equal(tcell.NewHexColor(0x00FF00), target.Origin())
		})

		t.Run("should try to decode hex color without hash prefix", func(t *testing.T) {
			var target Color
			err := (&target).UnmarshalJSON([]byte(`"2C2C2C"`))
			assert.NoError(err)
			assert.Equal(tcell.NewHexColor(0x2C2C2C), target.Origin())
		})

		t.Run("should return an error for invalid color value", func(t *testing.T) {
			var target Color
			err := (&target).UnmarshalJSON([]byte(`"foo"`))
			assert.Error(err)
		})
	})
}
