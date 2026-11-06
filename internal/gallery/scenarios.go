package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
)

// Scenario records a visible state beyond a component's default appearance.
// Tests and the screenshot suite share these actual Gallery models/actions.
type Scenario struct {
	ID, Component, Expected string
	Prepare                 func(*App)
	Open                    func(*ui.Tester) error
}

func Scenarios() []Scenario {
	return []Scenario{
		{ID: "long-confirm", Component: "AlertDialog", Expected: "删除项目", Prepare: func(a *App) { a.Advanced.LongConfirm = true; a.ConfirmOpen = true }},
		{ID: "avatar-members", Component: "AvatarGroupOverflow", Expected: "协作成员名单", Prepare: func(a *App) { a.Advanced.MembersOpen = true }},
		{ID: "mobile-drawer", Component: "MobileNavToggle", Expected: "移动站点导航", Prepare: func(a *App) { a.Advanced.MobileOpen = true }},
		{ID: "sheet-review", Component: "BottomSheet", Expected: "核对资料", Prepare: func(a *App) { a.Advanced.SheetFlowOpen = true; a.Advanced.SheetStep = "review" }},
		{ID: "date-rule", Component: "DateRangeInput", Expected: "日期范围包含不可用日期", Prepare: func(a *App) { a.Advanced.ExcludeWeekends = true }},
		{ID: "multi-selector", Component: "MultiSelector", Expected: "完成", Open: func(tt *ui.Tester) error { return tt.Click("语言多选") }},
		{ID: "command-empty", Component: "CommandPalette", Expected: "没有匹配的命令", Prepare: func(a *App) { a.Advanced.PaletteOpen = true; a.Advanced.Palette.Query = "__no_command_matches__" }},
		{ID: "file-busy", Component: "FileInput", Expected: "文件处理进度", Prepare: func(a *App) { a.Advanced.Files.Busy = true; a.Advanced.Files.Progress = .55 }},
		{ID: "chat-newline", Component: "Chat", Expected: "Ctrl/⌘+Enter 发送", Prepare: func(a *App) { a.Advanced.ChatMultiline = true }},
		{ID: "mega-scroll", Component: "TopNavMegaMenuItem", Expected: "扩展示例 24", Prepare: func(a *App) { a.Advanced.LongMenu = true }, Open: func(tt *ui.Tester) error {
			r, ok := tt.Find("探索产品")
			if !ok {
				return fmt.Errorf("missing mega menu trigger")
			}
			bounds := tt.Image().Bounds()
			if r.H <= 0 || r.Y+r.H > float32(bounds.Dy()) {
				tt.Scroll(float32(bounds.Dx()-40), float32(bounds.Dy()-40), 0, 3000)
			}
			return tt.Click("探索产品")
		}},
	}
}
