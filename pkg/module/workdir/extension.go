package workdir

import "github.com/jumale/gooster/pkg/dirtree"

type Extension interface {
	OnRefresh() dirtree.NodesHook
}

func (m *Module) AddExtension(ext Extension) *Module {
	m.ext = append(m.ext, ext)
	return m
}

func (m *Module) getTreeExtensionHooks() dirtree.NodesHook {
	return func(config dirtree.Config, nodes []*dirtree.Node) []*dirtree.Node {
		for _, ext := range m.ext {
			nodes = ext.OnRefresh()(config, nodes)
		}
		return nodes
	}
}
