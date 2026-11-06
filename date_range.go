package astryx

import (
	"github.com/egoist/mygo/ui"
	"time"
)

type DateRange struct {
	Start, End time.Time
	Empty      bool
}
type DateRangePreset struct {
	Label      string
	Start, End time.Time
}
type DateRangeOptions struct {
	Label    string
	Min, Max time.Time
	MaxDays  int
	Disabled bool
	Presets  []DateRangePreset
	Allowed  func(time.Time) bool
	OnChange func(DateRange)
}

func dayUTC(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
func ValidateDateRange(value DateRange, o DateRangeOptions) string {
	if value.Empty || value.Start.IsZero() && value.End.IsZero() {
		return ""
	}
	if value.Start.IsZero() || value.End.IsZero() {
		return "请选择开始和结束日期"
	}
	start, end := dayUTC(value.Start), dayUTC(value.End)
	if end.Before(start) {
		return "结束日期不能早于开始日期"
	}
	if !o.Min.IsZero() && start.Before(dayUTC(o.Min)) {
		return "开始日期早于允许范围"
	}
	if !o.Max.IsZero() && end.After(dayUTC(o.Max)) {
		return "结束日期超出允许范围"
	}
	days := end.Unix()/86400 - start.Unix()/86400 + 1
	if o.MaxDays > 0 && days > int64(o.MaxDays) {
		return "所选日期范围过长"
	}
	if o.Allowed != nil {
		if days > 36600 {
			return "带日期规则的范围不能超过 100 年"
		}
		for d := dateOnly(value.Start); !dayUTC(d).After(end); d = d.AddDate(0, 0, 1) {
			if !o.Allowed(d) {
				return "日期范围包含不可用日期"
			}
		}
	}
	return ""
}

// dateOnly preserves the caller's calendar region while removing clock time.
func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

type dateRangeInputState struct {
	Editing bool
	Draft   DateRange
}

func (k *UI) DateRangeInput(value *DateRange, options ...DateRangeOptions) ui.Element {
	o := option(options)
	if o.Label == "" {
		o.Label = "日期范围"
	}
	root := ui.Column(k.c).Gap(8).Label(o.Label).Disabled(o.Disabled)
	state := ui.Local(root, "date-range-editor", func() dateRangeInputState { return dateRangeInputState{} })
	commit := func(candidate DateRange) {
		if o.Disabled || ValidateDateRange(candidate, o) != "" {
			return
		}
		*value = candidate
		state.Editing = false
		root.SetChanged()
		if o.OnChange != nil {
			o.OnChange(candidate)
		}
	}
	editors := func(target *DateRange, onChange func(), clear func()) {
		ui.Row(k.c).AlignItems(ui.Center).Gap(8).Wrap().Children(func() {
			k.Key("start").DateInput(&target.Start).Label(o.Label + "开始").Disabled(o.Disabled).OnChange(onChange)
			k.Caption("至")
			k.Key("end").DateInput(&target.End).Label(o.Label + "结束").Disabled(o.Disabled).OnChange(onChange)
			if clear != nil {
				k.Button("清除"+o.Label, ButtonOptions{Variant: Ghost, Size: Small, Disabled: o.Disabled}).OnClick(clear)
			}
		})
	}
	root.Children(func() {
		empty := value.Empty || value.Start.IsZero() && value.End.IsZero()
		if empty {
			if !state.Editing {
				k.Button("选择日期范围", ButtonOptions{Variant: Outline, Disabled: o.Disabled}).AlignSelf(ui.Start).OnClick(func() {
					start := dateOnly(k.c.Now())
					if !o.Min.IsZero() && dayUTC(start).Before(dayUTC(o.Min)) {
						start = dateOnly(o.Min)
					}
					if !o.Max.IsZero() && dayUTC(start).After(dayUTC(o.Max)) {
						start = dateOnly(o.Max)
					}
					days := 7
					if o.MaxDays > 0 {
						days = min(days, o.MaxDays)
					}
					end := start.AddDate(0, 0, days-1)
					if !o.Max.IsZero() && dayUTC(end).After(dayUTC(o.Max)) {
						end = dateOnly(o.Max)
					}
					candidate := DateRange{Start: start, End: end}
					if ValidateDateRange(candidate, o) == "" {
						commit(candidate)
					} else {
						state.Draft = candidate
						state.Editing = true
					}
				})
			} else {
				ui.Column(k.c.Key("range-draft")).Gap(8).Children(func() {
					editors(&state.Draft, nil, nil)
					problem := ValidateDateRange(state.Draft, o)
					if problem != "" {
						root.Error(problem)
						k.Caption("当前日期范围不可用，请选择其他日期").TextColor(k.tokens.Native.Danger)
						k.Caption(problem).TextColor(k.tokens.Native.Danger)
					}
					k.Inline().Children(func() {
						k.Button("应用日期范围", ButtonOptions{Variant: Primary, Disabled: o.Disabled || problem != ""}).OnClick(func() { commit(state.Draft) })
						k.Button("取消选择日期", ButtonOptions{Variant: Ghost, Disabled: o.Disabled}).OnClick(func() { state.Editing = false })
					})
				})
			}
		} else {
			state.Editing = false
			editors(value, func() {
				root.SetChanged()
				if ValidateDateRange(*value, o) == "" && o.OnChange != nil {
					o.OnChange(*value)
				}
			}, func() { commit(DateRange{Empty: true}) })
			if problem := ValidateDateRange(*value, o); problem != "" {
				root.Error(problem)
				k.Caption(problem).TextColor(k.tokens.Native.Danger)
			}
		}
		if len(o.Presets) > 0 {
			k.Inline().Children(func() {
				for _, preset := range o.Presets {
					candidate := DateRange{Start: dateOnly(preset.Start), End: dateOnly(preset.End)}
					k.Button(preset.Label, ButtonOptions{Size: Small, Variant: Outline, Disabled: o.Disabled || ValidateDateRange(candidate, o) != ""}).OnClick(func() { commit(candidate) })
				}
			})
		}
	})
	return root
}
func (k *UI) DateTimeInput(value *time.Time, options ...InputOptions) ui.Element {
	o := option(options)
	root := ui.Row(k.c).AlignItems(ui.Center).Gap(8).Wrap().Disabled(o.Disabled).Label(o.Label)
	return root.Children(func() {
		date := k.Key("date").DateInput(value).Label(o.Label + "日期").Disabled(o.Disabled)
		clock := k.Key("time").TimeInput(value).Label(o.Label + "时间").Disabled(o.Disabled)
		if date.Changed() || clock.Changed() {
			root.SetChanged()
		}
	})
}
