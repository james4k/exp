package ui

import "strings"

// Key represents a key and any modifier keys, in a succinct text
// format...which is to be documented, but is mostly inspired by Cocoa's
// keybinding format with some modifications to be print and programmer
// friendly. An example: "^f" represents Control + F.
type Key string

func (k Key) contains(c rune) bool {
	return strings.ContainsRune(string(k), c)
}

// Ctrl checks if Key contains the Control modifier key, represented by '^'.
func (k Key) Ctrl() bool {
	return k.contains('^')
}

// Alt checks if Key contains the Alt modifier key, represented by '~'.
func (k Key) Alt() bool {
	return k.contains('~')
}

// Shift checks if Key contains the Shift modifier key, represented by '~'.
func (k Key) Shift() bool {
	return k.contains('$')
}

// Keypad extracts a numeric keypad value from Key, preceded by '#'.
func (k Key) Keypad() (c byte, ok bool) {
	idx := strings.IndexRune(string(k), '#')
	if idx < 0 || idx >= len(k) {
		return 0, false
	}
	return k[idx+1], true
}

const (
	Escape    Key = "(esc)"
	Enter         = "(enter)"
	Tab           = "(tab)"
	Backspace     = "(bs)"
	Insert        = "(ins)"
	Delete        = "(del)"
	Left          = "(left)"
	Up            = "(up)"
	Right         = "(right)"
	Down          = "(down)"
	PageUp        = "(pgup)"
	PageDown      = "(pgdn)"
	PageHome      = "(home)"
	PageEnd       = "(end)"
	F1            = "(f1)"
	F2            = "(f2)"
	F3            = "(f3)"
	F4            = "(f4)"
	F5            = "(f5)"
	F6            = "(f6)"
	F7            = "(f7)"
	F8            = "(f8)"
	F9            = "(f9)"
	F10           = "(f10)"
	F11           = "(f11)"
	F12           = "(f12)"
	F13           = "(f13)"
	F14           = "(f14)"
	F15           = "(f15)"
	F16           = "(f16)"
	F17           = "(f17)"
	F18           = "(f18)"
	F19           = "(f19)"
	F20           = "(f20)"
	F21           = "(f21)"
	F22           = "(f22)"
	F23           = "(f23)"
	F24           = "(f24)"
	F25           = "(f25)"
)
