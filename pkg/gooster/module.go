package gooster

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Module interface {
	ModuleView

	// Config returns a basic module config
	// which specifies how ans where the module is displayed in the app.
	Config() ModuleConfig

	// Init initializes the module based on the provided AppContext.
	// It's like a constructor, and it's called once, when the module
	// is registered in the app.
	Init(*AppContext) error
}

type ModuleView interface {
	Box
	tview.Primitive
}

// ------------------------------------------------------------ //

type ModuleConfig struct {
	Position `json:",inline"`
	Enabled  bool `json:"enabled"`
	Focused  bool `json:"focused"`
	FocusKey tcell.Key
}

type Position struct {
	Col    int `json:"col"`
	Row    int `json:"row"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ------------------------------------------------------------ //

type InputAware interface {
	// GetInputCapture returns the function installed with SetInputCapture() or nil
	// if no such function has been installed.
	GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey

	// SetInputCapture installs a function which captures key events before they are
	// forwarded to the primitive's default key event handler. This function can
	// then choose to forward that key event (or a different one) to the default
	// handler by returning it. If nil is returned, the default handler will not
	// be called.
	//
	// Providing a nil handler will remove a previously existing handler.
	//
	// Note that this function will not have an effect on primitives composed of
	// other primitives, such as Form, Flex, or Grid. Key events are only captured
	// by the primitives that have focus (e.g. InputField) and only one primitive
	// can have focus at a time. Composing primitives such as Form pass the focus on
	// to their contained primitives and thus never receive any key events
	// themselves. Therefore, they cannot intercept key events.
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey)
}

type Box interface {
	InputAware

	// SetBorderPadding sets the size of the borders around the box content.
	SetBorderPadding(top, bottom, left, right int)

	// SetBackgroundColor sets the box's background color.
	SetBackgroundColor(color tcell.Color)

	// SetBorder sets the flag indicating whether or not the box should have a
	// border.
	SetBorder(show bool)

	// SetBorderColor sets the box's border color.
	SetBorderColor(color tcell.Color)

	// SetBorderAttributes sets the border's style attributes. You can combine
	// different attributes using bitmask operations:
	//
	//   box.SetBorderAttributes(tcell.AttrUnderline | tcell.AttrBold)
	SetBorderAttributes(attr tcell.AttrMask)

	// SetTitle sets the box's title.
	SetTitle(title string)

	// SetTitleColor sets the box's title color.
	SetTitleColor(color tcell.Color)

	// SetTitleAlign sets the alignment of the title, one of AlignLeft, AlignCenter,
	// or AlignRight.
	SetTitleAlign(align int)
}

// ------------------------------------------------------------ //

type BaseModule struct {
	*BoxAdaptor
	*AppContext
	tview.Primitive
	cfg ModuleConfig
}

func NewBaseModule(
	cfg ModuleConfig,
	ctx *AppContext,
	view tview.Primitive,
	box *tview.Box,
) *BaseModule {
	return &BaseModule{
		cfg:        cfg,
		AppContext: ctx,
		Primitive:  view,
		BoxAdaptor: &BoxAdaptor{Box: box},
	}
}

func (m *BaseModule) Config() ModuleConfig {
	return m.cfg
}

type KeyEventHandler func(event *tcell.EventKey) *tcell.EventKey
type KeyEventHandlers map[tcell.Key]KeyEventHandler

func (m *BaseModule) HandleKeyEvents(handlers KeyEventHandlers) {
	handleKeyEvents(m.BoxAdaptor, handlers)
}

func handleKeyEvents(target InputAware, handlers KeyEventHandlers) {
	prev := target.GetInputCapture()
	target.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if prev != nil {
			if ev = prev(ev); ev == nil {
				return nil
			}
		}

		if handler, ok := handlers[ev.Key()]; ok {
			if ev = handler(ev); ev == nil {
				return nil
			}
		}

		return ev
	})
}

// ------------------------------------------------------------ //

type appInputAdaptor struct {
	*tview.Application
}

func (a *appInputAdaptor) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	a.Application.SetInputCapture(capture)
}

// ------------------------------------------------------------ //

type BoxAdaptor struct {
	Box *tview.Box
}

func (b BoxAdaptor) GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return b.Box.GetInputCapture()
}

func (b BoxAdaptor) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) {
	b.Box.SetInputCapture(capture)
}

func (b BoxAdaptor) SetBorderPadding(top, bottom, left, right int) {
	b.Box.SetBorderPadding(top, bottom, left, right)
}

func (b BoxAdaptor) SetBackgroundColor(color tcell.Color) {
	b.Box.SetBackgroundColor(color)
}

func (b BoxAdaptor) SetBorder(show bool) {
	b.Box.SetBorder(show)
}

func (b BoxAdaptor) SetBorderColor(color tcell.Color) {
	b.Box.SetBorderColor(color)
}

func (b BoxAdaptor) SetBorderAttributes(attr tcell.AttrMask) {
	b.Box.SetBorderAttributes(attr)
}

func (b BoxAdaptor) SetTitle(title string) {
	b.Box.SetTitle(title)
}

func (b BoxAdaptor) SetTitleColor(color tcell.Color) {
	b.Box.SetTitleColor(color)
}

func (b BoxAdaptor) SetTitleAlign(align int) {
	b.Box.SetTitleAlign(align)
}

// ------------------------------------------------------------ //
