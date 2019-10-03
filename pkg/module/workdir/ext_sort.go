package workdir

import (
	"github.com/jumale/gooster/pkg/dirtree"
	"path"
	"sort"
	"strings"
)

type SortExtension struct {
	Mode SortMode
}

type SortMode uint8

const (
	SortByType SortMode = 1 << iota
	SortDesc
)

func (s SortExtension) Sort(nodes []*dirtree.Node) []*dirtree.Node {
	byType := s.Mode&SortByType != 0
	ASC := s.Mode&SortDesc == 0

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

func (s SortExtension) OnRefresh() dirtree.NodesHook {
	return func(config dirtree.Config, nodes []*dirtree.Node) []*dirtree.Node {
		return s.Sort(nodes)
	}
}
