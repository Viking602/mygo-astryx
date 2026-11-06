package astryx_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	kit "github.com/Viking602/mygo-astryx"
	"github.com/egoist/mygo/ui"
)

func TestLongAlertKeepsDecisionButtonsInWindow(t *testing.T) {
	open, confirmed := true, false
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.AlertDialog("确认长说明", strings.Repeat("请阅读这段说明，再决定是否继续。\n", 60), &open,
			kit.AlertDialogOptions{Confirm: "执行", Cancel: "保留", Destructive: true, OnConfirm: func() { confirmed = true }})
	}), 360, 340)
	for _, label := range []string{"执行", "保留"} {
		r, ok := tt.Find(label)
		if !ok || r.Y < 0 || r.Y+r.H > 340 || r.X < 0 || r.X+r.W > 360 {
			t.Fatalf("decision %q is outside viewport: %+v", label, r)
		}
	}
	click(t, tt, "执行")
	if open || !confirmed {
		t.Fatal("long decision cannot be completed")
	}
}

func TestDesktopMegaMenuScrollsLongContent(t *testing.T) {
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.TopNavMegaMenu("产品菜单", func() {
			for i := range 40 {
				k.Button(fmt.Sprintf("产品 %d", i))
			}
		}, nil)
	}), 900, 400)
	click(t, tt, "产品菜单")
	r, _ := tt.Find("产品 0")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 2000)
	r, ok := tt.Find("产品 39")
	if !ok || r.Y < 0 || r.Y+r.H > 400 {
		t.Fatalf("last mega menu item is unreachable: %+v", r)
	}
	click(t, tt, "产品 39")
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("产品 39") {
		t.Fatal("long mega menu did not dismiss")
	}
}

func TestMobileNavigationExposesChildrenAndSkipsDisabled(t *testing.T) {
	selected := "home"
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.MobileNav(&selected, []kit.NavItem{{ID: "home", Label: "首页"}, {ID: "work", Label: "空间", Children: []kit.NavItem{{ID: "locked", Label: "不可用", Disabled: true}, {ID: "project", Label: "原生项目"}}}})
	}), 420, 400)
	click(t, tt, "空间")
	if !tt.HasText("原生项目") || selected != "home" {
		t.Fatal("mobile parent selected itself instead of opening children")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if selected != "project" {
		t.Fatalf("mobile submenu selection=%q", selected)
	}
}

func TestDateRangeAllowedIsDateOnlyInOriginalLocation(t *testing.T) {
	zone := time.FixedZone("user-zone", 8*3600)
	start := time.Date(2026, 11, 6, 15, 42, 0, 0, zone)
	err := kit.ValidateDateRange(kit.DateRange{Start: start, End: start.AddDate(0, 0, 2)}, kit.DateRangeOptions{Allowed: func(d time.Time) bool { return d.Hour() == 0 && d.Location() == zone }})
	if err != "" {
		t.Fatalf("date-only predicate received time-of-day or a different region: %s", err)
	}
}

func TestEmptyDateRangeDoesNotCommitDisallowedDefault(t *testing.T) {
	s, changes := kit.DateRange{Empty: true}, 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.DateRangeInput(&s, kit.DateRangeOptions{Allowed: func(time.Time) bool { return false }, OnChange: func(kit.DateRange) { changes++ }})
	}), 600, 400)
	click(t, tt, "选择日期范围")
	if !s.Empty || changes != 0 {
		t.Fatal("unavailable default range changed the controlled value")
	}
	if !tt.HasText("当前日期范围不可用，请选择其他日期") {
		t.Fatal("unavailable default has no visible explanation")
	}
}

func TestPowerSearchReadOnlyCannotClearConditions(t *testing.T) {
	s := kit.PowerSearchState{Filters: []kit.SearchFilter{{ID: "one", Field: "name", Operator: "=", Value: "Go"}}}
	changes := 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.PowerSearch(&s, []kit.FilterField{{ID: "name", Label: "名称"}}, kit.PowerSearchOptions{ReadOnly: true, OnChange: func([]kit.SearchFilter) { changes++ }})
	}), 600, 400)
	if _, ok := tt.Find("清除可编辑条件"); ok {
		t.Fatal("read-only search exposes a clear action")
	}
	if len(s.Filters) != 1 || changes != 0 {
		t.Fatal("read-only search changed conditions")
	}
}

func TestAvatarPresenceAndInteractiveOverflow(t *testing.T) {
	count := 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Avatar("Ada", nil, kit.AvatarOptions{Size: 64, Status: "在线", StatusTone: kit.Success})
		k.AvatarStatusDot("离线", kit.Neutral)
		k.AvatarGroup([]kit.Person{{Name: "甲"}, {Name: "乙"}, {Name: "丙"}}, 1, kit.AvatarGroupOptions{OverflowLabel: "更多", OnOverflow: func() { count++ }})
	}), 500, 400)
	r, ok := tt.Find("Ada，在线")
	if !ok || r.W != 64 || r.H != 64 {
		t.Fatalf("avatar did not expose composed status/size: %+v", r)
	}
	if !tt.HasText("更多") {
		t.Fatal("custom overflow text was not rendered")
	}
	click(t, tt, "另外 2 人")
	tt.Key(0, ui.KeyEnter)
	if count != 2 {
		t.Fatalf("overflow is not keyboard actionable: %d", count)
	}
}

func TestMobileNavToggleSharesDrawerDismissalState(t *testing.T) {
	open, disabled := false, false
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.MobileNavToggle("导航菜单", &open, kit.ButtonOptions{Disabled: disabled})
		k.Drawer("移动页面", &open, func() { k.Text("页面内容") })
	}), 420, 500)
	click(t, tt, "导航菜单")
	if !open || !tt.HasText("页面内容") {
		t.Fatal("mobile toggle did not open its navigation surface")
	}
	tt.Key(0, ui.KeyEscape)
	if open || !tt.Focused("导航菜单") {
		t.Fatal("Escape did not update toggle state and restore focus")
	}
	disabled = true
	tt.Frame()
	click(t, tt, "导航菜单")
	if open {
		t.Fatal("disabled mobile toggle opened navigation")
	}
}

func TestTokenOnlyChatPreservesContextAndMessageCallbacks(t *testing.T) {
	s := kit.ChatState{Tokens: []kit.ChatToken{{ID: "doc", Label: "README.md", Kind: "文档"}}}
	retries := 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Chat(&s, kit.ChatOptions{MessageOptions: kit.ChatMessageOptions{OnRetry: func(string) { retries++ }}}).Height(420)
	}), 600, 500)
	click(t, tt, "发送消息")
	if len(s.Messages) != 1 || len(s.Messages[0].Tokens) != 1 || s.Messages[0].Tokens[0].Label != "README.md" {
		t.Fatal("token-only message lost its context")
	}
	if len(s.Tokens) != 0 || !tt.HasText("README.md") {
		t.Fatal("sent context did not move from composer to message")
	}
	s.Messages[0].Status = "error"
	tt.Frame()
	click(t, tt, "重试")
	if retries != 1 {
		t.Fatal("high-level Chat did not forward message callbacks")
	}
}

func TestDatePresetsRemainAvailableWhenValueIsEmpty(t *testing.T) {
	s := kit.DateRange{Empty: true}
	day := time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC)
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.DateRangeInput(&s, kit.DateRangeOptions{Presets: []kit.DateRangePreset{{Label: "固定排期", Start: day, End: day.AddDate(0, 0, 1)}}})
	}), 600, 400)
	click(t, tt, "固定排期")
	if s.Empty || !s.Start.Equal(day) {
		t.Fatal("empty range could not choose a preset")
	}
}

func TestDialogDescriptionsScrollWithoutHidingFooter(t *testing.T) {
	for _, sheet := range []bool{false, true} {
		open, done := true, false
		tt := ui.NewTester(view(func(k *kit.UI) {
			o := kit.DialogOptions{Description: strings.Repeat("完整的说明文字需要在正文区域滚动。\n", 60), Footer: func() { k.Button("完成阅读").OnClick(func() { done = true; open = false }) }}
			if sheet {
				k.BottomSheet("长底部面板", &open, func() { k.Text("正文结束") }, o)
			} else {
				k.Dialog("长对话框", &open, func() { k.Text("正文结束") }, o)
			}
		}), 360, 360)
		r, ok := tt.Find("完成阅读")
		if !ok || r.H <= 0 || r.Y < 0 || r.Y+r.H > 360 {
			t.Fatalf("footer hidden (sheet=%v): %+v", sheet, r)
		}
		click(t, tt, "完成阅读")
		if open || !done {
			t.Fatal("footer action did not complete")
		}
	}
}

func TestMegaMenuSelectionDismissesAndReturnsFocus(t *testing.T) {
	for _, width := range []int{900, 420} {
		calls := 0
		tt := ui.NewTester(view(func(k *kit.UI) {
			k.TopNavMegaMenu("应用目录", func() { k.TopNavMegaMenuItem("打开原生项目", kit.ItemOptions{OnClick: func() { calls++ }}) }, nil)
		}), width, 500)
		click(t, tt, "应用目录")
		click(t, tt, "打开原生项目")
		tt.Advance(time.Second)
		if calls != 1 || tt.HasText("打开原生项目") || !tt.Focused("应用目录") {
			t.Fatalf("mega selection did not dismiss/restore focus (width %d)", width)
		}
	}
}

func TestMegaMenuKeyboardEntersContentAndEscapeRestoresTrigger(t *testing.T) {
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.TopNavMegaMenu("键盘目录", func() { k.TopNavMegaMenuItem("第一项目"); k.TopNavMegaMenuItem("第二项目") }, nil)
	}), 900, 500)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("第一项目") {
		t.Fatalf("keyboard opening did not enter menu content: trigger=%v first=%v second=%v", tt.Focused("键盘目录"), tt.Focused("第一项目"), tt.Focused("第二项目"))
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if !tt.Focused("键盘目录") || tt.HasText("第一项目") {
		t.Fatal("Escape did not restore trigger or the menu reopened")
	}
}

func TestCustomPopupsKeepLongContentScrollable(t *testing.T) {
	for _, custom := range []bool{false, true} {
		open, calls := false, 0
		tt := ui.NewTester(view(func(k *kit.UI) {
			body := func() {
				for i := range 35 {
					k.Button(fmt.Sprintf("面板项 %d", i)).OnClick(func() { calls++ })
				}
			}
			if custom {
				k.ComplexSelector("定制面板", &open, "选择", func(func()) { body() })
			} else {
				anchor := k.Button("普通面板").OnClick(func() { open = !open })
				k.Popover(anchor, &open, body)
			}
		}), 600, 420)
		label := "普通面板"
		if custom {
			label = "定制面板"
		}
		click(t, tt, label)
		first, _ := tt.Find("面板项 0")
		tt.Scroll(first.X+5, first.Y+5, 0, 2000)
		last, ok := tt.Find("面板项 34")
		if !ok || last.H <= 0 || last.Y < 0 || last.Y+last.H > 420 {
			t.Fatalf("popup content is unreachable: %+v", last)
		}
		click(t, tt, "面板项 34")
		if calls != 1 {
			t.Fatal("scrolled popup action did not run")
		}
	}
}

func TestScrollBothAxesReachesContentCorner(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		k := kit.New(c)
		k.Scroll(kit.ScrollOptions{Both: true}).Size(280, 160).Children(func() {
			ui.Box(c).Size(900, 600).Children(func() { k.Text("双轴末尾").Absolute().End(8).Bottom(8) })
		})
	}, 500, 400)
	tt.Scroll(10, 10, 2000, 2000)
	r, ok := tt.Find("双轴末尾")
	if !ok || r.H <= 0 || r.X < 0 || r.Y < 0 || r.X+r.W > 280 || r.Y+r.H > 160 {
		t.Fatalf("both-axis scroll did not reach corner: %+v", r)
	}
}

func TestLinkBoundsStayOnClickableTextInStretchedColumn(t *testing.T) {
	tt := ui.NewTester(view(func(k *kit.UI) { k.Link("查看组件源码", "https://example.com/source") }), 900, 300)
	click(t, tt, "查看组件源码")
	if urls := tt.OpenedURLs(); len(urls) != 1 || urls[0] != "https://example.com/source" {
		t.Fatalf("link center did not activate its text: %v", urls)
	}
}

func TestMobileMegaMenuReachesItemsBeforeFeaturedPanel(t *testing.T) {
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.TopNavMegaMenu("完整移动菜单", func() {
			for i := range 30 {
				k.TopNavMegaMenuItem(fmt.Sprintf("移动链接 %d", i))
			}
		}, func() { k.Text("最后的推荐内容") })
	}), 420, 600)
	click(t, tt, "完整移动菜单")
	r, _ := tt.Find("移动链接 0")
	tt.Scroll(r.X+5, r.Y+5, 0, 3000)
	r, ok := tt.Find("移动链接 29")
	if !ok || r.H <= 0 || r.Y < 0 || r.Y+r.H > 600 {
		t.Fatalf("last mobile item cannot be reached: %+v", r)
	}
}

func TestDesignLanguageConsistencyInStretchedColumn(t *testing.T) {
	segment, number := 0, 8.0
	var multi []string
	choice, open := "a", false
	options := []kit.Option{{Value: "a", Label: "选项甲"}, {Value: "b", Label: "选项乙"}}
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Segmented(&segment, []string{"日", "周", "月"})
		k.Token("独立标签", nil)
		k.Citation("[1] 规范", kit.CitationOptions{URL: "https://example.com"})
		k.NumberInput(&number, 0, 20, 1, kit.InputOptions{Label: "步进数量"})
		k.Select(&choice, options, kit.InputOptions{Label: "单选触发器"}).FillWidth()
		k.MultiSelector(&multi, options, kit.SelectionOptions{Label: "多选触发器"}).FillWidth()
		k.ComplexSelector("定制触发器", &open, "已选方案", nil).FillWidth()
	}), 900, 500)
	for _, compact := range []string{"分段选择", "独立标签", "[1] 规范"} {
		r, ok := tt.Find(compact)
		if !ok || r.W <= 0 || r.W > 320 {
			t.Fatalf("%s stretched across column instead of hugging content: %+v", compact, r)
		}
	}
	s1, ok1 := tt.Find("单选触发器")
	s2, ok2 := tt.Find("多选触发器")
	s3, ok3 := tt.Find("定制触发器")
	if !ok1 || !ok2 || !ok3 || s1.W < 800 || s1.W != s2.W || s1.W != s3.W || s1.H != s2.H || s1.H != s3.H {
		t.Fatalf("selector triggers do not share unified full-width geometry: select=%+v multi=%+v complex=%+v", s1, s2, s3)
	}
	click(t, tt, "增加")
	if number != 9 {
		t.Fatalf("number input increment failed: %v", number)
	}
}
