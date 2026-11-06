package astryx

import (
	"github.com/egoist/mygo/ui"
	"time"
)

// Table preserves MyGo's virtualized rows, selection, sorting and resizing.
// Configure ListState on the app model before calling Table.
func (k *UI) Table(state *ui.ListState, columns []ui.TableColumn, count int, cell func(row, col int)) ui.Element {
	return ui.Table(k.c, state, columns, count, cell).Radius(k.tokens.Radius).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
}
func (k *UI) List(state *ui.ListState, count int, row func(int)) ui.Element {
	return ui.List(k.c, state, count, row).Radius(k.tokens.Radius).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
}
func (k *UI) Tree(build func()) ui.Element {
	return ui.Tree(k.c, build).Padding(8).Radius(k.tokens.Radius).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
}
func (k *UI) TreeItem(label string, open *bool, build func()) ui.Element {
	return ui.TreeItem(k.c, label, open, build)
}
func (k *UI) Calendar(value *time.Time) ui.Element {
	return ui.Calendar(k.c, value).AlignSelf(ui.Start).Padding(12).Radius(k.tokens.RadiusLarge).Border(1, k.tokens.Native.Border).Background(k.tokens.Elevated)
}
func (k *UI) DateInput(value *time.Time) ui.Element {
	height, _, _ := k.metrics(Medium)
	return ui.DateInput(k.c, value).MinHeight(height).Padding(6, 12).Radius(k.tokens.Radius).Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border)
}
func (k *UI) TimeInput(value *time.Time) ui.Element {
	height, _, _ := k.metrics(Medium)
	return ui.TimeInput(k.c, value).MinHeight(height).Padding(6, 12).Radius(k.tokens.Radius).Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border)
}
func (k *UI) ColorPicker(value *ui.Color) ui.Element {
	return ui.ColorPicker(k.c, value).AlignSelf(ui.Start)
}
