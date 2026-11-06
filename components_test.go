package astryx_test

import (
	kit "github.com/Viking602/mygo-astryx"
	"github.com/Viking602/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
	"testing"
	"time"
)

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}
func view(build func(*kit.UI)) func(*ui.Context) {
	return func(c *ui.Context) {
		kit.Scope(c, theme.Neutral(false, c.Preferences()), func(k *kit.UI) { k.Stack().Padding(16).FillWidth().Children(func() { build(k) }) }).Fill()
	}
}

func TestButtonActionsAndDisabledStates(t *testing.T) {
	count := 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Key("save").Button("保存", kit.ButtonOptions{Variant: kit.Primary}).OnClick(func() { count++ })
		k.Key("disabled").Button("禁用", kit.ButtonOptions{Disabled: true}).OnClick(func() { count += 100 })
		k.Key("loading").Button("加载", kit.ButtonOptions{Loading: true}).OnClick(func() { count += 100 })
	}), 400, 250)
	click(t, tt, "保存")
	if count != 1 {
		t.Fatalf("click count=%d", count)
	}
	tt.Frame()
	if count != 1 {
		t.Fatalf("action repeated during rebuild: %d", count)
	}
	tt.Key(0, ui.KeyEnter)
	if count != 2 {
		t.Fatalf("Enter count=%d", count)
	}
	click(t, tt, "禁用")
	click(t, tt, "加载")
	if count != 2 {
		t.Fatalf("disabled button acted: %d", count)
	}
}

func TestChoiceKeyboardAndPointer(t *testing.T) {
	checked, on := false, false
	selected := "one"
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Key("check").Checkbox(&checked, "接收通知")
		k.Key("switch").Switch(&on, "同步")
		k.RadioGroup("类型", func() { k.Key("one").Radio(&selected, "one", "个人"); k.Key("two").Radio(&selected, "two", "团队") })
	}), 400, 280)
	click(t, tt, "接收通知")
	if !checked {
		t.Fatal("checkbox did not update")
	}
	tt.Key(0, ui.KeySpace)
	if checked {
		t.Fatal("Space did not toggle checkbox")
	}
	click(t, tt, "同步")
	if !on {
		t.Fatal("switch did not update")
	}
	click(t, tt, "团队")
	if selected != "two" {
		t.Fatal("radio did not update")
	}
	tt.Key(0, ui.KeyUp)
	if selected != "one" {
		t.Fatalf("radio arrow selection=%s", selected)
	}
}

func TestNativeTextCompositionAndReadOnly(t *testing.T) {
	text, locked := "", "不可修改"
	composing := false
	tt := ui.NewTester(view(func(k *kit.UI) {
		in := k.Key("input").TextInput(&text, kit.InputOptions{Label: "名称"})
		composing = in.Composing()
		k.Key("locked").TextInput(&locked, kit.InputOptions{Label: "只读", ReadOnly: true})
	}), 400, 200)
	click(t, tt, "名称")
	tt.Compose("ni", 2)
	if !composing || text != "" {
		t.Fatalf("composition leaked into model: %q, composing=%v", text, composing)
	}
	tt.Type("你")
	if composing || text != "你" {
		t.Fatalf("commit=%q composing=%v", text, composing)
	}
	tt.Type("好")
	tt.Command("selectAll")
	tt.Command("copy")
	if tt.Clipboard() != "你好" {
		t.Fatalf("clipboard=%q", tt.Clipboard())
	}
	click(t, tt, "只读")
	tt.Type("x")
	if locked != "不可修改" {
		t.Fatal("read-only input changed")
	}
}

func TestSelectSkipsDisabledOptionsAndCommitsDisplayValue(t *testing.T) {
	selected := "one"
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Key("select").Select(&selected, []kit.Option{{Value: "one", Label: "个人"}, {Value: "blocked", Label: "不可用", Disabled: true}, {Value: "two", Label: "团队"}}, kit.InputOptions{Label: "方案"})
	}), 400, 300)
	click(t, tt, "方案")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if selected != "two" {
		t.Fatalf("expected keyboard to skip blocked option, got %q", selected)
	}
	if !tt.HasText("团队") {
		t.Fatalf("selected label was not rebuilt: %v", tt.Texts())
	}
}

func TestComboboxChoosesFilteredSuggestion(t *testing.T) {
	value := ""
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Key("combo").Combobox(&value, []string{"Go", "Rust", "Swift"}, kit.InputOptions{Label: "语言"})
	}), 400, 300)
	click(t, tt, "语言")
	tt.Type("Ru")
	tt.Key(0, ui.KeyEnter)
	if value != "Rust" {
		t.Fatalf("chosen suggestion not committed: %q", value)
	}
}

func TestNestedOverlaysEscapeAndRestoreFocus(t *testing.T) {
	open := false
	choice := "one"
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Button("打开").OnClick(func() { open = true })
		k.Dialog("编辑", &open, func() {
			k.Key("dialog-choice").Select(&choice, []kit.Option{{Value: "one", Label: "个人"}, {Value: "two", Label: "团队"}}, kit.InputOptions{Label: "弹窗方案"})
		})
	}), 640, 480)
	click(t, tt, "打开")
	if !open {
		t.Fatal("dialog did not open")
	}
	click(t, tt, "弹窗方案")
	if !tt.HasText("团队") {
		t.Fatal("nested popup did not open")
	}
	tt.Key(0, ui.KeyEscape)
	if !open {
		t.Fatal("Escape closed dialog before its popup")
	}
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape did not close dialog")
	}
	if !tt.Focused("打开") {
		t.Fatal("dialog did not restore trigger focus")
	}
}

func TestKeysPreserveInputIdentityAfterReorder(t *testing.T) {
	items := []struct{ ID, Text string }{{"alpha", "A"}, {"beta", "B"}}
	order := []int{0, 1}
	tt := ui.NewTester(view(func(k *kit.UI) {
		for _, i := range order {
			k.Key(items[i].ID).TextInput(&items[i].Text, kit.InputOptions{Label: items[i].ID})
		}
	}), 400, 200)
	click(t, tt, "beta")
	order = []int{1, 0}
	tt.Frame()
	if !tt.Focused("beta") {
		t.Fatal("focus did not follow the key after reorder")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Type("!")
	if items[0].Text != "A" || items[1].Text != "B!" {
		t.Fatalf("editing followed position instead of key: %+v", items)
	}
}

func TestToastActionAndExpiry(t *testing.T) {
	undone := 0
	tt := ui.NewTester(func(c *ui.Context) {
		kit.Scope(c, theme.Neutral(false, c.Preferences()), func(k *kit.UI) {
			k.Button("归档").OnClick(func() {
				k.Toast(ui.Toast{Title: "已归档", Action: "撤销", OnAction: func() { undone++ }, Timeout: time.Second})
			})
			k.ToastViewport()
		}).Fill()
	}, 640, 480)
	click(t, tt, "归档")
	if !tt.HasText("已归档") {
		t.Fatal("toast not visible")
	}
	click(t, tt, "撤销")
	if undone != 1 {
		t.Fatalf("undo calls=%d", undone)
	}
	tt.Advance(2 * time.Second)
	if tt.HasText("已归档") {
		t.Fatal("toast did not leave")
	}
}

func TestTabsPaginationAndDisclosure(t *testing.T) {
	tab, page, open := 0, 1, false
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Key("tabs").Tabs(&tab, []string{"概览", "活动"})
		k.Key("pagination").Pagination(&page, 3)
		k.Key("disclosure").Collapsible("更多内容", &open, func() { k.Text("展开内容") })
	}), 500, 300)
	click(t, tt, "活动")
	if tab != 1 {
		t.Fatal("tab did not change")
	}
	click(t, tt, "下一页")
	if page != 2 {
		t.Fatal("page did not change")
	}
	click(t, tt, "更多内容")
	if !open || !tt.HasText("展开内容") {
		t.Fatal("disclosure did not expand")
	}
}

func TestDisabledSelectDoesNotReopenWhenEnabled(t *testing.T) {
	value, disabled := "one", false
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Key("select").Select(&value, []kit.Option{{Value: "one", Label: "个人"}, {Value: "two", Label: "团队"}}, kit.InputOptions{Label: "方案", Disabled: disabled})
	}), 400, 300)
	click(t, tt, "方案")
	if !tt.HasText("团队") {
		t.Fatal("popup did not open")
	}
	disabled = true
	tt.Frame()
	if tt.HasText("团队") {
		t.Fatal("disabled popup still visible")
	}
	disabled = false
	tt.Frame()
	if tt.HasText("团队") {
		t.Fatal("popup reopened without user input")
	}
	click(t, tt, "方案")
	click(t, tt, "团队")
	if value != "two" {
		t.Fatal("reenabled select could not choose")
	}
}

func TestAlertDialogRequiresActionAndDestructiveDefaultsToCancel(t *testing.T) {
	open, count := false, 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Button("删除入口").OnClick(func() { open = true })
		k.AlertDialog("确定删除？", "不可撤销。", &open, kit.AlertDialogOptions{Confirm: "确认删除", Cancel: "保留", Destructive: true, OnConfirm: func() { count++ }})
	}), 640, 480)
	click(t, tt, "删除入口")
	tt.ClickAt(2, 2)
	if !open {
		t.Fatal("outside click dismissed alert")
	}
	if !tt.Focused("保留") {
		t.Fatal("destructive alert did not focus its cancel action")
	}
	tt.Key(0, ui.KeyEnter)
	if open || count != 0 {
		t.Fatal("Enter performed destructive action by default")
	}
	click(t, tt, "删除入口")
	click(t, tt, "确认删除")
	if open || count != 1 {
		t.Fatalf("confirm open=%v actions=%d", open, count)
	}
}
