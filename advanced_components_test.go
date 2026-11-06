package astryx_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	kit "mygo-astryx"
)

func TestLongSelectScrollsItsOptionsAndChoosesLast(t *testing.T) {
	value := "0"
	items := []kit.Option{}
	for i := 0; i < 100; i++ {
		items = append(items, kit.Option{Value: fmt.Sprint(i), Label: fmt.Sprintf("选项 %d", i)})
	}
	tt := ui.NewTester(view(func(k *kit.UI) { k.Select(&value, items, kit.InputOptions{Label: "长列表"}) }), 480, 600)
	click(t, tt, "长列表")
	tt.Key(0, ui.KeyEnd)
	r, ok := tt.Find("选项 99")
	if !ok || r.Y < 0 || r.Y+r.H > 600 {
		t.Fatalf("last option was not scrolled into view: %+v", r)
	}
	tt.Key(0, ui.KeyEnter)
	if value != "99" || tt.HasText("选项 98") {
		t.Fatal("last option did not commit and dismiss")
	}
}

func TestMenusCommitCheckedRadioAndNestedActions(t *testing.T) {
	checked, radio, calls := false, "a", 0
	entries := []kit.MenuEntry{
		{Label: "勾选项", Checked: &checked},
		{Label: "方案乙", Radio: &radio, Value: "b"},
		{Label: "禁用项", Disabled: true, OnSelect: func() { calls += 100 }},
		{Label: "子菜单", Children: []kit.MenuEntry{{Label: "子动作", OnSelect: func() { calls++ }}}},
	}
	tt := ui.NewTester(view(func(k *kit.UI) { k.DropdownMenu("选项菜单", entries) }), 600, 500)
	click(t, tt, "选项菜单")
	click(t, tt, "勾选项")
	if !checked {
		t.Fatal("checked menu binding did not change")
	}
	click(t, tt, "选项菜单")
	click(t, tt, "方案乙")
	if radio != "b" {
		t.Fatal("radio menu binding did not change")
	}
	click(t, tt, "选项菜单")
	click(t, tt, "禁用项")
	if calls != 0 {
		t.Fatal("disabled menu action fired")
	}
	tt.Key(0, ui.KeyEscape)
	click(t, tt, "选项菜单")
	click(t, tt, "子菜单")
	click(t, tt, "子动作")
	tt.Frame()
	if calls != 1 {
		t.Fatalf("nested action ran %d times", calls)
	}
	// The system context-menu adapter uses the same controlled item model.
	tt = ui.NewTester(view(func(k *kit.UI) {
		k.ContextMenu(k.Button("右键目标"), entries, kit.ContextMenuOptions{Native: true})
	}), 400, 250)
	if err := tt.RightClick("右键目标"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("子菜单", "子动作"); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatal("native submenu did not invoke action")
	}
}

func TestCommandPaletteFiltersSkipsDisabledAndDoesNotRepeat(t *testing.T) {
	open, calls := true, 0
	s := kit.CommandPaletteState{}
	commands := []kit.Command{{ID: "a", Label: "Alpha"}, {ID: "x", Label: "Blocked", Disabled: true}, {ID: "b", Label: "Beta", Description: "project", Run: func() { calls++ }}}
	tt := ui.NewTester(view(func(k *kit.UI) { k.CommandPalette(&open, &s, commands) }), 640, 520)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if open || s.Chosen != "b" || calls != 1 {
		t.Fatalf("keyboard command: open=%v chosen=%q calls=%d", open, s.Chosen, calls)
	}
	tt.Frame()
	if calls != 1 {
		t.Fatal("command repeated on rebuild")
	}
	open, s.Query = true, ""
	tt.Frame()
	click(t, tt, "搜索命令")
	tt.Type("project")
	if !tt.HasText("Beta") || tt.HasText("Alpha") {
		t.Fatal("command filtering failed")
	}
	tt.Compose("ni", 2)
	tt.Key(0, ui.KeyEnter)
	if calls != 1 || !open {
		t.Fatal("IME confirmation executed a command")
	}
	tt.Compose("", 0)
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape did not dismiss command palette")
	}
}

func TestMultiSelectorLimitClearAndDisabledLifecycle(t *testing.T) {
	var values []string
	disabled, changes := false, 0
	items := []kit.Option{{Value: "a", Label: "设计"}, {Value: "x", Label: "受限", Disabled: true}, {Value: "b", Label: "工程"}}
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.MultiSelector(&values, items, kit.SelectionOptions{Label: "团队选择", Disabled: disabled, MaxSelected: 1, OnChange: func(v []string) { changes++ }})
	}), 480, 420)
	click(t, tt, "团队选择")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(values, []string{"a"}) {
		t.Fatalf("keyboard selection=%v", values)
	}
	click(t, tt, "工程")
	click(t, tt, "受限")
	if len(values) != 1 || changes != 1 {
		t.Fatal("selection limits were ignored")
	}
	click(t, tt, "清空")
	if len(values) != 0 || changes != 2 {
		t.Fatal("clear did not notify")
	}
	disabled = true
	tt.Frame()
	if tt.HasText("完成") {
		t.Fatal("disabled multi-selector kept its popup")
	}
	disabled = false
	tt.Frame()
	if tt.HasText("完成") {
		t.Fatal("multi-selector reopened without input")
	}
}

func TestTypeaheadCommitsStableIDAndCustomSelectorCloses(t *testing.T) {
	s, open, count := kit.TypeaheadState{}, false, 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Typeahead(&s, []kit.Suggestion{{ID: "blocked", Label: "Rust locked", Disabled: true}, {ID: "lang-rust", Label: "Rust", Description: "systems"}}, kit.TypeaheadOptions{InputOptions: kit.InputOptions{Label: "建议"}, OnSelect: func(kit.Suggestion) { count++ }})
		k.ComplexSelector("自定义选择", &open, "配置", func(close func()) { k.Button("应用配置").OnClick(close) })
	}), 500, 400)
	click(t, tt, "建议")
	tt.Type("Rust")
	tt.Key(0, ui.KeyEnter)
	if s.Value != "lang-rust" || s.Query != "Rust" || count != 1 {
		t.Fatalf("typeahead result=%+v count=%d", s, count)
	}
	click(t, tt, "自定义选择")
	click(t, tt, "应用配置")
	if open {
		t.Fatal("custom selector close callback did not close")
	}
}

func TestPowerSearchValidatesEditsAndPreservesReadOnlyFilters(t *testing.T) {
	s := kit.PowerSearchState{Field: "amount", Value: "NaN", Filters: []kit.SearchFilter{{ID: "fixed", Field: "amount", Operator: ">", Value: "0", ReadOnly: true}}}
	changes, readonly := 0, false
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.PowerSearch(&s, []kit.FilterField{{ID: "amount", Label: "数量", Kind: "number"}}, kit.PowerSearchOptions{ReadOnly: readonly, OnChange: func([]kit.SearchFilter) { changes++ }})
	}), 650, 500)
	click(t, tt, "添加条件")
	if s.Error == "" || len(s.Filters) != 1 || changes != 0 {
		t.Fatal("invalid numeric filter was accepted")
	}
	click(t, tt, "筛选值")
	tt.Command("selectAll")
	tt.Type("12")
	click(t, tt, "添加条件")
	if len(s.Filters) != 2 || changes != 1 {
		t.Fatal("valid filter was not added")
	}
	id := s.Filters[1].ID
	click(t, tt, "编辑条件"+id)
	click(t, tt, "筛选值")
	tt.Command("selectAll")
	tt.Type("24")
	click(t, tt, "更新条件")
	if s.Filters[1].Value != "24" || changes != 2 {
		t.Fatal("editing did not update the existing filter")
	}
	click(t, tt, "清除可编辑条件")
	if len(s.Filters) != 1 || s.Filters[0].ID != "fixed" {
		t.Fatal("clear removed a read-only filter")
	}
	readonly = true
	tt.Frame()
	if _, ok := tt.Find("添加条件"); ok {
		t.Fatal("read-only search still has an editor")
	}
}

func TestFileInputBrowseDropValidationAndRemoval(t *testing.T) {
	dir := t.TempDir()
	a, b, bad := filepath.Join(dir, "one.txt"), filepath.Join(dir, "two.txt"), filepath.Join(dir, "wrong.png")
	for _, p := range []string{a, b, bad} {
		if err := os.WriteFile(p, []byte("content"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	s, changes := kit.FileInputState{}, 0
	o := kit.FileInputOptions{Label: "附件", Accept: ".txt", Multiple: true, MaxFiles: 2, MaxSize: 20, Browse: func() ([]string, error) { return []string{a, bad}, nil }, OnChange: func([]kit.UploadFile) { changes++ }}
	tt := ui.NewTester(view(func(k *kit.UI) { k.FileInput(&s, o) }), 500, 420)
	click(t, tt, "浏览文件")
	if len(s.Files) != 1 || s.Error == "" || changes != 1 {
		t.Fatalf("browse result=%+v changes=%d", s, changes)
	}
	r, _ := tt.Find("附件")
	if !tt.DropFiles(r.X+8, r.Y+8, b) {
		t.Fatal("file drop was not consumed")
	}
	if len(s.Files) != 2 || changes != 2 {
		t.Fatal("valid drop was not appended")
	}
	if kit.AddUploadPaths(&s, []string{a}, o) || changes != 2 {
		t.Fatal("duplicate file changed selection")
	}
	click(t, tt, "移除文件one.txt")
	if len(s.Files) != 1 || changes != 3 {
		t.Fatal("remove did not notify")
	}
	s.Busy = true
	tt.Frame()
	click(t, tt, "浏览文件")
	if kit.AddUploadPaths(&s, []string{a}, o) || len(s.Files) != 1 {
		t.Fatal("busy input accepted files")
	}
	s.Busy = false
	o.MaxSize = 2
	if kit.AddUploadPaths(&s, []string{a}, o) || !strings.Contains(s.Error, "过大") {
		t.Fatal("file size limit was not applied")
	}
}

func TestDateRangeValidatesCalendarDaysAndPresets(t *testing.T) {
	zone, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(2026, 3, 7, 23, 0, 0, 0, zone)
	rangeValue := kit.DateRange{Start: start, End: start.AddDate(0, 0, 2)}
	if kit.ValidateDateRange(rangeValue, kit.DateRangeOptions{MaxDays: 3}) != "" {
		t.Fatal("DST affected inclusive calendar-day count")
	}
	if kit.ValidateDateRange(rangeValue, kit.DateRangeOptions{MaxDays: 2}) == "" {
		t.Fatal("inclusive maximum was not enforced")
	}
	if kit.ValidateDateRange(rangeValue, kit.DateRangeOptions{Allowed: func(d time.Time) bool { return d.Day() != 8 }}) == "" {
		t.Fatal("disallowed day was accepted")
	}
	s, changes := kit.DateRange{Empty: true}, 0
	min := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.DateRangeInput(&s, kit.DateRangeOptions{Label: "出行", Min: min, Max: min.AddDate(0, 0, 2), MaxDays: 2, OnChange: func(kit.DateRange) { changes++ }, Presets: []kit.DateRangePreset{{Label: "一天", Start: min, End: min}}})
	}), 650, 350)
	click(t, tt, "选择日期范围")
	if s.Start.Before(min) || kit.ValidateDateRange(s, kit.DateRangeOptions{Min: min, Max: min.AddDate(0, 0, 2), MaxDays: 2}) != "" {
		t.Fatal("initial range violated constraints")
	}
	click(t, tt, "一天")
	if !s.Start.Equal(s.End) || changes != 2 {
		t.Fatal("preset did not commit")
	}
	click(t, tt, "清除出行")
	if !s.Empty || changes != 3 {
		t.Fatal("range clear did not notify")
	}
}

func TestHoverCardDelayClickDismissalAndEscape(t *testing.T) {
	tt := ui.NewTester(view(func(k *kit.UI) {
		anchor := k.Button("资料")
		k.HoverCard(anchor, func() { k.Text("账户详情") }, kit.HoverCardOptions{Delay: 300 * time.Millisecond, OpenOnClick: true})
	}), 600, 400)
	r, _ := tt.Find("资料")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(200 * time.Millisecond)
	if tt.HasText("账户详情") {
		t.Fatal("hover card opened before delay")
	}
	tt.Advance(110 * time.Millisecond)
	if !tt.HasText("账户详情") {
		t.Fatal("hover card did not open")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if tt.HasText("账户详情") {
		t.Fatal("dismissed hover card reopened under stationary pointer")
	}
	tt.Move(590, 390)
	tt.Advance(time.Second)
	click(t, tt, "资料")
	if !tt.HasText("账户详情") {
		t.Fatal("click did not open hover card")
	}
	click(t, tt, "资料")
	tt.Advance(time.Second)
	if tt.HasText("账户详情") {
		t.Fatal("click-to-close immediately reopened hover card")
	}
}

func TestCarouselAndLightboxActions(t *testing.T) {
	index, image, open := 0, 0, false
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Carousel(&index, []kit.CarouselSlide{{ID: "a", Label: "甲", Content: func() { k.Text("内容甲") }}, {ID: "b", Label: "乙", Content: func() { k.Text("内容乙") }}}, kit.CarouselOptions{Loop: true, Height: 60})
		k.Button("预览").OnClick(func() { open = true })
		k.Lightbox(&open, &image, []kit.ImageAsset{{ID: "a", Label: "图片甲"}, {ID: "b", Label: "图片乙"}})
	}), 700, 600)
	click(t, tt, "上一项")
	if index != 1 || !tt.HasText("内容乙") {
		t.Fatal("carousel loop failed")
	}
	click(t, tt, "预览")
	click(t, tt, "放大")
	if !tt.HasText("1 / 2 · 图片甲 · 125%") {
		t.Fatal("zoom did not update")
	}
	tt.Key(0, ui.KeyRight)
	if image != 1 || !tt.HasText("2 / 2 · 图片乙 · 100%") {
		t.Fatal("lightbox navigation did not reset zoom")
	}
	click(t, tt, "旋转")
	tt.Key(0, ui.KeyEscape)
	if open || !tt.Focused("预览") {
		t.Fatal("lightbox did not dismiss and restore focus")
	}
}

func TestTimerPauseResumeAndReset(t *testing.T) {
	s := kit.TimerState{}
	tt := ui.NewTester(func(c *ui.Context) {
		k := kit.New(c)
		now := c.Now()
		k.Stack().Children(func() {
			k.Timer(&s, kit.TimerOptions{Clock: true})
			k.Button("开始").OnClick(func() { s.Start(now) })
			k.Button("暂停").OnClick(func() { s.Pause(now) })
			k.Button("重置").OnClick(s.Reset)
		})
	}, 400, 240)
	click(t, tt, "开始")
	tt.Advance(3100 * time.Millisecond)
	if !tt.HasText("0:03") {
		t.Fatal("running timer did not tick")
	}
	click(t, tt, "暂停")
	tt.Advance(5 * time.Second)
	if !tt.HasText("0:03") {
		t.Fatal("paused timer kept advancing")
	}
	click(t, tt, "开始")
	tt.Advance(time.Second)
	if !tt.HasText("0:04") {
		t.Fatal("resume lost elapsed time")
	}
	click(t, tt, "重置")
	if s.Running || s.Elapsed != 0 || !tt.HasText("0:00") {
		t.Fatal("timer reset failed")
	}
}

func TestChatComposerIMEEnterNewlineErrorsBusyAndTokens(t *testing.T) {
	s := kit.ChatState{Messages: []kit.ChatMessageData{{ID: "local-1", Role: "user", Text: "existing"}}}
	fail, attempts, stop := false, 0, 0
	custom := false
	tt := ui.NewTester(view(func(k *kit.UI) {
		o := kit.ChatOptions{OnStop: func() { stop++; s.Busy = false }}
		if custom {
			o.OnSend = func(d kit.ChatDraft) error {
				attempts++
				if fail {
					return errors.New("连接中断")
				}
				return nil
			}
		}
		k.ChatComposer(&s, o)
	}), 600, 400)
	click(t, tt, "消息")
	tt.Compose("ni", 2)
	tt.Key(0, ui.KeyEnter)
	if len(s.Messages) != 1 {
		t.Fatal("composition submitted chat")
	}
	tt.Type("你好")
	tt.Key(ui.Shift, ui.KeyEnter)
	tt.Type("Go")
	if !strings.Contains(s.Draft, "\n") {
		t.Fatal("Shift Enter did not insert newline")
	}
	tt.Key(0, ui.KeyEnter)
	if len(s.Messages) != 2 || s.Messages[1].ID == "local-1" || s.Draft != "" {
		t.Fatalf("send result=%+v", s)
	}
	tt.Key(0, ui.KeyEnter)
	if len(s.Messages) != 2 {
		t.Fatal("empty message was sent")
	}
	custom, fail = true, true
	tt.Frame()
	click(t, tt, "消息")
	tt.Type("retry")
	tt.Key(0, ui.KeyEnter)
	if s.Error != "连接中断" || s.Draft != "retry" || attempts != 1 {
		t.Fatal("send failure lost draft")
	}
	fail = false
	click(t, tt, "发送消息")
	if s.Error != "" || s.Draft != "" || attempts != 2 {
		t.Fatal("retry did not clear error and draft")
	}
	s.Busy = true
	tt.Frame()
	click(t, tt, "消息")
	tt.Type("later")
	tt.Key(0, ui.KeyEnter)
	if attempts != 2 {
		t.Fatal("busy chat sent a message")
	}
	click(t, tt, "停止生成")
	if stop != 1 || s.Busy {
		t.Fatal("stop callback failed")
	}
	s.Draft = ""
	s.Tokens = []kit.ChatToken{{ID: "doc", Label: "context"}}
	tt.Frame()
	click(t, tt, "发送消息")
	if attempts != 3 || len(s.Tokens) != 0 {
		t.Fatal("token-only draft was not sent")
	}
}

func TestChatMessageToolsReactionCopyRetryAndDictation(t *testing.T) {
	m := kit.ChatMessageData{ID: "m", Role: "assistant", Text: "answer", Status: "error", Reasoning: "analysis", Tools: []kit.ChatToolCall{{ID: "tool", Name: "search", Status: "done", Output: "result"}}, Reactions: []kit.ChatReaction{{Label: "有用"}}}
	d, retry, start, stop := kit.DictationState{}, 0, 0, 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.ChatMessage(&m, kit.ChatMessageOptions{OnRetry: func(id string) {
			if id == "m" {
				retry++
			}
		}})
		k.ChatDictationButton(&d, kit.DictationOptions{Start: func() error { start++; return nil }, Stop: func() error { stop++; return nil }})
	}), 650, 800)
	click(t, tt, "思考过程")
	click(t, tt, "search · done")
	if !m.ReasoningOpen || !m.Tools[0].Expanded || !tt.HasText("result") {
		t.Fatal("message detail did not expand")
	}
	click(t, tt, "反馈有用")
	click(t, tt, "反馈有用")
	if m.Reactions[0].Active || m.Reactions[0].Count != 0 {
		t.Fatal("reaction did not toggle")
	}
	click(t, tt, "复制消息")
	if tt.Clipboard() != "answer" {
		t.Fatal("message was not copied")
	}
	click(t, tt, "重试")
	if retry != 1 {
		t.Fatal("retry callback failed")
	}
	click(t, tt, "开始听写")
	if start != 0 {
		t.Fatal("unavailable dictation started")
	}
	d.Available = true
	tt.Frame()
	click(t, tt, "开始听写")
	click(t, tt, "停止听写")
	if start != 1 || stop != 1 || d.Recording {
		t.Fatal("dictation callbacks did not toggle recording")
	}
}

func TestResizableKeyboardLimitsAndCollapse(t *testing.T) {
	s, changes := kit.ResizableState{Size: 110}, 0
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Resizable(&s, func() { k.Text("第一面板") }, func() { k.Text("第二面板") }, kit.ResizableOptions{Min: 100, Max: 130, Collapsible: true, OnResize: func(float32) { changes++ }}).Height(180).FillWidth()
	}), 600, 400)
	click(t, tt, "Divider")
	for range 5 {
		tt.Key(0, ui.KeyRight)
	}
	if s.Size != 130 || changes == 0 {
		t.Fatalf("max size=%v changes=%d", s.Size, changes)
	}
	for range 6 {
		tt.Key(0, ui.KeyLeft)
	}
	if s.Size != 100 {
		t.Fatalf("min size=%v", s.Size)
	}
	click(t, tt, "收起面板")
	if !s.Collapsed || tt.HasText("第一面板") || !tt.HasText("第二面板") {
		t.Fatal("collapse did not hide only first pane")
	}
	click(t, tt, "展开面板")
	if s.Collapsed || !tt.HasText("第一面板") || s.Size != 100 {
		t.Fatal("expand lost pane size")
	}
}

func TestStepperLinearDisabledAndTwoDigitSteps(t *testing.T) {
	current := 0
	steps := []kit.Step{}
	for i := range 12 {
		steps = append(steps, kit.Step{ID: fmt.Sprint(i), Label: fmt.Sprintf("步骤%d", i+1), Disabled: i == 1})
	}
	tt := ui.NewTester(view(func(k *kit.UI) { k.Stepper(&current, steps, kit.StepperOptions{Linear: true}) }), 900, 800)
	click(t, tt, "步骤3")
	click(t, tt, "步骤2")
	if current != 0 {
		t.Fatal("linear or disabled step was bypassed")
	}
	if !tt.HasText("10") || !tt.HasText("12") {
		t.Fatal("two digit step indicators were not formatted")
	}
}

func TestInputGroupForwardsChangesAndClearRestoresFocus(t *testing.T) {
	text, changes := "draft", 0
	var handle ui.Handle
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.InputGroup(&text, nil, nil, kit.InputOptions{Label: "组合输入"}).OnChange(func() { changes++ })
		k.TextInput(&text, kit.InputOptions{Label: "可清除输入"}).Bind(&handle)
		k.InputClearButton(&text, &handle)
	}), 500, 320)
	click(t, tt, "组合输入")
	tt.Type("!")
	if changes != 1 {
		t.Fatal("input group did not forward its child's change")
	}
	click(t, tt, "清除输入")
	if text != "" || !tt.Focused("可清除输入") {
		t.Fatal("clear action did not restore input focus")
	}
}

func TestReadOnlyAdvancedInputsDoNotCommitOrRemoveTokens(t *testing.T) {
	chat := kit.ChatState{Draft: "locked", Tokens: []kit.ChatToken{{ID: "ctx", Label: "context"}}}
	suggestion := kit.TypeaheadState{Query: "Go", Value: "go"}
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.Typeahead(&suggestion, []kit.Suggestion{{ID: "rust", Label: "Rust"}}, kit.TypeaheadOptions{InputOptions: kit.InputOptions{Label: "只读建议", ReadOnly: true}})
		k.ChatComposer(&chat, kit.ChatOptions{ReadOnly: true})
	}), 600, 500)
	click(t, tt, "只读建议")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if suggestion.Value != "go" || suggestion.Query != "Go" {
		t.Fatal("read-only typeahead committed a suggestion")
	}
	if _, ok := tt.Find("移除context"); ok {
		t.Fatal("read-only composer allows removing context")
	}
	click(t, tt, "消息")
	tt.Type("x")
	tt.Key(0, ui.KeyEnter)
	if chat.Draft != "locked" || len(chat.Messages) != 0 {
		t.Fatal("read-only composer edited or sent")
	}
}

func TestFileMIMEParametersAndEmptyAcceptPatterns(t *testing.T) {
	dir := t.TempDir()
	text, unknown := filepath.Join(dir, "note.txt"), filepath.Join(dir, "no-extension")
	for _, p := range []string{text, unknown} {
		if err := os.WriteFile(p, []byte("test"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	s := kit.FileInputState{}
	if !kit.AddUploadPaths(&s, []string{text}, kit.FileInputOptions{Accept: "text/plain"}) {
		t.Fatalf("text/plain MIME with charset was rejected: %s", s.Error)
	}
	if kit.AddUploadPaths(&s, []string{unknown}, kit.FileInputOptions{Accept: ".txt,", Multiple: true}) {
		t.Fatal("empty accept segment allowed an unknown file")
	}
}

func TestCollapsedSideNavRetainsNestedMenu(t *testing.T) {
	selected, collapsed := "home", true
	tt := ui.NewTester(view(func(k *kit.UI) {
		k.SideNav(&selected, []kit.NavItem{{ID: "home", Label: "主页"}, {ID: "work", Label: "工作空间", Children: []kit.NavItem{{ID: "project", Label: "子项目"}}}}, kit.NavOptions{Collapsed: &collapsed})
	}), 500, 400)
	if tt.HasText("子项目") {
		t.Fatal("collapsed navigation still exposes expanded children")
	}
	click(t, tt, "工作空间")
	click(t, tt, "子项目")
	if selected != "project" {
		t.Fatal("collapsed navigation lost nested selection")
	}
	click(t, tt, "展开导航")
	if collapsed || !tt.HasText("子项目") {
		t.Fatal("expanding navigation lost its nested items")
	}
}
