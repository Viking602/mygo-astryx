# 基于 MyGo 的原生 UI 组件库

日期：2026-11-06。方向已确认：**原生实现，参考 Astryx 重新构建完整 UI 库。**

MyGo 指 [`egoist/mygo`](https://github.com/egoist/mygo)，并非本项目新起的品牌名。正式模块路径和发布名称尚未设定。

## 源码基线与结论

已阅读官方文档，将公开源码克隆到临时目录，并核对以下固定版本：

- Tag：`v0.7.4`。
- Commit：`9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e`。
- 上游要求 Go `1.27.3`；本机原始版本为 `go1.25.10 darwin/arm64`，实施时已在临时 Go 缓存准备 `go1.27.3`。

项目已使用匹配工具链构建和测试，并完成 macOS 原生窗口冒烟验证。[固定版本 go.mod](https://github.com/egoist/mygo/blob/9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e/go.mod)

MyGo 已具备布局、绘制、控件、文字输入、无障碍和无窗口测试能力，适合在其上实现完整设计系统。优先使用公开的 Base 控件构建自有外观；高阶控件逐项验证定制边界，确有不足时再补公开扩展点。

这里的原生是 MyGo 自绘的 Go UI：Metal / Direct3D 11 / OpenGL 及软件渲染路径，并不表示每个控件都是 AppKit、WinUI 或 GTK 的系统控件实例。[渲染说明](https://mygo.egoist.dev/docs/ui/rendering)

## 可以直接复用的基础

### 无默认外观的交互控件

Base 控件负责鼠标、键盘、焦点和对应的可访问性语义，本库负责颜色、尺寸、布局与绘制：

- `ButtonBase`、`CheckboxBase`、`SwitchBase`、`ToggleBase`、`RadioBase`、`SliderBase`。
- `TabsBase`、`SegmentedBase`、`CollapsibleBase`。
- `SelectBase`、`ComboboxBase`。
- `PopoverBase`、`DialogBase`、`TooltipBase`、`ToastViewportBase`、`ToastBase`。
- `TextInputBase`、`TextAreaBase`。

依据：[自定义控件文档](https://mygo.egoist.dev/docs/ui/custom-widgets)、[固定版本公开 API](https://github.com/egoist/mygo/blob/9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e/ui/api_gen.go)、[内部实现对照](https://github.com/egoist/mygo/blob/9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e/ui/base.go)。本库只调用公开 API，不访问私有 node、engine 或 internal 包。

### 主题和系统偏好

`ui.Theme` 已有背景、表面、文字、强调色、边框、状态色、选区、焦点、圆角、间距和字体；窗口主题用 `SetTheme`，局部主题用 `ui.Themed`。

`c.Preferences()` 提供系统强调色、高对比度、文字缩放和减少动态效果。自有主题必须明确继承策略，不能认为调用 `SetTheme` 后所有偏好自然生效。

多级表面、弱化状态背景、不同组件尺寸与阴影等额外设计值放在本库自己的 Go 类型中，再映射到 MyGo 主题及元素样式，不能假定 `ui.Theme` 已有对应字段。

依据：[主题文档](https://mygo.egoist.dev/docs/ui/styling)、[Theme 源码](https://github.com/egoist/mygo/blob/9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e/ui/theme.go)、[Preferences 源码](https://github.com/egoist/mygo/blob/9bd34e1eb9525d9a1fe698aadc964c49c8ff8e7e/ui/preferences.go)。

### 高阶控件

MyGo 已有虚拟列表、表格、树、日期与时间、颜色选择、表单、Markdown、编辑器、导航和图表等能力。其定制入口不同，不能假设每项都有 Base 或可任意改造内部部件。[原生组件目录](https://mygo.egoist.dev/docs/ui)

## 本库的组成

1. **设计值与主题：** 参考 Astryx 建立语义颜色、排版、尺寸、间距、圆角、阴影和动效，再对接 MyGo。
2. **原生组件：** 在 Base 上实现统一外观，保留 `ui.Element` 和组合部件供调用方继续定制。
3. **页面模式：** 设置、列表、详情、表单和工作台，验证组件在真实组合中的一致性。
4. **原生 Gallery：** 展示组件、尺寸、变体与交互状态，可切换浅深色；作为主要设计和验收入口。
5. **文档与发布：** Go module、GoDoc、中文使用指南、示例、测试、版本和迁移说明。

外观先以 Astryx Neutral 的克制基础作为对照，建立可替换主题；具体色值和尺寸通过原生 Gallery 校准。未选择任何与框架同名的动漫或其他视觉品牌。[Astryx 对照调研](astryx-review.md)

## 组件映射

以下是最初的实施分工。当前已完成的组件族与差异以 [固定来源覆盖清单](astryx-coverage.md) 为准。

| 组件族 | 实现基础 | 本库需要完成 |
| --- | --- | --- |
| Button、IconButton、Toggle、ButtonGroup | `ButtonBase`、`ToggleBase`、布局 | 变体、尺寸、图标、加载、组合边界与可访问名称 |
| Checkbox、Switch、Radio | 对应 Base 和公开分组 API | 指示器、选中/混合/禁用状态、标签和分组 |
| Slider、Tabs、SegmentedControl | 对应 Base | 轨道、滑块、选中样式、方向和尺寸 |
| TextInput、TextArea、InputGroup | 输入 Base | 前后缀、清除、错误、只读；保留 IME、选区、光标与撤销 |
| Select、Combobox | 对应 Base 的组合部件 | 触发器、输入、弹层、选项、高亮、空状态和键盘行为 |
| Dialog、Popover、Tooltip、Toast | 对应 Base 与浮层 API | 内容结构、尺寸、遮罩、动作区、位置和动效 |
| Drawer、HoverCard、CommandPalette | Base、Overlay、AttachTo、Modal 等组合 | 定义交互并独立测试，遵守浮层与焦点规则 |
| Accordion、Collapsible | `CollapsibleBase` 与布局 | 标题、图标、内容展开和动画 |
| Text、Heading、Card、Badge、Divider、EmptyState | 文字、布局、绘制和角色 | 排版层级、表面、状态和内容组合 |
| Table、List、Tree、Outline | 现有高阶控件 | 保留虚拟化和导航，验证行、表头及内部部件定制边界 |
| Calendar、DateInput、TimeInput、ColorPicker | 现有高阶控件 | 样式、格式与校验适配，核查必要的公开扩展点 |
| AppShell、SideNav、Toolbar、Breadcrumbs、Pagination | 布局、导航与现有控件 | 窗口缩放布局、应用结构和页面组合 |
| Progress、Spinner、Skeleton、状态反馈 | 状态控件与动画/绘制 | 统一视觉、减少动态效果和停止重绘策略 |

文件选择、跨应用拖放和系统菜单按原生系统服务设计，不直接照搬浏览器行为。颜色和尺寸变体不重复计算为独立组件。

## API 和生命周期规则

- 组件以 `*ui.Context` 和 Go 函数组合；表单数据归应用状态所有。
- `ui.Element` 仅在当前构建过程有效，不能跨帧保存或传给后台 goroutine。持久身份用 `ui.Handle`，后台模型更新通过 `Window.Update`。
- 有状态控件构造前用 `c.Key(id)` 设置稳定身份，列表使用业务 ID，避免重排后串状态。
- 同一帧可能多次构建；构建函数不启动网络、文件写入、goroutine 或一次性任务。修改正在遍历的集合时优先使用事件回调。
- 普通文字输入复用输入 Base；自定义控件保留角色、Label、值、描述与错误信息。
- 主题按窗口和局部范围工作，避免进程级可变单例；测试局部主题内的弹层和多窗口。
- 使用 MyGo 的按需绘制和动画，不另建永久刷新循环；效果结束后应能静止。

依据：[视图](https://mygo.egoist.dev/docs/ui/views)、[API 迁移](https://mygo.egoist.dev/docs/ui/migration)、[浮层](https://mygo.egoist.dev/docs/ui/overlays)、[渲染](https://mygo.egoist.dev/docs/ui/rendering)。

## 建议的目录

先保持单个 Go module，统一组件包内按控件拆分文件，出现明确边界后再拆子包：

| 位置 | 职责 |
| --- | --- |
| `go.mod` | 固定 MyGo 版本与公开模块路径 |
| `theme/` | 设计值、主题及 `ui.Theme` 映射 |
| 根目录按组件拆分的 `.go` 文件 | 公开组件和组合部件 |
| `patterns/` | 设置、列表、详情等页面组合 |
| `cmd/gallery/` | Go 原生组件浏览器 |
| `internal/catalog/` | 组件分类、说明、状态和示例登记 |
| `docs/` | 安装、主题、API、迁移与使用说明 |

基础工程已按此结构创建；页面模式包尚未拆出。组件元数据支持 Gallery、单项选择、`-list` JSON 导出与演示覆盖检查。

## 实施顺序

1. **打通原生预览。** 准备工具链、固定版本，实现主题和 Button、Input、Checkbox、Select、Dialog 等代表控件，在原生 Gallery 验证交互和组合。这一步不代表全库完成。
2. **覆盖通用组件。** 按映射清单逐项登记“已实现 / 适配现有控件 / 需要扩展点 / 暂未支持”，覆盖焦点、禁用、加载、错误、选中、空内容和长文本。
3. **用完整页面验证。** 实现设置、数据列表、详情与表单示例，核查筛选、编辑、浮层嵌套、窗口尺寸与主题切换。
4. **交付独立可用的库。** 用独立应用消费公开 API，补齐 GoDoc、中文指南、可运行示例、版本、许可和变更说明。

## 验证方式

`ui.NewTester` 可无窗口测试点击、键盘、输入法合成、主题、缩放、焦点和部分可访问性输出，也可取得渲染快照。测试时钟应显式推进，不用 sleep。[测试文档](https://mygo.egoist.dev/docs/ui/testing)

已阅读上游 `base_test.go`、`themed_test.go`、`examples/counter-native/main_test.go` 的部分测试，确认入口和用例存在；本轮未运行它们。

重点验收实际行为：

- 点击、Enter/Space、Tab、方向键和禁用状态正确。
- 中文 IME 合成、编辑、撤销、复制和选区正确。
- Select 在 Dialog 中可用，Escape 逐层关闭并恢复焦点。
- 列表重排、过滤、增删后不串编辑状态和焦点。
- 局部主题、多窗口、高对比度和文字缩放可用。
- 减少动态效果下合理，动画结束后不持续请求绘制。
- `go test` 覆盖关键状态和组合交互，固定条件的渲染快照辅助视觉检查。
- macOS、Windows、Linux 的真实窗口和读屏仍需各自验证，无窗口测试和交叉编译不能代替。

## 当前状态

已完成 Go module、主题、11 类共 147 个公开组件入口、57 组原生 Gallery 演示、独立使用示例和关键交互测试。已核对固定 Astryx 提交根入口的 99 个 UI 组件族；复合部件以原生模型、插槽或 MyGo 适配表示，差异均登记在覆盖清单。

已通过每个组件入口的浅深色、宽窄窗口渲染测试，生成 228 张默认演示截图与 40 张状态截图。上一轮 macOS 真实窗口曾完成 143 个入口的遍历；本轮原生会话停留首帧，新增到 147 个入口的原生遍历未完成，最小 MyGo 程序也复现。当前结果以验证记录和失败报告为准。Windows/Linux 真实窗口、对应读屏工具与完整视觉回归仍需分别执行。

后续范围是更多业务页面、图表与编辑器等专门控件，以及 API 稳定化、服务接入和多平台发行验收。HoverCard、CommandPalette、聊天、文件和高级选择等已纳入本轮实现。详细 API 见 [组件文档](components.md)，测试证据见 [验证记录](verification.md)。
