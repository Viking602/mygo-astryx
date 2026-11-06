package astryx

import (
	"github.com/egoist/mygo/ui"
	"strconv"
)

type Step struct {
	ID, Label, Description string
	Disabled, Error        bool
}
type StepperOptions struct {
	Vertical, ReadOnly, Linear bool
	OnChange                   func(int)
}

func (k *UI) Stepper(current *int, steps []Step, options ...StepperOptions) ui.Element {
	o := option(options)
	clampChoice(current, len(steps))
	root := ui.Row(k.c).AlignItems(ui.Center).Gap(12).Wrap().Role(ui.RoleList).Label("步骤")
	if o.Vertical {
		root.Column().AlignItems(ui.Start)
	}
	root.Children(func() {
		for i, step := range steps {
			if i > 0 && !o.Vertical {
				lineColor := k.tokens.Native.Border
				if i <= *current {
					lineColor = k.tokens.Native.Accent.Alpha(0.45)
				}
				ui.Box(k.c).Size(24, 1).Background(lineColor)
			}
			label := step.Label
			state := "未开始"
			if i < *current {
				state = "已完成"
			} else if i == *current {
				state = "当前步骤"
			}
			if step.Error {
				state = "需要处理"
			}
			active := i == *current
			e := ui.ButtonBase(k.c.Key(step.ID)).Label(label).Description(state + " " + step.Description).
				Disabled(step.Disabled || o.ReadOnly || o.Linear && i > *current+1).
				Padding(8, 12).Gap(10).Radius(k.tokens.Radius).AlignItems(ui.Center)
			if active {
				e.Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border)
			} else if e.Hovered() && !e.IsDisabled() {
				e.Background(k.tokens.Subtle.Alpha(0.5))
			}
			e.Children(func() {
				ink, fill := k.tokens.Muted, k.tokens.Subtle
				mark := strconv.Itoa(i + 1)
				if i <= *current {
					ink, fill = k.tokens.Native.AccentText, k.tokens.Native.Accent
				}
				if step.Error {
					ink, fill = k.tokens.Native.Danger, k.tokens.DangerMuted
					mark = "!"
				}
				ui.Box(k.c).Size(26, 26).Radius(999).Center().Background(fill).Children(func() {
					if i < *current && !step.Error {
						k.Icon(IconCheck).FontSize(14).TextColor(ink).Role(ui.RoleNone)
					} else {
						k.Text(mark).FontSize(k.tokens.FontSmall).FontWeight(600).TextColor(ink)
					}
				})
				ui.Column(k.c).Gap(2).Children(func() {
					k.Text(label).FontWeight(550)
					if step.Description != "" {
						k.Caption(step.Description)
					}
				})
			})
			if e.Clicked() {
				*current = i
				root.SetChanged()
				if o.OnChange != nil {
					o.OnChange(i)
				}
			}
		}
	})
	return root
}

type NavOptions struct {
	Label     string
	Collapsed *bool
}

func navSelected(item NavItem, value string) bool {
	if item.ID == value {
		return true
	}
	for _, child := range item.Children {
		if navSelected(child, value) {
			return true
		}
	}
	return false
}
func menuForNav(items []NavItem, selected *string) []MenuEntry {
	var entries []MenuEntry
	for _, item := range items {
		entry := MenuEntry{ID: item.ID, Label: item.Label, Disabled: item.Disabled}
		if len(item.Children) > 0 {
			entry.Children = menuForNav(item.Children, selected)
		} else {
			entry.OnSelect = func() { *selected = item.ID }
		}
		entries = append(entries, entry)
	}
	return entries
}
func (k *UI) SideNav(selected *string, items []NavItem, options ...NavOptions) ui.Element {
	o := option(options)
	root := ui.Column(k.c).Gap(4).Padding(10).Radius(k.tokens.Radius).Background(k.tokens.Subtle.Alpha(0.5)).Border(1, k.tokens.Native.Border).Label(o.Label)
	root.Children(func() {
		if o.Collapsed != nil {
			label := "收起导航"
			if *o.Collapsed {
				label = "展开导航"
			}
			k.Button(label, ButtonOptions{Variant: Ghost, Size: Small}).AlignSelf(ui.Start).OnClick(func() { *o.Collapsed = !*o.Collapsed })
		}
		var render func([]NavItem, int)
		render = func(items []NavItem, depth int) {
			for _, item := range items {
				ui.Column(k.c.Key(item.ID)).Children(func() {
					collapsed := o.Collapsed != nil && *o.Collapsed
					if item.Section != "" && !collapsed {
						k.Caption(item.Section).FontWeight(600).Padding(10, 8, 4, 8)
					}
					if len(item.Children) > 0 {
						if collapsed {
							button := k.IconButton(item.Label, item.Icon, ButtonOptions{Variant: Ghost, Disabled: item.Disabled}).PopoverMenuButton(func(m *ui.Menu) { buildMenu(m, menuForNav(item.Children, selected)) })
							if navSelected(item, *selected) {
								button.Selected(true).Background(k.tokens.Native.SurfaceHover).TextColor(k.tokens.Native.Text)
							}
							k.Tooltip(button, item.Label)
							return
						}
						open := ui.Local(ui.Box(k.c), "nav-open", func() bool { return true })
						k.Collapsible(item.Label, open, func() { render(item.Children, depth+1) }).Disabled(item.Disabled).Padding(8, 12)
						return
					}
					active := *selected == item.ID
					row := ui.ButtonBase(k.c).Justify(ui.Start).Label(item.Label).Disabled(item.Disabled).Padding(8, 12+float32(depth)*14).Gap(10).Radius(k.tokens.RadiusSmall)
					if active {
						row.Selected(true).Background(k.tokens.Native.SurfaceHover).TextColor(k.tokens.Native.Text)
					} else if row.Hovered() {
						row.Background(k.tokens.Subtle)
					}
					row.Children(func() {
						if item.Icon != nil {
							k.Icon(item.Icon).FontSize(16)
						}
						if o.Collapsed != nil && *o.Collapsed {
							if item.Icon == nil {
								k.Text(item.Label).MaxLines(1).MaxWidth(28)
							}
						} else {
							txt := k.Text(item.Label).Grow(1)
							if active {
								txt.FontWeight(600)
							}
							if item.Badge != "" {
								k.Badge(item.Badge)
							}
						}
					})
					if row.Clicked() {
						*selected = item.ID
						root.SetChanged()
					}
					if collapsed {
						k.Tooltip(row, item.Label)
					}
				})
			}
		}
		render(items, 0)
	})
	return root
}
func (k *UI) TopNav(selected *string, items []NavItem, options ...NavOptions) ui.Element {
	o := option(options)
	return ui.Row(k.c).AlignItems(ui.Center).Gap(4).Wrap().FocusGroup(ui.Horizontal).Label(o.Label).Children(func() {
		for _, item := range items {
			var e ui.Element
			if len(item.Children) > 0 {
				e = k.NavMenu(item.Label, menuForNav(item.Children, selected)).Disabled(item.Disabled)
			} else {
				e = k.Button(item.Label, ButtonOptions{Variant: Ghost, Disabled: item.Disabled}).OnClick(func() { *selected = item.ID })
			}
			if navSelected(item, *selected) {
				e.Selected(true).Background(k.tokens.Native.SurfaceHover).TextColor(k.tokens.Native.Text)
			}
		}
	})
}
func (k *UI) MobileNav(selected *string, items []NavItem) ui.Element {
	return ui.Row(k.c).AlignItems(ui.Center).Padding(4).Gap(4).Radius(k.tokens.Radius).Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border).Wrap().FocusGroup(ui.Horizontal).Label("移动导航").Children(func() {
		for _, item := range items {
			var e ui.Element
			if len(item.Children) > 0 {
				e = k.Key(item.ID).NavMenu(item.Label, menuForNav(item.Children, selected)).Disabled(item.Disabled).Grow(1)
			} else {
				e = k.Key(item.ID).Button(item.Label, ButtonOptions{Variant: Ghost, Size: Small, Icon: item.Icon, Disabled: item.Disabled}).Grow(1).OnClick(func() { *selected = item.ID })
			}
			if navSelected(item, *selected) {
				e.Selected(true).Background(k.tokens.Elevated).TextColor(k.tokens.Native.Text).Shadow(0, 1, 2, 0, k.tokens.Shadow)
			}
		}
	})
}

// MobileNavToggle binds an application's drawer/sheet visibility. The same
// boolean must be passed to that surface so Escape and outside presses update
// the trigger's expanded state as well.
func (k *UI) MobileNavToggle(label string, open *bool, options ...ButtonOptions) ui.Element {
	o := option(options)
	if o.Variant == Default {
		o.Variant = Ghost
	}
	return k.Button(label, o).AlignSelf(ui.Start).Expanded(*open).OnClick(func() { *open = !*open })
}
func (k *UI) OverflowList(items []NavItem, visible int, onSelect func(string)) ui.Element {
	visible = max(visible, 0)
	return k.Inline().Children(func() {
		for _, item := range items[:min(visible, len(items))] {
			k.Button(item.Label, ButtonOptions{Variant: Ghost, Size: Small, Disabled: item.Disabled}).OnClick(func() {
				if onSelect != nil {
					onSelect(item.ID)
				}
			})
		}
		if len(items) > visible {
			entries := []MenuEntry{}
			for _, item := range items[visible:] {
				entries = append(entries, MenuEntry{ID: item.ID, Label: item.Label, Disabled: item.Disabled, OnSelect: func() {
					if onSelect != nil {
						onSelect(item.ID)
					}
				}})
			}
			k.MoreMenu("更多项目", entries)
		}
	})
}

type OutlineItem struct {
	ID, Title string
	Level     int
}

func (k *UI) Outline(selected *string, items []OutlineItem, onNavigate func(string)) ui.Element {
	return ui.Column(k.c).Gap(2).Label("内容目录").Children(func() {
		for _, item := range items {
			active := *selected == item.ID
			indent := float32(max(item.Level-1, 0))*14 + 12
			borderColor := k.tokens.Native.Border
			if active {
				borderColor = k.tokens.Native.Accent
			}
			e := ui.ButtonBase(k.c.Key(item.ID)).Justify(ui.Start).Padding(6, 12).PaddingStart(indent).
				BorderWidth(0, 0, 0, 2).BorderColor(borderColor).Radius(0, k.tokens.RadiusSmall, k.tokens.RadiusSmall, 0).Label(item.Title)
			if active {
				e.Selected(true).Background(k.tokens.Subtle).TextColor(k.tokens.Native.Text)
			} else if e.Hovered() {
				e.Background(k.tokens.Subtle.Alpha(0.5))
			}
			e.Children(func() {
				txt := k.Text(item.Title).FontSize(k.tokens.FontSmall).SingleLine()
				if active {
					txt.FontWeight(600)
				} else {
					txt.TextColor(k.tokens.Muted)
				}
			})
			e.OnClick(func() {
				*selected = item.ID
				if onNavigate != nil {
					onNavigate(item.ID)
				}
			})
		}
	})
}
func (k *UI) TreeList(state *ui.OutlineState[string], roots []string, children func(string) []string, label func(string) string) ui.Element {
	if children == nil {
		children = func(string) []string { return nil }
	}
	return ui.Outline(k.c, state, roots, children, func(id string) {
		text := id
		if label != nil {
			text = label(id)
		}
		k.Text(text).Padding(6, 8)
	}).Padding(4).Background(k.tokens.Elevated).Border(1, k.tokens.Native.Border).Radius(k.tokens.Radius)
}
