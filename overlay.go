package astryx

import "github.com/egoist/mygo/ui"

type DialogOptions struct {
	Width       float32
	Description string
	Footer      func()
	Header      func()
}

func (k *UI) backdropColor() ui.Color {
	if k.tokens.Native.Dark {
		return k.tokens.Native.Background.Alpha(0.72)
	}
	return k.tokens.Native.Background.Alpha(0.76)
}

func (k *UI) dialog(title string, open *bool, content func(), o DialogOptions, drawer bool) ui.Element {
	return ui.DialogBase(k.c, open, func(backdrop, panel ui.Element) {
		backdrop.Background(k.backdropColor())
		w, h := k.c.Size()
		bar := k.c.TitleBar()
		width := o.Width
		if width <= 0 {
			width = 440
		}
		width = min(width, max(w-32, 0))
		panel.Width(width).MaxHeight(max(h-32, 0)).Padding(24).Gap(18).Radius(k.tokens.RadiusLarge).
			Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).Shadow(0, 12, 40, 0, k.tokens.Shadow).Label(title)
		if drawer {
			backdrop.Justify(ui.Start).AlignItems(ui.End)
			panel.Height(h).MaxHeight(h).Padding(max(24, bar.Height+12), 24, 24, 24).Radius(0)
		}
		if o.Header != nil {
			o.Header()
		} else {
			ui.Row(k.c).AlignItems(ui.Center).Gap(12).Children(func() {
				k.Heading(title, 3).Grow(1)
				k.IconButton("关闭"+title, IconClose, ButtonOptions{Variant: Ghost, Size: Small}).OnClick(func() { *open = false })
			})
		}
		if o.Description != "" {
			panel.Description(o.Description)
		}
		ui.Scroll(k.c).Shrink(1).MinHeight(0).Children(func() {
			ui.Column(k.c).Gap(12).Children(func() {
				if o.Description != "" {
					k.Caption(o.Description)
				}
				if content != nil {
					content()
				}
			})
		})
		if o.Footer != nil {
			ui.Row(k.c).Justify(ui.End).Gap(8).Shrink(0).Wrap().Children(o.Footer)
		}
	})
}
func (k *UI) Dialog(title string, open *bool, content func(), options ...DialogOptions) ui.Element {
	return k.dialog(title, open, content, option(options), false)
}
func (k *UI) Drawer(title string, open *bool, content func(), options ...DialogOptions) ui.Element {
	return k.dialog(title, open, content, option(options), true)
}

type AlertDialogOptions struct {
	Confirm, Cancel string
	Destructive     bool
	OnConfirm       func()
}

func (k *UI) AlertDialog(title, description string, open *bool, options ...AlertDialogOptions) ui.Element {
	o := option(options)
	if o.Confirm == "" {
		o.Confirm = "确认"
	}
	if o.Cancel == "" {
		o.Cancel = "取消"
	}
	if !*open {
		return ui.Element{}
	}
	var panel ui.Element
	ui.Overlay(k.c, func() {
		back := ui.Box(k.c).Absolute().Top(0).Right(0).Bottom(0).Left(0).Center().Background(k.backdropColor()).Modal()
		back.Children(func() {
			w, h := k.c.Size()
			panel = k.Card(CardOptions{Padding: 24, Elevated: true}).Width(min(420, max(w-32, 0))).MaxHeight(max(h-32, 0)).Role(ui.RoleAlertDialog).Label(title).Description(description)
			panel.Children(func() {
				k.Heading(title, 3).Shrink(0)
				ui.Scroll(k.c).Shrink(1).MinHeight(0).Children(func() { k.Text(description) })
				ui.Row(k.c).Justify(ui.End).Gap(8).Shrink(0).Wrap().Margin(8, 0, 0, 0).Children(func() {
					cancel := k.Button(o.Cancel, ButtonOptions{Variant: Outline}).OnClick(func() { *open = false })
					variant := Primary
					if o.Destructive {
						variant = Destructive
						cancel.AutoFocus()
					}
					confirm := k.Button(o.Confirm, ButtonOptions{Variant: variant}).OnClick(func() {
						*open = false
						if o.OnConfirm != nil {
							o.OnConfirm()
						}
					})
					if !o.Destructive {
						confirm.AutoFocus()
					}
				})
			})
		})
		// Outside presses deliberately do not dismiss a decision dialog.
		if back.OverlayShortcut(0, ui.KeyEscape) {
			*open = false
		}
	})
	return panel
}
func (k *UI) Popover(anchor ui.Element, open *bool, content func()) ui.Element {
	return ui.PopoverBase(k.c, anchor, open, func(panel ui.Element) {
		k.popupContent(panel, 16, 10, content)
	})
}
func (k *UI) Tooltip(anchor ui.Element, description string) ui.Element {
	anchor.Description(description)
	return ui.TooltipBase(k.c, anchor, func(tip ui.Element) {
		tip.Padding(7, 10).Radius(k.tokens.RadiusSmall).Background(k.tokens.Native.Text).TextColor(k.tokens.Native.Background)
		k.Text(description).TextColor(k.tokens.Native.Background).FontSize(k.tokens.FontSmall)
	})
}
