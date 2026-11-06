package astryx

import (
	"fmt"
	"github.com/egoist/mygo/ui"
)

func clampChoice(selected *int, count int) {
	if count > 0 {
		*selected = min(max(*selected, 0), count-1)
	}
}

// Tabs returns the tab list; the app renders its selected panel separately.
func (k *UI) Tabs(selected *int, labels []string) ui.Element {
	clampChoice(selected, len(labels))
	t := ui.TabsBase(k.c, selected, len(labels))
	t.List.Gap(4).Padding(0, 0, 6, 0).BorderWidth(0, 0, 1, 0).BorderColor(k.tokens.Native.Border).Label("选项卡")
	t.List.Children(func() {
		for i, label := range labels {
			e := t.Tab(i).Label(label).Padding(7, 14).Radius(k.tokens.RadiusSmall)
			active := i == *selected
			if active {
				e.Background(k.tokens.Subtle).TextColor(k.tokens.Native.Text)
			} else if e.Hovered() {
				e.Background(k.tokens.Subtle.Alpha(0.5)).TextColor(k.tokens.Native.Text)
			} else {
				e.TextColor(k.tokens.Muted)
			}
			e.Children(func() {
				txt := k.Text(label).SingleLine()
				if active {
					txt.FontWeight(600)
				} else {
					txt.TextColor(k.tokens.Muted)
				}
			})
		}
	})
	return t.List
}
func (k *UI) Segmented(selected *int, labels []string) ui.Element {
	clampChoice(selected, len(labels))
	s := ui.SegmentedBase(k.c, selected, len(labels))
	s.Track.AlignSelf(ui.Start).Padding(3).Gap(3).Radius(k.tokens.Radius).
		Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border).Label("分段选择")
	s.Track.Children(func() {
		for i, label := range labels {
			e := s.Segment(i).Label(label).Padding(5, 13).Radius(k.tokens.RadiusSmall)
			active := i == *selected
			if active {
				e.Background(k.tokens.Elevated).Shadow(0, 1, 3, 0, k.tokens.Shadow)
			} else if e.Hovered() {
				e.Background(k.tokens.Native.SurfaceHover)
			}
			e.Children(func() {
				txt := k.Text(label).SingleLine()
				if active {
					txt.FontWeight(600)
				} else {
					txt.TextColor(k.tokens.Muted)
				}
			})
		}
	})
	return s.Track
}

type Breadcrumb struct {
	Label   string
	OnClick func()
	Menu    []MenuEntry
}

func (k *UI) Breadcrumbs(items []Breadcrumb) ui.Element {
	return ui.Row(k.c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(4).Wrap().Children(func() {
		for i, item := range items {
			if i > 0 {
				k.Icon(IconChevronRight).FontSize(12).TextColor(k.tokens.Muted).Margin(0, 2)
			}
			if len(item.Menu) > 0 {
				k.DropdownMenu(item.Label, item.Menu, ButtonOptions{Variant: Ghost, Size: Small})
			} else if item.OnClick != nil {
				k.Button(item.Label, ButtonOptions{Variant: Ghost, Size: Small}).OnClick(item.OnClick)
			} else {
				k.Text(item.Label).FontSize(k.tokens.FontSmall).FontWeight(600).Padding(5, 8)
			}
		}
	})
}

// Pagination uses one-based page numbers and clamps its model when totals shrink.
func (k *UI) Pagination(page *int, pages int) ui.Element {
	pages = max(pages, 1)
	*page = min(max(*page, 1), pages)
	return ui.Row(k.c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(6).Children(func() {
		if k.Button("上一页", ButtonOptions{Variant: Outline, Size: Small, Disabled: *page == 1}).Clicked() {
			*page--
		}
		start := max(1, min(*page-2, pages-4))
		end := min(pages, start+4)
		for p := start; p <= end; p++ {
			variant := Ghost
			if p == *page {
				variant = Primary
			}
			if k.Button(fmt.Sprint(p), ButtonOptions{Variant: variant, Size: Small}).MinWidth(32).Label(fmt.Sprintf("第 %d 页", p)).Clicked() {
				*page = p
			}
		}
		if k.Button("下一页", ButtonOptions{Variant: Outline, Size: Small, Disabled: *page == pages}).Clicked() {
			*page++
		}
	})
}

type NavItem struct {
	ID, Label string
	Badge     string
	Icon      *ui.SVG
	Disabled  bool
	Section   string
	Children  []NavItem
}

func (k *UI) Sidebar(selected *string, items []NavItem) ui.Element {
	return ui.Column(k.c).Gap(3).Children(func() {
		for _, item := range items {
			active := *selected == item.ID
			e := ui.ButtonBase(k.c.Key(item.ID)).Justify(ui.SpaceBetween).Label(item.Label).Padding(8, 12).Radius(k.tokens.Radius).Gap(10)
			if active {
				e.Selected(true).Background(k.tokens.Native.SurfaceHover).TextColor(k.tokens.Native.Text)
			} else if e.Hovered() {
				e.Background(k.tokens.Subtle)
			}
			e.Children(func() {
				txt := k.Text(item.Label).Grow(1).SingleLine()
				if active {
					txt.FontWeight(600)
				}
				if item.Badge != "" {
					k.Caption(item.Badge).FontFeatures("tnum").Padding(1, 7).Radius(999).Background(k.tokens.Elevated)
				}
			})
			if e.Clicked() {
				*selected = item.ID
			}
		}
	})
}

func (k *UI) Collapsible(label string, open *bool, build func()) ui.Element {
	d := ui.CollapsibleBase(k.c, open)
	d.Trigger.AlignSelf(ui.Stretch).Justify(ui.SpaceBetween).Padding(10, 4).Gap(12).Radius(k.tokens.RadiusSmall).Label(label).Children(func() {
		k.Text(label).FontWeight(550).Grow(1)
		ic := k.Icon(IconChevronDown).FontSize(14).TextColor(k.tokens.Muted)
		if *open {
			ic.Rotate(180)
		}
	})
	d.Panel(func() { ui.Column(k.c).Padding(4, 4, 12, 4).Gap(8).Children(build) })
	return d.Trigger
}

type Disclosure struct {
	ID, Title string
	Open      *bool
	Content   func()
}

func (k *UI) Accordion(items []Disclosure) ui.Element {
	return ui.Column(k.c).Children(func() {
		for _, item := range items {
			ui.Column(k.c.Key(item.ID)).BorderWidth(0, 0, 1, 0).BorderColor(k.tokens.Native.Border).Children(func() { k.Collapsible(item.Title, item.Open, item.Content) })
		}
	})
}
