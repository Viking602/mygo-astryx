package astryx

import "github.com/egoist/mygo/ui"

// ButtonOptions configures a button at construction. Native element methods
// remain available on the returned element for layout and event handling.
type ButtonOptions struct {
	Variant           Variant
	Size              Size
	Disabled, Loading bool
	Icon, EndIcon     *ui.SVG
}

func (k *UI) styleButton(e ui.Element, o ButtonOptions) ui.Element {
	height, pad, font := k.metrics(o.Size)
	b := k.tokens.Native
	fill, hover, pressed, ink, border := b.Surface, b.SurfaceHover, b.SurfacePressed, b.Text, b.Border.Alpha(0.75)
	switch o.Variant {
	case Primary:
		fill, hover, pressed, ink, border = b.Accent, b.AccentHover, b.AccentPressed, b.AccentText, ui.Transparent
	case Outline:
		fill, hover, pressed, border = k.tokens.Elevated, b.Surface, b.SurfaceHover, b.Border
	case Ghost:
		fill, hover, pressed, border = ui.Transparent, b.Surface, b.SurfaceHover, ui.Transparent
	case Destructive:
		fill, hover, pressed, ink, border = k.tokens.DangerMuted, k.tokens.DangerMuted.Mix(b.Danger, 0.1), k.tokens.DangerMuted.Mix(b.Danger, 0.2), b.Danger, b.Danger.Alpha(0.28)
	}
	e.Disabled(o.Disabled || o.Loading).MinHeight(height).Padding(6, pad).Gap(8).Radius(k.tokens.Radius).
		FontSize(font).FontWeight(550).TextColor(ink).Background(fill)
	if border.A > 0 {
		e.Border(1, border)
	}
	if !e.IsDisabled() {
		if e.Pressed() {
			e.Background(pressed)
		} else if e.Hovered() {
			e.Background(hover)
		}
	}
	return e
}

func (k *UI) Button(label string, options ...ButtonOptions) ui.Element {
	o := option(options)
	e := k.styleButton(ui.ButtonBase(k.c).Label(label), o)
	e.Children(func() {
		if o.Loading {
			ui.Spinner(k.c).Size(14, 14).Label("加载中")
		} else if o.Icon != nil {
			k.Icon(o.Icon).FontSize(16)
		}
		if label != "" {
			k.Text(label).FontSize(buttonFont(k, o.Size)).FontWeight(550).SingleLine()
		}
		if o.EndIcon != nil {
			k.Icon(o.EndIcon).FontSize(16)
		}
	})
	return e
}
func buttonFont(k *UI, size Size) float32 { _, _, font := k.metrics(size); return font }

// IconButton requires a descriptive label, which is exposed to assistive tools.
func (k *UI) IconButton(label string, icon *ui.SVG, options ...ButtonOptions) ui.Element {
	o := option(options)
	e := k.styleButton(ui.ButtonBase(k.c).Label(label), o)
	height, _, _ := k.metrics(o.Size)
	iconSize := float32(18)
	if o.Size == Small {
		iconSize = 14
	} else if o.Size == Large {
		iconSize = 20
	}
	e.Size(height, height).Padding(6).Children(func() {
		if o.Loading {
			ui.Spinner(k.c).Size(iconSize, iconSize)
		} else if icon != nil {
			k.Icon(icon).FontSize(iconSize)
		} else {
			k.Icon(IconMore).FontSize(iconSize)
		}
	})
	return e
}
func (k *UI) Link(label, url string) ui.Element {
	e := ui.Link(k.c, label, url).FontSize(k.tokens.FontBody).FontWeight(500).AlignSelf(ui.Start)
	if e.Hovered() {
		e.Underline()
	}
	return e
}
func (k *UI) Toggle(on *bool, label string, options ...ButtonOptions) ui.Element {
	o := option(options)
	e := ui.ToggleBase(k.c, on).Label(label)
	height, pad, font := k.metrics(o.Size)
	if k.inGroup {
		height = max(height-6, 26)
		pad = max(pad-2, 10)
	}
	b := k.tokens.Native
	fill, hover, pressed, ink, border := k.tokens.Subtle, b.SurfaceHover, b.SurfacePressed, k.tokens.Muted, b.Border
	weight := 500
	if k.inGroup {
		fill, border = ui.Transparent, ui.Transparent
	}
	if *on {
		fill, hover, pressed, ink, border = b.Accent, b.AccentHover, b.AccentPressed, b.AccentText, ui.Transparent
		weight = 600
	}
	radius := k.tokens.Radius
	if k.inGroup {
		radius = k.tokens.RadiusSmall
	}
	e.Disabled(o.Disabled || o.Loading).MinHeight(height).Padding(4, pad).Gap(6).Radius(radius).
		FontSize(font).FontWeight(weight).TextColor(ink).Background(fill)
	if border.A > 0 {
		e.Border(1, border)
	}
	if *on && k.inGroup {
		e.Shadow(0, 1, 2, 0, k.tokens.Shadow)
	}
	if !e.IsDisabled() {
		if e.Pressed() {
			e.Background(pressed)
		} else if e.Hovered() {
			e.Background(hover)
			if !*on {
				ink = b.Text
				e.TextColor(ink)
			}
		}
	}
	e.Children(func() {
		k.Text(label).FontSize(font).FontWeight(weight).TextColor(ink).SingleLine()
	})
	return e
}
func (k *UI) ButtonGroup(build func()) ui.Element {
	prev := k.inGroup
	k.inGroup = true
	defer func() { k.inGroup = prev }()
	return ui.ToggleGroup(k.c, build).AlignSelf(ui.Start).Padding(3).Gap(2).Radius(k.tokens.Radius).Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border)
}
