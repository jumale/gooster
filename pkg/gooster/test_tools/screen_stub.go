package testtools

import (
	"fmt"
	"github.com/gdamore/tcell"
	"strings"
	"time"
)

type ScreenStub struct {
	width  int
	height int
	data   [][]rune
}

func NewScreenStub(width, height int) *ScreenStub {
	var data [][]rune
	for i := 0; i < height; i++ {
		data = append(data, make([]rune, width))
	}
	return &ScreenStub{
		width:  width,
		height: height,
		data:   data,
	}
}

func (scr *ScreenStub) GetDisplayString() string {
	var data []string
	for _, row := range scr.data {
		data = append(data, string(row))
	}
	return strings.Join(data, "\n")
}

func (scr *ScreenStub) GetDisplay() []byte {
	return []byte(scr.GetDisplayString())
}

func (scr *ScreenStub) Init() error {
	return nil
}

func (scr *ScreenStub) Fini() {

}

func (scr *ScreenStub) Clear() {

}

func (scr *ScreenStub) Fill(r rune, s tcell.Style) {
	for i, row := range scr.data {
		for j := range row {
			scr.data[i][j] = r
		}
	}
}

func (scr *ScreenStub) SetCell(x int, y int, style tcell.Style, ch ...rune) {
	if x < len(scr.data) && y < len(scr.data[x]) && len(ch) > 0 {
		scr.data[x][y] = ch[0]
	}
}

func (scr *ScreenStub) GetContent(x, y int) (mainc rune, combc []rune, style tcell.Style, width int) {
	if x < len(scr.data) && y < len(scr.data[x]) {
		return scr.data[x][y], nil, tcell.StyleDefault, 1
	} else {
		return 0, nil, tcell.StyleDefault, 0
	}
}

func (scr *ScreenStub) SetContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
	if x < len(scr.data) && y < len(scr.data[x]) {
		scr.data[x][y] = mainc
	}
}

func (scr *ScreenStub) SetStyle(style tcell.Style) {

}

func (scr *ScreenStub) ShowCursor(x int, y int) {

}

func (scr *ScreenStub) HideCursor() {

}

func (scr *ScreenStub) Size() (int, int) {
	return scr.width, scr.height
}

type stubEvent struct{}

func (s stubEvent) When() time.Time {
	return time.Now()
}

func (scr *ScreenStub) PollEvent() tcell.Event {
	return stubEvent{}
}

func (scr *ScreenStub) PostEvent(ev tcell.Event) error {
	return nil
}

func (scr *ScreenStub) PostEventWait(ev tcell.Event) {
}

func (scr *ScreenStub) EnableMouse() {

}

func (scr *ScreenStub) DisableMouse() {

}

func (scr *ScreenStub) HasMouse() bool {
	return false
}

func (scr *ScreenStub) Colors() int {
	return 256
}

func (scr *ScreenStub) Show() {
	fmt.Print(scr.GetDisplayString())
}

func (scr *ScreenStub) Sync() {
}

func (scr *ScreenStub) CharacterSet() string {
	return "UTF-8"
}

func (scr *ScreenStub) RegisterRuneFallback(r rune, subst string) {
}

func (scr *ScreenStub) UnregisterRuneFallback(r rune) {
}

func (scr *ScreenStub) CanDisplay(r rune, checkFallbacks bool) bool {
	return true
}

func (scr *ScreenStub) Resize(int, int, int, int) {

}

func (scr *ScreenStub) HasKey(tcell.Key) bool {
	return true
}
