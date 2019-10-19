package help

//
//func NewKeyNamesModule(cfg gooster.ModuleConfig) gooster.Module {
//	return &KeyNames{cfg: cfg}
//}
//
//type KeyNames struct {
//	cfg  gooster.ModuleConfig
//	view *tview.TextView
//	*gooster.AppContext
//}
//
//func (w *KeyNames) Name() string {
//	return "help_keys"
//}
//
//func (w *KeyNames) Init(ctx *gooster.AppContext) (tview.Primitive, gooster.ModuleConfig, error) {
//	w.AppContext = ctx
//
//	w.view = tview.NewTextView()
//	w.view.SetTitle("Available keys")
//	w.view.SetBorder(true)
//	w.view.SetWordWrap(true)
//	w.view.SetBackgroundColor(tcell.ColorDefault)
//
//	var names []string
//	for _, name := range tcell.KeyNames {
//		names = append(names, name)
//	}
//	sort.Strings(names)
//
//	text := ""
//	for _, name := range names {
//		text += fmt.Sprintf("%s  ", name)
//	}
//
//	w.view.SetText(text)
//
//	return w.view, w.cfg, nil
//}
