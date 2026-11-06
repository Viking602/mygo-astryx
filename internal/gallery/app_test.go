package gallery

import (
	"fmt"
	"github.com/Viking602/mygo-astryx/internal/catalog"
	"github.com/egoist/mygo/ui"
	"strings"
	"testing"
)

func TestEveryComponentRendersInBothThemesAndWindowSizes(t *testing.T) {
	for _, mode := range []int{1, 2} {
		for _, width := range []int{1280, 320} {
			for _, component := range catalog.Components {
				t.Run(fmt.Sprintf("%s/mode%d/%d", component.Name, mode, width), func(t *testing.T) {
					a := New()
					a.Page, a.Component, a.Mode = component.Category, component.Name, mode
					var demo Demo
					for _, d := range a.Demos(nil) {
						if slicesContain(d.Names, component.Name) {
							demo = d
							break
						}
					}
					tt := ui.NewTester(a.View, width, 900)
					tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
					if !tt.HasText(demo.Title) || !tt.HasText(strings.Join(demo.Names, " · ")) {
						t.Fatal("component selection did not render its live demo")
					}
					if tt.Image().Bounds().Dx() != width {
						t.Fatal("incorrect output width")
					}
					r, ok := tt.Find("查看组件")
					if !ok || r.X < 0 || r.X+r.W > float32(width)+1 {
						t.Fatalf("component picker is outside window: %+v", r)
					}
				})
			}
		}
	}
}

func slicesContain(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func TestEveryCatalogComponentHasOneLiveDemo(t *testing.T) {
	app := New()
	covered := map[string]string{}
	for _, demo := range app.Demos(nil) {
		if demo.Build == nil || demo.Code == "" {
			t.Fatalf("incomplete demo %s", demo.ID)
		}
		for _, name := range demo.Names {
			if previous := covered[name]; previous != "" {
				t.Fatalf("duplicate %s in %s and %s", name, previous, demo.ID)
			}
			covered[name] = demo.Category
		}
	}
	for _, component := range catalog.Components {
		if covered[component.Name] != component.Category {
			t.Errorf("%s has no demo in %s", component.Name, component.Category)
		}
		delete(covered, component.Name)
	}
	if len(covered) > 0 {
		t.Fatalf("uncatalogued demos: %v", covered)
	}
}

func TestGalleryCategoriesThemesAndNarrowWindow(t *testing.T) {
	for _, mode := range []int{1, 2} {
		for _, cat := range catalog.Categories {
			t.Run(cat.ID+[]string{"", "-light", "-dark"}[mode], func(t *testing.T) {
				app := New()
				app.Page = cat.ID
				app.Mode = mode
				tt := ui.NewTester(app.View, 1100, 850)
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				if !tt.HasText(cat.Title) {
					t.Fatalf("missing category heading: %v", tt.Texts())
				}
				if len(tt.Texts()) < 15 {
					t.Fatal("category did not render its demos")
				}
				tt.SetSize(420, 780)
				if _, ok := tt.Find("组件分类"); !ok {
					t.Fatal("narrow layout lost category picker")
				}
				if tt.Image().Bounds().Dx() != 420 {
					t.Fatal("resize did not render")
				}
			})
		}
	}
}

func TestGallerySearchNavigationAndTheme(t *testing.T) {
	a := New()
	tt := ui.NewTester(a.View, 1280, 900)
	if err := tt.Click("浏览组件"); err != nil {
		t.Fatal(err)
	}
	if a.Page != "actions" {
		t.Fatal("overview navigation did not work")
	}
	if err := tt.Click("主要动作"); err != nil {
		t.Fatal(err)
	}
	if a.Clicks != 1 {
		t.Fatal("live button did not update model")
	}
	if err := tt.Click("深色"); err != nil {
		t.Fatal(err)
	}
	if a.Mode != 2 {
		t.Fatal("theme selector did not change theme")
	}
	if err := tt.Click("搜索组件"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Dialog")
	if !tt.HasText("搜索结果") || !tt.HasText("打开对话框") {
		t.Fatalf("search did not find overlay demos: %v", tt.Texts())
	}
	if err := tt.Click("打开对话框"); err != nil {
		t.Fatal(err)
	}
	if !a.DialogOpen {
		t.Fatal("gallery dialog did not open")
	}
	tt.Key(0, ui.KeyEscape)
	if a.DialogOpen {
		t.Fatal("gallery dialog did not close")
	}
}
