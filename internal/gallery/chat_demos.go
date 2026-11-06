package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	kit "mygo-astryx"
	"strings"
)

func (a *App) ChatDemos(c *ui.Context) []Demo {
	if c != nil && strings.TrimSpace(a.Query) == "" && a.Page != "chat" {
		return nil
	}
	s := &a.Advanced
	return []Demo{
		{ID: "chat-workspace", Category: "chat", Title: "完整原生对话", Description: "本地消息演示，未连接 AI 服务。Enter 发送，Shift+Enter 换行。", Names: []string{"Chat", "ChatLayout", "ChatMessageList", "ChatComposer", "ChatComposerInput", "ChatSendButton", "ChatLayoutScrollButton"}, Code: `k.Chat(&state.Conversation, astryx.ChatOptions{OnSend: send})`, Build: func(k *kit.UI) {
			k.Switch(&s.ChatMultiline, "Enter 换行，Ctrl/⌘ Enter 发送")
			k.Key("conversation").Chat(&s.Chat, kit.ChatOptions{MultilineEnter: s.ChatMultiline, Placeholder: "在本地演示中发送消息…", OnSend: func(draft kit.ChatDraft) error {
				n := len(s.Chat.Messages)
				s.Chat.Messages = append(s.Chat.Messages, kit.ChatMessageData{ID: fmt.Sprintf("user-%d", n), Role: "user", Text: draft.Text, Tokens: draft.Tokens}, kit.ChatMessageData{ID: fmt.Sprintf("reply-%d", n), Role: "assistant", Text: "已收到：" + draft.Text + "\n\n这是本地交互演示，没有请求模型或外部服务。"})
				return nil
			}, OnStop: func() { s.Chat.Busy = false }}).Height(460)
			k.ChatLayoutScrollButton(&s.Chat.List, "回到最新消息")
		}},
		{ID: "chat-message-parts", Category: "chat", Title: "消息、反馈与工具调用", Description: "支持复制、可展开说明、工具输入输出和独立反馈。", Names: []string{"ChatMessage", "ChatMessageBubble", "ChatMessageMetadata", "ChatSystemMessage", "ChatToolCalls"}, Code: `k.ChatMessage(&state.Message)`, Build: func(k *kit.UI) {
			k.Key("rich-message").ChatMessage(&s.Message)
			k.ChatSystemMessage("工具调用与状态由应用模型提供。")
		}},
		{ID: "chat-context", Category: "chat", Title: "上下文、扩展区与听写边界", Description: "听写按钮由宿主服务驱动；此处仅模拟本地回调，不访问麦克风。", Names: []string{"ChatComposerTokenElement", "ChatTokenizedText", "ChatComposerDrawer", "ChatDictationButton"}, Code: `k.ChatDictationButton(&state.Dictation, astryx.DictationOptions{Start: start, Stop: stop})`, Build: func(k *kit.UI) {
			k.Inline().Children(func() {
				if s.ContextTokenVisible {
					k.ChatComposerTokenElement(kit.ChatToken{ID: "design", Label: "@设计规范", Kind: "文档上下文"}, func() { s.ContextTokenVisible = false })
				} else {
					k.Button("添加上下文", kit.ButtonOptions{Variant: kit.Outline, Size: kit.Small}).OnClick(func() { s.ContextTokenVisible = true })
				}
			})
			token := kit.ChatToken{ID: "readme", Label: "README.md", Kind: "文件"}
			k.ChatTokenizedText([]kit.ChatTextRun{{Text: "请参考 "}, {Token: &token}, {Text: " 中的组件说明。"}})
			k.ChatComposerDrawer("更多输入选项", &s.ComposerDrawer, func() {
				k.Checkbox(&a.Checked, "包含项目上下文")
				k.Caption("扩展区可以放置模型、文件或工具配置。")
			})
			k.Inline().Children(func() {
				k.ChatDictationButton(&s.Dictation, kit.DictationOptions{Start: func() error { s.MenuAction = "听写模拟已开始"; return nil }, Stop: func() error {
					s.Chat.Draft += strings.TrimSpace(" 这是模拟听写输入。")
					s.MenuAction = "听写模拟已停止，草稿已更新"
					return nil
				}})
				k.Badge("本地模拟", kit.Info)
			})
			if s.MenuAction != "" {
				k.Caption(s.MenuAction)
			}
		}},
	}
}
