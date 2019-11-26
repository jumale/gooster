package complete

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/command"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/gooster/module/prompt"
	"github.com/rivo/tview"
	"math"
)

func (m *Module) handleSetCompletion(event gooster.EventSetCompletion) {
	m.current = event
	m.view.Clear()

	list := event.Completion
	_, _, width, _ := m.view.GetRect()
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

	if len(event.Completion) > 1 {
		m.Events().Dispatch(gooster.EventSetFocus{Target: m.view})
	}
}

func (m *Module) handleNextItem(event *tcell.EventKey) *tcell.EventKey {
	row, col := m.view.GetSelection()
	col += 1
	if col >= m.view.GetColumnCount() {
		col = 0
		row += 1
	}
	if row >= m.view.GetRowCount() {
		row = 0
	}
	m.Log().DebugF("Select %d %d", row, col)
	m.view.Select(row, col)
	return event
}

func (m *Module) handleSelectItem(event *tcell.EventKey) *tcell.EventKey {
	selected := m.view.GetCell(m.view.GetSelection()).Text
	m.view.Clear()
	m.Events().Dispatch(prompt.EventSetPrompt{
		Input: command.ApplyCompletion(m.current.Input, selected+" "),
		Focus: true,
	})
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
