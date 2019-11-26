package config

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAlign(t *testing.T) {
	assert := require.New(t)

	t.Run("UnmarshalJSON", func(t *testing.T) {
		testCases := []struct {
			strVal   string
			constVal int
			hasErr   bool
		}{
			{"left", tview.AlignLeft, false},
			{"right", tview.AlignRight, false},
			{"center", tview.AlignCenter, false},
			{"", tview.AlignCenter, false},
			{"foo", 0, true},
		}

		for _, testCase := range testCases {
			var target Align
			err := (&target).UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, testCase.strVal)))
			assert.Equal(testCase.constVal, target.Origin())
			if testCase.hasErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		}
	})
}
