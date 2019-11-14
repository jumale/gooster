package complete

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/rivo/tview"
	"math"
)

func (m *Module) handleSetCompletion(event gooster.EventSetCompletion) {
	m.view.Clear()

	list := event.Completion
	_, _, width, _ := m.GetRect()
	cols := numColsForList(list, width)
	if cols == 0 {
		return
	}

	row := -1
	for i := range list {
		col := i % cols
		if col == 0 {
			row += 1
		}
		m.view.SetCell(row, col, tview.NewTableCell(list[i]))
	}
}

func (m *Module) handleSelectNextItem(event *tcell.EventKey) *tcell.EventKey {

	return event
}

func (m *Module) handleMoveDown(event *tcell.EventKey) *tcell.EventKey {

	return event
}

func (m *Module) handleMoveLeft(event *tcell.EventKey) *tcell.EventKey {

	return event
}

func (m *Module) handleMoveRight(event *tcell.EventKey) *tcell.EventKey {

	return event
}

func numColsForList(list []string, boxWidth int) (numCols int) {
	if len(list) == 0 {
		return 0
	}

	maxItemWidth := 0
	for i := range list {
		if len(list[i]) > maxItemWidth {
			maxItemWidth = len(list[i])
		}
	}

	if maxItemWidth == 0 {
		return 0
	}

	return int(math.Floor(float64(boxWidth / maxItemWidth)))
}
