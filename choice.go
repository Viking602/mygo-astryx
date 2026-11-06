package astryx

import "github.com/egoist/mygo/ui"

type ChoiceOptions struct {
	Disabled, Mixed bool
	Size            Size
}

func (k *UI) choiceMetrics(size Size) (float32, float32) {
	_, _, font := k.metrics(size)
	switch size {
	case Small:
		return 16, font
	case Large:
		return 22, font
	default:
		return 18, font
	}
}

func (k *UI) Checkbox(value *bool, label string, options ...ChoiceOptions) ui.Element {
	o := option(options)
	_, font := k.choiceMetrics(o.Size)
	e := ui.CheckboxBase(k.c, value).Disabled(o.Disabled).Label(label).Gap(9).MinHeight(28)
	if o.Mixed {
		e.Mixed()
	}
	e.Children(func() {
		state := IndicatorUnchecked
		if *value {
			state = IndicatorChecked
		}
		if o.Mixed {
			state = IndicatorMixed
		}
		k.CheckboxIndicator(state, IndicatorOptions{Size: o.Size})
		if label != "" {
			k.Text(label).FontSize(font)
		}
	})
	return e
}
func (k *UI) Switch(value *bool, label string, options ...ChoiceOptions) ui.Element {
	o := option(options)
	size, font := k.choiceMetrics(o.Size)
	height := size + 4
	width := height + 16
	knob := height - 6
	e := ui.SwitchBase(k.c, value).Disabled(o.Disabled).Label(label).Gap(10).MinHeight(28)
	e.Children(func() {
		fill := k.tokens.Native.SurfacePressed
		border := k.tokens.Native.Border
		if *value {
			fill = k.tokens.Native.Accent
			border = k.tokens.Native.Accent
		}
		ui.Box(k.c).Size(width, height).Radius(999).Background(fill).Border(1, border).Children(func() {
			x := float32(2)
			if *value {
				x = width - knob - 4
			}
			ink := ui.Hex("#ffffff")
			if *value {
				ink = k.tokens.Native.AccentText
			}
			ui.Box(k.c).Absolute().Top(2).Left(x).Size(knob, knob).Radius(999).Background(ink).Shadow(0, 1, 2, 0, k.tokens.Shadow)
		})
		if label != "" {
			k.Text(label).FontSize(font)
		}
	})
	return e
}
func (k *UI) RadioGroup(label string, build func()) ui.Element {
	return ui.RadioGroup(k.c, build).Gap(10).Label(label)
}
func (k *UI) Radio(value *string, choice, label string, options ...ChoiceOptions) ui.Element {
	o := option(options)
	_, font := k.choiceMetrics(o.Size)
	e := ui.RadioBase(k.c, value, choice).Disabled(o.Disabled).Label(label).Gap(9).MinHeight(28)
	e.Children(func() {
		state := IndicatorUnchecked
		if *value == choice {
			state = IndicatorChecked
		}
		k.RadioIndicator(state, IndicatorOptions{Size: o.Size})
		k.Text(label).FontSize(font)
	})
	return e
}
func (k *UI) Slider(value *float64, low, high float64, options ...InputOptions) ui.Element {
	o := option(options)
	e := ui.Slider(k.c, value, low, high).Disabled(o.Disabled)
	if o.Label != "" {
		e.Label(o.Label)
	}
	return e
}
func (k *UI) RangeSlider(lowValue, highValue *float64, low, high, step float64) ui.Element {
	return ui.RangeSlider(k.c, lowValue, highValue, low, high, step)
}
func (k *UI) Rating(value *int, maximum int) ui.Element { return ui.Rating(k.c, value, maximum) }
