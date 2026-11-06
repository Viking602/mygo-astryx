package astryx

import "github.com/egoist/mygo/ui"

// Stack arranges children vertically with the design system's default gap.
func (k *UI) Stack() ui.Element { return ui.Column(k.c).Gap(k.tokens.Space * 3) }

// Inline is a wrapping row, suitable for actions and compact controls.
func (k *UI) Inline() ui.Element { return ui.Row(k.c).AlignItems(ui.Center).Gap(k.tokens.Space * 2).Wrap() }
func (k *UI) Grid(columns int) ui.Element {
	return ui.Grid(k.c).Columns(max(columns, 1)).Gap(k.tokens.Space * 4)
}
func (k *UI) Container() ui.Element {
	return ui.Column(k.c).FillWidth().MaxWidth(1120).Padding(24).Gap(20)
}

type ScrollOptions struct{ Horizontal, Both bool }

func (k *UI) Scroll(options ...ScrollOptions) ui.Element {
	o := option(options)
	if o.Both {
		return ui.ScrollBoth(k.c)
	}
	if o.Horizontal {
		return ui.ScrollHorizontal(k.c)
	}
	return ui.Scroll(k.c)
}
func (k *UI) Spacer() ui.Element  { return ui.Spacer(k.c) }
func (k *UI) Divider() ui.Element { return ui.Divider(k.c) }

type CardOptions struct {
	Padding  float32
	Elevated bool
}

func (k *UI) Card(options ...CardOptions) ui.Element {
	o := option(options)
	pad := o.Padding
	if pad == 0 {
		pad = k.tokens.Space * 5
	}
	e := ui.Column(k.c).Padding(pad).Gap(k.tokens.Space*3).Radius(k.tokens.RadiusLarge).
		Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border)
	if o.Elevated {
		e.Shadow(0, 6, 20, 0, k.tokens.Shadow)
	}
	return e
}

// AppShell creates an adaptive navigation/content shell with an immersive
// window title bar region. At narrow widths the navigation becomes a compact
// top section; neither callback is omitted.
func (k *UI) AppShell(navigation, content func()) ui.Element {
	w, _ := k.c.Size()
	bar := k.c.TitleBar()
	root := ui.Row(k.c).Fill().AlignItems(ui.Stretch)
	if w < 760 {
		root.Column()
	}
	return root.Children(func() {
		topPad := float32(20)
		if bar.Height > 0 && (bar.Left > 0 || w < 760) {
			topPad = max(20, bar.Height+10)
		}
		if w >= 760 {
			nav := ui.Column(k.c).Padding(topPad, 20, 20, 20).Gap(12).Width(220).FillHeight().
				Background(k.tokens.Subtle).BorderWidth(0, 1, 0, 0).BorderColor(k.tokens.Native.Border).DragWindow()
			nav.Children(navigation)
		} else {
			leftPad := max(float32(20), bar.Left+12)
			rightPad := max(float32(20), bar.Right+12)
			nav := ui.Column(k.c).Padding(topPad, rightPad, 8, leftPad).Gap(10).FillWidth().DragWindow()
			nav.Children(navigation)
		}
		ui.Column(k.c).Grow(1).MinWidth(0).MinHeight(0).Children(content)
	})
}

// Split adapts the native resizable split view to this component set.
func (k *UI) Split(position *float32, first, second func()) ui.Element {
	return ui.Split(k.c, position, first, second)
}
