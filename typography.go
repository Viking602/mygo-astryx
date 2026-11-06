package astryx

import "github.com/egoist/mygo/ui"

func (k *UI) Text(text string) ui.Element {
	return ui.Text(k.c, text).FontSize(k.tokens.FontBody).LineHeight(1.45)
}
func (k *UI) Caption(text string) ui.Element {
	return ui.Text(k.c, text).FontSize(k.tokens.FontSmall).LineHeight(1.4).TextColor(k.tokens.Muted)
}
func (k *UI) Heading(text string, level int) ui.Element {
	level = min(max(level, 1), 6)
	sizes := []float32{2.2, 1.65, 1.32, 1.15, 1.05, 1}
	weights := []int{700, 650, 650, 600, 600, 600}
	return ui.Text(k.c, text).FontSize(k.tokens.FontBody * sizes[level-1]).FontWeight(weights[level-1]).LineHeight(1.25).Role(ui.RoleHeading).Level(level)
}

// Code displays selectable monospace text, without interpreting it as markup.
func (k *UI) Code(text string) ui.Element {
	return ui.Text(k.c, text).Font("monospace").FontSize(k.tokens.FontSmall).LineHeight(1.5).Selectable().
		Padding(12, 16).Radius(k.tokens.Radius).Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border)
}
func (k *UI) Icon(icon *ui.SVG) ui.Element { return ui.Icon(k.c, icon).FontSize(18) }
func (k *UI) Badge(text string, tones ...Tone) ui.Element {
	t := option(tones)
	ink, fill := k.colors(t)
	border := k.tokens.Native.Border
	if t != Neutral {
		border = ink.Alpha(0.22)
	}
	return ui.Text(k.c, text).SingleLine().FontSize(k.tokens.FontSmall).FontWeight(550).
		TextColor(ink).Padding(3, 10).Radius(999).Background(fill).Border(1, border)
}
func (k *UI) EmptyState(title, description string, action func()) ui.Element {
	return ui.Column(k.c).Padding(28, 24).Gap(8).Radius(k.tokens.RadiusLarge).
		Background(k.tokens.Subtle.Alpha(0.35)).Border(1, k.tokens.Native.Border).BorderStyle(ui.BorderDashed).
		AlignItems(ui.Center).Children(func() {
		ui.Box(k.c).Size(40, 40).Radius(999).Center().Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).Margin(0, 0, 2, 0).Children(func() {
			k.Icon(IconPlus).FontSize(18).TextColor(k.tokens.Muted)
		})
		k.Heading(title, 3)
		k.Caption(description).MaxWidth(380).TextAlign(ui.Center)
		if action != nil {
			ui.Row(k.c).Margin(6, 0, 0, 0).Gap(8).Children(action)
		}
	})
}
