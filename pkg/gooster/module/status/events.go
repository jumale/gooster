package status

type EventShowInStatus struct {
	Col   int    // column number where to show the status value
	Value string // string value, may include string formatting
	Align int    // tview.Align* constants
}
