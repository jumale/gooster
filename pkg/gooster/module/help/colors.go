package help

//
//import (
//	"fmt"
//	"github.com/gdamore/tcell"
//	"github.com/jumale/gooster/pkg/gooster"
//	"github.com/rivo/tview"
//	"sort"
//)
//
//func NewColorNamesModule(cfg gooster.ModuleConfig) gooster.Module {
//	return &ColorNames{cfg: cfg}
//}
//
//type ColorNames struct {
//	cfg  gooster.ModuleConfig
//	view *tview.TextView
//	*gooster.AppContext
//}
//
//func (w *ColorNames) Name() string {
//	return "help_colors"
//}
//
//func (w *ColorNames) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
//	w.AppContext = ctx
//
//	w.view = tview.NewTextView()
//	w.view.SetTitle("Available colors")
//	w.view.SetBorder(true)
//	w.view.SetDynamicColors(true)
//	w.view.SetWordWrap(true)
//	w.view.SetBackgroundColor(tcell.ColorDefault)
//
//	type colType struct {
//		name  string
//		value tcell.Color
//	}
//	var colors []colType
//	for name, value := range tcell.ColorNames {
//		colors = append(colors, colType{name: name, value: value})
//	}
//
//	sort.SliceStable(colors, func(i, j int) bool {
//		return colors[i].name < colors[j].name
//	})
//
//	text := ""
//	for _, color := range colors {
//		text += fmt.Sprintf("[:%s]  %s  [:-]", color.name, color.name)
//	}
//
//	w.view.SetText(text)
//
//	return w.view, w.cfg, nil
//}
