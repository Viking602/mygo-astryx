# MyGo / Astryx

参考 Astryx 设计体系、基于 [MyGo](https://github.com/egoist/mygo) 实现的 Go 原生 UI 组件库。

提供 **11 类、147 个公开组件入口、57 组可操作演示**，覆盖固定版本 Astryx 根入口的 **99 个 UI 组件族**。组件包括自有外观、原生组合与 MyGo 高阶适配，尺寸与配色不重复计数。复合部件通过 Go 模型或插槽表示；具体映射和平台差异见 [Astryx 覆盖清单](docs/astryx-coverage.md)。

## 运行 Gallery

要求 Go **1.27.3+**，MyGo 固定为 **v0.7.4**。支持 MyGo 所支持的桌面平台，真实平台验证情况见下文。

```sh
go run ./cmd/gallery
```

也可以先构建再运行：

```sh
go build -o bin/gallery ./cmd/gallery
./bin/gallery
```

Gallery 包含分类导航、组件搜索与单项选择、系统/浅色/深色主题、代码示例、表单、菜单、浮层、数据和聊天控件。窄窗口自动改为分类选择器，Ctrl/⌘ K 打开命令面板。可用 Astryx 的原始名称搜索等价组件，并展开各演示的固定源码对照。

```sh
go run ./cmd/gallery -theme dark -page forms
go run ./cmd/gallery -theme dark -component Chat
go run ./cmd/gallery -focus
go run ./cmd/gallery -list
go run ./cmd/gallery -snapshot artifacts/gallery-light.png -theme light
go run ./cmd/gallery -snapshot artifacts/gallery-dark.png -theme dark
go run ./cmd/gallery -snapshot artifacts/gallery-narrow.png -width 420 -height 780
```

可用分类：`layout`、`typography`、`actions`、`forms`、`navigation`、`feedback`、`overlays`、`data`、`menus`、`media`、`chat`，以及 `overview`。

## 使用组件

通过 Go module 安装 `v0.1.0`：

```sh
go get github.com/Viking602/mygo-astryx@v0.1.0
```

在代码中引入并使用：

```go
import (
    astryx "github.com/Viking602/mygo-astryx"
    "github.com/Viking602/mygo-astryx/theme"
    "github.com/egoist/mygo/ui"
)

func (app *State) view(c *ui.Context) {
    tokens := theme.Neutral(c.Theme().Dark, c.Preferences())
    astryx.Scope(c, tokens, func(k *astryx.UI) {
        k.Container().Children(func() {
            k.Heading("项目设置", 2)
            k.Field("名称", func() {
                k.Key("project-name").TextInput(&app.Name).FillWidth()
            })
            k.Button("保存", astryx.ButtonOptions{Variant: astryx.Primary}).
                OnClick(app.Save)
        })
        k.ToastViewport()
    }).Fill()
}
```

`State`、`Name` 和 `Save` 属于使用方自己的应用模型。完整可运行示例：

```sh
go run ./examples/basic
```

每次视图构建重新创建 builder。`UI`、`ui.Element` 和组合部件不跨构建保存；需要持久控件身份时使用 `ui.Handle`。有状态控件使用 `k.Key(id).Control(...)`，它等价于在构造前调用 `c.Key(id)`。

## 组件范围

| 分类 | 数量 | 内容 |
| --- | ---: | --- |
| 布局 | 16 | Stack、Grid、Layout 插槽、Resizable、AppShell 等 |
| 文字与展示 | 21 | Markdown、CodeBlock、Citation、AvatarGroup、ClickableCard 等 |
| 动作 | 5 | Button、IconButton、Toggle、ButtonGroup、Link |
| 输入与选择 | 34 | MultiSelector、Typeahead、PowerSearch、FileInput、Tokenizer、可替换指示器等 |
| 导航 | 16 | TopNav、SideNav、MobileNav、MegaMenu、Outline、Stepper 等 |
| 反馈 | 10 | Banner、StatusDot、Timer、Timestamp、Toast 等 |
| 浮层 | 9 | Dialog、Drawer、HoverCard、BottomSheet、Overlay 等 |
| 数据与日期 | 11 | Table、List、TreeList、DateRangeInput、DateTimeInput 等 |
| 菜单与命令 | 6 | DropdownMenu、MoreMenu、ContextMenu、NavMenu、Toolbar、CommandPalette |
| 图片与媒体 | 3 | Thumbnail、Carousel、Lightbox |
| 聊天与上下文 | 16 | Chat、消息、输入、工具详情、上下文、反馈与听写服务接口 |

[完整 API 与实现方式](docs/components.md) · [Astryx 覆盖清单](docs/astryx-coverage.md) · [验证记录](docs/verification.md) · [架构](docs/native-ui-plan.md)

## 检查与预览

```sh
go test ./...
go vet ./...
make snapshots
make snapshots-all
```

测试覆盖关键交互、中文 IME、禁用选项、弹窗嵌套与焦点恢复、文件拖放与校验、筛选增删改、聊天发送与失败保留，以及每个公开组件在浅深色、1280/320 宽窗口下的渲染。`snapshots-all` 生成 57 组演示的 228 张截图，以及 10 种交互状态的 40 张截图和 JSON 报告。中性主题的主要文字颜色组合有对比度检查。

原生窗口冒烟验证需 Python 3 和可用的桌面会话：

```sh
make smoke-native
```

该命令仅在 `e2e` 构建中启用 MyGo 本地自动化，逐项选择全部 147 个入口并操作代表性控件，验证后关闭它启动的 Gallery。普通构建不包含此服务。结果、窗口截图和日志写入 `artifacts/`。

AI 后端、文件上传和系统听写由使用方接入；Gallery 的聊天与听写明确标为本地演示。本轮无窗口测试与截图通过，macOS 原生复验遇到首帧后不刷新的运行环境问题，最小 MyGo 程序也复现；原生遍历尚未完成。Windows/Linux 真实窗口、各平台读屏和完整视觉回归仍需分别验收。组件族覆盖不等于与 Astryx 的 React API 完全兼容，TabMenu 与 TableFooter 的保留差异已明确登记。
