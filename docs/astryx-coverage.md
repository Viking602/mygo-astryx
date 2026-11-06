# Astryx 原生覆盖清单

对照来源固定为 [`meta/astryx@e8ecde4ff82d`](https://github.com/meta/astryx/tree/e8ecde4ff82d424d76de8d6d8b122421a184c4a6) 的 `packages/core/src/index.ts`。

根入口共有 **99 个 UI 组件族**。本库为每个组件族提供原生入口或明确的原生等价组合。Gallery 的 **147 个公开入口**还包括组合部件及 MyGo 扩展；两种数字的统计范围不同。尺寸、配色、React hook、Provider、Context 与同一组件的别名不重复计数。

这份清单说明组件族覆盖，不代表 React props、DOM 行为、全部插件或每个复合部件都具有同名 Go API。具体部件的表示方法、源码依据、状态和差异保存在 [机器可读清单](../internal/catalog/astryx.json)。Gallery 可按 Astryx 名称搜索，并展开每组演示的源码对照。

## 部件核对方式

每个导出部件都明确登记为原生方法、等价命名、模型、插槽、组合、平台适配、弃用迁移或部分支持。未知导出会让生成脚本报错，禁止使用泛化的父组件说明自动视为已覆盖。以下数量按导出记录计，包含 Layout 等重导出，不是独立控件数量。

| 状态 | 导出记录 |
| --- | ---: |
| `deprecated-migration` | 1 |
| `native-api` | 121 |
| `native-composition` | 1 |
| `native-equivalent` | 10 |
| `native-model` | 62 |
| `native-slot` | 6 |
| `partial` | 2 |
| `platform-adaptation` | 12 |

## 明确保留的差异

| 部件 | 处理方式 |
| --- | --- |
| BottomSheetSwitcher | 上游已弃用；Gallery 用一个受控 BottomSheet 按步骤替换内容，不映射为 Header/Footer。 |
| TabMenu | 可在标签栏旁组合菜单；尚无绑定单个 Tab 的内建操作菜单。 |
| TableFooter | 可在表格之后组合内容；尚无与虚拟列自动对齐的内建固定页脚。 |

## 组件族映射

| Astryx 组件族 | MyGo 原生入口 |
| --- | --- |
| AlertDialog | `AlertDialog` |
| AppShell | `AppShell` |
| AspectRatio | `AspectRatio` |
| Avatar | `Avatar`、`AvatarStatusDot` |
| AvatarGroup | `AvatarGroup`、`AvatarGroupOverflow` |
| Badge | `Badge` |
| Banner | `Banner` |
| Blockquote | `Blockquote` |
| BottomSheet | `BottomSheet` |
| Breadcrumbs | `Breadcrumbs` |
| Button | `Button` |
| ButtonGroup | `ButtonGroup` |
| Calendar | `Calendar` |
| Card | `Card` |
| Carousel | `Carousel` |
| Center | `Center` |
| Chat | `Chat`、`ChatLayout`、`ChatMessageList`、`ChatComposer`、`ChatComposerInput`、`ChatSendButton`、`ChatLayoutScrollButton`、`ChatMessage`、`ChatMessageBubble`、`ChatMessageMetadata`、`ChatSystemMessage`、`ChatToolCalls`、`ChatComposerTokenElement`、`ChatTokenizedText`、`ChatComposerDrawer`、`ChatDictationButton` |
| CheckboxInput | `Checkbox` |
| CheckboxList | `CheckboxList` |
| Citation | `Citation` |
| ClickableCard | `ClickableCard` |
| CodeBlock | `CodeBlock`、`Code` |
| Collapsible | `Collapsible` |
| CommandPalette | `CommandPalette` |
| ComplexSelector | `ComplexSelector` |
| ContextMenu | `ContextMenu` |
| DateInput | `DateInput` |
| DateRangeInput | `DateRangeInput` |
| DateTimeInput | `DateTimeInput` |
| Dialog | `Dialog` |
| Divider | `Divider` |
| Drawer | `Drawer` |
| DropdownMenu | `DropdownMenu` |
| EmptyState | `EmptyState` |
| Field | `Field`、`FieldStatus`、`InputClearButton` |
| FileInput | `FileInput` |
| FormLayout | `FormLayout` |
| Grid | `Grid`、`GridSpan` |
| HoverCard | `HoverCard` |
| Icon | `Icon` |
| IconButton | `IconButton` |
| Indicator | `CheckboxIndicator`、`RadioIndicator`、`CheckIndicator` |
| InfoTip | `InfoTip` |
| InputGroup | `InputGroup` |
| Item | `Item` |
| Kbd | `Kbd` |
| Layout | `Layout`、`HStack`、`Stack`、`Section`、`Card` |
| Lightbox | `Lightbox` |
| Link | `Link` |
| List | `List` |
| Markdown | `Markdown` |
| MetadataList | `MetadataList` |
| MobileNav | `MobileNav`、`MobileNavToggle` |
| MoreMenu | `MoreMenu` |
| MultiSelector | `MultiSelector` |
| NavIcon | `Icon` |
| NavMenu | `NavMenu` |
| NumberInput | `NumberInput` |
| Outline | `Outline` |
| OverflowList | `OverflowList` |
| Overlay | `Overlay` |
| Pagination | `Pagination` |
| Popover | `Popover` |
| PowerSearch | `PowerSearch` |
| ProgressBar | `Progress` |
| RadioList | `RadioList` |
| Resizable | `Resizable` |
| ScrollableArea | `Scroll` |
| Section | `Section` |
| SegmentedControl | `Segmented` |
| SelectableCard | `SelectableCard` |
| Selector | `Select` |
| SideNav | `SideNav` |
| Skeleton | `Skeleton` |
| Slider | `Slider` |
| Spinner | `Spinner` |
| Stack | `Stack` |
| StatusDot | `StatusDot` |
| Stepper | `Stepper` |
| Switch | `Switch` |
| TabList | `Tabs` |
| Table | `Table` |
| Text | `Text` |
| TextArea | `TextArea` |
| TextInput | `TextInput` |
| Thumbnail | `Thumbnail` |
| TimeInput | `TimeInput` |
| Timer | `Timer` |
| Timestamp | `Timestamp` |
| Toast | `Toast` |
| ToggleButton | `Toggle` |
| Token | `Token` |
| Tokenizer | `Tokenizer` |
| Toolbar | `Toolbar` |
| Tooltip | `Tooltip` |
| TopNav | `TopNav`、`TopNavMegaMenu`、`TopNavMegaMenuItem` |
| TreeList | `TreeList` |
| Typeahead | `Typeahead` |
| VisuallyHidden | `VisuallyHidden` |

## 原生转换边界

- 菜单子项、分隔线、分组、勾选、单选与子菜单统一为 `MenuEntry`；Layout、Dialog 等部件使用插槽回调；表格内部结构由 MyGo 管理。
- `CheckboxInput → Checkbox`、`Selector → Select`、`TabList → Tabs`、`ToggleButton → Toggle`、`ProgressBar → Progress` 等采用 Go 命名。`VStack` 使用 `Stack`。
- `LayerProvider` 使用 `Scope`、`Overlay`、原生浮层和 `ToastViewport`；React hook 与 DOM 上下文不直接移植。
- 文件传输、AI 请求、流式数据和系统听写由宿主连接服务；Gallery 使用本地模型，不产生 AI 或上传请求。文件格式校验根据扩展名/MIME 元数据，不做内容嗅探。
- 日期范围采用两端原生日期输入；Carousel 没有触摸惯性滑动，Lightbox 不内置视频播放器。OverflowList 由宿主传入可见数量。
- ScrollableArea 采用 MyGo Scroll，保留可配置原生滚动；DOM stickyContainment、overscroll、浏览器 contenteditable 内嵌 token 与 Astryx Table React 插件不属于逐项兼容范围。

## 更新来源

```sh
python3 scripts/update_astryx_manifest.py /path/to/astryx
go test ./internal/catalog ./internal/gallery
```

脚本只读取指定的上游 checkout，登记导出名和映射，不复制上游组件实现。目录测试检查所有组件族映射到真实且有 Gallery 演示的公开 API。
