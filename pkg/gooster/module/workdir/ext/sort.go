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

type SortTreeConfig struct {
	Mode SortMode `json:"mode"`
}

type SortTree struct {
	cfg SortTreeConfig
}

func NewSortTree() gooster.Extension {
	return &SortTree{cfg: SortTreeConfig{
		Mode: SortByType,
	}}
}

func (ext *SortTree) Name() string {
	return "sort"
}

func (ext *SortTree) Init(_ gooster.Module, ctx gooster.Context) error {
	if err := ctx.LoadConfig(&ext.cfg); err != nil {
		return err
	}

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
