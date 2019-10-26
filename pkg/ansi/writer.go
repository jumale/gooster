package ansi

import (
	"fmt"
	"github.com/gdamore/tcell"
	"io"
	"strconv"
	"strings"
)

//
//var (
//	colorTagRegExp = regexp.MustCompile(`\033\[[\d;]+m`)
//)

type ColorId = int

type WriterConfig struct {
	ColorMap  map[ColorId]tcell.Color
	DefaultFg tcell.Color
	DefaultBg tcell.Color
}

type writer struct {
	target    io.Writer
	colorMap  map[ColorId]colorValue
	resetFg   string
	resetBg   string
	resetFl   string
	state     state
	currCode  []byte
	tagValues [3]string
}

func NewWriter(target io.Writer, cfg WriterConfig) io.Writer {
	colorMap := make(map[ColorId]colorValue)
	for id, value := range defaultColorMap {
		colorMap[30+id] = value
		colorMap[40+id] = value
	}
	if cfg.ColorMap != nil {
		for id, color := range cfg.ColorMap {
			colorMap[id] = colorVal(color)
		}
	}

	resetFg := "-"
	if cfg.DefaultFg != tcell.ColorDefault {
		resetFg = colorVal(cfg.DefaultFg)
	}
	resetBg := "-"
	if cfg.DefaultBg != tcell.ColorDefault {
		resetBg = colorVal(cfg.DefaultBg)
	}

	return &writer{
		target:   target,
		colorMap: colorMap,
		resetFg:  resetFg,
		resetBg:  resetBg,
		resetFl:  "-",
	}
}

// Converts all ASCII colors in the text to corresponding "tview" tags.
func (f *writer) Write(data []byte) (int, error) {
	var buf, tmp []byte

	reset := func() {
		f.state = closed
		buf = append(buf, tmp...)
		tmp = []byte{}
	}
	parseCode := func() {
		f.tagValues = f.applyColorCode(f.currCode, f.tagValues)
		f.currCode = []byte{}
	}

	for _, b := range data {
		if f.state == open {
			tmp = append(tmp, b)
			if b == chEscape {
				f.state = escaped
			} else {
				reset()
			}

		} else if f.state == escaped {
			tmp = append(tmp, b)
			if b == chDiv {
				if len(f.currCode) > 0 {
					parseCode()
				}
			} else if b == chClose {
				if len(f.currCode) > 0 {
					parseCode()
					buf = append(buf, f.createTag(f.tagValues[:]...)...)
					f.tagValues = [3]string{}
					f.state = closed
				} else {
					reset()
				}
			} else if 48 <= b && b <= 57 {
				f.currCode = append(f.currCode, b)
			} else {
				reset()
			}

		} else if b == chOpen {
			tmp = append(tmp, b)
			f.state = open

		} else {
			buf = append(buf, b)
		}

	}

	return f.target.Write(buf)

	//return f.target.Write(colorTagRegExp.ReplaceAllFunc(data, func(data []byte) []byte {
	//	var tagValues [3]string
	//	for _, value := range bytes.Split(data[2:len(data)-1], []byte(";")) {
	//		f.applyColorCode(value, tagValues)
	//	}
	//	return f.createTag(tagValues[:]...)
	//}))
}

func (f *writer) applyColorCode(code []byte, tagVals [3]string) [3]string {
	id, err := strconv.Atoi(string(code))
	if err != nil {
		return tagVals
	}

	switch id {
	case 0:
		tagVals[0] = f.resetFg
		tagVals[1] = f.resetBg
		tagVals[2] = f.resetFl
	case 39:
		tagVals[0] = f.resetFg
	case 49:
		tagVals[1] = f.resetFg
	case 21, 22, 24, 25, 27:
		tagVals[2] = f.resetFl
	case 1, 2, 4, 5, 7:
		tagVals[2] += flagsMap[id]
	case 30, 31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96, 97:
		tagVals[0] = f.colorMap[id]
	case 40, 41, 42, 43, 44, 45, 46, 47, 100, 101, 102, 103, 104, 105, 106, 107:
		tagVals[1] = f.colorMap[id]
	}

	return tagVals
}

func (f *writer) createTag(values ...string) []byte {
	if values[2] == "" {
		values = values[:len(values)-1]
		if values[1] == "" {
			values = values[:len(values)-1]
		}
	}
	return []byte("[" + strings.Join(values, ":") + "]")
}

type colorValue = string

func colorVal(c tcell.Color) colorValue {
	for name, value := range tcell.ColorNames {
		if value == c {
			return name
		}
	}
	return fmt.Sprintf("#%06x", c.Hex())
}

const (
	chOpen   = 27
	chEscape = 91
	chDiv    = 59
	chClose  = 109
)

type state int

const (
	closed state = iota
	open
	escaped
)

var defaultColorMap = map[int]colorValue{
	0:  colorVal(tcell.ColorBlack),
	1:  colorVal(tcell.ColorRed),
	2:  colorVal(tcell.ColorGreen),
	3:  colorVal(tcell.ColorYellow),
	4:  colorVal(tcell.ColorBlue),
	5:  colorVal(tcell.ColorDarkMagenta), // magenta
	6:  colorVal(tcell.ColorDarkCyan),
	7:  colorVal(tcell.ColorLightGray),
	60: colorVal(tcell.ColorDarkGray),
	61: colorVal(tcell.ColorIndianRed), // light red
	62: colorVal(tcell.ColorLightGreen),
	63: colorVal(tcell.ColorLightYellow),
	64: colorVal(tcell.ColorLightBlue),
	65: colorVal(tcell.ColorMistyRose), // light magenta
	66: colorVal(tcell.ColorLightCyan),
	67: colorVal(tcell.ColorWhite),
}

var flagsMap = map[ColorId]string{
	1: "b",
	2: "d",
	4: "u",
	5: "l",
	7: "r",
	// 8: "???", 'hidden' is not implemented in tview
}
