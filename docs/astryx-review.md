# Astryx 对照调研

调研日期：2026-11-06。

**用户已确认：基于 `egoist/mygo` 的 Go 原生 UI，参考 Astryx 重新实现完整组件库。** 此前面向 React 扩展包的建议已被本方案替代。

原生架构、组件映射和交付顺序见 [MyGo 原生 UI 实施方案](native-ui-plan.md)。

## 可以借鉴的设计

| Astryx 的做法 | 本项目的采用方式 | 官方来源 |
| --- | --- | --- |
| 设计基础、组件、页面模式组成统一体系 | 同时建设 Token、组件和完整页面组合 | [仓库 README](https://github.com/meta/astryx#readme) |
| 官方宣称提供 150+ 可访问组件，目前仍为 Beta | 组件目录作为范围参考，按实际功能验收，不以名称数量证明完整度 | [仓库 README](https://github.com/meta/astryx#readme) |
| Core 使用 React 19+、StyleX 与预构建 CSS | 借鉴设计意图，运行依赖采用 Go 与 MyGo | [Core README](https://github.com/meta/astryx/blob/main/packages/core/README.md)、[包清单](https://github.com/meta/astryx/blob/main/packages/core/package.json) |
| 主题统一管理浅深色、字体、圆角、阴影、动效和组件变体 | 转为 Go 类型、MyGo 主题和原生元素样式 | [Neutral 主题源码](https://github.com/meta/astryx/blob/main/packages/themes/neutral/src/neutralTheme.ts) |
| 开放组合，支持源码导出 | 保留 `ui.Element` 和组合部件，让调用方继续定制 | [Core README](https://github.com/meta/astryx/blob/main/packages/core/README.md) |
| CLI 提供组件文档、模板、主题及迁移 | 建立统一组件元数据，为原生 Gallery、文档和后续检索工具复用 | [CLI README](https://github.com/meta/astryx/blob/main/packages/cli/README.md) |
| Button 处理焦点、图标按钮名称与减少动态效果 | 保留对应的原生交互和可访问性目标 | [Button 源码](https://github.com/meta/astryx/blob/main/packages/core/src/Button/Button.tsx) |

## 原生转换

| Web 概念 | MyGo 原生实现 |
| --- | --- |
| CSS variables | Go 语义化 Token 与 `ui.Theme` 映射 |
| React 组件与插槽 | Go 函数、`ui.Element`、`Children(func())` 和组合部件 |
| 状态选择器 | 根据控件状态设置元素样式 |
| DOM 焦点、ARIA 与输入事件 | MyGo Base 控件、原生角色、焦点和输入系统 |
| CSS 动效 | MyGo 的动画及转换 API，尊重系统偏好 |
| Storybook | Go 原生组件 Gallery 与状态预览 |
| npm 与 TypeScript 声明 | Go module、公开类型、GoDoc 与语义化版本 |

这是设计意图的映射，不保证两个平台逐项等价；生命周期、样式和浮层处理不能直接照搬。

## 实现边界

- 采用 MyGo 原生 UI，不依赖 Astryx npm 包，不 fork Astryx，不使用 WebView 承载组件。
- 复用 MyGo 的渲染、文字输入、键盘、焦点和无障碍基础。
- 参考设计体系后以 MyGo 公开 API 实现，不将现有控件重复导出描述为独立实现。
- 如果后续复制上游代码或实质性文档内容，应保留适用的许可和版权声明。Astryx 标为 MIT；字体、图标和图片需分别核查来源。[官方 LICENSE](https://github.com/meta/astryx/blob/main/LICENSE)

实现阶段另行固定来源为 [`e8ecde4ff82d424d76de8d6d8b122421a184c4a6`](https://github.com/meta/astryx/tree/e8ecde4ff82d424d76de8d6d8b122421a184c4a6)，逐项读取根入口及组件 barrel 导出。99 个组件族、组合部件与原生等价关系登记在 [覆盖清单](astryx-coverage.md)。没有运行 Astryx React 应用，验证对象是本项目的 MyGo 原生实现。
