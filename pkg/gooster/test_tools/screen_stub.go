package testtools

import (
	"fmt"
	"github.com/gdamore/tcell"
	"strconv"
	"strings"
	"time"
)

type cell struct {
	r rune
	s tcell.Style
}

type screenStub struct {
	cells  [][]cell
	style  tcell.Style
	width  int
	height int
}

func NewScreenStub(width, height int) *screenStub {
	cells := make([][]cell, height)
	for i := range cells {
		cells[i] = make([]cell, width)
	}
	return &screenStub{
		width:  width,
		height: height,
		cells:  cells,
	}
}

func (scr *screenStub) GetView() string {
	return strings.Join(scr.GetViewLines(), "\n")
}

func (scr *screenStub) GetViewLines() []string {
	fg := tcell.ColorDefault
	bg := tcell.ColorDefault
	attr := tcell.AttrNone

	styleTag := func(s tcell.Style) string {
		styleTag := make([]string, 3)
		cFg, cBg, cAttr := s.Decompose()
		if fg != cFg {
			fg = cFg
			styleTag[0] = colorName(fg)
		}
		if bg != cBg {
			bg = cBg
			styleTag[1] = colorName(bg)
		}
		if cAttr != attr {
			attr = cAttr
			styleTag[2] = attrValues(attr)
		}
		return string(createTag(styleTag...))
	}

	var lines []string

	for _, row := range scr.cells {
		var line string

		for _, c := range row {
			line += styleTag(c.s)
			line += string(c.r)
		}

		lines = append(lines, line)
	}

	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	for i, line := range lines {
		lines[i] = fmt.Sprintf("%-"+strconv.Itoa(maxWidth)+"s", line)
	}

	return lines
}

func (scr *screenStub) Fill(r rune, s tcell.Style) {
	for i, row := range scr.cells {
		for j := range row {
			scr.cells[i][j] = cell{r: r, s: s}
		}
	}
}

func (scr *screenStub) SetCell(x int, y int, style tcell.Style, ch ...rune) {
	if len(ch) > 0 {
		scr.SetContent(x, y, ch[0], nil, style)
	}
}

func (scr *screenStub) SetContent(x int, y int, mainc rune, _ []rune, style tcell.Style) {
	if y >= len(scr.cells) || x >= len(scr.cells[y]) {
		return
	}

	if style == tcell.StyleDefault {
		style = scr.style
	}
	scr.cells[y][x] = cell{r: mainc, s: style}
}

func (scr *screenStub) GetContent(x, y int) (mainc rune, combc []rune, style tcell.Style, width int) {
	if y < len(scr.cells) && x < len(scr.cells[y]) {
		c := scr.cells[y][x]
		return c.r, nil, c.s, 1
	} else {
		return 0, nil, tcell.StyleDefault, 0
	}
}

func (scr *screenStub) SetStyle(style tcell.Style) {
	scr.style = style
}

func (scr *screenStub) Size() (int, int) {
	return scr.width, scr.height
}

func (scr *screenStub) PollEvent() tcell.Event {
	return stubEvent{}
}

func (scr *screenStub) PostEvent(tcell.Event) error {
	return nil
}

func (scr *screenStub) Init() error {
	return nil
}

func (scr *screenStub) HasMouse() bool {
	return false
}

func (scr *screenStub) Colors() int {
	return 256
}

func (scr *screenStub) Show() {
	fmt.Print(scr.GetView())
}

func (scr *screenStub) CharacterSet() string {
	return "UTF-8"
}

func (scr *screenStub) CanDisplay(rune, bool) bool {
	return true
}

func (scr *screenStub) HasKey(tcell.Key) bool {
	return true
}

func (scr *screenStub) Fini()                             {}
func (scr *screenStub) Clear()                            {}
func (scr *screenStub) ShowCursor(int, int)               {}
func (scr *screenStub) HideCursor()                       {}
func (scr *screenStub) PostEventWait(tcell.Event)         {}
func (scr *screenStub) EnableMouse()                      {}
func (scr *screenStub) DisableMouse()                     {}
func (scr *screenStub) Sync()                             {}
func (scr *screenStub) RegisterRuneFallback(rune, string) {}
func (scr *screenStub) UnregisterRuneFallback(rune)       {}
func (scr *screenStub) Resize(int, int, int, int)         {}

type stubEvent struct{}

func (s stubEvent) When() time.Time {
	return time.Now()
}

func colorName(c tcell.Color) string {
	if c == tcell.ColorDefault {
		return "-"
	}
	for name, val := range tcell.ColorNames {
		if val == c {
			return name
		}
	}
	return fmt.Sprintf("#%06x", c.Hex())
}

func attrValues(a tcell.AttrMask) string {
	if a == tcell.AttrNone {
		return "-"
	}

	var v string
	if a&tcell.AttrBold != 0 {
		v += "b"
	}
	if a&tcell.AttrBlink != 0 {
		v += "l"
	}
	if a&tcell.AttrReverse != 0 {
		v += "r"
	}
	if a&tcell.AttrUnderline != 0 {
		v += "u"
	}
	if a&tcell.AttrDim != 0 {
		v += "d"
	}
	return v
}

func createTag(values ...string) []byte {
	if values[2] == "" {
		values = values[:len(values)-1]
		if values[1] == "" {
			values = values[:len(values)-1]
			if values[0] == "" {
				return nil
			}
		}
	}
	return []byte("[" + strings.Join(values, ":") + "]")
}
