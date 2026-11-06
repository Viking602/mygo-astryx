// Package catalog is the shared index for the Gallery and component docs.
package catalog

import "strings"

type Category struct{ ID, Title, Description string }

var Categories = []Category{
	{"layout", "布局与容器", "从间距、栅格到完整应用结构。"},
	{"typography", "文字与展示", "清晰的排版层级与内容表面。"},
	{"actions", "按钮与动作", "一致的动作层级、状态与键盘交互。"},
	{"forms", "输入与选择", "原生输入法、数据绑定与表单语义。"},
	{"navigation", "导航与组织", "页面切换、层级和渐进展开。"},
	{"feedback", "反馈与状态", "呈现进度、结果与等待状态。"},
	{"overlays", "浮层与对话", "焦点管理、Escape 与嵌套交互。"},
	{"data", "数据与日期", "保留 MyGo 虚拟化与原生编辑能力。"},
	{"menus", "菜单与命令", "菜单状态、子菜单和可搜索命令。"},
	{"media", "图片与媒体", "缩略图、轮播和可操作的图片预览。"},
	{"chat", "聊天与上下文", "原生消息、Markdown、输入和工具状态。"},
}

type Component struct{ Name, Category, Implementation string }

// Components contains public components, not their sizes or color variants.
var Components = build()

func build() []Component {
	var result []Component
	add := func(category, kind, names string) {
		for _, name := range strings.Fields(names) {
			result = append(result, Component{name, category, kind})
		}
	}
	add("layout", "composition", "Container Stack Inline Grid Scroll Spacer Divider AppShell Split")
	add("typography", "composition", "Text Caption Heading Code Icon Avatar Badge Card EmptyState")
	add("actions", "base", "Button IconButton Toggle ButtonGroup")
	add("actions", "adapter", "Link")
	add("forms", "base", "TextInput TextArea InputGroup Checkbox Switch RadioGroup Radio Select Combobox")
	add("forms", "adapter", "NumberInput SearchField Field Form Fieldset Slider RangeSlider Rating")
	add("navigation", "base", "Tabs Segmented Collapsible")
	add("navigation", "composition", "Breadcrumbs Pagination Sidebar Accordion")
	add("feedback", "composition", "Alert Skeleton")
	add("feedback", "adapter", "Progress Spinner")
	add("feedback", "base", "Toast ToastViewport")
	add("overlays", "base", "Dialog AlertDialog Drawer Popover Tooltip")
	add("data", "adapter", "Table List Tree TreeItem Calendar DateInput TimeInput ColorPicker")
	add("layout", "composition", "AspectRatio Center HStack GridSpan Layout Section Resizable")
	add("typography", "composition", "VisuallyHidden Kbd AvatarGroup Blockquote MetadataList Citation Item ClickableCard")
	add("typography", "adapter", "Markdown CodeBlock")
	add("typography", "composition", "AvatarStatusDot AvatarGroupOverflow")
	add("forms", "base", "FormLayout FieldStatus InputClearButton SelectableCard Indicator CheckboxIndicator RadioIndicator CheckIndicator CheckboxList RadioList MultiSelector Typeahead ComplexSelector PowerSearch Token Tokenizer FileInput")
	add("data", "adapter", "DateRangeInput DateTimeInput TreeList")
	add("feedback", "composition", "Banner StatusDot Timer Timestamp")
	add("menus", "adapter", "DropdownMenu MoreMenu ContextMenu NavMenu Toolbar CommandPalette")
	add("overlays", "base", "HoverCard InfoTip BottomSheet Overlay")
	add("navigation", "composition", "TopNav SideNav MobileNav MobileNavToggle OverflowList Outline Stepper TopNavMegaMenu TopNavMegaMenuItem")
	add("media", "composition", "Thumbnail Carousel Lightbox")
	add("chat", "composition", "Chat ChatLayout ChatMessageList ChatComposer ChatComposerInput ChatSendButton ChatLayoutScrollButton ChatMessage ChatMessageBubble ChatMessageMetadata ChatSystemMessage ChatToolCalls ChatComposerTokenElement ChatTokenizedText ChatComposerDrawer ChatDictationButton")
	return result
}

func Count(category string) int {
	n := 0
	for _, c := range Components {
		if c.Category == category {
			n++
		}
	}
	return n
}
func Match(c Component, query string) bool {
	return strings.Contains(strings.ToLower(c.Name+" "+c.Category), strings.ToLower(strings.TrimSpace(query)))
}
