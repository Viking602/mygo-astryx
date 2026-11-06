package astryx

import "github.com/egoist/mygo/ui"

func (k *UI) AspectRatio(ratio float32) ui.Element {
	if ratio <= 0 {
		ratio = 1
	}
	return ui.Box(k.c).AspectRatio(ratio).Clip()
}
func (k *UI) Center() ui.Element { return ui.Column(k.c).Center() }
func (k *UI) HStack() ui.Element { return ui.Row(k.c).AlignItems(ui.Center).Gap(k.tokens.Space * 3) }
func (k *UI) Section(title, description string, build func()) ui.Element {
	return k.Stack().Gap(12).Children(func() {
		if title != "" || description != "" {
			ui.Column(k.c).Gap(4).Children(func() {
				if title != "" {
					k.Heading(title, 3)
				}
				if description != "" {
					k.Caption(description)
				}
			})
		}
		if build != nil {
			build()
		}
	})
}

type LayoutSlots struct {
	Header, Content, Panel, Footer func()
	PanelWidth                     float32
}

func (k *UI) Layout(slots LayoutSlots) ui.Element {
	return ui.Column(k.c).Fill().MinHeight(0).Background(k.tokens.Elevated).Children(func() {
		if slots.Header != nil {
			ui.Row(k.c).AlignItems(ui.Center).Padding(12, 16).DragWindow().Children(slots.Header)
		}
		body := ui.Row(k.c).Grow(1).MinHeight(0).AlignItems(ui.Stretch)
		w, _ := k.c.Size()
		if w < 700 {
			body.Column()
		}
		body.Children(func() {
			ui.Column(k.c).Grow(1).MinWidth(0).MinHeight(0).Padding(8, 16, 16, 16).Children(slots.Content)
			if slots.Panel != nil {
				width := slots.PanelWidth
				if width <= 0 {
					width = 240
				}
				panel := ui.Column(k.c).Padding(16).Gap(8).Background(k.tokens.Subtle)
				if w >= 700 {
					panel.Width(width).BorderWidth(0, 0, 0, 1).BorderColor(k.tokens.Native.Border)
				} else {
					panel.BorderWidth(1, 0, 0, 0).BorderColor(k.tokens.Native.Border)
				}
				panel.Children(slots.Panel)
			}
		})
		if slots.Footer != nil {
			ui.Row(k.c).Padding(10, 16).BorderWidth(1, 0, 0, 0).BorderColor(k.tokens.Native.Border).Children(slots.Footer)
		}
	})
}
func (k *UI) FormLayout(columns int) ui.Element { return k.Grid(columns).Gap(20) }

// VisuallyHidden keeps text available to assistive technology without taking
// layout room or pointer/keyboard input.
func (k *UI) VisuallyHidden(text string) ui.Element {
	return ui.Text(k.c, text).Label(text).Absolute().Top(0).Left(0).Size(1, 1).Clip().Opacity(0).IgnoreInput(true)
}

type ResizableState struct {
	Size      float32
	Collapsed bool
}
type ResizableOptions struct {
	Min, Max                       float32
	Vertical, Percent, Collapsible bool
	OnResize                       func(float32)
}

func (k *UI) Resizable(state *ResizableState, first, second func(), options ...ResizableOptions) ui.Element {
	o := option(options)
	if o.Min <= 0 {
		if o.Percent {
			o.Min = 5
		} else {
			o.Min = 40
		}
	}
	if o.Max <= 0 {
		o.Max = 10000
	}
	if o.Percent {
		o.Max = min(o.Max, 95)
		o.Min = min(o.Min, 95)
	}
	o.Max = max(o.Max, o.Min)
	state.Size = min(max(state.Size, o.Min), o.Max)
	root := ui.Column(k.c).Gap(6).MinHeight(0)
	root.Children(func() {
		if o.Collapsible {
			k.Button(map[bool]string{false: "收起面板", true: "展开面板"}[state.Collapsed], ButtonOptions{Size: Small, Variant: Ghost}).AlignSelf(ui.Start).OnClick(func() { state.Collapsed = !state.Collapsed })
		}
		if state.Collapsed {
			ui.Column(k.c).Grow(1).MinHeight(0).Children(second)
			return
		}
		var e ui.Element
		switch {
		case o.Vertical && o.Percent:
			e = ui.SplitVerticalPercent(k.c, &state.Size, first, second)
		case o.Vertical:
			e = ui.SplitVertical(k.c, &state.Size, first, second)
		case o.Percent:
			e = ui.SplitPercent(k.c, &state.Size, first, second)
		default:
			e = ui.Split(k.c, &state.Size, first, second)
		}
		e.Grow(1).MinHeight(0).OnChange(func() {
			state.Size = min(max(state.Size, o.Min), o.Max)
			if o.OnResize != nil {
				o.OnResize(state.Size)
			}
			root.SetChanged()
		})
	})
	return root
}

type OverlayOptions struct {
	NonModal, DismissOnOutside bool
	Label                      string
}

func (k *UI) Overlay(open *bool, build func(), options ...OverlayOptions) ui.Element {
	if !*open {
		return ui.Element{}
	}
	o := option(options)
	var root ui.Element
	ui.Overlay(k.c, func() {
		root = ui.Column(k.c).Absolute().Top(0).Right(0).Bottom(0).Left(0)
		if !o.NonModal {
			root.Modal().Background(k.backdropColor())
		}
		if o.Label != "" {
			root.Label(o.Label)
		}
		root.Children(build)
		if root.OverlayShortcut(0, ui.KeyEscape) || o.DismissOnOutside && root.Clicked() {
			*open = false
		}
	})
	return root
}
