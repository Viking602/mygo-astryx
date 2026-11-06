package astryx

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"strings"
)

type Person struct {
	Name, Role, Status string
	StatusTone         Tone
	Image              *ui.Bitmap
}

func (k *UI) AvatarGroup(people []Person, maximum int, options ...AvatarGroupOptions) ui.Element {
	maximum = max(maximum, 1)
	return ui.Row(k.c).AlignSelf(ui.Start).Gap(-8).AlignItems(ui.Center).Children(func() {
		for _, p := range people[:min(len(people), maximum)] {
			k.Avatar(p.Name, p.Image).Border(2, k.tokens.Elevated)
		}
		if len(people) > maximum {
			k.AvatarGroupOverflow(len(people)-maximum, option(options))
		}
	})
}
func (k *UI) Blockquote(text, attribution string) ui.Element {
	return ui.Column(k.c).Padding(10, 16).Gap(6).Radius(0, k.tokens.RadiusSmall, k.tokens.RadiusSmall, 0).
		Background(k.tokens.Subtle.Alpha(0.45)).BorderWidth(0, 0, 0, 3).BorderColor(k.tokens.Native.Accent).Children(func() {
		k.Text(text).Selectable()
		if attribution != "" {
			k.Caption("— " + attribution)
		}
	})
}
func (k *UI) Kbd(keys string) ui.Element {
	return ui.Text(k.c, keys).Font("monospace").FontSize(k.tokens.FontSmall).FontWeight(550).SingleLine().
		Padding(2, 7).Radius(5).Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border).Shadow(0, 1, 0, 0, k.tokens.Native.Border)
}
func (k *UI) StatusDot(label string, tone Tone) ui.Element {
	ink, _ := k.colors(tone)
	if tone == Neutral {
		ink = k.tokens.Muted
	}
	return ui.Row(k.c).AlignItems(ui.Center).Gap(7).Role(ui.RoleStatus).Label(label).Children(func() {
		ui.Box(k.c).Size(8, 8).Radius(999).Background(ink)
		if label != "" {
			k.Text(label)
		}
	})
}

type ItemOptions struct {
	Description, Trailing, Avatar, Badge string
	BadgeTone                            Tone
	Icon                                 *ui.SVG
	Disabled, Selected                   bool
	OnClick                              func()
}

func (k *UI) Item(label string, options ...ItemOptions) ui.Element {
	o := option(options)
	var root ui.Element
	if o.OnClick != nil {
		root = ui.ButtonBase(k.c).Justify(ui.Start).Label(label).OnClick(o.OnClick)
	} else {
		root = ui.Row(k.c)
	}
	root.Padding(10, 12).Gap(12).Disabled(o.Disabled).Radius(k.tokens.Radius).AlignItems(ui.Center)
	if o.Selected {
		root.Selected(true).Background(k.tokens.Subtle).TextColor(k.tokens.Native.Text)
	} else if root.Hovered() && o.OnClick != nil && !o.Disabled {
		root.Background(k.tokens.Subtle)
	}
	root.Children(func() {
		if o.Avatar != "" {
			k.Avatar(o.Avatar, nil, AvatarOptions{Size: 36})
		} else if o.Icon != nil {
			ui.Box(k.c).Size(32, 32).Radius(k.tokens.RadiusSmall).Center().Background(k.tokens.Subtle).Children(func() {
				k.Icon(o.Icon).FontSize(16)
			})
		}
		ui.Column(k.c).Gap(2).Grow(1).MinWidth(0).Children(func() {
			k.Text(label).FontWeight(550)
			if o.Description != "" {
				k.Caption(o.Description)
			}
		})
		if o.Badge != "" {
			k.Badge(o.Badge, o.BadgeTone)
		}
		if o.Trailing != "" {
			k.Caption(o.Trailing)
		}
	})
	return root
}

type Metadata struct{ Label, Value string }

func (k *UI) MetadataList(items []Metadata) ui.Element {
	return ui.Column(k.c).Radius(k.tokens.Radius).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated).Clip().Children(func() {
		for i, item := range items {
			row := ui.Row(k.c).Padding(10, 14).Gap(16).Wrap()
			if i > 0 {
				row.BorderWidth(1, 0, 0, 0).BorderColor(k.tokens.Native.Border)
			}
			row.Children(func() {
				k.Caption(item.Label).Width(110).FontWeight(550)
				k.Text(item.Value).Grow(1).MinWidth(0).Selectable()
			})
		}
	})
}
func (k *UI) ClickableCard(label string, build func(), options ...CardOptions) ui.Element {
	o := option(options)
	pad := o.Padding
	if pad == 0 {
		pad = 20
	}
	root := ui.ButtonBase(k.c).Column().Justify(ui.Start).AlignItems(ui.Stretch).Label(label).Padding(pad).Gap(10).Radius(k.tokens.RadiusLarge).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
	if root.Hovered() {
		root.BorderColor(k.tokens.Native.Accent).Background(k.tokens.Subtle.Alpha(0.35))
	}
	if o.Elevated {
		root.Shadow(0, 3, 12, 0, k.tokens.Shadow)
	}
	return root.Children(build)
}
func (k *UI) SelectableCard(selected *bool, label string, build func()) ui.Element {
	root := ui.CheckboxBase(k.c, selected).Column().Justify(ui.Start).AlignItems(ui.Stretch).Label(label).Padding(16, 18).Gap(8).Radius(k.tokens.RadiusLarge).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
	state := IndicatorUnchecked
	if *selected {
		state = IndicatorChecked
		root.Border(2, k.tokens.Native.Accent).Background(k.tokens.Subtle.Alpha(0.55))
	} else if root.Hovered() {
		root.BorderColor(k.tokens.Native.Border.Mix(k.tokens.Native.Text, 0.25))
	}
	return root.Children(func() {
		ui.Row(k.c).AlignItems(ui.Center).Gap(10).Children(func() {
			k.Text(label).FontWeight(600).Grow(1)
			k.CheckboxIndicator(state)
		})
		if build != nil {
			build()
		}
	})
}

type CitationOptions struct {
	Detail, URL string
	OnClick     func()
}

func (k *UI) Citation(label string, options ...CitationOptions) ui.Element {
	o := option(options)
	e := k.Button(label, ButtonOptions{Variant: Outline, Size: Small}).AlignSelf(ui.Start).Radius(999)
	if o.OnClick != nil {
		e.OnClick(o.OnClick)
	} else if o.URL != "" {
		e.OnClick(func() { k.c.OpenURL(o.URL) })
	}
	if o.Detail != "" {
		k.Tooltip(e, o.Detail)
	}
	return e
}

type MarkdownOptions struct {
	Native    *markdown.Options
	Streaming bool
}

func (k *UI) markdownOptions(custom *markdown.Options) *markdown.Options {
	if custom != nil {
		return custom
	}
	return &markdown.Options{
		WrapCode: true,
		Theme: &markdown.Theme{
			CodeBlockBackground: k.tokens.Subtle,
			CodeBlockBorder:     k.tokens.Native.Border,
			CodeBlockRadius:     k.tokens.Radius,
			CodeBackground:      k.tokens.Subtle,
			QuoteBar:            k.tokens.Native.Accent,
			Rule:                k.tokens.Native.Border,
			TableBorder:         k.tokens.Native.Border,
			TableHeader:         k.tokens.Subtle,
		},
	}
}

func (k *UI) Markdown(source string, options ...MarkdownOptions) ui.Element {
	o := option(options)
	opts := k.markdownOptions(o.Native)
	if o.Streaming {
		return markdown.Stream(k.c, source, opts)
	}
	return markdown.View(k.c, source, opts)
}

// CodeBlock uses MyGo's Markdown renderer for native syntax highlighting,
// selection and the copy action. A longer fence prevents source injection.
func (k *UI) CodeBlock(source, language string) ui.Element {
	fence := "```"
	for strings.Contains(source, fence) {
		fence += "`"
	}
	language = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '`' {
			return -1
		}
		return r
	}, language)
	return markdown.View(k.c, fence+language+"\n"+source+"\n"+fence, k.markdownOptions(nil))
}

type BannerOptions struct {
	Tone                Tone
	Action, ActionLabel string
	OnAction            func()
	Dismissed           *bool
}

func (k *UI) Banner(title, description string, options ...BannerOptions) ui.Element {
	o := option(options)
	if o.Dismissed != nil && *o.Dismissed {
		return ui.Element{}
	}
	ink, fill := k.colors(o.Tone)
	border := k.tokens.Native.Border
	if o.Tone != Neutral {
		border = ink.Alpha(0.24)
	}
	return ui.Row(k.c).AlignItems(ui.Center).Padding(14, 16).Gap(12).Radius(k.tokens.Radius).
		Background(fill).Border(1, border).Role(ui.RoleStatus).Label(title).Children(func() {
		ui.Column(k.c).Gap(3).Grow(1).MinWidth(0).Children(func() {
			k.Text(title).FontWeight(600).TextColor(ink)
			if description != "" {
				k.Text(description)
			}
		})
		label := o.ActionLabel
		if label == "" {
			label = o.Action
		}
		if label != "" {
			k.Button(label, ButtonOptions{Size: Small, Variant: Outline}).OnClick(o.OnAction)
		}
		if o.Dismissed != nil {
			k.IconButton("关闭通知栏", IconClose, ButtonOptions{Size: Small, Variant: Ghost}).OnClick(func() { *o.Dismissed = true })
		}
	})
}
