package astryx

import "github.com/egoist/mygo/ui"

func (k *UI) Alert(title, description string, tones ...Tone) ui.Element {
	t := option(tones)
	ink, fill := k.colors(t)
	border := k.tokens.Native.Border
	if t != Neutral {
		border = ink.Alpha(0.24)
	}
	return ui.Column(k.c).Padding(14, 16).Gap(4).Radius(k.tokens.Radius).
		Background(fill).Border(1, border).Role(ui.RoleStatus).Label(title).Children(func() {
		k.Text(title).FontWeight(600).TextColor(ink)
		if description != "" {
			k.Text(description).TextColor(k.tokens.Native.Text)
		}
	})
}
func (k *UI) Progress(value float64) ui.Element {
	return ui.Progress(k.c, min(max(value, 0), 1)).Height(7).Radius(999)
}
func (k *UI) Spinner() ui.Element { return ui.Spinner(k.c).Label("加载中") }

// Skeleton is deliberately static: idle placeholders do not keep the window
// rendering and are suitable for users who request reduced motion.
func (k *UI) Skeleton(width, height float32) ui.Element {
	return ui.Box(k.c).Size(width, height).Radius(k.tokens.RadiusSmall).Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border.Alpha(0.5)).Role(ui.RoleNone)
}

// Toast queues a message. Call it in response to an action, not unconditionally
// in a view. Include ToastViewport once in the view for the custom appearance.
func (k *UI) Toast(message ui.Toast) string { return k.c.AddToast(message) }
func (k *UI) ToastViewport() ui.Element {
	w, _ := k.c.Size()
	return ui.ToastViewportBase(k.c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Padding(20).Gap(8).AlignItems(ui.End)
		for _, toast := range toasts {
			parts := ui.ToastBase(k.c, toast)
			parts.Root.Width(min(360, max(w-40, 100))).Padding(14, 16).Gap(12).Row().AlignItems(ui.Center).Radius(k.tokens.RadiusLarge).
				Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).Shadow(0, 6, 22, 0, k.tokens.Shadow)
			parts.Root.Children(func() {
				ui.Column(k.c).Grow(1).MinWidth(0).Gap(4).Children(func() {
					ink := k.tokens.Native.Text
					if toast.Type == "error" {
						ink = k.tokens.Native.Danger
					}
					k.Text(toast.Title).FontWeight(600).TextColor(ink)
					if toast.Description != "" {
						k.Caption(toast.Description)
					}
				})
				if toast.Action != "" {
					k.styleButton(parts.ActionButton(), ButtonOptions{Variant: Outline, Size: Small}).Label(toast.Action).Children(func() { k.Text(toast.Action).FontSize(k.tokens.FontSmall).FontWeight(550) })
				}
				k.styleButton(parts.CloseButton(), ButtonOptions{Variant: Ghost, Size: Small}).Size(28, 28).Padding(4).Label("关闭通知").Children(func() { k.Icon(IconClose).FontSize(14) })
			})
		}
	})
}
