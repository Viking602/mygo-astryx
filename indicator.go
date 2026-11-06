package astryx

import "github.com/egoist/mygo/ui"

type IndicatorState string

const (
	IndicatorUnchecked IndicatorState = "unchecked"
	IndicatorChecked   IndicatorState = "checked"
	IndicatorMixed     IndicatorState = "mixed"
)

type IndicatorOptions struct {
	Size     Size
	Disabled bool
	Kind     string
}

// Indicator is a decorative control visual, not a notification badge. Its
// parent remains responsible for role, state, keyboard and focus behavior.
func (k *UI) Indicator(state IndicatorState, options ...IndicatorOptions) ui.Element {
	o := option(options)
	kind := o.Kind
	if kind == "" {
		kind = "checkbox"
	}
	return k.indicator(kind, state, o)
}
func (k *UI) CheckboxIndicator(state IndicatorState, options ...IndicatorOptions) ui.Element {
	return k.indicator("checkbox", state, option(options))
}
func (k *UI) RadioIndicator(state IndicatorState, options ...IndicatorOptions) ui.Element {
	return k.indicator("radio", state, option(options))
}
func (k *UI) CheckIndicator(state IndicatorState, options ...IndicatorOptions) ui.Element {
	return k.indicator("check", state, option(options))
}
func (k *UI) indicator(kind string, state IndicatorState, o IndicatorOptions) ui.Element {
	size, _ := k.choiceMetrics(o.Size)
	if k.tokens.Indicator != nil {
		if e, ok := k.tokens.Indicator(k.c, kind, string(state), size); ok {
			return e.Role(ui.RoleNone).IgnoreInput(true)
		}
	}
	root := ui.Box(k.c).Size(size, size).Center().Role(ui.RoleNone).IgnoreInput(true)
	if o.Disabled {
		root.Opacity(.5)
	}
	iconSize := float32(int(size*0.78) &^ 1)
	if iconSize < 12 {
		iconSize = 12
	}
	switch kind {
	case "check":
		if state == IndicatorChecked {
			root.Children(func() { k.Icon(IconCheck).FontSize(iconSize).TextColor(k.tokens.Native.Accent).Role(ui.RoleNone) })
		}
	case "radio":
		dot := float32(int(size*0.46) &^ 1)
		root.Radius(999).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
		if state == IndicatorChecked {
			root.BorderColor(k.tokens.Native.Accent).Children(func() { ui.Box(k.c).Size(dot, dot).Radius(999).Background(k.tokens.Native.Accent) })
		}
	default:
		root.Radius(5).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
		if state != IndicatorUnchecked {
			root.Background(k.tokens.Native.Accent).BorderColor(k.tokens.Native.Accent).Children(func() {
				icon := IconCheck
				if state == IndicatorMixed {
					icon = IconMinus
				}
				k.Icon(icon).FontSize(iconSize).TextColor(k.tokens.Native.AccentText).Role(ui.RoleNone)
			})
		}
	}
	return root
}
func (k *UI) FieldStatus(message string, tone Tone) ui.Element {
	ink, _ := k.colors(tone)
	return k.Caption(message).TextColor(ink).Role(ui.RoleStatus)
}
func (k *UI) InputClearButton(value *string, focus *ui.Handle, options ...ButtonOptions) ui.Element {
	o := option(options)
	o.Variant = Ghost
	o.Size = Small
	o.Disabled = o.Disabled || *value == ""
	return k.IconButton("清除输入", IconClose, o).OnClick(func() {
		*value = ""
		if focus != nil {
			focus.Focus()
		}
	})
}
func (k *UI) GridSpan(columns int) ui.Element { return ui.Box(k.c).ColumnSpan(max(columns, 1)) }
