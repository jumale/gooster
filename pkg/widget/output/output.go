package output

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"strings"
	"sync"
)

type Config struct {
	gooster.WidgetConfig `json:",inline"`
}

func NewWidget(cfg Config) *Widget {
	return &Widget{
		cfg:   cfg,
		Mutex: &sync.Mutex{},
	}
}

type Widget struct {
	cfg  Config
	view *tview.TextView
	text string
	*sync.Mutex
}

func (w *Widget) Name() string {
	return "Console Output"
}

func (w *Widget) Init(ctx *gooster.AppContext) error {
	w.view = tview.NewTextView()
	w.view.SetBorder(false)
	w.view.SetDynamicColors(true)
	w.view.SetBorderPadding(0, 0, 1, 1)
	w.view.SetBackgroundColor(tcell.ColorDefault)

	ctx.EventManager.Subscribe(gooster.EventOutputMessage, func(event gooster.Event) {
		fmt.Printf("Output: %+v\n", event)
		w.addText(event.Data.(string))
	})

	return nil
}

func (w *Widget) addText(text string) {
	w.Lock()
	defer w.Unlock()

	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	w.text += text
	w.view.SetText(w.text)
}

func (w *Widget) View() tview.Primitive {
	return w.view
}

func (w *Widget) Config() gooster.WidgetConfig {
	return w.cfg.WidgetConfig
}
