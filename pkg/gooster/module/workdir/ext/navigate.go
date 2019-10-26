package ext

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"strings"
	"sync"
	"time"
)

type TypingSearchConfig struct {
	gooster.ExtensionConfig
	KeyPressInterval time.Duration
}

type TypingSearch struct {
	search   string
	children []*dirtree.Node
	timer    *time.Timer
	cfg      TypingSearchConfig
	workDir  workdir.Actions
	sync.Mutex
	*gooster.AppContext
}

func NewTypingSearch(cfg TypingSearchConfig) gooster.Extension {
	return &TypingSearch{cfg: cfg}
}

func (ext *TypingSearch) Config() gooster.ExtensionConfig {
	return ext.cfg.ExtensionConfig
}

func (ext *TypingSearch) Init(m gooster.Module, ctx *gooster.AppContext) error {
	ext.AppContext = ctx
	ext.workDir = workdir.Actions{AppContext: ctx}

	ext.workDir.ExtendSetChildren(func(nodes []*dirtree.Node) []*dirtree.Node {
		ext.Lock()
		ext.children = nodes
		ext.Unlock()
		return nodes
	})

	prev := m.GetInputCapture()
	m.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			ext.workDir.ActivateNode(child.Path)
			return
		}
	}
}
