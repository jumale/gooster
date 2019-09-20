package dialog

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/log"
	"github.com/rivo/tview"
	"math"
)

const textBottomPadding = 1

type Text struct {
	Title   string
	Text    string
	Border  bool
	Buttons []Button
	Log     log.Logger
}

func (t Text) CreateBox(cfg Config, onDone ActionHandler) tview.Primitive {
	content := tview.NewTextView()
	content.SetBackgroundColor(tcell.ColorDefault)
	if t.Text != "" {
		content.SetText(t.Text)
	}
	content.SetRect(t.getContentRect())

	return Base{
		Title:   t.Title,
		Border:  t.Border,
		Buttons: t.Buttons,
		Log:     t.Log,
	}.CreateBox(content, cfg, onDone)
}

func (t Text) getContentRect() (int, int, int, int) {
	width := 0
	height := 0

	contentLen := len(t.Text)
	if len(t.Title) > contentLen {
		contentLen = len(t.Title)
	}

	if contentLen > width && contentLen < maxWidth {
		width = contentLen
		height = 1
	}
	if contentLen > maxWidth {
		width = maxWidth
		height = int(math.Ceil(float64(contentLen / maxWidth)))
	}

	if width < minWidth {
		width = minWidth
	}
	height += textBottomPadding

	t.debug("text size: %dx%d", width, height)

	return 0, 0, width, height
}

func (t Text) debug(msg string, args ...interface{}) {
	if t.Log != nil {
		t.Log.DebugF(msg, args...)
	}
}
