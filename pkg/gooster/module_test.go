package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	_assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBaseModule(t *testing.T) {
	assert := _assert.New(t)

	t.Run("should create a base module and extract a Box from the view", func(t *testing.T) {
		ctx := &AppContext{}
		cfg := ModuleConfig{}
		view := tview.NewTextView()

		base := NewBaseModule(cfg, ctx, view, view.Box)
		base.SetBackgroundColor(tcell.ColorRed) // just should not fail when calling Box methods

		assert.Implements((*ModuleView)(nil), base)
		assert.Equal(cfg, base.Config())
		assert.Equal(ctx, base.AppContext)
		assert.Equal(view, base.Primitive)
		assert.Equal(view.Box, base.BoxAdaptor.Box)
	})

}
