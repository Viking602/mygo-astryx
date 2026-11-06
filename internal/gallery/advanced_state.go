package gallery

import (
	kit "github.com/Viking602/mygo-astryx"
	"github.com/egoist/mygo/ui"
	"time"
)

type AdvancedState struct {
	Multi, Checks, Tokens                                                                                                       []string
	Radio, Navigation, Outline, ComplexDraft, ComplexValue, ClearText                                                           string
	ClearHandle                                                                                                                 ui.Handle
	ComplexOpen, SheetOpen, OverlayOpen, LightboxOpen, PaletteOpen, BannerDismissed, SelectedCard, CollapsedNav, ComposerDrawer bool
	MenuChecked                                                                                                                 bool
	RemovableToken, ContextTokenVisible                                                                                         bool
	MenuRadio, MenuAction                                                                                                       string
	Step, CarouselIndex, ImageIndex, TreeSelected                                                                               int
	Resize                                                                                                                      kit.ResizableState
	Suggest                                                                                                                     kit.TypeaheadState
	Filters                                                                                                                     kit.PowerSearchState
	Files                                                                                                                       kit.FileInputState
	Range                                                                                                                       kit.DateRange
	Timer                                                                                                                       kit.TimerState
	Palette                                                                                                                     kit.CommandPaletteState
	Tree                                                                                                                        ui.OutlineState[string]
	Images                                                                                                                      []kit.ImageAsset
	Chat                                                                                                                        kit.ChatState
	Message                                                                                                                     kit.ChatMessageData
	Dictation                                                                                                                   kit.DictationState
	MembersOpen, MobileOpen, LongConfirm, LongMenu, ExcludeWeekends, ChatMultiline                                              bool
	SheetStep                                                                                                                   string
	SheetFlowOpen                                                                                                               bool
}

func newAdvanced() AdvancedState {
	a := AdvancedState{RemovableToken: true, ContextTokenVisible: true, Step: 1, Multi: []string{"go"}, Checks: []string{"go"}, Tokens: []string{"Go", "原生"}, Radio: "go", Navigation: "home", Outline: "start", ComplexDraft: "标准", ComplexValue: "标准", ClearText: "可以清除的文字", MenuRadio: "stable", Resize: kit.ResizableState{Size: 180}, Range: kit.DateRange{Start: time.Date(2026, 11, 6, 0, 0, 0, 0, time.Local), End: time.Date(2026, 11, 12, 0, 0, 0, 0, time.Local)}, Timer: kit.TimerState{Elapsed: 65 * time.Second}, Dictation: kit.DictationState{Available: true}}
	a.Tree.Open.Add("root")
	for i, palette := range [][3]string{{"#e8d6a4", "#ba6b4a", "#58766a"}, {"#283449", "#7b90ad", "#c9b5a0"}, {"#d7e0d1", "#6b8b78", "#dab575"}} {
		name := []string{"晨光", "夜航", "旷野"}[i]
		svg, err := ui.ParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 160"><rect width="240" height="160" fill="` + palette[0] + `"/><circle cx="176" cy="48" r="24" fill="` + palette[1] + `"/><path d="M0 120 70 54 144 128 206 88 240 120V160H0Z" fill="` + palette[2] + `"/></svg>`))
		if err != nil {
			panic(err)
		}
		a.Images = append(a.Images, kit.ImageAsset{ID: name, Label: name, Source: svg})
	}
	a.Message = kit.ChatMessageData{ID: "rich-message", Role: "assistant", Author: "MyGo 助手", Text: "原生消息支持 **Markdown**、代码、复制和反馈。", Reasoning: "这是可展开的演示说明，不包含模型内部推理。", Status: "已完成", Tools: []kit.ChatToolCall{{ID: "tool-read", Name: "读取项目", Status: "success", Input: `{"path":"README.md"}`, Output: `{"components":"native"}`}}, Reactions: []kit.ChatReaction{{Label: "有帮助", Count: 2}, {Label: "需改进"}}}
	a.Chat.Messages = []kit.ChatMessageData{{ID: "welcome", Role: "system", Text: "本地演示 · 未连接 AI 服务"}, {ID: "question", Role: "user", Author: "你", Text: "MyGo UI 支持什么？"}, {ID: "answer", Role: "assistant", Author: "MyGo 助手", Text: "支持原生组件、主题和键盘交互。试试下方输入框，**Enter 发送**，Shift+Enter 换行。"}}
	a.Chat.Tokens = []kit.ChatToken{{ID: "context", Label: "@项目", Kind: "上下文"}}
	return a
}
