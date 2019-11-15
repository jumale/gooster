package ext

import (
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/workdir"
	"path"
	"sort"
	"strings"
)

type SortMode uint8

const (
	SortByType SortMode = 1 << iota
	SortDesc
)

type SortTreeConfig struct {
	gooster.ExtensionConfig
	Mode SortMode
}

type SortTree struct {
	cfg SortTreeConfig
}

func NewSortTree(cfg SortTreeConfig) gooster.Extension {
	return &SortTree{cfg: cfg}
}

func (ext *SortTree) Config() gooster.ExtensionConfig {
	return ext.cfg.ExtensionConfig
}

func (ext *SortTree) Init(m gooster.Module, ctx *gooster.AppContext) error {
	ctx.Events().Subscribe(events.HandleWithPrio(100, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case workdir.EventSetChildren:
			event.Children = ext.sort(event.Children)
			return event
		}
		return e
	}))
	return nil
}

func (ext SortTree) sort(nodes []*dirtree.Node) []*dirtree.Node {
	byType := ext.cfg.Mode&SortByType != 0
	ASC := ext.cfg.Mode&SortDesc == 0

	sort.SliceStable(nodes, func(i, j int) bool {
		a := nodes[i].Info
		b := nodes[j].Info

		if byType && a.IsDir() != b.IsDir() {
			return a.IsDir() == ASC
		}

		aName := a.Name()
		bName := b.Name()
		aIsDot := strings.HasPrefix(aName, ".")
		bIsDot := strings.HasPrefix(bName, ".")
		if byType && aIsDot != bIsDot {
			return aIsDot == ASC
		}

		aExt := path.Ext(aName)
		bExt := path.Ext(bName)
		if byType && aExt != bExt {
			return aExt < bExt == ASC
		}

		return aName < bName == ASC
	})
	return nodes
}
