package astryx

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"slices"
	"strings"
)

type SelectionOptions struct {
	Label       string
	Disabled    bool
	MaxSelected int
	OnChange    func([]string)
}

func changeSelection(values *[]string, value string, on bool) {
	if on {
		if !slices.Contains(*values, value) {
			*values = append(*values, value)
		}
	} else {
		*values = slices.DeleteFunc(*values, func(v string) bool { return v == value })
	}
}
func (k *UI) CheckboxList(values *[]string, items []Option, options ...SelectionOptions) ui.Element {
	o := option(options)
	root := ui.Column(k.c).Gap(8).Role(ui.RoleGroup).Label(o.Label).Disabled(o.Disabled)
	return root.Children(func() {
		for _, item := range items {
			on := slices.Contains(*values, item.Value)
			disabled := o.Disabled || item.Disabled || o.MaxSelected > 0 && len(*values) >= o.MaxSelected && !on
			box := k.Key(item.Value).Checkbox(&on, item.Label, ChoiceOptions{Disabled: disabled})
			if box.Changed() {
				changeSelection(values, item.Value, on)
				root.SetChanged()
				if o.OnChange != nil {
					o.OnChange(slices.Clone(*values))
				}
			}
		}
	})
}
func (k *UI) RadioList(value *string, items []Option, options ...InputOptions) ui.Element {
	o := option(options)
	return k.RadioGroup(o.Label, func() {
		for _, item := range items {
			k.Key(item.Value).Radio(value, item.Value, item.Label, ChoiceOptions{Disabled: o.Disabled || item.Disabled})
		}
	}).Disabled(o.Disabled)
}

type multiSelectorState struct {
	Open     bool
	Query    string
	First    ui.Handle
	HasFirst bool
}

func (k *UI) MultiSelector(values *[]string, items []Option, options ...SelectionOptions) ui.Element {
	o := option(options)
	root := ui.Column(k.c).MinWidth(0)
	state := ui.Local(root, "multi-selector", func() multiSelectorState { return multiSelectorState{} })
	var trigger ui.Element
	root.Children(func() {
		label := o.Label
		if label == "" {
			label = "多项选择"
		}
		summary := "请选择"
		if len(*values) > 0 {
			summary = fmt.Sprintf("已选择 %d 项", len(*values))
		}
		trigger = k.styleTrigger(ui.ButtonBase(k.c).Role(ui.RolePopUpButton).Value(summary).Expanded(state.Open),
			summary, len(*values) == 0, InputOptions{Label: label, Disabled: o.Disabled})
		if trigger.Clicked() {
			state.Open = !state.Open
		}
		if o.Disabled {
			state.Open = false
		}
		ui.PopoverBase(k.c, trigger, &state.Open, func(panel ui.Element) {
			k.popup(panel)
			panel.MinWidth(260).Padding(10).Gap(8)
			k.Key("query").TextInput(&state.Query, InputOptions{Label: "筛选" + label, Placeholder: "搜索选项…"}).AutoFocus().FillWidth().OnTextKey(func(key ui.TextKey) bool {
				if key.Key == ui.KeyDown && state.HasFirst {
					state.First.Focus()
					return true
				}
				return false
			})
			count, first := 0, true
			state.HasFirst = false
			ui.Scroll(k.c).MaxHeight(220).Shrink(1).MinHeight(0).Children(func() {
				ui.Column(k.c).Gap(4).FocusGroup(ui.Vertical).Children(func() {
					for _, item := range items {
						if !strings.Contains(strings.ToLower(item.Label), strings.ToLower(state.Query)) {
							continue
						}
						count++
						on := slices.Contains(*values, item.Value)
						disabled := item.Disabled || o.MaxSelected > 0 && len(*values) >= o.MaxSelected && !on
						e := k.Key(item.Value).Checkbox(&on, item.Label, ChoiceOptions{Disabled: disabled}).Padding(6, 8).Radius(k.tokens.RadiusSmall)
						if e.Hovered() && !disabled {
							e.Background(k.tokens.Subtle)
						}
						if first && !disabled {
							e.Bind(&state.First)
							first = false
							state.HasFirst = true
						}
						if e.Changed() {
							changeSelection(values, item.Value, on)
							root.SetChanged()
							trigger.SetChanged()
							if o.OnChange != nil {
								o.OnChange(slices.Clone(*values))
							}
						}
					}
				})
			})
			if count == 0 {
				k.Caption("没有匹配的选项").Padding(8)
			}
			ui.Row(k.c).Gap(8).Padding(4, 0, 0, 0).BorderWidth(1, 0, 0, 0).BorderColor(k.tokens.Native.Border).Children(func() {
				k.Caption(fmt.Sprintf("已选 %d 项", len(*values))).Grow(1)
				k.Button("清空", ButtonOptions{Variant: Ghost, Size: Small, Disabled: len(*values) == 0}).OnClick(func() {
					*values = nil
					root.SetChanged()
					trigger.SetChanged()
					if o.OnChange != nil {
						o.OnChange(nil)
					}
				})
				k.Button("完成", ButtonOptions{Variant: Primary, Size: Small}).OnClick(func() { state.Open = false })
			})
		})
	})
	return trigger
}

// ComplexSelector delegates its popup body to the caller. The close callback
// captures only application state, so the caller may safely attach it to actions.
func (k *UI) ComplexSelector(label string, open *bool, summary string, content func(close func()), options ...InputOptions) ui.Element {
	o := option(options)
	placeholder := summary == ""
	if summary == "" {
		summary = o.Placeholder
	}
	if summary == "" {
		summary = "请选择"
	}
	o.Label = label
	trigger := k.styleTrigger(ui.ButtonBase(k.c).Role(ui.RolePopUpButton).Value(summary).Expanded(*open), summary, placeholder, o)
	if trigger.Clicked() {
		*open = !*open
	}
	if o.Disabled {
		*open = false
	}
	ui.PopoverBase(k.c, trigger, open, func(panel ui.Element) {
		k.popupContent(panel, 16, 12, func() {
			if content != nil {
				content(func() { *open = false })
			}
		})
	})
	return trigger
}
func (k *UI) Token(label string, onRemove func()) ui.Element {
	return ui.Row(k.c).AlignSelf(ui.Start).AlignItems(ui.Center).Padding(4, 10).Gap(6).Radius(k.tokens.RadiusSmall).
		Background(k.tokens.Subtle).Border(1, k.tokens.Native.Border).Role(ui.RoleGroup).Label(label).Children(func() {
		k.Text(label).FontSize(k.tokens.FontSmall).FontWeight(500).SingleLine()
		if onRemove != nil {
			btn := ui.ButtonBase(k.c).Size(18, 18).Radius(999).Label("移除" + label).OnClick(onRemove)
			if btn.Hovered() {
				btn.Background(k.tokens.Native.SurfaceHover)
			}
			btn.Children(func() { k.Icon(IconClose).FontSize(12).TextColor(k.tokens.Muted) })
		}
	})
}
func (k *UI) Tokenizer(values *[]string, suggestions []string, options ...InputOptions) ui.Element {
	o := option(options)
	height, _, _ := k.metrics(o.Size)
	fill := k.tokens.Elevated
	if o.ReadOnly {
		fill = k.tokens.Subtle
	}
	e := ui.TokenField(k.c, values, suggestions).MinHeight(height).Padding(5, 10).Radius(k.tokens.Radius).
		Background(fill).Disabled(o.Disabled)
	if o.Label != "" {
		e.Label(o.Label)
	}
	return e
}

type Suggestion struct {
	ID, Label, Description, Group string
	Disabled                      bool
}
type TypeaheadState struct {
	Query, Value string
	Loading      bool
	Error        string
}
type TypeaheadOptions struct {
	InputOptions
	OnSelect func(Suggestion)
	OnQuery  func(string)
}

func (k *UI) Typeahead(state *TypeaheadState, items []Suggestion, options ...TypeaheadOptions) ui.Element {
	o := option(options)
	s := ui.ComboboxBase(k.c, &state.Query)
	k.styleInput(s.Input, o.InputOptions)
	if o.Disabled || o.ReadOnly {
		s.SetOpen(false)
	}
	s.Popup(func(panel ui.Element) {
		k.popup(panel)
		ui.Scroll(k.c).MaxHeight(260).MinHeight(0).Children(func() {
			ui.Column(k.c).Gap(2).Children(func() {
				count := 0
				group := ""
				for _, item := range items {
					if s.Filtering() && !strings.Contains(strings.ToLower(item.Label+" "+item.Description), strings.ToLower(state.Query)) {
						continue
					}
					count++
					if item.Group != "" && item.Group != group {
						group = item.Group
						k.Caption(group).Padding(8, 10)
					}
					if item.Disabled {
						k.Caption(item.Label+" · 不可用").Padding(8, 10)
						continue
					}
					row := s.Item(item.ID).Label(item.Label).Padding(8, 10).Radius(k.tokens.RadiusSmall)
					if row.Highlighted() {
						row.Background(k.tokens.Subtle)
					}
					row.Children(func() {
						k.Text(item.Label)
						if item.Description != "" {
							k.Caption(item.Description)
						}
					})
				}
				if state.Loading {
					k.Spinner()
				}
				if state.Error != "" {
					k.Alert("无法获取选项", state.Error, Danger)
				} else if count == 0 && !state.Loading {
					k.Caption("没有匹配的选项").Padding(8)
				}
			})
		})
	})
	if id, ok := s.Chosen(); ok {
		for _, item := range items {
			if item.ID == id && !item.Disabled {
				state.Value, state.Query = item.ID, item.Label
				s.Input.SetChanged()
				if o.OnSelect != nil {
					o.OnSelect(item)
				}
				break
			}
		}
	} else if s.Input.Changed() && o.OnQuery != nil {
		s.Input.OnChange(func() { o.OnQuery(state.Query) })
	}
	return s.Input
}
