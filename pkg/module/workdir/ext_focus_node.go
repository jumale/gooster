package workdir

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"strings"
	"sync"
	"time"
)

type FocusNodeExtensionConfig struct {
	KeyPressInterval time.Duration
	Log              log.Logger
}

type FocusNodeExtension struct {
	search string
	timer  *time.Timer
	cfg    FocusNodeExtensionConfig
	log    log.Logger
	sync.Mutex
}

func NewFocusNodeExtension(cfg FocusNodeExtensionConfig) *FocusNodeExtension {
	var logger log.Logger = log.EmptyLogger{}
	if cfg.Log != nil {
		logger = cfg.Log
	}

	return &FocusNodeExtension{
		cfg: cfg,
		log: logger,
	}
}

func (ext FocusNodeExtension) OnRefresh(ExtendableTree) dirtree.NodesHook {
	return nil
}

func (ext *FocusNodeExtension) OnInput(tree ExtendableTree) gooster.InputHandler {
	return func(event *tcell.EventKey) (newEvent *tcell.EventKey, handled bool) {
		if event.Key() != tcell.KeyRune {
			return event, false
		}

		if ext.timer == nil {
			ext.timer = time.NewTimer(ext.cfg.KeyPressInterval)
		} else {
			ext.timer.Reset(ext.cfg.KeyPressInterval)
		}
		go ext.clearSearch()

		ext.Lock()
		ext.search += string(event.Rune())
		ext.focusNode(tree, ext.search)
		ext.Unlock()

		return event, true
	}
}

func (ext *FocusNodeExtension) clearSearch() {
	<-ext.timer.C
	ext.Lock()
	ext.log.Debug("clearing focus search")
	ext.search = ""
	ext.Unlock()
}

func (ext FocusNodeExtension) focusNode(tree ExtendableTree, search string) {
	//ext.log.DebugF("focus node `%s`", search)
	for _, child := range tree.GetRoot().GetChildren() {
		if strings.Contains(strings.ToLower(child.GetText()), search) {
			tree.SetCurrentNode(child)
			return
		}
	}
}
