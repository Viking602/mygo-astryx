package gallery

import (
	"fmt"
	kit "github.com/Viking602/mygo-astryx"
	"github.com/egoist/mygo/ui"
	"strings"
	"time"
)

var cachedNames = map[string][]string{}

func splitNames(names string) []string {
	if s, ok := cachedNames[names]; ok {
		return s
	}
	s := strings.Fields(names)
	cachedNames[names] = s
	return s
}

func (a *App) AdvancedDemos(c *ui.Context) []Demo {
	s := &a.Advanced
	filterPage := c != nil && strings.TrimSpace(a.Query) == ""
	var demos []Demo
	add := func(id, category, title, description, names, code string, build func(*kit.UI)) {
		if filterPage && category != a.Page {
			return
		}
		demos = append(demos, Demo{ID: id, Category: category, Title: title, Description: description, Names: splitNames(names), Code: code, Build: build})
	}
	add("geometry", "layout", "比例、居中与跨列", "保持比例并在栅格中分配空间。", "AspectRatio Center HStack GridSpan", `k.AspectRatio(16.0/9).Width(240).Children(content)`, func(k *kit.UI) {
		k.HStack().Children(func() {
			k.AspectRatio(16.0 / 9).Width(220).Radius(8).Background(k.Tokens().Subtle).Border(1, k.Tokens().Native.Border).Children(func() { k.Center().Fill().Children(func() { k.Text("16 : 9").Bold() }) })
			k.Caption("HStack 保持单行，Inline 可自动换行。").MaxWidth(220)
		})
		k.Grid(3).Children(func() {
			k.GridSpan(2).Padding(16).Radius(k.Tokens().RadiusLarge).Background(k.Tokens().Subtle).Border(1, k.Tokens().Native.Border).Children(func() { k.Text("跨两列") })
			k.Card().Children(func() { k.Text("一列") })
		})
	})
	add("layout-slots", "layout", "完整页面插槽", "Header、Content、Panel 和 Footer 组合在同一布局中。", "Layout Section", `k.Layout(astryx.LayoutSlots{Header: header, Content: content, Panel: panel, Footer: footer})`, func(k *kit.UI) {
		ui.Box(c).Height(240).Border(1, k.Tokens().Native.Border).Radius(8).Clip().Children(func() {
			k.Layout(kit.LayoutSlots{Header: func() { k.Text("项目详情").Bold() }, Content: func() { k.Section("概览", "统一的页面结构。", func() { k.Text("主要内容") }) }, Panel: func() { k.Caption("附加信息"); k.Badge("已同步", kit.Success).AlignSelf(ui.Start) }, Footer: func() { k.Caption("最后保存 · 刚刚") }})
		})
	})
	add("resizable", "layout", "约束与折叠分栏", "拖动分隔线，聚焦后用方向键调整；支持折叠。", "Resizable", `k.Resizable(&state.Panes, left, right, astryx.ResizableOptions{Min: 100, Max: 320, Collapsible: true})`, func(k *kit.UI) {
		k.Key("advanced-resize").Resizable(&s.Resize, func() { k.Text("侧栏").Padding(12) }, func() {
			k.Card().Margin(8).Children(func() { k.Text("工作区域"); k.Caption(fmt.Sprintf("面板宽度 %.0f DIP", s.Resize.Size)) })
		}, kit.ResizableOptions{Min: 100, Max: 320, Collapsible: true}).Height(200).Border(1, k.Tokens().Native.Border)
	})
	add("hidden-keyboard", "typography", "键盘说明与隐藏文本", "隐藏文字不占布局空间，仍保留可访问性语义。", "VisuallyHidden Kbd", `k.Kbd("⌘ K"); k.VisuallyHidden("仅供辅助技术阅读")`, func(k *kit.UI) {
		k.Inline().Children(func() { k.Text("打开命令面板"); k.Kbd("⌘ / Ctrl K") })
		k.VisuallyHidden("Gallery 辅助说明：这是原生组件展示")
	})
	add("avatar-presence", "typography", "头像与在线状态", "颜色与形状共同表达状态，状态也会读入头像名称。", "AvatarStatusDot", `k.Avatar("Ada", nil, astryx.AvatarOptions{Status: "在线", StatusTone: astryx.Success})`, func(k *kit.UI) {
		k.Inline().Children(func() {
			k.Avatar("Ada", nil, kit.AvatarOptions{Status: "在线", StatusTone: kit.Success, Size: 48})
			k.Avatar("Grace", nil, kit.AvatarOptions{Status: "离线", StatusTone: kit.Neutral, Size: 48})
			k.Avatar("Linus", nil, kit.AvatarOptions{Status: "连接错误", StatusTone: kit.Danger, Size: 48})
		})
		k.Inline().Children(func() {
			k.AvatarStatusDot("独立离线状态", kit.Neutral, kit.AvatarStatusOptions{AvatarSize: 64})
			k.Caption("独立状态点：离线（圆环）")
		})
	})
	add("avatar-group", "typography", "协作成员", "重叠头像组与同尺寸圆形溢出指示器，点击可查看完整成员角色与在线状态。", "AvatarGroup AvatarGroupOverflow", `k.AvatarGroup(people, 3, astryx.AvatarGroupOptions{OnOverflow: showMembers})`, func(k *kit.UI) {
		people := []kit.Person{
			{Name: "林知夏", Role: "Engineering Lead · 架构负责人", Status: "在线", StatusTone: kit.Success},
			{Name: "陈叙", Role: "Product Designer · 设计系统", Status: "在线", StatusTone: kit.Success},
			{Name: "王敬之", Role: "Systems Engineer · 原生渲染", Status: "离线", StatusTone: kit.Neutral},
			{Name: "李牧", Role: "Runtime Engineer · 布局引擎", Status: "在线", StatusTone: kit.Success},
			{Name: "赵清", Role: "QA Lead · 自动化验证", Status: "请勿打扰", StatusTone: kit.Danger},
		}
		k.Inline().Gap(28).Children(func() {
			ui.Column(c).Gap(6).Children(func() {
				k.Caption("评审成员 (默认 +N 溢出)").FontWeight(550)
				k.AvatarGroup(people, 3, kit.AvatarGroupOptions{OnOverflow: func() { s.MembersOpen = true }})
			})
			ui.Column(c).Gap(6).Children(func() {
				k.Caption("自定义溢出文案").FontWeight(550)
				ui.Row(c).AlignItems(ui.Center).Gap(-8).Children(func() {
					for _, p := range people[:2] {
						k.Avatar(p.Name, p.Image).Border(2, k.Tokens().Elevated)
					}
					k.AvatarGroupOverflow(12, kit.AvatarGroupOptions{OverflowLabel: "更多成员", OnOverflow: func() { s.MembersOpen = true }})
				})
			})
		})
		k.Dialog("协作成员名单", &s.MembersOpen, func() {
			for _, person := range people {
				ui.Row(c.Key(person.Name)).AlignItems(ui.Center).Padding(8, 10).Gap(12).Radius(k.Tokens().Radius).Background(k.Tokens().Subtle.Alpha(0.45)).Border(1, k.Tokens().Native.Border).Children(func() {
					k.Avatar(person.Name, person.Image, kit.AvatarOptions{Size: 40, Status: person.Status, StatusTone: person.StatusTone})
					ui.Column(c).Gap(2).Grow(1).MinWidth(0).Children(func() {
						k.Text(person.Name).FontWeight(600)
						k.Caption(person.Role)
					})
					k.Badge(person.Status, person.StatusTone)
				})
			}
		}, kit.DialogOptions{Description: "点击溢出指示器查看全部协作成员、职责与实时在线状态。"})
	})
	add("rich-content", "typography", "原生富文本与代码", "Markdown、代码高亮、复制与元数据都在原生视图中。", "Markdown CodeBlock Blockquote MetadataList Citation", `k.Markdown("## 标题"); k.CodeBlock(source, "go")`, func(k *kit.UI) {
		k.Markdown("### 有结构的内容\n支持 **粗体**、列表、链接和表格。\n\n- 原生文字选择\n- 保留 Markdown 复制")
		k.CodeBlock("package main\n\nfunc main() {\n    println(\"Hello MyGo\")\n}", "go")
		k.Blockquote("一致性来自共享的设计与交互规则。", "设计规范")
		k.MetadataList([]kit.Metadata{{Label: "运行方式", Value: "Go 原生"}, {Label: "渲染器", Value: "MyGo"}, {Label: "依赖版本", Value: "v0.7.4"}})
		k.Inline().Children(func() {
			k.Citation("[1] MyGo 文档", kit.CitationOptions{Detail: "官方原生 UI 文档", URL: "https://mygo.egoist.dev/docs/ui"})
			k.Citation("[2] Astryx 规范", kit.CitationOptions{Detail: "Meta Astryx 设计体系", URL: "https://github.com/meta/astryx"})
			k.Citation("[3] WAI-ARIA APG", kit.CitationOptions{Detail: "无障碍交互模式指南", URL: "https://www.w3.org/WAI/ARIA/apg/"})
		})
	})
	add("interactive-content", "typography", "可操作的内容块", "嵌套按钮独立响应，避免整卡与按钮同时触发。", "Item ClickableCard", `k.ClickableCard("项目", content).OnClick(openProject)`, func(k *kit.UI) {
		ui.Column(c).Radius(k.Tokens().Radius).Border(1, k.Tokens().Native.Border).Background(k.Tokens().Elevated).Padding(4).Gap(2).Children(func() {
			k.Item("林知夏", kit.ItemOptions{Avatar: "林知夏", Description: "Engineering Lead · 架构负责人", Badge: "管理员", BadgeTone: kit.Info, OnClick: func() { s.MenuAction = "查看成员：林知夏" }})
			k.Item("项目设计", kit.ItemOptions{Description: "主题与组件规范", Trailing: "刚刚", Icon: kit.IconSettings, Selected: true, OnClick: func() { s.MenuAction = "打开项目设计" }})
		})
		k.ClickableCard("打开项目卡片", func() {
			k.Text("原生工作台").Bold()
			k.Caption("点击卡片打开详情。")
			k.Button("卡片内动作", kit.ButtonOptions{Variant: kit.Outline, Size: kit.Small}).AlignSelf(ui.Start).OnClick(func() { s.MenuAction = "只触发了内层按钮" })
		}).OnClick(func() { s.MenuAction = "触发了卡片动作" })
		if s.MenuAction != "" {
			k.Caption(s.MenuAction)
		}
	})
	add("selection-cards", "forms", "可选择卡片", "卡片保留复选语义和键盘操作。", "SelectableCard", `k.SelectableCard(&state.Selected, "团队方案", content)`, func(k *kit.UI) {
		k.SelectableCard(&s.SelectedCard, "团队方案", func() { k.Caption("共享项目、团队权限与统一设置。") })
		k.Caption(fmt.Sprintf("已选择：%v", s.SelectedCard))
	})
	add("form-grid", "forms", "表单布局与辅助动作", "跨列布局、字段状态和可恢复焦点的清除按钮。", "FormLayout FieldStatus InputClearButton", `k.InputGroup(&state.Text, nil, clearButton)`, func(k *kit.UI) {
		k.FormLayout(2).Children(func() {
			k.Field("项目", func() { k.TextInput(&a.GridName).FillWidth() })
			k.Field("邮箱", func() { k.TextInput(&a.GridEmail).FillWidth() })
		})
		k.InputGroup(&s.ClearText, nil, func() { k.InputClearButton(&s.ClearText, &s.ClearHandle) }, kit.InputOptions{Label: "可清除输入", Focus: &s.ClearHandle})
		k.FieldStatus("输入信息只保留在本地演示状态。", kit.Info)
	})
	add("control-indicators", "forms", "可替换的状态视觉", "装饰性视觉不接管父控件的角色、键盘或焦点。", "Indicator CheckboxIndicator RadioIndicator CheckIndicator", `k.CheckboxIndicator(astryx.IndicatorMixed)`, func(k *kit.UI) {
		k.Inline().Children(func() {
			k.Indicator(kit.IndicatorChecked)
			for _, state := range []kit.IndicatorState{kit.IndicatorUnchecked, kit.IndicatorChecked, kit.IndicatorMixed} {
				k.CheckboxIndicator(state)
			}
			k.RadioIndicator(kit.IndicatorUnchecked)
			k.RadioIndicator(kit.IndicatorChecked)
			k.CheckIndicator(kit.IndicatorChecked)
		})
		tokens := k.Tokens()
		tokens.Indicator = func(c *ui.Context, kind, state string, size float32) (ui.Element, bool) {
			if kind != "checkbox" {
				return ui.Element{}, false
			}
			e := ui.Box(c).Size(size, size).Radius(999).Border(1, tokens.Native.Accent).Center()
			if state != "unchecked" {
				e.Background(tokens.Native.Accent).Children(func() { ui.Text(c, "◆").TextColor(tokens.Native.AccentText).Role(ui.RoleNone) })
			}
			return e, true
		}
		kit.Scope(c, tokens, func(custom *kit.UI) { custom.Checkbox(&a.Checked, "自定义指示器，保留原生行为") }).Background(ui.Transparent)
	})
	add("multiple-selection", "forms", "列表选择与可搜索多选", "支持清空、选择上限、禁用项和键盘聚焦。", "CheckboxList RadioList MultiSelector", `k.MultiSelector(&state.Values, items, astryx.SelectionOptions{MaxSelected: 2})`, func(k *kit.UI) {
		items := []kit.Option{{Value: "go", Label: "Go"}, {Value: "rust", Label: "Rust"}, {Value: "swift", Label: "Swift"}, {Value: "blocked", Label: "暂未开放", Disabled: true}}
		k.Key("multiple").MultiSelector(&s.Multi, items, kit.SelectionOptions{Label: "语言多选", MaxSelected: 2}).FillWidth()
		cols := 2
		if w, _ := c.Size(); w < 640 {
			cols = 1
		}
		k.Grid(cols).Children(func() {
			ui.Column(c).Gap(6).Padding(12, 14).Radius(k.Tokens().Radius).Background(k.Tokens().Subtle.Alpha(0.45)).Border(1, k.Tokens().Native.Border).Children(func() {
				k.Caption("多选清单 (CheckboxList)").FontWeight(600)
				k.Key("checks").CheckboxList(&s.Checks, items, kit.SelectionOptions{Label: "语言清单"})
			})
			ui.Column(c).Gap(6).Padding(12, 14).Radius(k.Tokens().Radius).Background(k.Tokens().Subtle.Alpha(0.45)).Border(1, k.Tokens().Native.Border).Children(func() {
				k.Caption("单选清单 (RadioList)").FontWeight(600)
				k.Key("radios").RadioList(&s.Radio, items, kit.InputOptions{Label: "默认语言"})
			})
		})
		k.Caption("已选择：" + strings.Join(s.Multi, ", "))
	})
	add("custom-selectors", "forms", "搜索建议与自定义选择面板", "建议的稳定 ID 与展示文字分离；自定义面板显式应用。", "Typeahead ComplexSelector", `k.Typeahead(&state.Suggestions, items)`, func(k *kit.UI) {
		k.Key("suggestions").Typeahead(&s.Suggest, []kit.Suggestion{{ID: "go", Label: "Go", Description: "原生应用开发", Group: "语言"}, {ID: "rust", Label: "Rust", Description: "系统开发", Group: "语言"}, {ID: "mygo", Label: "MyGo", Description: "原生 UI 框架", Group: "框架"}}, kit.TypeaheadOptions{InputOptions: kit.InputOptions{Label: "搜索技术", Placeholder: "搜索语言或框架…"}}).FillWidth()
		k.Key("custom-selector").ComplexSelector("选择工作模式", &s.ComplexOpen, s.ComplexValue, func(close func()) {
			k.RadioGroup("工作模式", func() {
				k.Radio(&s.ComplexDraft, "标准", "标准模式")
				k.Radio(&s.ComplexDraft, "专注", "专注模式")
			})
			k.Button("应用选择", kit.ButtonOptions{Variant: kit.Primary}).AlignSelf(ui.Start).OnClick(func() { s.ComplexValue = s.ComplexDraft; close() })
		}).FillWidth()
		k.Caption("建议 ID：" + s.Suggest.Value)
	})
	add("structured-query", "forms", "结构化搜索", "添加、编辑与移除条件，保留只读筛选。", "PowerSearch", `k.PowerSearch(&state.Search, fields)`, func(k *kit.UI) {
		k.Key("structured-search").PowerSearch(&s.Filters, []kit.FilterField{{ID: "name", Label: "名称", Kind: "text"}, {ID: "count", Label: "数量", Kind: "number"}, {ID: "date", Label: "日期", Kind: "date"}, {ID: "status", Label: "状态", Kind: "enum", Options: []kit.Option{{Value: "active", Label: "活跃"}, {Value: "archived", Label: "归档"}}}})
		k.Caption(fmt.Sprintf("文字：%s · 活跃条件 %d 个", s.Filters.Query, len(s.Filters.Filters)))
	})
	add("tokens", "forms", "标签与文本标记", "支持添加、建议、移除与退格操作。", "Token Tokenizer", `k.Tokenizer(&state.Tags, suggestions).Label("标签")`, func(k *kit.UI) {
		k.Inline().Children(func() {
			k.Token("不可移除", nil)
			if s.RemovableToken {
				k.Token("可移除", func() { s.RemovableToken = false })
			} else {
				k.Button("恢复标签", kit.ButtonOptions{Size: kit.Small, Variant: kit.Ghost}).OnClick(func() { s.RemovableToken = true })
			}
		})
		k.Key("token-field").Tokenizer(&s.Tokens, []string{"Go", "原生", "设计", "组件"}, kit.InputOptions{Label: "项目标签"}).FillWidth()
		k.Caption("当前标签：" + strings.Join(s.Tokens, " · "))
	})
	add("file-input", "forms", "文件选择与拖放", "使用真实系统选择器；校验格式、大小、数量，支持移除。", "FileInput", `k.FileInput(&state.Uploads, astryx.FileInputOptions{Accept: ".md,.txt,.png", Multiple: true})`, func(k *kit.UI) {
		k.Key("upload").FileInput(&s.Files, kit.FileInputOptions{Label: "附件", Accept: ".md,.txt,.png", Multiple: true, MaxFiles: 3, MaxSize: 1024 * 1024}).FillWidth()
		k.Inline().Children(func() {
			k.Caption("最多 3 个文件，每个 1 MB。这里不发送文件到服务器。").Grow(1)
			k.Button("切换处理状态", kit.ButtonOptions{Variant: kit.Outline, Size: kit.Small}).OnClick(func() { s.Files.Busy = !s.Files.Busy; s.Files.Progress = .55 })
		})
	})
	add("date-range", "data", "日期范围与日期时间", "支持预设、清空、范围限制和错误提示。", "DateRangeInput DateTimeInput", `k.DateRangeInput(&state.Range, astryx.DateRangeOptions{MaxDays: 31})`, func(k *kit.UI) {
		start := time.Date(2026, 11, 6, 0, 0, 0, 0, time.Local)
		k.Checkbox(&s.ExcludeWeekends, "排除周末")
		var allowed func(time.Time) bool
		if s.ExcludeWeekends {
			allowed = func(d time.Time) bool { return d.Weekday() != time.Saturday && d.Weekday() != time.Sunday }
		}
		k.Key("date-range").DateRangeInput(&s.Range, kit.DateRangeOptions{Label: "排期", MaxDays: 31, Allowed: allowed, Presets: []kit.DateRangePreset{{Label: "一周", Start: start, End: start.AddDate(0, 0, 6)}, {Label: "两周", Start: start, End: start.AddDate(0, 0, 13)}, {Label: "工作日一天", Start: start, End: start}}})
		k.Key("date-time").DateTimeInput(&a.Date, kit.InputOptions{Label: "发布"})
	})
	add("banner-status", "feedback", "通知栏与状态点", "动作、关闭和状态文字保持明确。", "Banner StatusDot", `k.Banner("更新可用", "查看新功能。", options)`, func(k *kit.UI) {
		k.Banner("组件库更新", "新增菜单、聊天与高级输入。", kit.BannerOptions{Tone: kit.Info, ActionLabel: "查看", OnAction: func() { s.MenuAction = "已查看更新" }, Dismissed: &s.BannerDismissed})
		if s.BannerDismissed {
			k.Button("重新显示通知栏", kit.ButtonOptions{Size: kit.Small, Variant: kit.Outline}).AlignSelf(ui.Start).OnClick(func() { s.BannerDismissed = false })
		}
		k.Inline().Children(func() {
			k.StatusDot("在线", kit.Success)
			k.StatusDot("处理中", kit.Warning)
			k.StatusDot("离线", kit.Neutral)
		})
	})
	add("time-controls", "feedback", "时间戳与经过时间", "计时按需更新，暂停后不再持续刷新。", "Timer Timestamp", `k.Timer(&state.Timer, astryx.TimerOptions{Clock: true})`, func(k *kit.UI) {
		k.Key("timer").Timer(&s.Timer, kit.TimerOptions{Clock: true, Label: "经过时间"}).FontSize(28)
		k.Inline().Children(func() {
			k.Button("开始计时", kit.ButtonOptions{Size: kit.Small, Variant: kit.Primary}).OnClick(func() { s.Timer.Start(c.Now()) })
			k.Button("暂停计时", kit.ButtonOptions{Size: kit.Small, Variant: kit.Outline}).OnClick(func() { s.Timer.Pause(c.Now()) })
			k.Button("重置计时", kit.ButtonOptions{Size: kit.Small, Variant: kit.Ghost}).OnClick(s.Timer.Reset)
		})
		k.Timestamp(a.Date, kit.TimestampOptions{Format: "absolute", Copy: true})
	})
	add("menus", "menus", "菜单、子菜单与状态项", "支持禁用项、分组、勾选、单选和键盘导航。", "DropdownMenu MoreMenu ContextMenu NavMenu Toolbar", `k.DropdownMenu("操作", entries)`, func(k *kit.UI) {
		entries := []kit.MenuEntry{{ID: "open", Label: "打开项目", Group: "项目", OnSelect: func() { s.MenuAction = "打开项目" }}, {ID: "blocked", Label: "不可用操作", Disabled: true}, {Separator: true}, {ID: "check", Label: "自动保存", Group: "偏好", Checked: &s.MenuChecked}, {ID: "versions", Label: "发行通道", Children: []kit.MenuEntry{{Label: "稳定版", Radio: &s.MenuRadio, Value: "stable"}, {Label: "预览版", Radio: &s.MenuRadio, Value: "preview"}}}}
		k.Inline().Children(func() {
			k.Key("dropdown").DropdownMenu("项目操作", entries)
			k.Key("more").MoreMenu("更多操作", entries)
			k.Key("nav-menu").NavMenu("工作空间菜单", entries)
			k.Toolbar("编辑工具", func() {
				k.Toggle(&a.ToggleA, "粗体")
				k.Toggle(&a.ToggleB, "斜体")
				k.Button("工具动作", kit.ButtonOptions{Variant: kit.Ghost, Size: kit.Small}).OnClick(func() { s.MenuAction = "工具动作" })
			})
		})
		target := ui.Row(c).Padding(14, 16).Radius(k.Tokens().Radius).Background(k.Tokens().Subtle.Alpha(0.45)).Border(1, k.Tokens().Native.Border).BorderStyle(ui.BorderDashed).Label("右键菜单区域").Children(func() { k.Caption("在此右键，或聚焦后按 Shift+F10 打开上下文菜单") }).Focusable()
		k.ContextMenu(target, entries)
		k.Caption(fmt.Sprintf("动作：%s · 自动保存：%v · %s", s.MenuAction, s.MenuChecked, s.MenuRadio))
	})
	add("command-palette", "menus", "可搜索命令面板", "用方向键、Enter 和 Escape 操作；Ctrl/⌘ K 可全局打开。", "CommandPalette", `k.CommandPalette(&state.Open, &state.Palette, commands)`, func(k *kit.UI) {
		k.Inline().Children(func() {
			k.Button("打开命令面板", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { s.PaletteOpen = true })
			k.Kbd("⌘ K")
		})
		commands := []kit.Command{{ID: "new", Label: "新建原生项目", Description: "创建工作空间", Group: "项目", Shortcut: "⌘ N", Run: func() { s.MenuAction = "新建原生项目" }}, {ID: "disabled", Label: "未开放命令", Disabled: true, Group: "项目"}, {ID: "settings", Label: "打开设置", Description: "主题与偏好", Group: "设置", Run: func() { s.MenuAction = "打开设置" }}}
		k.CommandPalette(&s.PaletteOpen, &s.Palette, commands)
		k.Caption("已执行：" + s.MenuAction)
	})
	add("hover-cards", "overlays", "延迟悬停与内容预览", "支持键盘聚焦、点击切换及延迟隐藏。", "HoverCard InfoTip", `k.HoverCard(anchor, preview, astryx.HoverCardOptions{OpenOnClick: true})`, func(k *kit.UI) {
		k.Inline().Children(func() {
			anchor := k.Button("悬停查看成员", kit.ButtonOptions{Variant: kit.Outline})
			k.HoverCard(anchor, func() {
				ui.Row(c).AlignItems(ui.Start).Gap(12).Children(func() {
					k.Avatar("MyGo", nil, kit.AvatarOptions{Size: 44, Status: "在线", StatusTone: kit.Success})
					ui.Column(c).Gap(4).Grow(1).MinWidth(0).Children(func() {
						ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
							k.Text("MyGo 原生团队").FontWeight(600)
							k.Badge("@mygo-ui")
						})
						k.Caption("构建跨平台、无障碍的 Go 原生桌面组件体系。可将指针移入面板继续操作。")
						k.Button("面板内动作", kit.ButtonOptions{Size: kit.Small, Variant: kit.Primary}).AlignSelf(ui.Start).Margin(4, 0, 0, 0).OnClick(func() { s.MenuAction = "已操作悬停面板" })
					})
				})
			}, kit.HoverCardOptions{OpenOnClick: true, Width: 320, Label: "成员预览"})
			k.InfoTip("关于悬停面板", "支持延迟开启、关闭和辅助描述。")
		})
	})
	add("sheets-overlay", "overlays", "底部面板与自定义覆盖层", "模态内容保留焦点约束和 Escape 关闭。", "BottomSheet Overlay", `k.BottomSheet("选项", &state.Open, content)`, func(k *kit.UI) {
		k.Inline().Children(func() {
			k.Button("打开底部面板", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { s.SheetOpen = true })
			k.Button("打开覆盖层", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { s.OverlayOpen = true })
			k.Button("分步底部面板", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { s.SheetStep = "details"; s.SheetFlowOpen = true })
		})
		if s.SheetFlowOpen {
			title := "填写资料"
			if s.SheetStep == "review" {
				title = "核对资料"
			}
			k.Key("sheet-flow").BottomSheet(title, &s.SheetFlowOpen, func() {
				if s.SheetStep == "details" {
					k.TextInput(&a.Name, kit.InputOptions{Label: "流程项目名称"}).FillWidth()
					k.Button("下一步：核对").AlignSelf(ui.Start).OnClick(func() { s.SheetStep = "review" })
				} else {
					k.Text("项目名称：" + a.Name)
					k.Inline().Children(func() {
						k.Button("返回填写").OnClick(func() { s.SheetStep = "details" })
						k.Button("完成流程", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { s.SheetFlowOpen = false })
					})
				}
			})
		}
		k.BottomSheet("底部选项", &s.SheetOpen, func() { k.Segmented(&a.Segment, []string{"详情", "设置"}); k.Checkbox(&a.Checked, "记住选择") }, kit.DialogOptions{Description: "用于窄窗口的模态选择。"})
		k.Overlay(&s.OverlayOpen, func() {
			k.Center().Fill().Children(func() {
				k.Card().Width(260).Children(func() {
					k.Heading("自定义覆盖层", 3)
					k.Button("关闭覆盖层", kit.ButtonOptions{Variant: kit.Primary}).AlignSelf(ui.Start).OnClick(func() { s.OverlayOpen = false })
				})
			})
		}, kit.OverlayOptions{Label: "自定义覆盖层"})
	})
	add("advanced-navigation", "navigation", "多层导航与目录", "横向、侧栏、溢出菜单和可滚动的 MegaMenu。", "TopNav SideNav OverflowList Outline TopNavMegaMenu TopNavMegaMenuItem", `k.TopNav(&state.Page, navigation)`, func(k *kit.UI) {
		items := []kit.NavItem{{ID: "home", Label: "主页", Icon: kit.IconPlus}, {ID: "work", Label: "工作空间", Children: []kit.NavItem{{ID: "projects", Label: "项目"}, {ID: "settings", Label: "设置"}}}, {ID: "docs", Label: "文档"}}
		k.TopNav(&s.Navigation, items, kit.NavOptions{Label: "顶部导航"})
		cols := 2
		if w, _ := c.Size(); w < 720 {
			cols = 1
		}
		k.Grid(cols).Children(func() {
			k.SideNav(&s.Navigation, items, kit.NavOptions{Label: "侧栏导航", Collapsed: &s.CollapsedNav})
			ui.Column(c).Gap(10).Children(func() {
				k.Caption("章节目录 (Outline)").FontWeight(600)
				k.Outline(&s.Outline, []kit.OutlineItem{{ID: "start", Title: "开始使用", Level: 1}, {ID: "theme", Title: "主题", Level: 2}, {ID: "components", Title: "组件", Level: 1}}, func(id string) { s.MenuAction = "导航到 " + id })
				k.Caption("溢出折叠 (OverflowList)").FontWeight(600).Padding(4, 0, 0, 0)
				k.OverflowList([]kit.NavItem{{ID: "one", Label: "一"}, {ID: "two", Label: "二"}, {ID: "three", Label: "三"}, {ID: "four", Label: "四"}}, 2, func(id string) { s.MenuAction = "选择溢出项 " + id })
			})
		})
		k.Inline().Children(func() {
			k.TopNavMegaMenu("探索产品", func() {
				k.TopNavMegaMenuItem("设计工具", kit.ItemOptions{Description: "组件与主题", OnClick: func() { s.MenuAction = "设计工具" }})
				k.TopNavMegaMenuItem("原生开发", kit.ItemOptions{Description: "Go 与 MyGo", OnClick: func() { s.MenuAction = "原生开发" }})
				if s.LongMenu {
					for i := range 24 {
						label := fmt.Sprintf("扩展示例 %02d", i+1)
						k.TopNavMegaMenuItem(label, kit.ItemOptions{OnClick: func() { s.MenuAction = label }})
					}
				}
			}, func() { k.Heading("从一个组件开始", 4); k.Caption("可交互的特色内容面板。") })
			k.Switch(&s.LongMenu, "显示长菜单内容")
		})
		k.Caption("当前位置：" + s.Navigation)
		k.Caption("导航动作：" + s.MenuAction)
	})
	add("mobile-navigation", "navigation", "移动导航与展开入口", "父项目打开子菜单；展开入口与侧栏共享关闭状态。", "MobileNav MobileNavToggle", `k.MobileNavToggle("导航菜单", &state.Open); k.Drawer("站点导航", &state.Open, content)`, func(k *kit.UI) {
		items := []kit.NavItem{{ID: "home", Label: "移动主页"}, {ID: "work", Label: "移动工作空间", Children: []kit.NavItem{{ID: "locked", Label: "未开放", Disabled: true}, {ID: "projects", Label: "移动项目"}, {ID: "settings", Label: "移动设置"}}}, {ID: "docs", Label: "移动文档"}}
		k.MobileNav(&s.Navigation, items)
		k.MobileNavToggle("打开站点导航", &s.MobileOpen, kit.ButtonOptions{Variant: kit.Outline})
		k.Drawer("移动站点导航", &s.MobileOpen, func() { k.SideNav(&s.Navigation, items).OnChange(func() { s.MobileOpen = false }) })
		k.Caption("移动当前位置：" + s.Navigation)
	})
	add("stepper", "navigation", "步骤进度", "支持当前、完成、错误、禁用和逐步前进。", "Stepper", `k.Stepper(&state.Step, steps, astryx.StepperOptions{Linear: true})`, func(k *kit.UI) {
		k.Stepper(&s.Step, []kit.Step{{ID: "basic", Label: "基础信息", Description: "填写名称"}, {ID: "options", Label: "配置", Description: "选择方案"}, {ID: "review", Label: "确认", Description: "检查并完成"}}, kit.StepperOptions{Linear: true})
		k.Caption(fmt.Sprintf("当前步骤 %d", s.Step+1))
	})
	add("virtual-tree", "data", "虚拟化树形列表", "稳定节点 ID、选择、展开和键盘导航。", "TreeList", `k.TreeList(&state.Tree, roots, children, label)`, func(k *kit.UI) {
		s.Tree.List.Selected = &s.TreeSelected
		k.TreeList(&s.Tree, []string{"root"}, func(id string) []string {
			if id == "root" {
				return []string{"设计", "运行时", "文档"}
			}
			return nil
		}, func(id string) string {
			if id == "root" {
				return "MyGo UI"
			}
			return id
		}).Height(180).Label("虚拟项目树")
	})
	add("media-viewers", "media", "图片与轮播", "本项目自定义 SVG 示例，支持预览、切换、缩放和旋转。", "Thumbnail Carousel Lightbox", `k.Lightbox(&state.Open, &state.Index, images)`, func(k *kit.UI) {
		k.Inline().Children(func() {
			for i, image := range s.Images {
				k.Key(image.ID).Thumbnail(image, kit.ThumbnailOptions{Selected: i == s.ImageIndex, OnClick: func() { s.ImageIndex = i; s.LightboxOpen = true }})
			}
		})
		slides := []kit.CarouselSlide{}
		for _, image := range s.Images {
			slides = append(slides, kit.CarouselSlide{ID: image.ID, Label: image.Label, Content: func() { ui.Image(c, image.Source).Fill().Fit(ui.FillBox).Label(image.Label) }})
		}
		k.Carousel(&s.CarouselIndex, slides, kit.CarouselOptions{Height: 180, Loop: true})
		k.Lightbox(&s.LightboxOpen, &s.ImageIndex, s.Images)
	})
	return append(demos, a.ChatDemos(c)...)
}
