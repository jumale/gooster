package ext

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"strings"
	"sync"
	"time"
)

type TypingSearchConfig struct {
	KeyPressInterval time.Duration `json:"key_press_interval"`
}

type TypingSearch struct {
	search   string
	children []*dirtree.Node
	timer    *time.Timer
	cfg      TypingSearchConfig
	sync.Mutex
	gooster.Context
}

func NewTypingSearch() gooster.Extension {
	return &TypingSearch{cfg: TypingSearchConfig{
		KeyPressInterval: 400 * time.Millisecond,
	}}
}

func (ext *TypingSearch) Name() string {
	return "navigate"
}

func (ext *TypingSearch) Init(m gooster.Module, ctx gooster.Context) error {
	ext.Context = ctx
	if err := ctx.LoadConfig(&ext.cfg); err != nil {
		return err
	}

	ext.Events().Subscribe(events.HandleWithPrio(-100, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case workdir.EventSetChildren:
			ext.Lock()
			ext.children = event.Children
			ext.Unlock()
		}
		return e
	}))

	prev := m.View().GetBox().GetInputCapture()
	m.View().GetBox().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ext.navigate(event)
		return prev(event)
	})

	return nil
}

func (ext *TypingSearch) navigate(event *tcell.EventKey) {
	if event.Key() != tcell.KeyRune {
		return
	}

	if ext.timer == nil {
		ext.timer = time.NewTimer(ext.cfg.KeyPressInterval)
	} else {
		ext.timer.Reset(ext.cfg.KeyPressInterval)
	}
	go ext.clearSearch()

	ext.Lock()
	ext.search += string(event.Rune())
	ext.focusNode(ext.children, ext.search)
	ext.Unlock()
}

func (ext *TypingSearch) clearSearch() {
	<-ext.timer.C
	ext.Lock()
	ext.search = ""
	ext.Unlock()
}

func (ext *TypingSearch) focusNode(nodes []*dirtree.Node, search string) {
	//ext.log.DebugF("focus node `%s`", search)
	for _, child := range nodes {
		if strings.Contains(strings.ToLower(child.GetText()), search) {
			ext.Events().Dispatch(workdir.EventActivateNode{Path: child.Path})
			return
		}
	}
}
