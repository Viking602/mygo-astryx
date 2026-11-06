package astryx

import (
	"github.com/Viking602/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
)

// UI builds components in one MyGo build pass. Do not store it in app state.
type UI struct {
	c         *ui.Context
	tokens    theme.Theme
	closeMenu func()
	inGroup   bool
}

// New follows the current host or scoped theme without changing it.
func New(c *ui.Context) *UI { return &UI{c: c, tokens: theme.System(c)} }

// Key gives the next outer control a stable identity before it initializes
// state. It is equivalent to c.Key(id), and is intended for immediate chaining.
func (k *UI) Key(id any) *UI { k.c.Key(id); return k }

// Scope builds a themed subtree. Its floating panels inherit this theme too.
// The returned container can be sized with ordinary MyGo element methods.
func Scope(c *ui.Context, tokens theme.Theme, build func(*UI)) ui.Element {
	return ui.Themed(c, &tokens.Native).Children(func() {
		if build != nil {
			build(&UI{c: c, tokens: tokens})
		}
	})
}

// Tokens returns a copy of this builder's resolved tokens.
func (k *UI) Tokens() theme.Theme { return k.tokens }

type Size uint8

const (
	Medium Size = iota
	Small
	Large
)

type Variant uint8

const (
	Default Variant = iota
	Primary
	Outline
	Ghost
	Destructive
)

type Tone uint8

const (
	Neutral Tone = iota
	Info
	Success
	Warning
	Danger
)

func option[T any](o []T) (v T) {
	if len(o) > 0 {
		return o[len(o)-1]
	}
	return
}

func (k *UI) metrics(size Size) (height, pad, font float32) {
	unit := k.tokens.Space
	if unit <= 0 {
		unit = 4
	}
	switch size {
	case Small:
		return max(30, k.tokens.FontBody+14), unit * 3, k.tokens.FontSmall
	case Large:
		return max(44, k.tokens.FontBody+26), unit * 4.5, k.tokens.FontBody + 1
	default:
		return max(36, k.tokens.FontBody+20), unit * 3.5, k.tokens.FontBody
	}
}

func (k *UI) colors(t Tone) (ink, fill ui.Color) {
	switch t {
	case Info:
		return k.tokens.Info, k.tokens.InfoMuted
	case Success:
		return k.tokens.Native.Success, k.tokens.SuccessMuted
	case Warning:
		return k.tokens.Native.Warning, k.tokens.WarningMuted
	case Danger:
		return k.tokens.Native.Danger, k.tokens.DangerMuted
	default:
		return k.tokens.Native.Text, k.tokens.Native.Surface
	}
}
