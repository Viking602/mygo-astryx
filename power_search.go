package astryx

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

type FilterField struct {
	ID, Label, Kind string
	Operators       []string
	Options         []Option
	Validate        func(string) error
	Editor          func(*UI, *string)
}
type SearchFilter struct {
	ID, Field, Operator, Value string
	ReadOnly                   bool
}
type PowerSearchState struct {
	Query                                  string
	Filters                                []SearchFilter
	Field, Operator, Value, Editing, Error string
	next                                   int
}
type PowerSearchOptions struct {
	Label              string
	Disabled, ReadOnly bool
	OnChange           func([]SearchFilter)
}

func filterOperators(field FilterField) []string {
	if len(field.Operators) > 0 {
		return field.Operators
	}
	switch field.Kind {
	case "number", "date":
		return []string{"=", ">", "<", ">=", "<="}
	default:
		return []string{"包含", "=", "!="}
	}
}
func validateFilter(field FilterField, operator, value string) string {
	if !slices.Contains(filterOperators(field), operator) {
		return "请选择有效运算符"
	}
	if strings.TrimSpace(value) == "" {
		return "请输入筛选值"
	}
	switch field.Kind {
	case "number":
		n, err := strconv.ParseFloat(value, 64)
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return "请输入有效数字"
		}
	case "date":
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return "日期格式为 YYYY-MM-DD"
		}
	case "enum":
		if !slices.ContainsFunc(field.Options, func(o Option) bool { return !o.Disabled && o.Value == value }) {
			return "请选择有效选项"
		}
	}
	if field.Validate != nil {
		if err := field.Validate(value); err != nil {
			return err.Error()
		}
	}
	return ""
}
func (k *UI) PowerSearch(state *PowerSearchState, fields []FilterField, options ...PowerSearchOptions) ui.Element {
	o := option(options)
	if o.Label == "" {
		o.Label = "结构化搜索"
	}
	root := ui.Column(k.c).Gap(12).Label(o.Label).Disabled(o.Disabled)
	changed := func() {
		root.SetChanged()
		if o.OnChange != nil {
			o.OnChange(slices.Clone(state.Filters))
		}
	}
	return root.Children(func() {
		k.Key("free-query").SearchField(&state.Query, InputOptions{Label: o.Label, Placeholder: "搜索文字…", Disabled: o.Disabled, ReadOnly: o.ReadOnly}).FillWidth().OnChange(changed)
		if len(state.Filters) > 0 {
			k.Inline().Children(func() {
				for _, filter := range state.Filters {
					label := filter.Field
					for _, field := range fields {
						if field.ID == filter.Field {
							label = field.Label
						}
					}
					var remove func()
					if !filter.ReadOnly && !o.ReadOnly && !o.Disabled {
						remove = func() {
							state.Filters = slices.DeleteFunc(state.Filters, func(f SearchFilter) bool { return f.ID == filter.ID })
							changed()
						}
					}
					k.Key(filter.ID).Token(label+" "+filter.Operator+" "+filter.Value, remove)
					if !filter.ReadOnly && !o.ReadOnly {
						k.Button("编辑", ButtonOptions{Size: Small, Variant: Ghost, Disabled: o.Disabled}).Label("编辑条件" + filter.ID).OnClick(func() {
							state.Editing, state.Field, state.Operator, state.Value = filter.ID, filter.Field, filter.Operator, filter.Value
						})
					}
				}
			})
		}
		if len(fields) == 0 {
			k.Caption("暂无可筛选字段")
			return
		}
		if o.ReadOnly {
			return
		}
		if state.Editing != "" && !slices.ContainsFunc(state.Filters, func(f SearchFilter) bool { return f.ID == state.Editing && !f.ReadOnly }) {
			state.Editing, state.Value, state.Error = "", "", ""
		}
		selected := 0
		for i, field := range fields {
			if field.ID == state.Field {
				selected = i
			}
		}
		if !slices.ContainsFunc(fields, func(f FilterField) bool { return f.ID == state.Field }) {
			state.Field = fields[0].ID
			selected = 0
			state.Operator = ""
			state.Value = ""
		}
		field := fields[selected]
		operators := filterOperators(field)
		if !slices.Contains(operators, state.Operator) {
			state.Operator = operators[0]
		}
		k.Inline().Children(func() {
			fieldOptions := make([]Option, 0, len(fields))
			for _, f := range fields {
				fieldOptions = append(fieldOptions, Option{Value: f.ID, Label: f.Label})
			}
			k.Key("field").Select(&state.Field, fieldOptions, InputOptions{Label: "筛选字段", Disabled: o.Disabled}).Width(150).OnChange(func() { state.Value = ""; state.Operator = ""; state.Error = "" })
			operatorOptions := make([]Option, 0, len(operators))
			for _, op := range operators {
				operatorOptions = append(operatorOptions, Option{Value: op, Label: op})
			}
			k.Key("operator").Select(&state.Operator, operatorOptions, InputOptions{Label: "筛选运算符", Disabled: o.Disabled}).Width(110)
			ui.Column(k.c.Key("filter-value")).MinWidth(140).Grow(1).Children(func() {
				if field.Editor != nil {
					field.Editor(k, &state.Value)
				} else if field.Kind == "enum" {
					k.Select(&state.Value, field.Options, InputOptions{Label: "筛选值", Disabled: o.Disabled}).FillWidth()
				} else {
					k.TextInput(&state.Value, InputOptions{Label: "筛选值", Placeholder: map[bool]string{true: "YYYY-MM-DD", false: "输入筛选值"}[field.Kind == "date"], Disabled: o.Disabled}).FillWidth()
				}
			})
			label := "添加条件"
			if state.Editing != "" {
				label = "更新条件"
			}
			k.Button(label, ButtonOptions{Variant: Primary, Disabled: o.Disabled}).OnClick(func() {
				current := field
				for _, f := range fields {
					if f.ID == state.Field {
						current = f
					}
				}
				state.Error = validateFilter(current, state.Operator, state.Value)
				if state.Error != "" {
					return
				}
				if state.Editing != "" {
					for i := range state.Filters {
						if state.Filters[i].ID == state.Editing && !state.Filters[i].ReadOnly {
							state.Filters[i].Field, state.Filters[i].Operator, state.Filters[i].Value = state.Field, state.Operator, state.Value
						}
					}
				} else {
					state.next++
					id := fmt.Sprintf("filter-%d", state.next)
					for slices.ContainsFunc(state.Filters, func(f SearchFilter) bool { return f.ID == id }) {
						state.next++
						id = fmt.Sprintf("filter-%d", state.next)
					}
					state.Filters = append(state.Filters, SearchFilter{ID: id, Field: state.Field, Operator: state.Operator, Value: state.Value})
				}
				state.Value, state.Editing = "", ""
				changed()
			})
		})
		if state.Error != "" {
			root.Error(state.Error)
			k.Caption(state.Error).TextColor(k.tokens.Native.Danger)
		}
		if len(state.Filters) > 0 {
			k.Button("清除可编辑条件", ButtonOptions{Variant: Ghost, Size: Small, Disabled: o.Disabled}).AlignSelf(ui.Start).OnClick(func() {
				state.Filters = slices.DeleteFunc(state.Filters, func(f SearchFilter) bool { return !f.ReadOnly })
				changed()
			})
		}
	})
}
