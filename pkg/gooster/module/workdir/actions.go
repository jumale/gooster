package workdir

import (
	"github.com/jumale/gooster/pkg/dirtree"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
)

const (
	ActionRefresh      events.EventId = "workdir:refresh"
	ActionChangeDir                   = "workdir:change_dir"
	ActionSetChildren                 = "workdir:set_children"
	ActionActivateNode                = "workdir:activate_node"
	ActionCreateFile                  = "workdir:create_file"
	ActionCreateDir                   = "workdir:create_dir"
	ActionViewFile                    = "workdir:view_file"
	ActionDelete                      = "workdir:delete"
	ActionOpen                        = "workdir:open"
)

type PayloadSetChildren struct {
	Target   *tview.TreeNode
	Children []*dirtree.Node
}

type PayloadActivateNode struct {
	Path string
	Mode dirtree.FindMode
}

type Actions struct {
	*gooster.AppContext
}

func (a Actions) Refresh() {
	a.Events().Dispatch(events.Event{Id: ActionRefresh})
}

func (a Actions) ChangeDir(path string) {
	a.Events().Dispatch(events.Event{Id: ActionChangeDir, Payload: path})
}

func (a Actions) SetChildren(target *tview.TreeNode, children []*dirtree.Node) {
	a.Events().Dispatch(events.Event{
		Id: ActionSetChildren,
		Payload: PayloadSetChildren{
			Target:   target,
			Children: children,
		},
	})
}

func (a Actions) ExtendSetChildren(callback func(nodes []*dirtree.Node) []*dirtree.Node) {
	a.Events().Extend(events.Extension{
		Id: ActionSetChildren,
		Fn: func(data events.EventPayload) (newData events.EventPayload) {
			payload := data.(PayloadSetChildren)
			return PayloadSetChildren{
				Target:   payload.Target,
				Children: callback(payload.Children),
			}
		},
	})
}

func (a Actions) ActivateNode(path string) {
	a.Events().Dispatch(events.Event{
		Id:      ActionActivateNode,
		Payload: PayloadActivateNode{Path: path, Mode: dirtree.FindCurrent},
	})
}

func (a Actions) ActivatePrevNode(path string) {
	a.Events().Dispatch(events.Event{
		Id:      ActionActivateNode,
		Payload: PayloadActivateNode{Path: path, Mode: dirtree.FindPrev},
	})
}

func (a Actions) ActivateNextNode(path string) {
	a.Events().Dispatch(events.Event{
		Id:      ActionActivateNode,
		Payload: PayloadActivateNode{Path: path, Mode: dirtree.FindNext},
	})
}

func (a Actions) CreateFile(name string) {
	a.Events().Dispatch(events.Event{Id: ActionCreateFile, Payload: name})
}

func (a Actions) CreateDir(dirPath string) {
	a.Events().Dispatch(events.Event{Id: ActionCreateDir, Payload: dirPath})
}

func (a Actions) ViewFile(path string) {
	a.Events().Dispatch(events.Event{Id: ActionViewFile, Payload: path})
}

func (a Actions) Delete(path string) {
	a.Events().Dispatch(events.Event{Id: ActionDelete, Payload: path})
}

func (a Actions) Open(path string) {
	a.Events().Dispatch(events.Event{Id: ActionOpen, Payload: path})
}
