package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	kit "mygo-astryx"
	"strings"
)

// Demos maps every catalog component to a live, interactive example.
func (a *App) Demos(c *ui.Context) []Demo {
	filterPage := c != nil && strings.TrimSpace(a.Query) == ""
	demos := []Demo{
		{"layout-primitives", "layout", "有节奏的布局", "容器、栅格、滚动和间距。", `k.Stack().Gap(16).Children(content)`, []string{"Container", "Stack", "Inline", "Grid", "Scroll", "Spacer", "Divider"}, func(k *kit.UI) {
			k.Container().Padding(0).Children(func() {
				k.Inline().FillWidth().Children(func() { k.Badge("Inline"); k.Spacer(); k.Caption("Spacer 将内容推向另一端") })
				k.Divider()
				k.Grid(3).Children(func() {
					for _, name := range []string{"01 / 对齐", "02 / 间距", "03 / 栅格"} {
						k.Card().Children(func() { k.Text(name) })
					}
				})
				k.Scroll().Height(108).Border(1, k.Tokens().Native.Border).Radius(8).Background(k.Tokens().Subtle.Alpha(0.35)).Children(func() {
					k.Stack().Gap(8).Padding(12).Children(func() {
						for i := 1; i <= 8; i++ {
							k.Caption(fmt.Sprintf("滚动内容 %02d", i))
						}
					})
				})
				k.Scroll(kit.ScrollOptions{Both: true}).Height(110).FillWidth().Border(1, k.Tokens().Native.Border).Radius(8).Background(k.Tokens().Subtle.Alpha(0.35)).Children(func() {
					ui.Box(c).Size(600, 240).Children(func() {
						k.Caption("可同时横向和纵向滚动").Padding(12)
						k.Text("右下角内容").Absolute().End(12).Bottom(12)
					})
				})
			})
		}},
		{"layout-shell", "layout", "应用骨架与分栏", "拖动分栏边界，调整两侧空间。", `k.Split(&state.Position, sidebar, content).Height(180)`, []string{"AppShell", "Split"}, func(k *kit.UI) {
			k.Split(&a.SplitPosition, func() {
				k.Stack().Padding(12).Children(func() { k.Text("资源").Bold(); k.Caption("项目 · 设计 · 文档") })
			}, func() {
				k.Card().Margin(8).Children(func() { k.Text("可调整的内容区域"); k.Caption("Split 保留原生拖动与键盘行为。") })
			}).Height(160).Border(1, k.Tokens().Native.Border).Radius(8)
			ui.Box(c).Height(150).Border(1, k.Tokens().Native.Border).Radius(8).Clip().Children(func() {
				k.AppShell(func() { k.Caption("AppShell 导航") }, func() { k.Text("内容区域").Padding(16) })
			})
		}},
		{"text-hierarchy", "typography", "排版层级", "系统字体、语义标题与可选择的代码。", `k.Heading("清晰的层级", 2)`, []string{"Text", "Caption", "Heading", "Code"}, func(k *kit.UI) {
			k.Heading("让内容自然地被阅读。", 2)
			k.Heading("一个清晰的章节标题", 3)
			k.Text("正文以系统字体呈现，适合长时间使用的原生桌面应用。")
			k.Caption("辅助信息 · 随系统文字缩放")
			k.Code(`k.Button("继续", astryx.ButtonOptions{Variant: astryx.Primary})`).FillWidth()
		}},
		{"display-surfaces", "typography", "内容表面", "状态标记、身份和图标。", `k.Badge("已同步", astryx.Success)`, []string{"Card", "Avatar", "Badge", "Icon"}, func(k *kit.UI) {
			k.Inline().Children(func() {
				k.Avatar("Go Native", nil)
				k.Icon(kit.IconSettings)
				for _, tone := range []kit.Tone{kit.Neutral, kit.Info, kit.Success, kit.Warning, kit.Danger} {
					k.Badge([]string{"默认", "信息", "成功", "警告", "错误"}[tone], tone)
				}
			})
			k.Card(kit.CardOptions{Elevated: true}).Children(func() {
				k.Text("独立的内容区块").Bold()
				k.Caption("边框、表面、圆角和阴影来自同一主题。")
			})
		}},
		{"empty-state", "typography", "空状态", "解释状态，给出可执行的下一步。", `k.EmptyState("还没有项目", "创建第一个项目。", actions)`, []string{"EmptyState"}, func(k *kit.UI) {
			k.EmptyState("还没有项目", "从一个小想法开始，创建第一个原生应用。", func() {
				k.Button("导入配置", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { k.Toast(ui.Toast{Title: "已打开导入向导"}) })
				k.Button("新建项目", kit.ButtonOptions{Variant: kit.Primary, Icon: kit.IconPlus}).OnClick(func() { k.Toast(ui.Toast{Title: "已创建示例项目"}) })
			})
		}},
		{"buttons", "actions", "动作层级", "主动作、次动作与危险动作。", `k.Button("继续", astryx.ButtonOptions{Variant: astryx.Primary})`, []string{"Button", "IconButton"}, func(k *kit.UI) {
			k.Inline().Children(func() {
				k.Button("主要动作", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { a.Clicks++ })
				k.Button("默认动作")
				k.Button("描边按钮", kit.ButtonOptions{Variant: kit.Outline})
				k.Button("轻量按钮", kit.ButtonOptions{Variant: kit.Ghost})
				k.Button("删除", kit.ButtonOptions{Variant: kit.Destructive}).OnClick(func() { a.ConfirmOpen = true })
				k.IconButton("添加", kit.IconPlus, kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { a.Clicks++ })
			})
			k.Inline().Children(func() {
				k.Button("小尺寸", kit.ButtonOptions{Size: kit.Small})
				k.Button("标准尺寸")
				k.Button("大尺寸", kit.ButtonOptions{Size: kit.Large})
				k.Button("不可用", kit.ButtonOptions{Disabled: true})
				k.Button("加载中", kit.ButtonOptions{Variant: kit.Primary, Loading: true})
			})
			k.Caption(fmt.Sprintf("主要动作和添加已触发 %d 次", a.Clicks))
			k.AlertDialog("删除示例项目", "演示操作，不会删除文件。", &a.ConfirmOpen, kit.AlertDialogOptions{Confirm: "确认删除", Destructive: true, OnConfirm: func() { k.Toast(ui.Toast{Title: "示例项目已删除"}) }})
		}},
		{"toggle-links", "actions", "组合动作", "选择状态、组内键盘导航和链接。", `k.ButtonGroup(func() { k.Toggle(&state.Bold, "粗体") })`, []string{"Toggle", "ButtonGroup", "Link"}, func(k *kit.UI) {
			k.Inline().Gap(16).Children(func() {
				k.ButtonGroup(func() {
					k.Key("bold").Toggle(&a.ToggleA, "粗体")
					k.Key("italic").Toggle(&a.ToggleB, "斜体")
					k.Key("underline").Toggle(&a.ToggleC, "下划线")
				}).Label("文字格式")
				preview := k.Text("排版预览 Aa").Padding(6, 12).Radius(k.Tokens().Radius).Background(k.Tokens().Subtle)
				if a.ToggleA {
					preview.Bold()
				}
				if a.ToggleB {
					preview.Italic()
				}
				if a.ToggleC {
					preview.Underline()
				}
				k.Spacer()
				k.Link("阅读 MyGo 文档 ↗", "https://mygo.egoist.dev/docs/ui")
			})
		}},
		{"text-inputs", "forms", "文字输入", "原生输入法、选择、撤销与剪贴板。", `k.Field("名称", func() { k.Key("name").TextInput(&state.Name) })`, []string{"TextInput", "TextArea", "Field"}, func(k *kit.UI) {
			k.Field("名称", func() { k.Key("name").TextInput(&a.Name, kit.InputOptions{Placeholder: "请输入名称"}).FillWidth() }, kit.InputOptions{Description: "试试中文输入法。"})
			k.Field("说明", func() { k.Key("bio").TextArea(&a.Bio).FillWidth() })
			k.Inline().Children(func() {
				k.TextInput(&a.ReadOnlyEmail, kit.InputOptions{Label: "只读邮箱", ReadOnly: true}).Width(220)
				k.TextInput(&a.InvalidEmail, kit.InputOptions{Label: "无效邮箱", Error: "请检查格式"}).Width(220)
			})
		}},
		{"input-composition", "forms", "输入组合", "前后缀、搜索与数字步进。", `k.InputGroup(&state.Query, prefix, suffix, astryx.InputOptions{Label: "地址"})`, []string{"InputGroup", "NumberInput", "SearchField"}, func(k *kit.UI) {
			k.Key("group").InputGroup(&a.GroupText, func() { k.Caption("https://") }, func() { k.Badge("URL") }, kit.InputOptions{Label: "项目地址", Placeholder: "example.com"}).FillWidth()
			k.Key("search-demo").SearchField(&a.Search, kit.InputOptions{Label: "示例搜索", Placeholder: "输入以搜索…"}).FillWidth()
			k.Field("数量", func() { k.NumberInput(&a.Number, 0, 100, 1).Label("数量") })
		}},
		{"choices", "forms", "开关与选项", "鼠标、空格和方向键都能操作。", `k.Checkbox(&state.Accepted, "同意")`, []string{"Checkbox", "Switch", "RadioGroup", "Radio"}, func(k *kit.UI) {
			k.Inline().Children(func() {
				k.Key("check").Checkbox(&a.Checked, "接收通知")
				k.Key("mixed").Checkbox(&a.Mixed, "部分选中", kit.ChoiceOptions{Mixed: true})
				k.Key("switch").Switch(&a.Enabled, "启用同步")
			})
			k.RadioGroup("工作空间类型", func() {
				k.Key("personal").Radio(&a.Radio, "personal", "个人")
				k.Key("team").Radio(&a.Radio, "team", "团队")
				k.Key("enterprise").Radio(&a.Radio, "enterprise", "企业", kit.ChoiceOptions{Disabled: true})
			})
			k.Caption("当前选择：" + a.Radio)
		}},
		{"selects", "forms", "选择与搜索选择", "键盘跳过不可用选项，文字可用于筛选。", `k.Select(&state.Plan, []astryx.Option{{Value: "studio", Label: "工作室"}})`, []string{"Select", "Combobox"}, func(k *kit.UI) {
			k.Field("方案", func() {
				k.Key("plan").Select(&a.Choice, []kit.Option{{Value: "personal", Label: "个人方案"}, {Value: "studio", Label: "工作室方案"}, {Value: "unavailable", Label: "暂未开放", Disabled: true}, {Value: "team", Label: "团队方案"}}, kit.InputOptions{Label: "方案"}).FillWidth()
			})
			k.Field("语言", func() {
				k.Key("language").Combobox(&a.Combo, []string{"Go", "Rust", "Swift", "Kotlin", "TypeScript"}, kit.InputOptions{Label: "语言", Placeholder: "输入以筛选"}).FillWidth()
			})
			k.Caption("方案值：" + a.Choice + " · 语言：" + a.Combo)
		}},
		{"sliders", "forms", "连续与离散的值", "双端范围、键盘调节与星级。", `k.Slider(&state.Volume, 0, 100).Label("音量")`, []string{"Slider", "RangeSlider", "Rating"}, func(k *kit.UI) {
			ui.Column(c).Gap(6).Children(func() {
				ui.Row(c).Justify(ui.SpaceBetween).Children(func() { k.Text("音量").FontWeight(500); k.Caption(fmt.Sprintf("%.0f%%", a.Volume)).FontFeatures("tnum") })
				k.Key("volume").Slider(&a.Volume, 0, 100, kit.InputOptions{Label: "音量"}).FillWidth()
			})
			ui.Column(c).Gap(6).Children(func() {
				ui.Row(c).Justify(ui.SpaceBetween).Children(func() { k.Text("范围").FontWeight(500); k.Caption(fmt.Sprintf("%.0f — %.0f", a.Low, a.High)).FontFeatures("tnum") })
				k.Key("range").RangeSlider(&a.Low, &a.High, 0, 100, 1).Label("范围").FillWidth()
			})
			k.Inline().Gap(10).Children(func() {
				k.Text("评分").FontWeight(500)
				k.Key("rating").Rating(&a.Stars, 5).Label("评分")
			})
		}},
		{"forms-layout", "forms", "有语义的表单", "标签、说明、错误和分组。", `k.Form(func() { k.Fieldset("账户", fields) })`, []string{"Form", "Fieldset"}, func(k *kit.UI) {
			k.Form(func() {
				k.Fieldset("账户", func() {
					k.Field("邮箱", func() { k.Key("form-email").TextInput(&a.FormEmail).FillWidth() })
					k.Field("同步", func() { k.Key("form-sync").Switch(&a.Enabled, "自动同步更改") })
				})
			})
		}},
		{"tabs-segmented", "navigation", "切换视图", "方向键移动选择与焦点。", `k.Tabs(&state.Tab, []string{"概览", "活动", "设置"})`, []string{"Tabs", "Segmented"}, func(k *kit.UI) {
			k.Key("tabs").Tabs(&a.Tab, []string{"概览", "活动", "设置"}).Label("项目视图")
			k.Text([]string{"项目概览内容", "最近的活动记录", "项目设置选项"}[a.Tab]).Padding(4, 0)
			k.Key("segments").Segmented(&a.Segment, []string{"日", "周", "月"}).Label("时间范围")
		}},
		{"page-navigation", "navigation", "层级与页码", "当前位置与可用动作。", `k.Pagination(&state.Page, 12)`, []string{"Breadcrumbs", "Pagination", "Sidebar"}, func(k *kit.UI) {
			k.Breadcrumbs([]kit.Breadcrumb{{Label: "工作空间", OnClick: func() { k.Toast(ui.Toast{Title: "已返回工作空间"}) }}, {Label: "项目", Menu: []kit.MenuEntry{{Label: "设计系统", OnSelect: func() { k.Toast(ui.Toast{Title: "已选择设计系统"}) }}, {Label: "原生应用", OnSelect: func() { k.Toast(ui.Toast{Title: "已选择原生应用"}) }}}}, {Label: "原生组件"}})
			k.Pagination(&a.PageNumber, 12)
			k.Sidebar(&a.Nav, []kit.NavItem{{ID: "projects", Label: "全部项目", Badge: "12"}, {ID: "recent", Label: "最近使用"}, {ID: "archive", Label: "归档"}}).MaxWidth(260)
		}},
		{"disclosures", "navigation", "渐进展开", "只在需要时显示细节。", `k.Collapsible("了解更多", &state.Open, content)`, []string{"Collapsible", "Accordion"}, func(k *kit.UI) {
			k.Key("single-disclosure").Collapsible("这个库如何绘制界面？", &a.OpenA, func() { k.Text("由 MyGo 在原生窗口中绘制，组件通过 Go 函数组合。") })
			k.Accordion([]kit.Disclosure{{ID: "faq", Title: "支持哪些平台？", Open: &a.OpenB, Content: func() { k.Text("MyGo 支持 macOS、Windows 和 Linux。各平台仍需分别验证。") }}})
		}},
		{"alerts", "feedback", "状态与结果", "文字表达状态，不仅依赖颜色。", `k.Alert("已保存", "所有更改已同步。", astryx.Success)`, []string{"Alert"}, func(k *kit.UI) {
			k.Alert("已保存", "所有更改已同步。", kit.Success)
			k.Alert("需要注意", "继续前请检查项目设置。", kit.Warning)
			k.Alert("连接中断", "稍后重试，当前输入不会丢失。", kit.Danger)
		}},
		{"progress", "feedback", "等待与进度", "已知进度、未知等待和静态占位。", `k.Progress(0.64).Label("进度")`, []string{"Progress", "Spinner", "Skeleton"}, func(k *kit.UI) {
			ui.Column(c).Gap(6).Children(func() {
				ui.Row(c).Justify(ui.SpaceBetween).Children(func() { k.Text("上传进度").FontWeight(500); k.Caption("已完成 64%").FontFeatures("tnum") })
				k.Progress(0.64).FillWidth().Label("上传进度")
			})
			k.Inline().Gap(16).Children(func() {
				k.Inline().Children(func() { k.Spinner(); k.Caption("正在处理…") })
				k.Stack().Gap(6).Children(func() { k.Skeleton(220, 12); k.Skeleton(160, 12) })
			})
		}},
		{"toasts", "feedback", "短暂通知", "可关闭、可带动作，原生管理计时。", `k.Toast(ui.Toast{Title: "已保存"}); k.ToastViewport()`, []string{"Toast", "ToastViewport"}, func(k *kit.UI) {
			k.Inline().Children(func() {
				k.Button("显示通知", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { k.Toast(ui.Toast{Title: "设置已保存", Description: "所有更改均已应用。"}) })
				k.Button("带撤销动作", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() {
					k.Toast(ui.Toast{Title: "项目已归档", Action: "撤销", OnAction: func() { k.Toast(ui.Toast{Title: "已撤销归档"}) }})
				})
				k.Button("错误通知", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() {
					k.Toast(ui.Toast{Title: "暂时无法连接", Description: "请检查网络后重试。", Type: "error"})
				})
			})
		}},
		{"dialogs", "overlays", "对话与确认", "试试输入、Tab 和 Escape。", `k.Dialog("编辑项目", &state.Open, content)`, []string{"Dialog", "AlertDialog"}, func(k *kit.UI) {
			k.Switch(&a.Advanced.LongConfirm, "显示长确认说明")
			k.Inline().Children(func() {
				k.Button("打开对话框", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { a.DialogOpen = true })
				k.Button("危险操作确认", kit.ButtonOptions{Variant: kit.Destructive}).OnClick(func() { a.ConfirmOpen = true })
			})
			k.Dialog("编辑项目", &a.DialogOpen, func() {
				k.Field("名称", func() { k.Key("dialog-name").TextInput(&a.Name).FillWidth() })
				k.Field("方案", func() {
					k.Key("dialog-select").Select(&a.Choice, []kit.Option{{Value: "personal", Label: "个人方案"}, {Value: "studio", Label: "工作室方案"}, {Value: "team", Label: "团队方案"}}).FillWidth()
				})
			}, kit.DialogOptions{Description: "包含原生输入和嵌套选择弹层。", Footer: func() {
				k.Button("取消", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { a.DialogOpen = false })
				k.Button("保存更改", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { a.DialogOpen = false; k.Toast(ui.Toast{Title: "更改已保存"}) })
			}})
			description := "演示操作，不会修改文件。"
			if a.Advanced.LongConfirm {
				description = strings.Repeat("这是一段可以滚动阅读的确认说明，操作按钮始终保留在下方。\n", 30)
			}
			k.AlertDialog("删除项目", description, &a.ConfirmOpen, kit.AlertDialogOptions{Confirm: "删除项目", Destructive: true, OnConfirm: func() { k.Toast(ui.Toast{Title: "已完成演示操作"}) }})
		}},
		{"drawers", "overlays", "侧边面板", "容纳设置，将焦点留在面板内。", `k.Drawer("筛选条件", &state.Open, content)`, []string{"Drawer"}, func(k *kit.UI) {
			k.Button("打开侧边面板", kit.ButtonOptions{Variant: kit.Outline}).AlignSelf(ui.Start).OnClick(func() { a.DrawerOpen = true })
			k.Drawer("筛选条件", &a.DrawerOpen, func() {
				k.Field("关键词", func() {
					k.TextInput(&a.Search, kit.InputOptions{Label: "项目关键词", Placeholder: "输入项目名称或标签…"}).FillWidth()
				})
				k.Checkbox(&a.Checked, "仅显示活跃项目")
				k.Switch(&a.Enabled, "包含共享项目")
			}, kit.DialogOptions{Width: 360, Description: "调整当前视图的项目筛选范围。", Footer: func() {
				k.Button("取消", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { a.DrawerOpen = false })
				k.Button("应用筛选", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { a.DrawerOpen = false; k.Toast(ui.Toast{Title: "筛选条件已更新"}) })
			}})
		}},
		{"popovers", "overlays", "就近呈现信息", "依附触发器的弹层与辅助提示。", `k.Popover(anchor, &state.Open, content)`, []string{"Popover", "Tooltip"}, func(k *kit.UI) {
			k.Inline().Children(func() {
				anchor := k.Button("项目详情", kit.ButtonOptions{Variant: kit.Outline})
				if anchor.Clicked() {
					a.PopoverOpen = !a.PopoverOpen
				}
				k.Popover(anchor, &a.PopoverOpen, func() {
					k.Text("MyGo 原生项目").FontWeight(600)
					k.Caption("Go · 自动主题 · 可组合组件")
					k.Divider()
					k.Switch(&a.Enabled, "实时状态同步")
				})
				tip := k.IconButton("显示说明", kit.IconSettings, kit.ButtonOptions{Variant: kit.Ghost})
				k.Tooltip(tip, "悬停或通过键盘聚焦时显示提示。")
			})
		}},
		{"collections", "data", "虚拟化集合", "选中后，使用方向键继续导航。", `k.Table(&state.Table, columns, len(rows), cell).Height(200)`, []string{"Table", "List"}, func(k *kit.UI) {
			a.TableState.Selected = &a.SelectedRow
			rows := [][]string{{"MyGo UI", "进行中", "Go"}, {"Gallery", "可预览", "原生"}, {"设计规范", "审核中", "主题"}, {"交互测试", "已接入", "Tester"}}
			k.Table(&a.TableState, []ui.TableColumn{{Title: "项目"}, {Title: "状态", Width: 100}, {Title: "类型", Width: 80}}, len(rows), func(row, col int) { k.Text(rows[row][col]) }).Height(180).Label("示例项目表格")
			k.List(&a.ListState, 20, func(i int) { k.Text(fmt.Sprintf("列表项 %02d", i+1)).Padding(7, 12) }).Height(120).Border(1, k.Tokens().Native.Border).Label("示例列表")
		}},
		{"tree", "data", "树形层级", "原生树导航与展开能力。", `k.Tree(func() { k.TreeItem("项目", &state.Open, children) })`, []string{"Tree", "TreeItem"}, func(k *kit.UI) {
			k.Tree(func() {
				k.TreeItem("mygo-astryx", &a.TreeOpen, func() {
					k.TreeItem("theme", nil, nil)
					k.TreeItem("components", nil, nil)
					k.TreeItem("gallery", nil, nil)
				})
			}).MinHeight(120).Label("项目结构")
		}},
		{"date-time", "data", "日期与时间", "语言、格式与编辑由 MyGo 提供。", `k.DateInput(&state.Date); k.TimeInput(&state.Date)`, []string{"Calendar", "DateInput", "TimeInput"}, func(k *kit.UI) {
			k.Inline().Children(func() { k.DateInput(&a.Date).Label("日期"); k.TimeInput(&a.Date).Label("时间") })
			k.Calendar(&a.Date).Label("日历")
		}},
		{"colors", "data", "颜色选择", "高阶控件继承统一主题。", `k.ColorPicker(&state.Color).Label("强调色")`, []string{"ColorPicker"}, func(k *kit.UI) {
			k.Inline().AlignItems(ui.Start).Gap(20).Children(func() {
				k.ColorPicker(&a.Color).Label("示例颜色")
				ui.Row(c).AlignItems(ui.Center).Gap(12).Padding(10, 14).Radius(k.Tokens().Radius).Background(k.Tokens().Subtle).Border(1, k.Tokens().Native.Border).Children(func() {
					ui.Box(c).Size(36, 36).Radius(k.Tokens().RadiusSmall).Background(a.Color).Border(1, k.Tokens().Native.Border)
					ui.Column(c).Gap(2).Children(func() {
						k.Text("当前强调色").FontWeight(600)
						k.Caption(fmt.Sprintf("RGB(%d, %d, %d)", a.Color.R, a.Color.G, a.Color.B)).Font("monospace")
					})
				})
			})
		}},
	}
	if filterPage {
		filtered := demos[:0]
		for _, d := range demos {
			if d.Category == a.Page {
				filtered = append(filtered, d)
			}
		}
		demos = filtered
	}
	return append(demos, a.AdvancedDemos(c)...)
}
