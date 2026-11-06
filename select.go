package astryx

import (
	"github.com/egoist/mygo/ui"
	"strings"
)

// Option separates a choice's stable value from its display label.
type Option struct {
	Value, Label string
	Disabled     bool
}

func (k *UI) popup(panel ui.Element) {
	panel.Padding(5).Gap(2).MaxHeight(280).Radius(k.tokens.RadiusLarge).
		Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).Shadow(0, 6, 22, 0, k.tokens.Shadow).Margin(5, 0, 0, 0)
}

func (k *UI) popupContent(panel ui.Element, padding, gap float32, content func()) {
	w, h := k.c.Size()
	maxHeight := min(float32(360), max(h-32, 32))
	k.popup(panel)
	panel.Padding(padding).MaxWidth(max(w-32, 0)).MaxHeight(maxHeight)
	if content != nil {
		ui.Scroll(k.c).MinHeight(0).Shrink(1).MaxHeight(max(maxHeight-padding*2, 0)).Children(func() { ui.Column(k.c).Gap(gap).Children(content) })
	}
}

func (k *UI) Select(value *string, items []Option, options ...InputOptions) ui.Element {
	o := option(options)
	// SelectParts has no public close method. A scoped generation resets its
	// private popup state when availability changes, without touching internals.
	root := ui.Column(k.c).MinWidth(0)
	state := ui.Local(root, "astryx-select-availability", func() selectAvailability { return selectAvailability{Disabled: o.Disabled} })
	if state.Disabled != o.Disabled {
		state.Disabled = o.Disabled
		state.Generation++
	}
	var trigger ui.Element
	root.Children(func() { k.c.Key(state.Generation); trigger = k.selectTrigger(value, items, o) })
	return trigger
}

type selectAvailability struct {
	Disabled   bool
	Generation int
}

func (k *UI) styleTrigger(trigger ui.Element, text string, placeholder bool, o InputOptions) ui.Element {
	height, _, font := k.metrics(o.Size)
	b := k.tokens.Native
	fill := k.tokens.Elevated
	if o.ReadOnly {
		fill = k.tokens.Subtle
	}
	trigger.MinHeight(height).Padding(7, 12).Gap(12).Radius(k.tokens.Radius).
		Background(fill).Border(1, b.Border).Disabled(o.Disabled).Justify(ui.SpaceBetween)
	if o.Label != "" {
		trigger.Label(o.Label)
	}
	if o.Description != "" {
		trigger.Description(o.Description)
	}
	if o.Error != "" {
		trigger.Error(o.Error).BorderColor(b.Danger)
	} else if !trigger.IsDisabled() {
		if trigger.Focused() || trigger.FocusWithin() {
			trigger.BorderColor(b.Focus)
		} else if trigger.Hovered() {
			trigger.BorderColor(b.Border.Mix(b.Text, 0.24))
		}
	}
	trigger.Children(func() {
		t := k.Text(text).FontSize(font).Grow(1).SingleLine().TextAlign(ui.Start)
		if placeholder {
			t.TextColor(k.tokens.Muted)
		}
		k.Icon(IconChevronDown).FontSize(14).TextColor(k.tokens.Muted)
	})
	return trigger
}

func (k *UI) selectTrigger(value *string, items []Option, o InputOptions) ui.Element {
	s := ui.SelectBase(k.c, value)
	label := o.Placeholder
	if label == "" {
		label = "请选择"
	}
	found := false
	for _, item := range items {
		if item.Value == *value {
			label = item.Label
			found = true
			break
		}
	}
	k.styleTrigger(s.Trigger, label, !found, o)
	if o.Disabled {
		s.Trigger.Expanded(false)
		return s.Trigger
	}
	s.Popup(func(panel ui.Element) {
		k.popup(panel)
		ui.Scroll(k.c).MaxHeight(260).MinHeight(0).Children(func() {
			ui.Column(k.c).Gap(2).Children(func() {
				if len(items) == 0 {
					k.Caption("暂无选项").Padding(8)
				}
				for _, item := range items {
					// SelectBase registers keyboard choices when Item is called. Render
					// unavailable choices separately so arrows/Enter cannot select them.
					if item.Disabled {
						ui.Row(k.c).Padding(8, 10).Disabled(true).Label(item.Label + "，不可用").Children(func() { k.Caption(item.Label) })
						continue
					}
					row := s.Item(item.Value).Disabled(item.Disabled).Padding(8, 10).Gap(12).Radius(k.tokens.RadiusSmall).Label(item.Label)
					if row.Highlighted() {
						row.Background(k.tokens.Subtle)
					}
					row.Children(func() {
						k.Text(item.Label).Grow(1)
						if item.Value == *value {
							k.CheckIndicator(IndicatorChecked)
						}
					})
				}
			})
		})
	})
	return s.Trigger
}

// Combobox returns the editable input. It keeps MyGo's native text editing and
// keyboard selection; the model contains the selected or newly typed string.
func (k *UI) Combobox(value *string, items []string, options ...InputOptions) ui.Element {
	o := option(options)
	s := ui.ComboboxBase(k.c, value)
	k.styleInput(s.Input, o)
	if o.Disabled || o.ReadOnly {
		s.SetOpen(false)
	}
	s.Popup(func(panel ui.Element) {
		k.popup(panel)
		ui.Scroll(k.c).MaxHeight(260).MinHeight(0).Children(func() {
			ui.Column(k.c).Gap(2).Children(func() {
				count := 0
				for _, item := range items {
					if s.Filtering() && !strings.Contains(strings.ToLower(item), strings.ToLower(*value)) {
						continue
					}
					row := s.Item(item).Padding(8, 10).Radius(k.tokens.RadiusSmall).Label(item)
					if row.Highlighted() {
						row.Background(k.tokens.Subtle)
					}
					row.Children(func() { k.Text(item) })
					count++
				}
				if count == 0 {
					k.Caption("没有匹配的选项").Padding(8)
				}
			})
		})
	})
	if chosen, ok := s.Chosen(); ok {
		*value = chosen
		s.Input.SetChanged()
	}
	return s.Input
}
