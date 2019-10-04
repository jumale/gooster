package gooster

import "github.com/gdamore/tcell"

type InputCapture func(event *tcell.EventKey) (newEvent *tcell.EventKey)
type InputHandler func(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool)
