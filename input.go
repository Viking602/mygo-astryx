package astryx

import (
	"github.com/egoist/mygo/ui"
	"math"
	"strconv"
	"strings"
)

type InputOptions struct {
	Size                                   Size
	Label, Placeholder, Description, Error string
	Disabled, ReadOnly                     bool
	Focus                                  *ui.Handle
}

func (k *UI) styleInput(e ui.Element, o InputOptions) ui.Element {
	height, _, font := k.metrics(o.Size)
	b := k.tokens.Native
	fill := k.tokens.Elevated
	ink := b.Text
	if o.ReadOnly {
		fill = k.tokens.Subtle
		ink = k.tokens.Muted
	}
	e.MinHeight(height).Padding(7, 12).FontSize(font).Radius(k.tokens.Radius).
		Background(fill).TextColor(ink).Border(1, b.Border).
		Disabled(o.Disabled).ReadOnly(o.ReadOnly).Placeholder(o.Placeholder)
	if o.Focus != nil {
		e.Bind(o.Focus)
	}
	if o.Label != "" {
		e.Label(o.Label)
	}
	if o.Description != "" {
		e.Description(o.Description)
	}
	if o.Error != "" {
		e.Error(o.Error).BorderColor(b.Danger)
	} else if !e.IsDisabled() {
		if e.Focused() || e.FocusWithin() {
			e.BorderColor(b.Focus)
		} else if e.Hovered() {
			e.BorderColor(b.Border.Mix(b.Text, 0.24))
		}
	}
	return e
}
func (k *UI) TextInput(value *string, options ...InputOptions) ui.Element {
	return k.styleInput(ui.TextInputBase(k.c, value), option(options))
}
func (k *UI) TextArea(value *string, options ...InputOptions) ui.Element {
	return k.styleInput(ui.TextAreaBase(k.c, value), option(options)).MinHeight(100).Padding(10, 12)
}

// Field preserves MyGo's label-to-control and error-description relationships.
func (k *UI) Field(label string, build func(), options ...InputOptions) ui.Element {
	o := option(options)
	return ui.Field(k.c, label, build).Description(o.Description).Error(o.Error).Disabled(o.Disabled)
}
func (k *UI) Form(build func()) ui.Element { return ui.Form(k.c, build).Gap(16) }
func (k *UI) Fieldset(legend string, build func()) ui.Element {
	return ui.Fieldset(k.c, legend, build).Padding(16).Gap(12).Radius(k.tokens.Radius).
		Border(1, k.tokens.Native.Border).Background(k.tokens.Subtle.Alpha(0.4))
}

// InputGroup adds leading/trailing content around an IME-aware native input.
// The root is returned for layout; the input can be customized inside Field.
func (k *UI) InputGroup(value *string, leading, trailing func(), options ...InputOptions) ui.Element {
	o := option(options)
	height, _, _ := k.metrics(o.Size)
	b := k.tokens.Native
	fill := k.tokens.Elevated
	if o.ReadOnly {
		fill = k.tokens.Subtle
	}
	root := ui.Row(k.c).MinHeight(height).AlignItems(ui.Center).Gap(8).Padding(2, 12).Radius(k.tokens.Radius).
		Border(1, b.Border).Background(fill).Disabled(o.Disabled)
	root.Children(func() {
		if leading != nil {
			leading()
		}
		in := k.TextInput(value, o).Grow(1).MinWidth(0).MinHeight(max(height-4, 24)).Padding(4, 0).
			Border(0, ui.Transparent).Background(ui.Transparent).Radius(0).FocusRing(false)
		if in.Changed() {
			root.SetChanged()
		}
		if trailing != nil {
			trailing()
		}
	})
	if o.Error != "" {
		root.BorderColor(b.Danger)
	} else if !root.IsDisabled() {
		if root.FocusWithin() {
			root.BorderColor(b.Focus)
		} else if root.Hovered() {
			root.BorderColor(b.Border.Mix(b.Text, 0.24))
		}
	}
	return root
}

// NumberInput and SearchField are themed adapters to MyGo's complete controls.
func (k *UI) NumberInput(value *float64, low, high, step float64, options ...InputOptions) ui.Element {
	o := option(options)
	if step <= 0 {
		step = 1
	}
	row := ui.Row(k.c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(6).Disabled(o.Disabled).ReadOnly(o.ReadOnly)
	if o.Label != "" {
		row.Label(o.Label)
	}
	if o.Error != "" {
		row.Error(o.Error)
	}
	decimals := 0
	if s := strconv.FormatFloat(step, 'f', -1, 64); strings.Contains(s, ".") {
		decimals = len(s) - strings.IndexByte(s, '.') - 1
	}
	format := func(v float64) string { return strconv.FormatFloat(v, 'f', decimals, 64) }
	text := ui.Local(row, "number-text", func() string { return format(*value) })
	set := func(v float64) {
		v = math.Max(low, math.Min(high, v))
		if decimals >= 0 {
			p := math.Pow(10, float64(decimals))
			v = math.Round(v*p) / p
		}
		if v != *value {
			*value = v
			row.SetChanged()
		}
	}
	row.Children(func() {
		in := k.TextInput(text, o).Width(84).FontFeatures("tnum").Range(low, high, *value).Step(step)
		if in.Changed() {
			if v, err := strconv.ParseFloat(strings.TrimSpace(*text), 64); err == nil && v >= low && v <= high {
				set(v)
			}
		}
		if !o.Disabled && !o.ReadOnly {
			if in.Shortcut(0, ui.KeyUp) {
				set(*value + step)
				*text = format(*value)
			}
			if in.Shortcut(0, ui.KeyDown) {
				set(*value - step)
				*text = format(*value)
			}
		}
		if !in.Focused() {
			*text = format(*value)
		}
		k.IconButton("减少", IconMinus, ButtonOptions{Variant: Outline, Size: o.Size, Disabled: o.Disabled || o.ReadOnly || *value <= low}).OnClick(func() {
			set(*value - step)
			*text = format(*value)
		})
		k.IconButton("增加", IconPlus, ButtonOptions{Variant: Outline, Size: o.Size, Disabled: o.Disabled || o.ReadOnly || *value >= high}).OnClick(func() {
			set(*value + step)
			*text = format(*value)
		})
	})
	return row
}
func (k *UI) SearchField(value *string, options ...InputOptions) ui.Element {
	o := option(options)
	height, _, _ := k.metrics(o.Size)
	fill := k.tokens.Elevated
	if o.ReadOnly {
		fill = k.tokens.Subtle
	}
	e := ui.SearchField(k.c, value).MinHeight(height).Padding(6, 12).Radius(k.tokens.Radius).
		Background(fill).Disabled(o.Disabled).ReadOnly(o.ReadOnly)
	if o.Placeholder != "" {
		e.Placeholder(o.Placeholder)
	}
	if o.Label != "" {
		e.Label(o.Label)
	}
	return e
}
