package astryx

import (
	"fmt"
	"github.com/egoist/mygo/ui"
)

type AvatarOptions struct {
	Size       float32
	Status     string
	StatusTone Tone
	OnClick    func()
}

// Avatar composes the status into its accessible name. The status lives above
// the circular image's clipping boundary, so it stays visible at every size.
func (k *UI) Avatar(name string, image *ui.Bitmap, options ...AvatarOptions) ui.Element {
	if len(options) == 0 {
		return ui.Avatar(k.c, name, image)
	}
	o := option(options)
	base := k.tokens.Native.FontSize * 2.25
	if o.Size <= 0 {
		o.Size = base
	}
	label := name
	if o.Status != "" {
		label += "，" + o.Status
	}
	var root ui.Element
	if o.OnClick != nil {
		root = ui.ButtonBase(k.c).OnClick(o.OnClick)
	} else {
		root = ui.Box(k.c).Role(ui.RoleImage)
	}
	return root.Size(o.Size, o.Size).Shrink(0).Center().Radius(o.Size / 2).Label(label).Children(func() {
		ui.Avatar(k.c, name, image).Size(o.Size, o.Size).Radius(o.Size / 2).FontSize(float32(max(11, int(o.Size*0.4))))
		if o.Status != "" {
			k.AvatarStatusDot(o.Status, o.StatusTone, AvatarStatusOptions{AvatarSize: o.Size}).Absolute().End(0).Bottom(0).Role(ui.RoleNone)
		}
	})
}

type AvatarStatusOptions struct {
	AvatarSize float32
	Icon       *ui.SVG
}

// AvatarStatusDot distinguishes states by shape as well as color: filled for
// success, a ring for neutral, and a minus for danger.
func (k *UI) AvatarStatusDot(label string, tone Tone, options ...AvatarStatusOptions) ui.Element {
	o := option(options)
	size, border := float32(10), float32(1)
	if o.AvatarSize > 36 {
		size, border = 20, 2
	}
	if o.AvatarSize > 72 {
		size, border = 32, 4
	}
	ink, _ := k.colors(tone)
	root := ui.Box(k.c).Size(size, size).Shrink(0).Radius(999).Border(border, k.tokens.Elevated).Background(ink).Center().Role(ui.RoleImage).Label(label)
	if label == "" {
		root.Role(ui.RoleNone)
	}
	return root.Children(func() {
		if o.Icon != nil && size >= 20 {
			k.Icon(o.Icon).Size(size*.6, size*.6).TextColor(k.tokens.Native.AccentText).Role(ui.RoleNone)
			return
		}
		if tone == Neutral {
			root.Background(k.tokens.Muted)
			ui.Box(k.c).Size(size-border*4, size-border*4).Radius(999).Background(k.tokens.Elevated)
		}
		if tone == Danger || tone == Warning {
			ui.Box(k.c).Size(size*.5, max(1, size*.09)).Radius(1).Background(k.tokens.Native.AccentText)
		}
	})
}

type AvatarGroupOptions struct {
	OverflowLabel string
	OnOverflow    func()
	Disabled      bool
}

// AvatarGroupOverflow keeps the true hidden count in the accessible label,
// even when the visible text is customized (for example, "99+").
func (k *UI) AvatarGroupOverflow(count int, options ...AvatarGroupOptions) ui.Element {
	o := option(options)
	count = max(count, 0)
	text := o.OverflowLabel
	if text == "" {
		text = fmt.Sprintf("+%d", count)
	}
	label := fmt.Sprintf("另外 %d 人", count)
	size := k.tokens.Native.FontSize*2.25 + 4
	padX := float32(6)
	if len([]rune(text)) > 3 {
		padX = 12
	}
	var root ui.Element
	if o.OnOverflow != nil {
		root = ui.ButtonBase(k.c).Disabled(o.Disabled).OnClick(o.OnOverflow)
	} else {
		root = ui.Box(k.c).Role(ui.RoleImage)
	}
	root.AlignSelf(ui.Start).MinWidth(size).Height(size).Padding(0, padX).Center().
		Radius(999).Background(k.tokens.Subtle).Border(2, k.tokens.Elevated).Label(label)
	if o.OnOverflow != nil && !root.IsDisabled() {
		if root.Pressed() {
			root.Background(k.tokens.Native.SurfacePressed)
		} else if root.Hovered() {
			root.Background(k.tokens.Native.SurfaceHover)
		}
	}
	root.Children(func() {
		k.Text(text).FontSize(k.tokens.FontSmall).FontWeight(600).TextColor(k.tokens.Muted).SingleLine()
	})
	return root
}
