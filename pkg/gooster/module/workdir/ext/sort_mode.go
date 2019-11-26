package ext

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

type SortMode uint8

const (
	SortByType SortMode = 1 << iota
	SortDesc
)

var sortModeMap = map[string]SortMode{
	"sort_by_type": SortByType,
	"sort_desc":    SortDesc,
}

func (s *SortMode) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}

	vals := strings.Split(val, "|")
	for i := range vals {
		key := strings.Trim(vals[i], " ")
		if val, ok := sortModeMap[key]; ok {
			*s = *s | val
		} else {
			return errors.Errorf("Unknown sort mode '%s'", key)
		}
	}
	return nil
}
