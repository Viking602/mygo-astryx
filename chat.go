package astryx

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"slices"
	"strings"
	"time"
)

type ChatToken struct{ ID, Label, Kind string }
type ChatTextRun struct {
	Text  string
	Token *ChatToken
}
type ChatReaction struct {
	Label  string
	Count  int
	Active bool
}
type ChatToolCall struct {
	ID, Name, Status, Input, Output string
	Expanded                        bool
}
type ChatMessageData struct {
	ID, Role, Author, Text, Status, Reasoning string
	Time                                      time.Time
	Streaming, ReasoningOpen                  bool
	Tools                                     []ChatToolCall
	Reactions                                 []ChatReaction
	Tokens                                    []ChatToken
}
type ChatMessageOptions struct {
	OnRetry    func(string)
	OnReaction func(string, int)
}
type ChatDraft struct {
	Text   string
	Tokens []ChatToken
}
type ChatState struct {
	Messages []ChatMessageData
	Draft    string
	Tokens   []ChatToken
	Busy     bool
	Error    string
	List     ui.ListState
	next     int
}
type ChatOptions struct {
	Placeholder                        string
	Disabled, ReadOnly, MultilineEnter bool
	OnSend                             func(ChatDraft) error
	OnStop                             func()
	Extras                             func()
	MessageOptions                     ChatMessageOptions
}

func (k *UI) ChatLayout(messages, composer func()) ui.Element {
	return ui.Column(k.c).MinHeight(0).Gap(12).Children(func() {
		ui.Column(k.c).Grow(1).MinHeight(0).Children(messages)
		if composer != nil {
			ui.Column(k.c).BorderWidth(1, 0, 0, 0).BorderColor(k.tokens.Native.Border).Padding(12, 0, 0, 0).Children(composer)
		}
	})
}
func (k *UI) ChatMessageMetadata(message ChatMessageData) ui.Element {
	author := message.Author
	if author == "" {
		author = map[string]string{"user": "你", "assistant": "助手", "system": "系统"}[message.Role]
	}
	if author == "" {
		author = "消息"
	}
	return ui.Row(k.c).Gap(8).Wrap().Children(func() {
		k.Caption(author).FontWeight(600)
		if !message.Time.IsZero() {
			k.Timestamp(message.Time, TimestampOptions{Format: "time", Static: true})
		}
		if message.Streaming {
			k.StatusDot("生成中", Info)
		} else if message.Status != "" {
			k.Caption(message.Status)
		}
	})
}
func (k *UI) ChatMessageBubble(message ChatMessageData) ui.Element {
	user := message.Role == "user"
	fill := k.tokens.Elevated
	if user {
		fill = k.tokens.Subtle
	}
	return ui.Column(k.c).Padding(14, 16).Gap(8).Radius(k.tokens.RadiusLarge).Background(fill).Border(1, k.tokens.Native.Border).Children(func() {
		if len(message.Tokens) > 0 {
			k.Inline().Children(func() {
				for _, token := range message.Tokens {
					k.Key(token.ID).ChatComposerTokenElement(token, nil)
				}
			})
		}
		if user {
			k.Text(message.Text).Selectable()
		} else {
			k.Markdown(message.Text, MarkdownOptions{Streaming: message.Streaming})
		}
	})
}
func (k *UI) ChatMessage(message *ChatMessageData, options ...ChatMessageOptions) ui.Element {
	o := option(options)
	services := k.c.Services()
	root := ui.Column(k.c).Padding(8, 4).Gap(8)
	root.Children(func() {
		k.ChatMessageMetadata(*message)
		if message.Reasoning != "" {
			k.Key("reasoning").Collapsible("思考过程", &message.ReasoningOpen, func() { k.Caption(message.Reasoning) })
		}
		k.ChatMessageBubble(*message)
		if len(message.Tools) > 0 {
			k.ChatToolCalls(message.Tools)
		}
		k.Inline().Children(func() {
			k.Button("复制消息", ButtonOptions{Variant: Ghost, Size: Small}).OnClick(func() { services.WriteClipboard(message.Text) })
			if message.Status == "error" && o.OnRetry != nil {
				k.Button("重试", ButtonOptions{Variant: Outline, Size: Small}).OnClick(func() { o.OnRetry(message.ID) })
			}
			for i := range message.Reactions {
				reaction := &message.Reactions[i]
				variant := Ghost
				if reaction.Active {
					variant = Outline
				}
				k.Button(fmt.Sprintf("%s %d", reaction.Label, reaction.Count), ButtonOptions{Variant: variant, Size: Small}).Label("反馈" + reaction.Label).OnClick(func() {
					if reaction.Active {
						reaction.Count = max(reaction.Count-1, 0)
					} else {
						reaction.Count++
					}
					reaction.Active = !reaction.Active
					if o.OnReaction != nil {
						o.OnReaction(message.ID, i)
					}
				})
			}
		})
	})
	return root
}
func (k *UI) ChatMessageList(state *ui.ListState, messages []ChatMessageData, options ...ChatMessageOptions) ui.Element {
	state.Key = func(i int) any { return messages[i].ID }
	return ui.List(k.c, state, len(messages), func(i int) {
		if messages[i].Role == "system" {
			k.ChatSystemMessage(messages[i].Text)
		} else {
			k.ChatMessage(&messages[i], options...)
		}
	}).Label("聊天消息").Gap(10)
}
func (k *UI) ChatSystemMessage(text string) ui.Element {
	return ui.Row(k.c).Justify(ui.Center).Padding(8).Role(ui.RoleStatus).Label(text).Children(func() { k.Caption(text) })
}
func (k *UI) ChatToolCalls(calls []ChatToolCall) ui.Element {
	return ui.Column(k.c).Gap(8).Label("工具调用").Children(func() {
		for i := range calls {
			call := &calls[i]
			ui.Column(k.c.Key(call.ID)).Padding(8, 12).Background(k.tokens.Subtle.Alpha(0.5)).Border(1, k.tokens.Native.Border).Radius(k.tokens.Radius).Children(func() {
				k.Collapsible(call.Name+" · "+call.Status, &call.Expanded, func() {
					if call.Input != "" {
						k.Caption("输入")
						k.CodeBlock(call.Input, "json")
					}
					if call.Output != "" {
						k.Caption("输出")
						k.CodeBlock(call.Output, "text")
					}
				})
				if call.Status == "running" {
					k.Inline().Children(func() { k.Spinner().Size(14, 14); k.Caption("正在执行") })
				}
			})
		}
	})
}
func (k *UI) ChatTokenizedText(runs []ChatTextRun) ui.Element {
	spans := make([]ui.Span, 0, len(runs))
	for _, run := range runs {
		if run.Token != nil {
			spans = append(spans, ui.Span{Text: run.Token.Label, Color: k.tokens.Info, Background: k.tokens.InfoMuted, Radius: 4, Weight: 550})
		} else {
			spans = append(spans, ui.Span{Text: run.Text})
		}
	}
	return ui.RichText(k.c, spans...).Selectable()
}
func (k *UI) ChatComposerTokenElement(token ChatToken, onRemove func()) ui.Element {
	return k.Token(token.Label, onRemove).Description(token.Kind)
}

type ChatComposerInputOptions struct {
	InputOptions
	MultilineEnter bool
	OnSend         func(string) bool
}

func (k *UI) ChatComposerInput(value *string, options ...ChatComposerInputOptions) ui.Element {
	o := option(options)
	if o.Label == "" {
		o.Label = "消息"
	}
	if o.Placeholder == "" {
		o.Placeholder = "输入消息…"
	}
	input := k.TextArea(value, o.InputOptions).MinHeight(72).MaxHeight(180)
	composing := input.Composing()
	input.OnTextKey(func(key ui.TextKey) bool {
		if o.Disabled || o.ReadOnly || composing || o.OnSend == nil || key.Key != ui.KeyEnter {
			return false
		}
		send := key.Mods == 0 && !o.MultilineEnter || key.Mods&ui.Cmd != 0
		if !send {
			return false
		}
		text := key.Text()
		if o.OnSend(text) {
			key.SetText("")
		}
		return true
	})
	return input
}
func (k *UI) ChatSendButton(busy, canSend bool, onSend, onStop func()) ui.Element {
	if busy {
		return k.Button("停止生成", ButtonOptions{Variant: Destructive, Disabled: onStop == nil, Icon: IconClose}).OnClick(onStop)
	}
	return k.Button("发送消息", ButtonOptions{Variant: Primary, Disabled: !canSend, EndIcon: IconArrow}).OnClick(onSend)
}
func (k *UI) ChatComposerDrawer(label string, open *bool, content func()) ui.Element {
	return k.Collapsible(label, open, content)
}
func (k *UI) ChatComposer(state *ChatState, options ...ChatOptions) ui.Element {
	o := option(options)
	now := k.c.Now()
	submit := func(text string) bool {
		if o.Disabled || o.ReadOnly || state.Busy || strings.TrimSpace(text) == "" && len(state.Tokens) == 0 {
			return false
		}
		draft := ChatDraft{Text: text, Tokens: slices.Clone(state.Tokens)}
		if o.OnSend != nil {
			if err := o.OnSend(draft); err != nil {
				state.Error = err.Error()
				return false
			}
		} else {
			state.next++
			id := fmt.Sprintf("local-%d", state.next)
			for slices.ContainsFunc(state.Messages, func(m ChatMessageData) bool { return m.ID == id }) {
				state.next++
				id = fmt.Sprintf("local-%d", state.next)
			}
			state.Messages = append(state.Messages, ChatMessageData{ID: id, Role: "user", Text: text, Time: now, Tokens: draft.Tokens})
		}
		state.Draft = ""
		state.Tokens = nil
		state.Error = ""
		return true
	}
	card := k.Card(CardOptions{Padding: 14})
	card.Children(func() {
		if len(state.Tokens) > 0 {
			k.Inline().Children(func() {
				for _, token := range state.Tokens {
					var remove func()
					if !o.Disabled && !o.ReadOnly {
						remove = func() {
							state.Tokens = slices.DeleteFunc(state.Tokens, func(t ChatToken) bool { return t.ID == token.ID })
						}
					}
					k.Key(token.ID).ChatComposerTokenElement(token, remove)
				}
			})
		}
		k.Key("draft").ChatComposerInput(&state.Draft, ChatComposerInputOptions{InputOptions: InputOptions{Label: "消息", Placeholder: o.Placeholder, Disabled: o.Disabled, ReadOnly: o.ReadOnly}, MultilineEnter: o.MultilineEnter, OnSend: submit}).FillWidth().MinHeight(60).Padding(4, 2).Border(0, ui.Transparent).Background(ui.Transparent).FocusRing(false)
		k.Inline().Children(func() {
			if o.Extras != nil {
				o.Extras()
			}
			hint := "Shift+Enter 换行"
			if o.MultilineEnter {
				hint = "Ctrl/⌘+Enter 发送"
			}
			k.Caption(hint).Grow(1)
			k.ChatSendButton(state.Busy, !o.Disabled && !o.ReadOnly && (strings.TrimSpace(state.Draft) != "" || len(state.Tokens) > 0), func() { submit(state.Draft) }, o.OnStop)
		})
		if state.Error != "" {
			k.Caption(state.Error).TextColor(k.tokens.Native.Danger)
		}
	})
	if !o.Disabled && !o.ReadOnly && card.FocusWithin() {
		card.BorderColor(k.tokens.Native.Focus)
	}
	return card
}
func (k *UI) ChatLayoutScrollButton(state *ui.ListState, label string) ui.Element {
	if label == "" {
		label = "回到最新消息"
	}
	return k.Button(label, ButtonOptions{Variant: Outline, Size: Small}).AlignSelf(ui.Start).OnClick(state.ScrollToEnd)
}
func (k *UI) Chat(state *ChatState, options ...ChatOptions) ui.Element {
	o := option(options)
	state.List.FollowEnd = true
	return k.ChatLayout(func() {
		if len(state.Messages) == 0 {
			k.EmptyState("开始对话", "发送第一条消息。", nil)
		} else {
			k.ChatMessageList(&state.List, state.Messages, o.MessageOptions).Grow(1).MinHeight(0)
		}
	}, func() { k.ChatComposer(state, o) })
}

// Dictation is a UI/service boundary. Native speech recognition is supplied by
// the host application; without both callbacks this control stays unavailable.
type DictationState struct {
	Available, Recording bool
	Error                string
}
type DictationOptions struct{ Start, Stop func() error }

func (k *UI) ChatDictationButton(state *DictationState, options ...DictationOptions) ui.Element {
	o := option(options)
	available := state.Available && o.Start != nil && o.Stop != nil
	label := "开始听写"
	variant := Ghost
	if state.Recording {
		label = "停止听写"
		variant = Destructive
	}
	e := k.Button(label, ButtonOptions{Variant: variant, Disabled: !available}).OnClick(func() {
		var err error
		if state.Recording {
			err = o.Stop()
		} else {
			err = o.Start()
		}
		if err != nil {
			state.Error = err.Error()
			return
		}
		state.Recording = !state.Recording
		state.Error = ""
	})
	if !available {
		e.Description("需要由宿主应用连接本机听写服务")
	}
	if state.Error != "" {
		e.Error(state.Error)
	}
	return e
}
