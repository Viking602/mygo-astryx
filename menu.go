package astryx

import (
	"github.com/egoist/mygo/ui"
	"strings"
)

// MenuEntry is a stable, application-owned menu model. Checked and Radio
// bindings allow controlled checkbox/radio items; Children creates a submenu.
type MenuEntry struct {
	ID, Label           string
	Group               string
	Disabled, Separator bool
	Checked             *bool
	Radio               *string
	Value               string
	Children            []MenuEntry
	Modifiers           ui.Modifiers
	Shortcut            ui.Key
	OnSelect            func()
}

func buildMenu(m *ui.Menu, entries []MenuEntry) {
	group := ""
	for _, entry := range entries {
		if entry.Group != "" && entry.Group != group {
			if group != "" {
				m.Separator()
			}
			group = entry.Group
			m.Item(group).Disabled(true)
		}
		if entry.Separator {
			m.Separator()
			continue
		}
		if len(entry.Children) > 0 {
			m.Submenu(entry.Label, func(sub *ui.Menu) { buildMenu(sub, entry.Children) }).Disabled(entry.Disabled)
			continue
		}
		item := m.Item(entry.Label).Disabled(entry.Disabled)
		if entry.Checked != nil {
			item.Checked(*entry.Checked)
		}
		if entry.Radio != nil {
			item.Radio(*entry.Radio == entry.Value)
		}
		if entry.Shortcut != 0 {
			item.Shortcut(entry.Modifiers, entry.Shortcut)
		}
		if item.Chosen() && !entry.Disabled {
			if entry.Checked != nil {
				*entry.Checked = !*entry.Checked
			}
			if entry.Radio != nil {
				*entry.Radio = entry.Value
			}
			if entry.OnSelect != nil {
				entry.OnSelect()
			}
		}
	}
}
func (k *UI) DropdownMenu(label string, entries []MenuEntry, options ...ButtonOptions) ui.Element {
	o := option(options)
	if o.Variant == Default {
		o.Variant = Outline
	}
	o.EndIcon = IconChevronDown
	return k.Button(label, o).Label(label).PopoverMenuButton(func(m *ui.Menu) { buildMenu(m, entries) })
}
func (k *UI) MoreMenu(label string, entries []MenuEntry) ui.Element {
	return k.IconButton(label, IconMore, ButtonOptions{Variant: Ghost}).PopoverMenuButton(func(m *ui.Menu) { buildMenu(m, entries) })
}

type ContextMenuOptions struct{ Native bool }

func (k *UI) ContextMenu(target ui.Element, entries []MenuEntry, options ...ContextMenuOptions) ui.Element {
	if option(options).Native {
		return target.ContextMenu(func(m *ui.Menu) { buildMenu(m, entries) })
	}
	return target.PopoverMenu(func(m *ui.Menu) { buildMenu(m, entries) })
}
func (k *UI) NavMenu(label string, entries []MenuEntry) ui.Element {
	return k.Button(label, ButtonOptions{Variant: Ghost, EndIcon: IconChevronDown}).Label(label).PopoverMenuButton(func(m *ui.Menu) { buildMenu(m, entries) })
}
func (k *UI) Toolbar(label string, build func()) ui.Element {
	prev := k.inGroup
	k.inGroup = true
	defer func() { k.inGroup = prev }()
	return ui.Toolbar(k.c, build).AlignSelf(ui.Start).Label(label).Padding(4).Gap(4).Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border).Radius(k.tokens.Radius)
}

type Command struct {
	ID, Label, Description, Group, Shortcut string
	Disabled                                bool
	Run                                     func()
}
type CommandPaletteState struct {
	Query      string
	Active     int
	Chosen     string
	visible    bool
	lastQuery  string
	lastActive int
}
type CommandPaletteOptions struct {
	Title, Placeholder, Empty string
	Width                     float32
	Inline                    bool
	Footer                    func()
}

func commandMatches(items []Command, query string) []Command {
	words := strings.Fields(strings.ToLower(query))
	var matches []Command
	for _, item := range items {
		text := strings.ToLower(item.Label + " " + item.Description + " " + item.Group)
		ok := true
		for _, word := range words {
			if !strings.Contains(text, word) {
				ok = false
				break
			}
		}
		if ok {
			matches = append(matches, item)
		}
	}
	return matches
}
func commandIndex(items []Command, start, direction int) int {
	if len(items) == 0 {
		return -1
	}
	if start < 0 && direction < 0 {
		start = 0
	}
	for i := 0; i < len(items); i++ {
		start = (start + direction + len(items)) % len(items)
		if !items[start].Disabled {
			return start
		}
	}
	return -1
}
func (k *UI) CommandPalette(open *bool, state *CommandPaletteState, items []Command, options ...CommandPaletteOptions) ui.Element {
	o := option(options)
	if !*open && !o.Inline {
		state.visible = false
	}
	if o.Title == "" {
		o.Title = "命令面板"
	}
	if o.Placeholder == "" {
		o.Placeholder = "搜索命令…"
	}
	if o.Empty == "" {
		o.Empty = "没有匹配的命令"
	}
	build := func() {
		matches := commandMatches(items, state.Query)
		if state.Active < 0 || state.Active >= len(matches) || len(matches) > 0 && matches[state.Active].Disabled {
			state.Active = commandIndex(matches, -1, 1)
		}
		choose := func(item Command) {
			if item.Disabled {
				return
			}
			state.Chosen = item.ID
			*open = false
			if item.Run != nil {
				item.Run()
			}
		}
		input := k.Key("command-query").TextInput(&state.Query, InputOptions{Label: "搜索命令", Placeholder: o.Placeholder}).AutoFocus().FillWidth()
		composing := input.Composing()
		input.OnTextKey(func(key ui.TextKey) bool {
			now := commandMatches(items, key.Text())
			switch key.Key {
			case ui.KeyDown:
				state.Active = commandIndex(now, state.Active, 1)
				return true
			case ui.KeyUp:
				state.Active = commandIndex(now, state.Active, -1)
				return true
			case ui.KeyEnter:
				if !composing && state.Active >= 0 && state.Active < len(now) {
					choose(now[state.Active])
					return true
				}
			case ui.KeyEscape:
				*open = false
				return true
			}
			return false
		})
		if input.Changed() {
			state.Active = commandIndex(commandMatches(items, state.Query), -1, 1)
		}
		matches = commandMatches(items, state.Query)
		scrollActive := !state.visible || state.lastQuery != state.Query || state.lastActive != state.Active
		ui.Scroll(k.c).MaxHeight(360).MinHeight(40).Children(func() {
			ui.Column(k.c).Gap(3).Children(func() {
				if len(matches) == 0 {
					k.Caption(o.Empty).Padding(18)
				}
				group := ""
				for i, item := range matches {
					if item.Group != "" && group != item.Group {
						group = item.Group
						k.Caption(group).FontWeight(600).Padding(10, 8, 4, 8)
					}
					row := ui.ButtonBase(k.c.Key(item.ID)).Label(item.Label).Disabled(item.Disabled).Padding(10, 12).Gap(12).Radius(k.tokens.RadiusSmall)
					if i == state.Active {
						row.Background(k.tokens.Subtle)
					}
					row.Children(func() {
						ui.Column(k.c).Grow(1).Gap(3).Children(func() {
							k.Text(item.Label)
							if item.Description != "" {
								k.Caption(item.Description)
							}
						})
						if item.Shortcut != "" {
							k.Kbd(item.Shortcut)
						}
					})
					if i == state.Active && scrollActive {
						row.ScrollIntoView()
					}
					row.OnClick(func() { choose(item) })
				}
			})
		})
		state.visible, state.lastQuery, state.lastActive = true, state.Query, state.Active
		if o.Footer != nil {
			o.Footer()
		} else {
			k.Caption("↑ ↓ 移动 · Enter 执行 · Esc 关闭")
		}
	}
	if o.Inline {
		return k.Card().Children(build)
	}
	return k.Dialog(o.Title, open, build, DialogOptions{Width: o.Width})
}
