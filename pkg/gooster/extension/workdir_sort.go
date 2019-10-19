package extension

import (
	"github.com/jumale/gooster/pkg/dirtree"
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

type WorkDirSortConfig struct {
	gooster.ExtensionConfig
	Mode SortMode
}

type WorkDirSort struct {
	cfg WorkDirSortConfig
}

func NewWorkDirSort(cfg WorkDirSortConfig) gooster.Extension {
	return &WorkDirSort{cfg: cfg}
}

func (ext *WorkDirSort) Config() gooster.ExtensionConfig {
	return ext.cfg.ExtensionConfig
}

func (ext *WorkDirSort) Init(m gooster.Module, ctx *gooster.AppContext) error {
	workdir.Actions{AppContext: ctx}.ExtendSetChildren(ext.sort)
	return nil
}

func (ext WorkDirSort) sort(nodes []*dirtree.Node) []*dirtree.Node {
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
