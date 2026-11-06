package astryx

import (
	"fmt"
	"github.com/egoist/mygo/ui"
)

type ImageAsset struct {
	ID, Label string
	Source    ui.ImageSource
}
type ThumbnailOptions struct {
	Width, Height float32
	Selected      bool
	OnClick       func()
}

func (k *UI) Thumbnail(image ImageAsset, options ...ThumbnailOptions) ui.Element {
	o := option(options)
	if o.Width <= 0 {
		o.Width = 120
	}
	if o.Height <= 0 {
		o.Height = 80
	}
	var root ui.Element
	if o.OnClick != nil {
		root = ui.ButtonBase(k.c).Label(image.Label).OnClick(o.OnClick)
	} else {
		root = ui.Box(k.c).Role(ui.RoleImage).Label(image.Label)
	}
	root.Size(o.Width, o.Height).Radius(k.tokens.Radius).Clip().Border(1, k.tokens.Native.Border).Background(k.tokens.Subtle)
	if o.Selected {
		root.Border(2, k.tokens.Native.Accent)
	}
	return root.Children(func() {
		if image.Source != nil {
			fit := ui.Cover
			if _, isSVG := image.Source.(*ui.SVG); isSVG {
				fit = ui.FillBox
			}
			ui.Image(k.c, image.Source).Fill().Fit(fit).Label(image.Label)
		} else {
			k.Caption(image.Label).Padding(8)
		}
	})
}

type CarouselSlide struct {
	ID, Label string
	Content   func()
}
type CarouselOptions struct {
	Height float32
	Loop   bool
}

func (k *UI) Carousel(index *int, slides []CarouselSlide, options ...CarouselOptions) ui.Element {
	o := option(options)
	if o.Height <= 0 {
		o.Height = 180
	}
	clampChoice(index, len(slides))
	root := ui.Column(k.c).Gap(10).Label("轮播").Focusable()
	move := func(delta int) {
		if len(slides) == 0 {
			return
		}
		next := *index + delta
		if o.Loop {
			next = (next + len(slides)) % len(slides)
		}
		*index = min(max(next, 0), len(slides)-1)
		root.SetChanged()
	}
	root.Children(func() {
		if len(slides) == 0 {
			k.EmptyState("暂无内容", "添加轮播内容后显示。", nil)
			return
		}
		slide := slides[*index]
		ui.Box(k.c.Key(slide.ID)).Height(o.Height).Radius(k.tokens.RadiusLarge).Clip().Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border).Children(slide.Content)
		k.Inline().Children(func() {
			k.Button("上一项", ButtonOptions{Variant: Outline, Size: Small, Disabled: !o.Loop && *index == 0}).OnClick(func() { move(-1) })
			k.Caption(fmt.Sprintf("%d / %d · %s", *index+1, len(slides), slide.Label)).Grow(1).TextAlign(ui.Center)
			k.Button("下一项", ButtonOptions{Variant: Outline, Size: Small, Disabled: !o.Loop && *index == len(slides)-1}).OnClick(func() { move(1) })
		})
	})
	root.OnShortcut(0, ui.KeyLeft, func() { move(-1) }).OnShortcut(0, ui.KeyRight, func() { move(1) })
	return root
}
func (k *UI) Lightbox(open *bool, index *int, images []ImageAsset) ui.Element {
	if !*open {
		return ui.Element{}
	}
	clampChoice(index, len(images))
	return ui.DialogBase(k.c, open, func(back, panel ui.Element) {
		back.Background(k.backdropColor())
		w, h := k.c.Size()
		bar := k.c.TitleBar()
		topSafe := max(float32(32), bar.Height+16)
		pw := min(float32(880), max(w-64, 280))
		ph := min(float32(620), max(h-topSafe-32, 240))
		panel.Width(pw).Height(ph).Padding(20, 24).Gap(14).
			Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).
			Shadow(0, 12, 36, 0, k.tokens.Shadow).Radius(k.tokens.RadiusLarge).Label("图片预览")
		zoom := ui.Local(panel, "lightbox-zoom", func() float32 { return 1 })
		rotate := ui.Local(panel, "lightbox-rotate", func() float32 { return 0 })
		move := func(delta int) {
			if len(images) == 0 {
				return
			}
			*index = (*index + delta + len(images)) % len(images)
			*zoom = 1
			*rotate = 0
		}
		k.Inline().Children(func() {
			k.Heading("图片预览", 3).Grow(1)
			k.Button("缩小", ButtonOptions{Variant: Outline, Size: Small}).OnClick(func() { *zoom = max(*zoom-.25, .25) })
			k.Button("放大", ButtonOptions{Variant: Outline, Size: Small}).OnClick(func() { *zoom = min(*zoom+.25, 4) })
			k.Button("旋转", ButtonOptions{Variant: Outline, Size: Small}).OnClick(func() { *rotate += 90 })
			k.IconButton("关闭图片预览", IconClose, ButtonOptions{Variant: Ghost, Size: Small}).OnClick(func() { *open = false })
		})
		if len(images) == 0 {
			k.EmptyState("没有图片", "没有可预览的内容。", nil)
			return
		}
		image := images[*index]
		ui.ScrollBoth(k.c).Grow(1).MinHeight(0).Radius(k.tokens.Radius).Background(k.tokens.Subtle.Alpha(0.55)).Border(1, k.tokens.Native.Border).Children(func() {
			ui.Box(k.c).Size(max(pw-56, 100)*(*zoom), max(ph-148, 100)*(*zoom)). Center().Children(func() {
				if image.Source != nil {
					ui.Image(k.c, image.Source).Fill().Fit(ui.Contain).Rotate(*rotate).Label(image.Label)
				} else {
					k.Text(image.Label)
				}
			})
		})
		k.Inline().Children(func() {
			k.Button("上一张", ButtonOptions{Variant: Outline, Size: Small}).OnClick(func() { move(-1) })
			k.Caption(fmt.Sprintf("%d / %d · %s · %.0f%%", *index+1, len(images), image.Label, *zoom*100)).Grow(1).TextAlign(ui.Center)
			k.Button("下一张", ButtonOptions{Variant: Outline, Size: Small}).OnClick(func() { move(1) })
		})
		if panel.OverlayShortcut(0, ui.KeyLeft) {
			move(-1)
		}
		if panel.OverlayShortcut(0, ui.KeyRight) {
			move(1)
		}
	})
}
