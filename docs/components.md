# 组件 API

所有组件通过当前构建的 `*astryx.UI` 创建。使用 `astryx.New(c)` 跟随当前 MyGo 主题；使用 `astryx.Scope(c, tokens, build)` 应用完整自有主题。调用方继续使用返回元素的 MyGo 方法设置布局和事件。

## 主题与配置

- `theme.System(c)`：保留 MyGo 已解析的系统外观、强调色、文字和偏好。
- `theme.Neutral(dark, c.Preferences())`：浅/深中性预设，处理高对比度和文字缩放；动画仍由 MyGo 的系统偏好控制。
- `theme.FromNative(native)`：基于现有 `ui.Theme` 添加语义设计值，不修改传入主题。
- `k.Tokens()`：取得当前设计值的副本。自有主题按窗口或局部作用域工作，不使用全局可变单例。
- 一个组件可传零个或一个 Options；若传多个，最后一个完整配置生效，不做字段合并。

主题字段包含 MyGo 的 `Native`，以及 `Muted`、`Subtle`、`Elevated`、状态弱背景、阴影、间距、字体和圆角。直接修改 Native 配色后，宜重新通过 `FromNative` 计算相关语义背景，避免旧派生值与新配色混用。

## 布局与展示

| 方法 | 参数 / 返回 | 实现方式 |
| --- | --- | --- |
| `Container` | `() ui.Element` | 有默认边距、宽度上限的内容容器 |
| `Stack`、`Inline` | `() ui.Element` | 纵向堆叠 / 可换行横排 |
| `Grid` | `(columns int) ui.Element` | MyGo 栅格组合，默认统一间距 |
| `Scroll` | `(...ScrollOptions) ui.Element` | 默认纵向，Horizontal / Both 支持横向和双轴 |
| `Spacer`、`Divider` | `() ui.Element` | 原生布局适配 |
| `AppShell` | `(navigation, content func()) ui.Element` | 窄窗上下布局，宽窗左右布局 |
| `Split` | `(position *float32, first, second func()) ui.Element` | MyGo 原生可调整分栏 |
| `Text`、`Caption` | `(text string) ui.Element` | 正文 / 辅助文字 |
| `Heading` | `(text string, level int) ui.Element` | 1–6 级标题与可访问性等级 |
| `Code` | `(text string) ui.Element` | 可选择的等宽文字 |
| `Icon` | `(icon *ui.SVG) ui.Element` | 原生 SVG；不为动作图标自动提供名称 |
| `Avatar` | `(name string, image *ui.Bitmap, ...AvatarOptions) ui.Element` | 原生头像，可设置尺寸、状态与动作 |
| `Badge` | `(text string, tone ...Tone) ui.Element` | 语义颜色标记 |
| `Card` | `(...CardOptions) ui.Element` | 表面、边框、圆角及可选阴影 |
| `EmptyState` | `(title, description string, action func()) ui.Element` | 空状态组合 |

`CardOptions`：`Padding`（零值使用默认边距）、`Elevated`。

## 动作

| 方法 | 参数 |
| --- | --- |
| `Button` | `(label string, ...ButtonOptions)` |
| `IconButton` | `(label string, icon *ui.SVG, ...ButtonOptions)` |
| `Toggle` | `(on *bool, label string, ...ButtonOptions)` |
| `ButtonGroup` | `(build func())` |
| `Link` | `(label, url string)` |

均返回 `ui.Element`。`ButtonOptions` 包含 `Variant`、`Size`、`Disabled`、`Loading`、`Icon`、`EndIcon`。变体为 `Default`、`Primary`、`Outline`、`Ghost`、`Destructive`；尺寸为 `Small`、`Medium`、`Large`，零值为 Medium。加载中的按钮同时禁用动作。

优先用返回元素的 `OnClick` 注册动作，或使用 `Clicked()` 的原生轮询语义。修改正在遍历的集合时使用回调。图标按钮必须给出可理解的 Label。

## 输入与选择

| 方法 | 参数 |
| --- | --- |
| `TextInput`、`TextArea` | `(value *string, ...InputOptions)` |
| `InputGroup` | `(value *string, leading, trailing func(), ...InputOptions)` |
| `Field` | `(label string, build func(), ...InputOptions)` |
| `Form` | `(build func())` |
| `Fieldset` | `(legend string, build func())` |
| `NumberInput` | `(value *float64, low, high, step float64, ...InputOptions)` |
| `SearchField` | `(value *string, ...InputOptions)` |
| `Checkbox`、`Switch` | `(value *bool, label string, ...ChoiceOptions)` |
| `RadioGroup` | `(label string, build func())` |
| `Radio` | `(value *string, choice, label string, ...ChoiceOptions)` |
| `Select` | `(value *string, items []Option, ...InputOptions)` |
| `Combobox` | `(value *string, items []string, ...InputOptions)` |
| `Slider` | `(value *float64, low, high float64, ...InputOptions)` |
| `RangeSlider` | `(lowValue, highValue *float64, low, high, step float64)` |
| `Rating` | `(value *int, maximum int)` |

均返回 `ui.Element`。模型指针由应用持有，不能为 nil。普通输入建立在 TextInput/TextArea Base 上，保留原生 IME、选择、撤销和剪贴板。

`InputOptions` 包含 `Size`、`Label`、`Placeholder`、`Description`、`Error`、`Disabled`、`ReadOnly`。文字输入使用全部配置；选择器使用适用的标签、说明、错误、尺寸和禁用配置。MyGo 高阶适配控件保留其原生尺寸及内部结构，具体支持见源码；不要假设每个配置都适用于所有控件。

`ChoiceOptions` 包含 `Disabled`、`Mixed`（Checkbox）、`Size`。`Option` 包含稳定的 `Value`、展示的 `Label` 与 `Disabled`。禁用选项作为不可操作行显示，不注册为 MyGo 的键盘候选，因此不会被方向键选中；其无障碍表现为带不可用说明的行，而非完整可选择选项。

`Combobox` 的模型包含已提交的选项或用户自由输入的文字，不强制只能使用预设值。`Field(...).Description(...).Error(...)` 用于显示说明和错误，并关联原生可访问性信息。

## 导航

| 方法 | 参数 |
| --- | --- |
| `Tabs`、`Segmented` | `(selected *int, labels []string)` |
| `Breadcrumbs` | `(items []Breadcrumb)` |
| `Pagination` | `(page *int, pages int)` |
| `Sidebar` | `(selected *string, items []NavItem)` |
| `Collapsible` | `(label string, open *bool, build func())` |
| `Accordion` | `(items []Disclosure)` |

Tabs / Segmented 的索引从 0 开始；Tabs 只构建标签栏，应用按选中值构建内容。组名可通过返回元素的 `.Label(...)` 设置。Pagination 页码从 1 开始，越界值会夹紧。Sidebar 使用稳定 ID；Disclosure 必须有唯一 ID、非 nil 的 Open 指针和可选 Content。

`Breadcrumb`：Label、OnClick、Menu（`[]MenuEntry`）。`NavItem`：ID、Label、Badge、Icon、Disabled、Section、Children。`Disclosure`：ID、Title、Open、Content。

Collapsible 返回触发元素；内容由其 Panel 在同一父容器中构建。Accordion 组合多个独立展开项，允许同时打开。

## 反馈与浮层

| 方法 | 参数 |
| --- | --- |
| `Alert` | `(title, description string, ...Tone)` |
| `Progress` | `(value float64)`，范围 0–1 |
| `Spinner` | `()` |
| `Skeleton` | `(width, height float32)`，静态占位 |
| `Toast` | `(message ui.Toast) string`，返回通知 ID |
| `ToastViewport` | `()`，每个视图调用一次 |
| `Dialog`、`Drawer` | `(title string, open *bool, content func(), ...DialogOptions)` |
| `AlertDialog` | `(title, description string, open *bool, ...AlertDialogOptions)` |
| `Popover` | `(anchor ui.Element, open *bool, content func())` |
| `Tooltip` | `(anchor ui.Element, description string)` |

除 Toast 外均返回 `ui.Element`，关闭中的浮层可能返回零元素，用 `Valid()` 判断。Tone 为 Neutral、Info、Success、Warning、Danger。

DialogOptions：Width、Description、Header、Footer。Dialog / BottomSheet 可自定义 Header，Dialog / Drawer 复用 Base 的模态焦点、关闭和恢复行为。AlertDialog 由公开 Overlay/Modal 组合：点击外部不会关闭，Escape 取消；危险动作默认聚焦取消，普通确认默认聚焦确认。AlertDialogOptions：Confirm、Cancel、Destructive、OnConfirm。

Tooltip 同时给 anchor 设置 Description。通知只在用户动作或模型事件中触发，不能无条件在构建函数中添加。Toast 的计时、动作、关闭和可访问性由 MyGo Toast Base 负责。

## 数据与日期适配

| 方法 | 参数 |
| --- | --- |
| `Table` | `(state *ui.ListState, columns []ui.TableColumn, count int, cell func(row, col int))` |
| `List` | `(state *ui.ListState, count int, row func(int))` |
| `Tree` | `(build func())` |
| `TreeItem` | `(label string, open *bool, build func())` |
| `Calendar`、`DateInput`、`TimeInput` | `(value *time.Time)` |
| `ColorPicker` | `(value *ui.Color)` |

这些是 MyGo 高阶控件的主题适配，保留原生行为，不声称重新实现其内部算法。Table/List 的排序、选择和稳定行 Key 由应用持有的 ListState 配置。树叶子可传 nil 的 open 和 build。

## 高级布局与内容

以下入口也都返回 `ui.Element`，可继续设置大小、事件与布局。

| 方法 | 参数与作用 |
| --- | --- |
| `AspectRatio`、`Center`、`HStack` | 比例参数 / 无参数；比例、居中和横向布局 |
| `GridSpan`、`FormLayout` | `(columns int)`；跨列元素 / 表单栅格 |
| `Layout` | `(LayoutSlots)`；Header、Content、Panel、Footer、PanelWidth |
| `Section` | `(title, description string, build func())` |
| `Resizable` | `(*ResizableState, first, second func(), ...ResizableOptions)` |
| `VisuallyHidden`、`Kbd` | `(text string)`；隐藏语义文本 / 按键显示 |
| `AvatarGroup` | `([]Person, maximum int, ...AvatarGroupOptions)`；超出部分合并计数 |
| `AvatarGroupOverflow` | `(count int, ...AvatarGroupOptions)`；真实人数、自定义文字和可选动作 |
| `AvatarStatusDot` | `(label string, tone Tone, ...AvatarStatusOptions)`；实心、圆环和横杠 |
| `Blockquote` | `(text, attribution string)` |
| `MetadataList` | `([]Metadata)`，Metadata 为 Label、Value |
| `Citation` | `(label string, ...CitationOptions)`；Detail、URL、OnClick |
| `Markdown` | `(source string, ...MarkdownOptions)`；Native 为 MyGo `markdown.Options`，Streaming 启用流式渲染 |
| `CodeBlock` | `(source, language string)`；原生语法高亮与复制 |
| `Item` | `(label string, ...ItemOptions)`；Description、Trailing、Icon、Disabled、OnClick |
| `ClickableCard` | `(label string, content func(), ...CardOptions)`，支持独立的嵌套动作 |
| `SelectableCard` | `(selected *bool, label string, content func())` |

`ResizableState` 保存 Size、Collapsed。Options 支持 Min、Max、Vertical、Percent、Collapsible 与 OnResize。分隔线复用 MyGo 的拖动、键盘和可访问性，应用配置的尺寸上下限在回写模型时夹紧；读屏看到的范围仍来自 MyGo Split。`VisuallyHidden` 保留文本语义，不占正常布局和输入区域。

AvatarOptions 支持 Size、Status、StatusTone、OnClick；状态合并到头像的可访问名称中，不仅依靠颜色表达。AvatarStatusOptions.AvatarSize 决定状态点的尺寸级别，Icon 可替换中大尺寸的内置标记。AvatarGroupOptions 支持 OverflowLabel、OnOverflow、Disabled；自定义显示文字时，可访问名称继续保留真实隐藏人数。

## 高级选择与文件

| 方法 | 参数 |
| --- | --- |
| `CheckboxList`、`MultiSelector` | `(values *[]string, items []Option, ...SelectionOptions)` |
| `RadioList` | `(value *string, items []Option, ...InputOptions)` |
| `Typeahead` | `(*TypeaheadState, []Suggestion, ...TypeaheadOptions)` |
| `ComplexSelector` | `(label string, open *bool, summary string, content func(close func()), ...InputOptions)` |
| `Token` | `(label string, onRemove func())`；nil 表示不可移除 |
| `Tokenizer` | `(values *[]string, suggestions []string, ...InputOptions)` |
| `PowerSearch` | `(*PowerSearchState, []FilterField, ...PowerSearchOptions)` |
| `FileInput` | `(*FileInputState, ...FileInputOptions)` |
| `FieldStatus` | `(message string, tone Tone)` |
| `InputClearButton` | `(value *string, focus *ui.Handle, ...ButtonOptions)` |
| `Indicator`、`CheckboxIndicator`、`RadioIndicator`、`CheckIndicator` | `(IndicatorState, ...IndicatorOptions)` |

`SelectionOptions`：Label、Disabled、MaxSelected、OnChange。MultiSelector 支持搜索、方向键进入候选、选择上限、清空和完成；长列表在弹层内滚动。回调拿到切片副本。

`TypeaheadState` 将 Query（展示文本）、Value（稳定 ID）与 Loading、Error 分开。Suggestion 包含 ID、Label、Description、Group、Disabled；OnSelect 回传完整建议，OnQuery 用于宿主加载候选。宿主负责请求取消、去重和结果更新，后台结果通过窗口更新机制应用。

`PowerSearchState` 保存自由 Query 和 Filters，FilterField 的 Kind 支持 text、number、date、enum，也能提供 Operators、Validate 与 Editor。SearchFilter.ReadOnly 保护固定条件；Options.ReadOnly 保护整个编辑界面，Disabled 禁用交互。自由文字变化和条件增删改均触发根元素变化通知；OnChange 的参数是当前条件副本，Query 从模型读取。

`FileInputOptions` 包含 Label、Accept（逗号分隔的扩展名/MIME）、Multiple、Disabled、MaxFiles、MaxSize、Browse、OnChange。默认 Browse 调用本机系统文件选择器，Browse 也可由宿主注入；拖放与浏览共用 `AddUploadPaths` 校验。FileInputState 保存 Files、Error、Busy、Progress，忙碌时禁用浏览、删除和拖放。该控件不上传文件；大小和格式限制来自文件元数据，不做内容识别。

状态指示器是装饰性图形，父 Checkbox / Radio / Select 保留角色、键盘和焦点。可用 `theme.Theme.Indicator` 回调替换指定 kind/state 的绘制；返回 `(element, true)` 表示接管。自定义内容也应使用 `RoleNone`，避免把装饰文字暴露成独立控件。

## 菜单、导航与高级浮层

| 方法 | 参数 |
| --- | --- |
| `DropdownMenu` | `(label string, []MenuEntry, ...ButtonOptions)` |
| `MoreMenu`、`NavMenu` | `(label string, []MenuEntry)` |
| `ContextMenu` | `(target ui.Element, []MenuEntry, ...ContextMenuOptions)` |
| `Toolbar` | `(label string, content func())` |
| `CommandPalette` | `(open *bool, *CommandPaletteState, []Command, ...CommandPaletteOptions)` |
| `TopNav`、`SideNav` | `(selected *string, []NavItem, ...NavOptions)` |
| `MobileNav` | `(selected *string, []NavItem)` |
| `MobileNavToggle` | `(label string, open *bool, ...ButtonOptions)`；与 Drawer/BottomSheet 共用状态 |
| `TopNavMegaMenu` | `(label string, items, featured func(), ...HoverCardOptions)` |
| `TopNavMegaMenuItem` | `(label string, ...ItemOptions)`；执行导航后关闭所属面板 |
| `OverflowList` | `([]NavItem, visible int, onSelect func(string))` |
| `Outline` | `(selected *string, []OutlineItem, onNavigate func(string))` |
| `Stepper` | `(current *int, []Step, ...StepperOptions)` |
| `HoverCard` | `(anchor ui.Element, content func(), ...HoverCardOptions)` |
| `InfoTip` | `(label, description string)` |
| `BottomSheet` | `(title string, open *bool, content func(), ...DialogOptions)` |
| `Overlay` | `(open *bool, content func(), ...OverlayOptions)` |

`MenuEntry` 统一表示 Group、Separator、Disabled、Checked、Radio、Value、Children、快捷键和 OnSelect。Checked 为 `*bool`，Radio 为共享 `*string`；分组标题使用禁用菜单行。ContextMenu 默认由 MyGo 绘制，Native=true 使用系统上下文菜单。菜单保留方向键、子菜单、禁用项跳过和 Escape 顺序。

Command 包含 ID、Label、Description、Group、Shortcut、Disabled、Run。面板支持 Inline 或 Dialog 展示、关键词筛选、方向键和 Enter，输入法合成中的 Enter 不执行命令。全局快捷键由应用注册，Gallery 已注册 Ctrl/⌘ K。

HoverCardOptions 支持 Delay、HideDelay、Disabled、OpenOnClick、FocusOnOpen、Label、Width、MaxHeight。长内容在视口约束内滚动。悬停延迟打开，键盘可见焦点立即打开，指针进入内容区可继续操作；Escape 或点击关闭后不会在静止指针下立刻重开。

TopNavMegaMenu 在窄窗使用 Drawer；桌面键盘打开后进入内容区，Escape 返回触发器。在 items/featured 回调中使用同一个 builder 创建 TopNavMegaMenuItem，它会在执行宿主动作时关闭面板并恢复焦点；普通 Item 可以保留面板。MobileNav 的父项目打开子菜单，叶子项目回写选择。MobileNavToggle 的 open 指针应与对应 Drawer/BottomSheet 共用，关闭面板会同步更新展开状态。

OverlayOptions 可设 NonModal、DismissOnOutside、Label；模态内容应提供明确标题和关闭动作。Dialog、BottomSheet 和 AlertDialog 的长说明滚动，操作区保持可见。Gallery 的分步底部面板替换同一个受控面板的内容，不新增已被 Astryx 弃用的 BottomSheetSwitcher 别名。

## 日期、时间与图片

| 方法 | 参数 |
| --- | --- |
| `DateRangeInput` | `(*DateRange, ...DateRangeOptions)` |
| `DateTimeInput` | `(*time.Time, ...InputOptions)` |
| `TreeList` | `(*ui.OutlineState[string], roots []string, children func(string) []string, label func(string) string)` |
| `Banner` | `(title, description string, ...BannerOptions)` |
| `StatusDot` | `(label string, tone Tone)` |
| `Timestamp` | `(time.Time, ...TimestampOptions)` |
| `Timer` | `(*TimerState, ...TimerOptions)` |
| `Thumbnail` | `(ImageAsset, ...ThumbnailOptions)` |
| `Carousel` | `(index *int, []CarouselSlide, ...CarouselOptions)` |
| `Lightbox` | `(open *bool, index *int, []ImageAsset)` |

`DateRange` 保存 Start、End、Empty；Options 支持 Min、Max、MaxDays（含两端的日历天数）、Allowed、Presets、Disabled 和 OnChange。`ValidateDateRange` 可独立校验，不受夏令时的小时差影响。编辑时保留无效值并展示错误，只有有效范围或清空才执行 OnChange。带 Allowed 回调的范围最多检查 100 年，以免异常范围长时间阻塞 UI。

Allowed 接收开始日期所在时区的日历午夜，避免同一天因时分秒不同而改变可用性。空值仍能选择预设；默认范围违反规则时保持受控值为空，显示可修改的临时日期，只有有效范围才能应用。

`TimerState.Start/Pause(now)` 使用调用方时钟，Reset 清零，Duration(now) 读取累计值；只在 Running 时安排下一秒更新。Timestamp 可显示 relative、date、time、absolute，支持 Location、Static 和 Copy。Banner 可绑定 Dismissed 并提供动作；状态点同时显示文字。

ImageAsset 使用 `ui.ImageSource`，支持 Bitmap / SVG。CarouselSlide 使用通用 Content 回调，支持循环和左右键；Lightbox 提供图片切换、缩放、旋转和滚动，关闭恢复触发器焦点。

## 聊天组件

`Chat(*ChatState, ...ChatOptions)` 组合虚拟消息列表与输入区。ChatState 保存 Messages、Draft、Tokens、Busy、Error、List；消息 ID 必须稳定且唯一。高阶 Chat 默认跟随列表末尾；需要独立控制滚动时使用 ChatMessageList。

| 方法 | 参数 |
| --- | --- |
| `ChatLayout` | `(messages, composer func())` |
| `ChatMessageList` | `(*ui.ListState, []ChatMessageData, ...ChatMessageOptions)` |
| `ChatComposer` | `(*ChatState, ...ChatOptions)` |
| `ChatComposerInput` | `(value *string, ...ChatComposerInputOptions)` |
| `ChatSendButton` | `(busy, canSend bool, onSend, onStop func())` |
| `ChatLayoutScrollButton` | `(*ui.ListState, label string)` |
| `ChatMessage` | `(*ChatMessageData, ...ChatMessageOptions)` |
| `ChatMessageBubble`、`ChatMessageMetadata` | `(ChatMessageData)` |
| `ChatSystemMessage` | `(text string)` |
| `ChatToolCalls` | `([]ChatToolCall)` |
| `ChatComposerTokenElement` | `(ChatToken, onRemove func())` |
| `ChatTokenizedText` | `([]ChatTextRun)` |
| `ChatComposerDrawer` | `(label string, open *bool, content func())` |
| `ChatDictationButton` | `(*DictationState, ...DictationOptions)` |

ChatOptions.OnSend 接收 `ChatDraft{Text, Tokens}` 并返回 error；失败保留草稿并显示错误，成功清空草稿与上下文。默认 Enter 发送、Shift+Enter 换行；MultilineEnter=true 时普通 Enter 换行，Ctrl/⌘ Enter 发送。合成中的输入、只读/禁用和 Busy 状态不发送。没有 OnSend 时仅把用户消息追加到本地模型，不产生模型请求。

ChatMessageData.Tokens 保存消息携带的上下文；默认发送和 Gallery 的本地回调都保留它，纯上下文消息不会变成空白气泡。ChatOptions.MessageOptions 将重试与反馈回调转发给高阶 Chat 内的消息列表。

消息支持 Markdown、流式展示、元数据、展开说明、工具输入/输出、复制、反馈和重试回调。宿主负责网络、流式更新、重试逻辑和服务状态。听写按钮只有 Available=true 且 Start/Stop 回调齐全时可用；Gallery 的听写回调是明确标注的本地模拟，不访问麦克风。

## 生命周期

`UI`、`ui.Element`、组件部件只在当前构建过程使用。不要把它们放进应用模型或后台 goroutine。跨帧保留模型和 `ui.Handle`，后台更新通过 `Window.Update`。以 `k.Key(id)` 为动态列表项目提供稳定身份，避免用索引代替业务 ID。
