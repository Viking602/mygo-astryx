package astryx

import (
	"github.com/egoist/mygo/ui"
	"time"
)

type HoverCardOptions struct {
	Delay, HideDelay      time.Duration
	Disabled, OpenOnClick bool
	FocusOnOpen           bool
	Label                 string
	Width                 float32
	MaxHeight             float32
}
type hoverCardState struct {
	Open              bool
	Suppressed        bool
	Keyboard          bool
	Since, LastActive time.Time
}

func hoverState(anchor ui.Element) *hoverCardState {
	return ui.Local(anchor, "astryx-hover-card", func() hoverCardState { return hoverCardState{} })
}

func (k *UI) HoverCard(anchor ui.Element, content func(), options ...HoverCardOptions) ui.Element {
	o := option(options)
	if o.Delay <= 0 {
		o.Delay = 300 * time.Millisecond
	}
	if o.HideDelay <= 0 {
		o.HideDelay = 200 * time.Millisecond
	}
	state := hoverState(anchor)
	if !state.Open {
		state.Keyboard = false
	}
	focus := ui.Local(anchor, "hover-return-focus", func() ui.Handle { return ui.Handle{} })
	anchor.Bind(focus)
	now := k.c.Now()
	if o.Disabled {
		state.Open = false
		return ui.Element{}
	}
	active := anchor.Hovered() || anchor.FocusVisible()
	if !active {
		state.Suppressed = false
	}
	clicked := o.OpenOnClick && anchor.Clicked()
	if clicked {
		state.Open = !state.Open
		state.Suppressed = !state.Open
		state.LastActive = now
	}
	if active && !state.Suppressed && !clicked {
		state.LastActive = now
		if state.Since.IsZero() {
			state.Since = now
		}
		if anchor.FocusVisible() || now.Sub(state.Since) >= o.Delay {
			state.Open = true
		} else {
			k.c.After(o.Delay - now.Sub(state.Since))
		}
	} else {
		state.Since = time.Time{}
	}
	var panel ui.Element
	if state.Open {
		if o.FocusOnOpen {
			state.Keyboard = true
		}
		panel = ui.PopoverBase(k.c, anchor, &state.Open, func(p ui.Element) {
			k.popup(p)
			w, h := k.c.Size()
			maxHeight := o.MaxHeight
			if maxHeight <= 0 {
				maxHeight = 360
			}
			maxHeight = min(maxHeight, max(h-32, 32))
			p.Padding(16).Gap(8).MaxWidth(min(360, max(w-32, 0))).MaxHeight(maxHeight)
			if o.Width > 0 {
				p.Width(min(o.Width, max(w-32, 0))).MaxWidth(max(w-32, 0))
			}
			if o.Label != "" {
				p.Label(o.Label)
			}
			if state.Keyboard {
				p.Focusable().AutoFocus()
			}
			if content != nil {
				ui.Scroll(k.c).MaxHeight(max(maxHeight-32, 0)).MinHeight(0).Shrink(1).Children(func() { ui.Column(k.c).Gap(8).Children(content) })
			}
			if p.Hovered() || p.FocusWithin() {
				active = true
				state.LastActive = now
			}
		})
		if panel.Dismissed() {
			state.Suppressed = true
			if panel.FocusWithin() {
				focus.Focus()
			}
		}
		if !active {
			left := o.HideDelay - now.Sub(state.LastActive)
			if left <= 0 {
				state.Open = false
			} else {
				k.c.After(left)
			}
		}
	}
	return panel
}
func (k *UI) InfoTip(label, description string) ui.Element {
	anchor := k.IconButton(label, IconInfo, ButtonOptions{Size: Small, Variant: Ghost}).Radius(999)
	k.Tooltip(anchor, description)
	return anchor
}
func (k *UI) BottomSheet(title string, open *bool, content func(), options ...DialogOptions) ui.Element {
	o := option(options)
	return ui.DialogBase(k.c, open, func(back, panel ui.Element) {
		w, h := k.c.Size()
		back.Background(k.backdropColor()).Justify(ui.End).AlignItems(ui.Center)
		width := w - 16
		if o.Width > 0 {
			width = min(width, o.Width)
		}
		panel.Width(max(width, 0)).MaxHeight(h*.85).Padding(20).Gap(14).Radius(k.tokens.RadiusLarge).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated).Shadow(0, 12, 36, 0, k.tokens.Shadow).Label(title)
		if o.Header != nil {
			o.Header()
		} else {
			ui.Row(k.c).Children(func() {
				k.Heading(title, 3).Grow(1)
				k.IconButton("关闭"+title, IconClose, ButtonOptions{Variant: Ghost, Size: Small}).OnClick(func() { *open = false })
			})
		}
		if o.Description != "" {
			panel.Description(o.Description)
		}
		ui.Scroll(k.c).MinHeight(0).Shrink(1).Children(func() {
			k.Stack().Children(func() {
				if o.Description != "" {
					k.Caption(o.Description)
				}
				if content != nil {
					content()
				}
			})
		})
		if o.Footer != nil {
			ui.Row(k.c).Justify(ui.End).Gap(8).Shrink(0).Wrap().Children(o.Footer)
		}
	})
}
