#!/usr/bin/env python3
"""Rebuild the pinned Astryx export audit from a local upstream checkout.

Usage: python3 scripts/update_astryx_manifest.py /path/to/astryx
Does not download dependencies or copy upstream implementations.
"""
import json
from collections import Counter
from pathlib import Path
import re
import subprocess
import sys

root = Path(__file__).resolve().parents[1]
upstream = Path(sys.argv[1]).resolve()
revision = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=upstream, text=True).strip()
source = upstream / "packages/core/src"
entry = (source / "index.ts").read_text()
excluded = {"hooks", "utils", "theme", "i18n", "InteractiveRoleContext", "SizeContext"}
families = sorted(set(re.findall(r"export \* from './([^']+)';", entry)) - excluded | {"Toast"})
native = set(re.findall(r"func \(\w+ \*UI\) ([A-Z]\w*)\(", "\n".join(p.read_text() for p in root.glob("*.go") if not p.name.endswith("_test.go"))))
aliases = {"CheckboxInput": ["Checkbox"], "NavIcon": ["Icon"], "ProgressBar": ["Progress"], "ScrollableArea": ["Scroll"], "SegmentedControl": ["Segmented"], "Selector": ["Select"], "TabList": ["Tabs"], "ToggleButton": ["Toggle"], "Indicator": ["CheckboxIndicator", "RadioIndicator", "CheckIndicator"]}
extras = {
    "Chat": ["ChatLayout", "ChatMessageList", "ChatComposer", "ChatComposerInput", "ChatSendButton", "ChatLayoutScrollButton", "ChatMessage", "ChatMessageBubble", "ChatMessageMetadata", "ChatSystemMessage", "ChatToolCalls", "ChatComposerTokenElement", "ChatTokenizedText", "ChatComposerDrawer", "ChatDictationButton"],
    "Layout": ["HStack", "Stack", "Section", "Card"], "Grid": ["GridSpan"], "Field": ["FieldStatus", "InputClearButton"], "CodeBlock": ["Code"], "TopNav": ["TopNavMegaMenu", "TopNavMegaMenuItem"],
    "Avatar": ["AvatarStatusDot"], "AvatarGroup": ["AvatarGroupOverflow"], "MobileNav": ["MobileNavToggle"],
}
notes = {
    "Chat": "消息保留文本与上下文，支持反馈、工具详情和输入；发送/停止/听写由宿主回调提供。Token 独立展示，不实现 contenteditable 光标内嵌标记；听写没有音量频段动画。",
    "Avatar": "AvatarOptions 支持尺寸、状态与动作。AvatarStatusDot 用实心/圆环/横杠区分状态，状态合并到头像的可访问名称。",
    "AvatarGroup": "AvatarGroupOverflow 支持真实隐藏人数、自定义文字、可选按钮动作与禁用状态。",
    "MobileNav": "父项目打开子菜单，叶子项目回写选择；MobileNavToggle 与 Drawer/BottomSheet 共用宿主 bool 状态。",
    "BottomSheet": "单个受控面板支持滚动正文和固定动作区；分步流程在同一面板中替换内容。弃用的 BottomSheetSwitcher 不作为 Header/Footer 别名。",
    "ScrollableArea": "ScrollOptions 支持纵向、横向或双轴 MyGo 滚动。DOM stickyContainment、overscroll 与 React hook 不逐项移植。",
    "Resizable": "MyGo Split 提供拖动和键盘分隔线；附加最小/最大、百分比、纵向和折叠模型。ResizeHandle 是内置分隔线。",
    "Table": "MyGo 虚拟表格；列、单元格与 ListState 表示行/头/体/选择。Astryx React 插件和 DOM 表格部件没有独立 Go 对象。",
    "List": "MyGo 虚拟列表，应用管理 ListState、稳定 Key 和行内容。",
    "TreeList": "MyGo OutlineState 与原生虚拟树；应用提供稳定节点 ID、子节点与内容。",
    "Layout": "Header/Content/Panel/Footer 是 LayoutSlots 回调；VStack 对应 Stack，StackItem 由元素布局方法表示。",
    "Grid": "GridSpan 设置 MyGo 栅格元素的跨列；可继续配置原生 Grid。",
    "Field": "Field(label, content) 表示标签和内容，FieldStatus 与 InputClearButton 独立展示。",
    "DropdownMenu": "条目、分组标题、分隔、勾选、单选、子菜单由 MenuEntry 统一建模；分组标题使用禁用菜单行。",
    "ContextMenu": "共享 MenuEntry；可选择 MyGo 绘制菜单或系统原生上下文菜单。",
    "Breadcrumbs": "Breadcrumb.Menu 支持展开下级菜单，共享 MenuEntry。",
    "CommandPalette": "Command 模型提供分组、筛选、禁用、方向键和 Enter；宿主注册全局快捷键。",
    "PowerSearch": "文字搜索和可增删改的类型化筛选；支持数字、日期、枚举、只读和自定义校验/编辑器。",
    "FileInput": "系统文件选择器与真实拖放；按扩展名/MIME 元数据、大小和数量校验。上传传输与进度由宿主负责；不做内容嗅探。",
    "DateRangeInput": "两端原生日历输入，预设、清空、范围与日期规则校验；不是双月并排日历的逐像素复刻。",
    "Calendar": "MyGo 原生日历；内部月标题、日期格使用 MyGo 实现，不暴露 React 部件。",
    "TopNav": "MegaMenu 长内容内部滚动、键盘进入与焦点恢复；TopNavMegaMenuItem 执行动作并关闭，窄窗改为 Drawer。",
    "SideNav": "分组、层级展开、选择、徽标与收起模式；布局交给宿主组合。",
    "OverflowList": "宿主传入可见项数量，其余进入菜单；不读取 DOM 自动测量项宽。",
    "Markdown": "MyGo 原生 Markdown 渲染及流式入口，不使用 WebView；渲染能力遵循 MyGo markdown.Options。",
    "CodeBlock": "MyGo Markdown fenced code 提供语言高亮与复制；Code 是可选择的等宽文字。",
    "Indicator": "装饰性 Checkbox/Radio/Check 视觉；主题 Indicator 回调可替换，父控件保留语义和输入。",
    "Overlay": "原生叠层，支持模态、Escape 和可选背景关闭；Scope 与主题取代浏览器 LayerProvider。",
    "Carousel": "通用 Content 回调、方向键和循环切换；未实现触摸惯性滑动。",
    "Lightbox": "原生图片预览、缩放、旋转、切换与滚动；不内置视频播放器。",
    "Tokenizer": "MyGo 原生 TokenField 提供建议、添加、删除与退格。",
    "Toast": "MyGo 通知队列，Toast 返回 ID，ToastViewport 呈现；动作、超时与焦点由原生 Base 管理。",
}

def exports(path):
    text = re.sub(r"/\*.*?\*/|//[^\n]*", "", path.read_text(), flags=re.S)
    result = set()
    for block in re.findall(r"export\s*\{([^}]+)\}\s*from", text, re.S):
        for item in block.split(","):
            item = item.strip()
            if item.startswith("type ") or not item:
                continue
            name = item.split(" as ")[-1].strip()
            if re.fullmatch(r"[A-Z][A-Za-z0-9_]+", name) and not name.isupper():
                result.add(name)
    return sorted(result)

docs = {p.name.removesuffix(".doc.mjs"): p for p in source.rglob("*.doc.mjs")}
code_files = {p.stem: p for suffix in ("*.tsx", "*.ts") for p in source.rglob(suffix)}

def part_record(family, part):
    # Every non-method export needs an explicit disposition. New upstream
    # exports must fail here instead of silently inheriting a generic claim.
    def mapped(status, methods, description):
        assert all(method in native for method in methods), (part, methods)
        path = docs.get(part) or code_files.get(part) or source / family / "index.ts"
        return {"export": part, "status": status, "native": methods,
                "representation": description, "source": str(path.relative_to(upstream)),
                "documented": part in docs}
    if part in native:
        return mapped("native-api", [part], part)
    renamed = {key: value[0] for key, value in aliases.items() if key != "Indicator"}
    renamed["VStack"] = "Stack"
    if part in renamed:
        return mapped("native-equivalent", [renamed[part]], f"原生命名为 {renamed[part]}，不重复新增别名方法。")
    contexts = {
        "AppShellMobileContext": (["AppShell", "MobileNavToggle"], "视口布局与宿主展开状态；无 React Context 对象。"),
        "DropdownMenuContext": (["DropdownMenu"], "MenuEntry 共享数据与 MyGo 菜单关闭/焦点；无自定义 React 菜单项 hook。"),
        "FormLayoutContext": (["FormLayout", "GridSpan"], "显式列数和跨列；没有 DOM 字段测量上下文。"),
        "LayoutAreaContext": (["Layout"], "LayoutSlots 显式指定 Header/Content/Panel/Footer。"),
        "LayoutDividerContext": (["Layout", "Divider"], "宿主明确组合边界与分隔线。"),
        "LinkProvider": (["Link"], "调用宿主 URL 服务；没有 React Router Provider 的自动链接替换。"),
        "NavHeadingCloseContext": (["NavMenu"], "原生菜单自动关闭并恢复触发器焦点。"),
        "NavHeadingMenuContext": (["NavMenu"], "MenuEntry 与原生菜单状态。"),
        "RadioListContext": (["RadioList", "RadioGroup"], "共享 string 模型与原生单选分组。"),
        "SideNavRenderContext": (["SideNav"], "NavItem 与 NavOptions 显式控制层级、选择和收起。"),
        "TableContext": (["Table"], "ListState、列定义和单元格回调，不暴露 DOM 表格上下文。"),
        "TopNavRenderContext": (["TopNav", "TopNavMegaMenuItem"], "NavItem、内容回调和当前 builder 的 MegaMenu 关闭动作。"),
    }
    if part in contexts:
        methods, desc = contexts[part]
        return mapped("platform-adaptation", methods, desc)
    if part == "BottomSheetSwitcher":
        assert "deprecated:" in docs[part].read_text(), "upstream deprecation changed"
        return mapped("deprecated-migration", ["BottomSheet"], "上游已弃用；Gallery 用一个受控 BottomSheet 按步骤替换内容，不映射为 Header/Footer。")
    if family in {"DropdownMenu", "ContextMenu", "Breadcrumbs"} and part.startswith(("DropdownMenu", "ContextMenu", "BreadcrumbMenu")):
        properties = {"CheckboxItem": "Checked", "RadioGroup": "共享 Radio 指针", "RadioItem": "Radio + Value", "Divider": "Separator", "Group": "Group（禁用的标题行）", "SubMenu": "Children", "Item": "Label/Disabled/OnSelect"}
        suffix = next((suffix for suffix in properties if part.endswith(suffix)), None)
        if suffix:
            return mapped("native-model", [family], "MenuEntry." + properties[suffix])
    model_parts = {
        "BreadcrumbItem": (["Breadcrumbs"], "Breadcrumb.Label/OnClick/Menu。"),
        "CheckboxListItem": (["CheckboxList"], "Option 与 SelectionOptions。"),
        "CollapsibleGroup": (["Accordion", "Collapsible"], "Disclosure 列表与各自的 Open 指针；允许多项展开。"),
        "FieldLabel": (["Field"], "Field(label, content) 的原生标签与输入关联。"),
        "InputGroupText": (["InputGroup", "Text"], "leading/trailing 回调中的 Text/Caption。"),
        "StackItem": (["Stack", "HStack"], "子元素的 Grow/Shrink/Align 等布局属性。"),
        "ListItem": (["List"], "虚拟列表的行内容回调；选择由 ListState 管理。"),
        "MetadataListItem": (["MetadataList"], "Metadata{Label,Value}。"),
        "NavHeadingMenu": (["NavMenu"], "标题位置由宿主组合 NavMenu；条目使用 MenuEntry。"),
        "NavHeadingMenuItem": (["NavMenu"], "MenuEntry 条目。"),
        "PowerSearchFilterEditor": (["PowerSearch"], "FilterField.Editor/Validate 与类型化内建编辑器。"),
        "PowerSearchToken": (["PowerSearch", "Token"], "SearchFilter 对应的可编辑/只读 Token。"),
        "RadioListItem": (["RadioList"], "Option + 共享 string 模型。"),
        "ResizeHandle": (["Resizable"], "MyGo Split 内建的可拖动/键盘分隔线，无独立 Go 部件。"),
        "SegmentedControlItem": (["Segmented"], "labels 中的索引项。"),
        "SelectorOption": (["Select"], "Option{Value,Label,Disabled}。"),
        "SideNavCollapseButton": (["SideNav"], "NavOptions.Collapsed 绑定的收起按钮。"),
        "SideNavHeading": (["SideNav", "Heading", "NavMenu"], "宿主在导航头部组合 Heading/Avatar/NavMenu。"),
        "SideNavItem": (["SideNav"], "NavItem，支持 Children/Disabled/Icon/Badge。"),
        "SideNavSection": (["SideNav"], "NavItem.Section/Children。"),
        "Step": (["Stepper"], "Step 的 Label/Description/Disabled/Error。"),
        "Tab": (["Tabs"], "labels 的索引项；宿主根据选中索引构建内容。"),
        "ToggleButtonGroup": (["ButtonGroup", "Toggle"], "ButtonGroup 的 ToggleGroup 分组及受控 Toggle。"),
        "TopNavHeading": (["TopNav", "Heading", "Avatar"], "宿主组合品牌标题与头像区域。"),
        "TopNavItem": (["TopNav"], "NavItem 数据。"),
        "TopNavMenu": (["TopNav", "NavMenu"], "带 Children 的 NavItem 使用原生菜单。"),
        "TopNavMegaMenuFeaturedCard": (["TopNavMegaMenu", "Card", "TopNavMegaMenuItem"], "featured 插槽中的 Card 与关闭菜单的操作条目。"),
        "BaseTypeahead": (["Typeahead"], "原生 MyGo ComboboxBase 管线；可直接使用公开 MyGo Base 定制。"),
        "TypeaheadItem": (["Typeahead"], "Suggestion{ID,Label,Description,Group,Disabled}。"),
    }
    if part in model_parts:
        methods, desc = model_parts[part]
        return mapped("native-model", methods, desc)
    if family == "CommandPalette" and part in {"CommandPaletteEmpty", "CommandPaletteFooter", "CommandPaletteGroup", "CommandPaletteInput", "CommandPaletteItem", "CommandPaletteList"}:
        return mapped("native-model", ["CommandPalette"], "Command / CommandPaletteOptions：内建搜索、分组、滚动列表与空提示；Footer 可由回调定制。")
    slots = {"DialogHeader": "Dialog", "DrawerHeader": "Drawer", "LayoutHeader": "Layout", "LayoutFooter": "Layout", "LayoutContent": "Layout", "LayoutPanel": "Layout"}
    if part in slots:
        return mapped("native-slot", [slots[part]], "DialogOptions.Header 或 LayoutSlots 对应的内容回调。")
    if family == "Table" and part in {"TableHeader", "TableHeaderCell", "TableBody", "TableRow", "TableCell"}:
        return mapped("native-model", ["Table"], "MyGo 内建表头、列定义和行/单元格回调；没有独立 DOM 部件对象。")
    if part == "TableSelectionToolbar":
        return mapped("native-composition", ["Table", "Toolbar"], "宿主根据 ListState 的选择构建 Toolbar。")
    partial = {
        "TableFooter": (["Table", "Stack"], "可在表格之后组合内容；尚无与虚拟列自动对齐的内建固定页脚。"),
        "TabMenu": (["Tabs", "DropdownMenu"], "可在标签栏旁组合菜单；尚无绑定单个 Tab 的内建操作菜单。"),
    }
    if part in partial:
        methods, desc = partial[part]
        return mapped("partial", methods, desc)
    raise RuntimeError(f"未审计的公开部件：{family}/{part}")

records = []
for family in families:
    methods = aliases.get(family, [family]) + extras.get(family, [])
    assert all(m in native for m in methods), (family, methods)
    parts = exports(source / family / "index.ts")
    records.append({"family": family, "native": methods, "mapping": "native-equivalent" if family in aliases else "native-api", "parts": [part_record(family, part) for part in parts], "notes": notes.get(family, "采用原生 Go API 和 MyGo 元素组合，样式跟随局部主题；不承诺 React props 完全一致。")})
manifest = {
    "source": {"repository": "https://github.com/meta/astryx", "revision": revision, "entry": "packages/core/src/index.ts"},
    "scope": "根入口的 UI 组件族（export * 排除 hooks/utils/theme/i18n/两个 Context，加命名导出的 Toast）。部件、别名、Context 单列，不重复计算为独立组件族。",
    "excluded": sorted(excluded | {"BaseProps", "useLayer", "LayerProvider", "createAnimator", "createPanGesture", "useToast"}),
    "provider_equivalents": {"LayerProvider": "Scope + Overlay/Popover/Dialog/ToastViewport", "SizeContext": "Options.Size", "InteractiveRoleContext": "原生 Element.Role", "theme/i18n": "局部 theme.Theme 与宿主文案/日期区域设置"},
    "families": records,
    "part_status_counts": dict(sorted(Counter(part["status"] for family in records for part in family["parts"]).items())),
}
output = root / "internal/catalog/astryx.json"
output.write_text(json.dumps(manifest,ensure_ascii=False,indent=2)+"\n")
doc = ["# Astryx 原生覆盖清单", "", f"对照来源固定为 [`meta/astryx@{revision[:12]}`](https://github.com/meta/astryx/tree/{revision}) 的 `packages/core/src/index.ts`。", "", f"根入口共有 **{len(records)} 个 UI 组件族**。本库为每个组件族提供原生入口或明确的原生等价组合。Gallery 的 **{len(native)-2} 个公开入口**还包括组合部件及 MyGo 扩展；两种数字的统计范围不同。尺寸、配色、React hook、Provider、Context 与同一组件的别名不重复计数。", "", "这份清单说明组件族覆盖，不代表 React props、DOM 行为、全部插件或每个复合部件都具有同名 Go API。具体部件的表示方法、源码依据、状态和差异保存在 [机器可读清单](../internal/catalog/astryx.json)。Gallery 可按 Astryx 名称搜索，并展开每组演示的源码对照。", "", "## 部件核对方式", "", "每个导出部件都明确登记为原生方法、等价命名、模型、插槽、组合、平台适配、弃用迁移或部分支持。未知导出会让生成脚本报错，禁止使用泛化的父组件说明自动视为已覆盖。以下数量按导出记录计，包含 Layout 等重导出，不是独立控件数量。", "", "| 状态 | 导出记录 |", "| --- | ---: |"]
for status,count in manifest["part_status_counts"].items():
    doc.append(f"| `{status}` | {count} |")
doc += ["", "## 明确保留的差异", "", "| 部件 | 处理方式 |", "| --- | --- |"]
for family in records:
    for part in family["parts"]:
        if part["status"] in {"partial", "deprecated-migration"}:
            doc.append(f"| {part['export']} | {part['representation']} |")
doc += ["", "## 组件族映射", "", "| Astryx 组件族 | MyGo 原生入口 |", "| --- | --- |"]
for record in records:
    doc.append("| " + record["family"] + " | " + "、".join(f"`{m}`" for m in record["native"]) + " |")
doc += ["", "## 原生转换边界", "", "- 菜单子项、分隔线、分组、勾选、单选与子菜单统一为 `MenuEntry`；Layout、Dialog 等部件使用插槽回调；表格内部结构由 MyGo 管理。", "- `CheckboxInput → Checkbox`、`Selector → Select`、`TabList → Tabs`、`ToggleButton → Toggle`、`ProgressBar → Progress` 等采用 Go 命名。`VStack` 使用 `Stack`。", "- `LayerProvider` 使用 `Scope`、`Overlay`、原生浮层和 `ToastViewport`；React hook 与 DOM 上下文不直接移植。", "- 文件传输、AI 请求、流式数据和系统听写由宿主连接服务；Gallery 使用本地模型，不产生 AI 或上传请求。文件格式校验根据扩展名/MIME 元数据，不做内容嗅探。", "- 日期范围采用两端原生日期输入；Carousel 没有触摸惯性滑动，Lightbox 不内置视频播放器。OverflowList 由宿主传入可见数量。", "- ScrollableArea 采用 MyGo Scroll，保留可配置原生滚动；DOM stickyContainment、overscroll、浏览器 contenteditable 内嵌 token 与 Astryx Table React 插件不属于逐项兼容范围。", "", "## 更新来源", "", "```sh", "python3 scripts/update_astryx_manifest.py /path/to/astryx", "go test ./internal/catalog ./internal/gallery", "```", "", "脚本只读取指定的上游 checkout，登记导出名和映射，不复制上游组件实现。目录测试检查所有组件族映射到真实且有 Gallery 演示的公开 API。", ""]
(root / "docs/astryx-coverage.md").write_text("\n".join(doc))
print(f"Pinned {len(records)} families, {sum(len(r['parts']) for r in records)} exported names; native catalog has {len(native)-2} component methods")
