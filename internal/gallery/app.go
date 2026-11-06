// Package gallery is the native, interactive component showcase.
package gallery

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	kit "mygo-astryx"
	"mygo-astryx/internal/catalog"
	"mygo-astryx/theme"
	"strings"
	"time"
)

type Demo struct {
	ID, Category, Title, Description, Code string
	Names                                  []string
	Build                                  func(*kit.UI)
}

// App owns demo models, never frame elements or builders.
type App struct {
	Page, Query                                                              string
	Component                                                                string
	Advanced                                                                 AdvancedState
	Mode                                                                     int
	ShowCode                                                                 bool
	Name, Bio, Email, Search, GroupText, Combo, Choice, Radio, Nav           string
	ReadOnlyEmail, InvalidEmail, FormEmail, GridName, GridEmail              string
	Checked, Mixed, Enabled, ToggleA, ToggleB, ToggleC                       bool
	Tab, Segment, PageNumber, Stars, SelectedRow                             int
	Volume, Low, High, Number                                                float64
	DialogOpen, ConfirmOpen, DrawerOpen, PopoverOpen, OpenA, OpenB, TreeOpen bool
	SplitPosition                                                            float32
	Date                                                                     time.Time
	Color                                                                    ui.Color
	TableState, ListState                                                    ui.ListState
	Clicks                                                                   int
}

func New() *App {
	s := newAdvanced()
	s.SelectedCard = true
	return &App{
		Advanced: s, Page: "overview", Name: "MyGo", Email: "hello@example.com",
		ReadOnlyEmail: "readonly@example.com", InvalidEmail: "invalid-email", FormEmail: "team@example.com",
		GridName: "Astryx Studio", GridEmail: "studio@example.com",
		Bio: "用 Go 构建原生桌面界面。", Checked: true, Enabled: true, ToggleA: true,
		Choice: "studio", Radio: "team", Combo: "Go", Nav: "projects", PageNumber: 1,
		Stars: 4, Volume: 42, Low: 20, High: 75, Number: 8, SplitPosition: 180,
		Date: time.Date(2026, 11, 6, 10, 30, 0, 0, time.Local), Color: ui.Hex("#6c82d9"), TreeOpen: true,
	}
}

func (a *App) View(c *ui.Context) {
	t := theme.System(c)
	if a.Mode != 0 {
		t = theme.Neutral(a.Mode == 2, c.Preferences())
	}
	kit.Scope(c, t, func(k *kit.UI) {
		c.OnShortcut(ui.Cmd, ui.KeyK, func() { a.Page = "menus"; a.Component = "CommandPalette"; a.Query = ""; a.Advanced.PaletteOpen = true })
		k.AppShell(func() { a.navigation(c, k) }, func() { a.content(c, k) })
		k.ToastViewport()
	}).Fill()
}

var (
	categoryOptions = func() []kit.Option {
		items := []kit.Option{{Value: "overview", Label: "总览"}}
		for _, cat := range catalog.Categories {
			items = append(items, kit.Option{Value: cat.ID, Label: cat.Title})
		}
		return items
	}()
	sidebarNavItems = func() []kit.NavItem {
		items := []kit.NavItem{{ID: "overview", Label: "总览", Badge: fmt.Sprint(len(catalog.Components))}}
		for _, cat := range catalog.Categories {
			items = append(items, kit.NavItem{ID: cat.ID, Label: cat.Title, Badge: fmt.Sprint(catalog.Count(cat.ID))})
		}
		return items
	}()
	categoryBadges = func() map[string]string {
		m := make(map[string]string, len(catalog.Categories))
		for _, cat := range catalog.Categories {
			m[cat.ID] = fmt.Sprint(catalog.Count(cat.ID))
		}
		return m
	}()
	componentPickers = func() map[string][]kit.Option {
		m := make(map[string][]kit.Option, len(catalog.Categories)+1)
		all := []kit.Option{{Value: "", Label: "全部组件"}}
		for _, component := range catalog.Components {
			opt := kit.Option{Value: component.Name, Label: component.Name}
			all = append(all, opt)
			if m[component.Category] == nil {
				m[component.Category] = []kit.Option{{Value: "", Label: "全部组件"}}
			}
			m[component.Category] = append(m[component.Category], opt)
		}
		m[""] = all
		m["overview"] = all
		return m
	}()
	demoSearchIndex = map[string]string{}
	demoNameIndex   = map[string]string{}
	demoHasSource   = map[string]bool{}
)

func ensureDemoIndex(d Demo) {
	if _, ok := demoSearchIndex[d.ID]; ok {
		return
	}
	names := strings.ToLower(strings.Join(d.Names, " ") + " " + strings.Join(catalog.AliasesFor(d.Names), " "))
	demoNameIndex[d.ID] = names
	demoSearchIndex[d.ID] = names + " " + strings.ToLower(d.Title+" "+d.Description+" "+d.Category)
	demoHasSource[d.ID] = len(catalog.FamiliesFor(d.Names)) > 0
}

func (a *App) navigation(c *ui.Context, k *kit.UI) {
	w, _ := c.Size()
	ui.Column(c).Gap(2).Children(func() {
		k.Text("MYGO / ASTRYX").FontSize(15).FontWeight(700).LetterSpacing(0.4)
		k.Caption("Go 原生 · 组件实验室")
	})
	if w < 760 {
		k.Key("category-picker").Select(&a.Page, categoryOptions, kit.InputOptions{Label: "组件分类"}).FillWidth().OnChange(func() { a.Component = ""; a.Query = "" })
		return
	}
	k.SearchField(&a.Query, kit.InputOptions{Label: "搜索组件", Placeholder: "搜索组件…"}).FillWidth().Margin(8, 0).OnChange(func() { a.Component = "" })
	old := a.Page
	ui.Scroll(c).Grow(1).MinHeight(0).Children(func() { k.Sidebar(&a.Page, sidebarNavItems) })
	if old != a.Page {
		a.Query = ""
		a.Component = ""
	}
	k.Divider()
	k.Caption("MyGo v0.7.4\n基于 Base 构建 · 原生绘制").LineHeight(1.5)
}

func (a *App) content(c *ui.Context, k *kit.UI) {
	w, _ := c.Size()
	bar := c.TitleBar()
	rightPad := max(float32(24), bar.Right+16)
	topMin := max(bar.Height, 40)
	topPad := float32(10)
	if w < 760 {
		topMin = 0
		topPad = 4
	}
	ui.Row(c).AlignItems(ui.Center).Gap(12).MinHeight(topMin).Padding(topPad, rightPad, 4, 24).FillWidth().DragWindow().Wrap().Children(func() {
		if w >= 760 {
			k.Spacer()
		}
		k.Key("theme-mode").Segmented(&a.Mode, []string{"系统", "浅色", "深色"}).Label("外观模式")
		k.Key("code-toggle").Switch(&a.ShowCode, "显示代码")
	})
	if w < 760 {
		ui.Column(c).FillWidth().Padding(4, 20, 8, 20).Children(func() {
			k.SearchField(&a.Query, kit.InputOptions{Label: "搜索组件", Placeholder: "搜索组件…"}).FillWidth().OnChange(func() { a.Component = "" })
		})
	}
	k.Key("gallery-scroll").Scroll().Grow(1).MinHeight(0).Children(func() {
		k.Container().MaxWidth(1140).Padding(12, 24, 24, 24).Children(func() {
			if a.Page == "overview" && strings.TrimSpace(a.Query) == "" && a.Component == "" {
				a.overview(c, k)
				return
			}
			title, description := "组件", "所有原生组件示例。"
			for _, cat := range catalog.Categories {
				if cat.ID == a.Page {
					title, description = cat.Title, cat.Description
				}
			}
			if a.Query != "" {
				title = "搜索结果"
				description = "匹配组件名称、分类或说明。"
			}
			picker := componentPickers[a.Page]
			if a.Query != "" || len(picker) == 0 {
				picker = componentPickers[""]
			}
			ui.Row(c).AlignItems(ui.End).Gap(16).Padding(0, 0, 4, 0).Wrap().Children(func() {
				ui.Column(c).Gap(4).Grow(1).MinWidth(200).Children(func() {
					k.Heading(title, 1)
					k.Caption(description).FontSize(k.Tokens().FontBody)
				})
				sel := k.Key("component-picker").Select(&a.Component, picker, kit.InputOptions{Label: "查看组件"})
				if w < 640 {
					sel.FillWidth()
				} else {
					sel.Width(220)
				}
			})
			found := 0
			q := strings.ToLower(strings.TrimSpace(a.Query))
			limit := 57
			if q != "" && a.Component == "" {
				limit = 6
			}
			demos := a.Demos(c)
			if q != "" && a.Component == "" {
				var primary, secondary []Demo
				for _, d := range demos {
					if !a.matches(d) {
						continue
					}
					if strings.Contains(demoNameIndex[d.ID], q) || strings.Contains(strings.ToLower(d.Title), q) {
						primary = append(primary, d)
					} else {
						secondary = append(secondary, d)
					}
				}
				demos = append(primary, secondary...)
			}
			for _, demo := range demos {
				if !a.matches(demo) {
					continue
				}
				found++
				if found > limit {
					continue
				}
				ui.Column(c.Key(demo.ID)).Children(func() {
					k.Card(kit.CardOptions{Padding: 20}).Gap(14).FillWidth().Children(func() {
						ui.Column(c).Gap(4).Children(func() {
							ui.Row(c).AlignItems(ui.Center).Gap(12).Wrap().Children(func() {
								k.Text(demo.Title).FontSize(16).FontWeight(650).Grow(1)
								k.Caption(strings.Join(demo.Names, " · ")).Font("monospace").Padding(2, 8).Radius(k.Tokens().RadiusSmall).Background(k.Tokens().Subtle)
							})
							k.Caption(demo.Description)
						})
						if demo.Build != nil {
							ui.Column(c).Gap(12).Children(func() {
								demo.Build(k)
							})
						}
						a.sourceDetails(c, k, demo)
						if a.ShowCode {
							k.Code(demo.Code).FillWidth()
						}
					})
				})
			}
			if found > limit {
				k.Caption(fmt.Sprintf("已显示前 %d 组最相关结果（共匹配 %d 组），可继续输入缩小范围。", limit, found))
			}
			if found == 0 {
				k.EmptyState("没有匹配的组件", "试试 Button、Input、Dialog，或清除搜索条件。", func() { k.Button("清除搜索").OnClick(func() { a.Query = "" }) })
			}
			k.Caption("MYGO ASTRYX · 原生组件，真实交互").Padding(8, 0)
		})
	})
}

func (a *App) matches(d Demo) bool {
	if a.Component != "" {
		for _, name := range d.Names {
			if name == a.Component {
				return true
			}
		}
		return false
	}
	q := strings.ToLower(strings.TrimSpace(a.Query))
	if q == "" {
		return d.Category == a.Page
	}
	ensureDemoIndex(d)
	return strings.Contains(demoSearchIndex[d.ID], q)
}

func (a *App) sourceDetails(c *ui.Context, k *kit.UI, d Demo) {
	ensureDemoIndex(d)
	if !demoHasSource[d.ID] {
		return
	}
	root := ui.Column(c.Key("source-details")).Padding(2, 0, 0, 0).BorderWidth(1, 0, 0, 0).BorderColor(k.Tokens().Native.Border.Alpha(0.6))
	open := ui.Local(root, "open", func() bool { return false })
	root.Children(func() {
		k.Collapsible("Astryx 组件对照", open, func() {
			names := d.Names
			if a.Component != "" {
				names = []string{a.Component}
			}
			for _, family := range catalog.FamiliesFor(names) {
				ui.Column(c).Gap(6).Padding(10, 14).Radius(k.Tokens().Radius).Background(k.Tokens().Subtle.Alpha(0.45)).Border(1, k.Tokens().Native.Border).Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(8).Wrap().Children(func() {
						k.Text("Astryx · " + family.Family).FontWeight(600).Grow(1)
						k.Badge(family.Mapping)
						k.Link("查看 "+family.Family+" 源码", catalog.SourceURL(family))
					})
					k.Caption(family.Notes)
					for _, part := range family.Parts {
						ui.Row(c).AlignItems(ui.Center).Gap(8).Wrap().Padding(2, 0).Children(func() {
							k.Caption(part.Export).Font("monospace").FontWeight(600).TextColor(k.Tokens().Native.Text)
							k.Caption("→ " + strings.Join(part.Native, ", ")).Font("monospace")
							k.Caption(part.Representation).Grow(1)
						})
					}
				})
			}
		})
	})
}

func (a *App) overview(c *ui.Context, k *kit.UI) {
	ui.Column(c).Gap(10).Children(func() {
		k.Caption("FOUNDATIONS, COMPONENTS & PATTERNS").FontSize(11).FontWeight(600).LetterSpacing(0.6)
		k.Heading("一套 Go 原生的界面语言。", 1).MaxWidth(700)
		k.Text("统一外观，保留原生交互。从一个按钮，到完整的桌面工作台。").TextColor(k.Tokens().Muted)
		k.Inline().Padding(4, 0, 0, 0).Children(func() {
			k.Badge(fmt.Sprintf("%d 个组件", len(catalog.Components)))
			k.Badge(fmt.Sprintf("%d 个分类", len(catalog.Categories)))
			k.Badge("浅色 / 深色")
			k.Badge("MyGo Base", kit.Info)
		})
	})
	k.Card(kit.CardOptions{Padding: 24, Elevated: true}).Children(func() {
		k.Inline().Children(func() {
			k.Avatar("MyGo", nil)
			k.Stack().Gap(2).Grow(1).Children(func() {
				k.Text("你的下一个原生应用").FontSize(17).FontWeight(650)
				k.Caption("由可组合的 Go 组件构建")
			})
			k.Badge("准备就绪", kit.Success)
		})
		k.Divider()
		k.Field("项目名称", func() {
			k.Key("overview-name").TextInput(&a.Name, kit.InputOptions{Placeholder: "为你的项目命名"}).FillWidth()
		})
		k.Inline().Children(func() {
			k.Key("overview-enabled").Switch(&a.Enabled, "启用同步")
			k.Spacer()
			k.Button("浏览组件", kit.ButtonOptions{Variant: kit.Outline}).OnClick(func() { a.Page = "actions" })
			k.Button("创建项目", kit.ButtonOptions{Variant: kit.Primary, EndIcon: kit.IconArrow}).OnClick(func() { k.Toast(ui.Toast{Title: "项目已创建", Description: a.Name}) })
		})
	})
	ui.Column(c).Gap(12).Children(func() {
		k.Heading("从基础开始探索", 3)
		cols := 2
		if w, _ := c.Size(); w < 920 {
			cols = 1
		}
		k.Grid(cols).Children(func() {
			for _, cat := range catalog.Categories {
				k.Key("category-"+cat.ID).ClickableCard(cat.Title, func() {
					ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
						k.Text(cat.Title).FontWeight(600).Grow(1)
						k.Badge(categoryBadges[cat.ID])
					})
					k.Caption(cat.Description)
				}, kit.CardOptions{Padding: 16}).OnClick(func() { a.Page = cat.ID })
			}
		})
	})
}
