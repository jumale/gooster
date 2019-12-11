package config

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	"regexp"
	"strings"
)

var keyValues = getKeyValues()

type Key struct {
	Type tcell.Key
	Rune rune
	Mod  tcell.ModMask
}

func NewKey(k tcell.Key) Key {
	var r rune
	var m tcell.ModMask
	if name, ok := tcell.KeyNames[k]; ok {
		if strings.HasPrefix(name, "Ctrl-") {
			m |= tcell.ModCtrl
			ch := strings.ReplaceAll(name, "Ctrl-", "")
			if ch == "Space" {
				ch = " "
			}
			r = []rune(ch)[0]
		}
	}
	return Key{Type: k, Rune: r, Mod: m}
}

func (k Key) String() string {
	return strings.ReplaceAll(tcell.NewEventKey(k.Type, k.Rune, k.Mod).Name(), "+", "-")
}

var keyMatcher = regexp.MustCompile(`^((?:[^-]+-)*)(-|[^\s-]+)$`) // @todo simplify

func (k *Key) UnmarshalJSON(b []byte) (err error) {
	var name string
	if err = json.Unmarshal(b, &name); err != nil {
		return err
	}

	*k, err = k.parseShortcut(name)
	return err
}

func (k *Key) parseShortcut(name string) (Key, error) {
	matches := keyMatcher.FindStringSubmatch(name)
	if matches == nil {
		return Key{}, KeyParseError{Key: name, Reason: "unexpected format, expected like 'Shift-Alt-Up'"}
	}

	mods := matches[1]
	var mod tcell.ModMask
	for _, modVal := range strings.Split(mods, "-") {
		if mask := getModMask(modVal); mask != tcell.ModNone {
			mod |= mask
		} else if modVal != "" {
			return Key{}, InvalidModError{Key: name, Mod: modVal, Reason: "unknown modifier"}
		}
	}

	keyName := matches[2]
	if mod&tcell.ModCtrl > 0 {
		if key, ok := keyValues["Ctrl-"+keyName]; ok {
			return k.createKey(key, rune(key), mod)
		} else {
			return Key{}, InvalidModError{Key: name, Mod: "Ctrl", Reason: noEffect}
		}
	}

	if key, ok := keyValues[keyName]; ok {
		return k.createKey(key, 0, mod)
	}

	keyRune := []rune(keyName)
	if len(keyRune) != 1 {
		return Key{}, KeyParseError{Key: name, Reason: "invalid target key, expected either one of the predefined keys, or a valid rune"}
	}

	return k.createKey(tcell.KeyRune, keyRune[0], mod)
}

func (k Key) createKey(keyType tcell.Key, char rune, mod tcell.ModMask) (Key, error) {
	key := Key{Type: keyType, Rune: char, Mod: mod}

	if char != 0 && mod&tcell.ModShift > 0 {
		if keyType == tcell.KeyRune {
			return key, InvalidModError{Key: key.String(), Mod: "Shift", Reason: "it produces another rune, use the resulting rune instead"}

		} else if int16(keyType) == int16(char) {
			return key, InvalidModError{Key: key.String(), Mod: "Shift", Reason: noEffect}
		}
	}

	return key, nil
}

func (k Key) Empty() bool {
	return k.Type == 0
}

func getModMask(val string) tcell.ModMask {
	switch strings.ToLower(val) {
	case "shift":
		return tcell.ModShift
	case "alt":
		return tcell.ModAlt
	case "meta":
		return tcell.ModMeta
	case "ctrl":
		return tcell.ModCtrl
	default:
		return tcell.ModNone
	}
}

func getKeyValues() map[string]tcell.Key {
	vals := make(map[string]tcell.Key)
	for key, name := range tcell.KeyNames {
		vals[name] = key
	}
	return vals
}

const noEffect = "it does not affect the result, skip it"

type InvalidModError struct {
	Key    string
	Mod    string
	Reason string
}

func (err InvalidModError) Error() string {
	return fmt.Sprintf("Invalid modifier %s in %s. Reason: %s", err.Mod, err.Key, err.Reason)
}

type KeyParseError struct {
	Key    string
	Reason string
}

func (err KeyParseError) Error() string {
	return fmt.Sprintf("Failed to parse key %s. Reason: %s", err.Key, err.Reason)
}
