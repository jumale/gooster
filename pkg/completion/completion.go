package completion

type Type uint8

const (
	space  rune = 32
	escape      = 92
)

const (
	TypeCommand Type = iota
	TypeArg
	TypeVar
	TypeDir
	TypeFile
	TypeCustom
)

type Completion struct {
	Type      Type
	Suggested []string
	Selected  string
}

func (e Completion) IsEmpty() bool {
	return len(e.Suggested) == 0
}

func (e Completion) IsUnique() bool {
	return len(e.Suggested) == 1
}

func (e Completion) Select(val string) Completion {
	e.Selected = val
	return e
}

func (e Completion) SelectByIndex(idx int) Completion {
	if idx < len(e.Suggested) {
		e.Selected = e.Suggested[idx]
	}
	return e
}

func (e Completion) SelectFirst() Completion {
	return e.SelectByIndex(0)
}

func (e Completion) ApplyTo(target string) string {
	if e.Selected == "" {
		return target
	}

	var suffix string
	switch e.Type {
	case TypeDir:
		suffix = "/"
	default:
		suffix = " "
	}

	for i := len(target) - 1; i >= 0; i -= 1 {
		if rune(target[i]) == space && (i-1 < 0 || rune(target[i-1]) != escape) {
			return target[:i+1] + e.Selected + suffix
		}
	}
	return e.Selected + suffix
}
