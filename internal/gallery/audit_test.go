package gallery

import (
	"fmt"
	"testing"

	"github.com/Viking602/mygo-astryx/internal/catalog"
	"github.com/egoist/mygo/ui"
)

func TestGalleryScenariosBothThemesAndSizes(t *testing.T) {
	for _, scenario := range Scenarios() {
		for _, mode := range []int{1, 2} {
			for _, width := range []int{320, 1280} {
				t.Run(fmt.Sprintf("%s/%d/%d", scenario.ID, mode, width), func(t *testing.T) {
					a := New()
					a.Component = scenario.Component
					a.Mode = mode
					found := false
					for _, component := range catalog.Components {
						if component.Name == scenario.Component {
							a.Page = component.Category
							found = true
							break
						}
					}
					if !found {
						t.Fatal("scenario has no public component")
					}
					if scenario.Prepare != nil {
						scenario.Prepare(a)
					}
					tt := ui.NewTester(a.View, width, 900)
					tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
					if scenario.Open != nil {
						if err := scenario.Open(tt); err != nil {
							t.Fatal(err)
						}
					}
					if _, ok := tt.Find(scenario.Expected); !ok {
						t.Fatalf("scenario did not render %q", scenario.Expected)
					}
				})
			}
		}
	}
}

func TestGallerySearchesAstryxNamesAndShowsPinnedReference(t *testing.T) {
	for _, query := range []string{"CheckboxInput", "ProgressBar", "ScrollableArea", "TabList"} {
		a := New()
		a.Query = query
		tt := ui.NewTester(a.View, 1280, 1000)
		if tt.HasText("没有匹配的组件") {
			t.Fatalf("upstream name %s cannot find its native equivalent", query)
		}
	}
	a := New()
	a.Page = "forms"
	a.Component = "Checkbox"
	tt := ui.NewTester(a.View, 1280, 1100)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Astryx 组件对照"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Astryx · CheckboxInput") {
		t.Fatal("reference did not identify upstream family")
	}
	if err := tt.Click("查看 CheckboxInput 源码"); err != nil {
		t.Fatal(err)
	}
	urls := tt.OpenedURLs()
	if len(urls) != 1 || urls[0] != catalog.Astryx.Source.Repository+"/tree/"+catalog.Astryx.Source.Revision+"/packages/core/src/CheckboxInput" {
		t.Fatalf("reference link is not pinned: %v", urls)
	}
}

func TestGalleryInputKeepsFocusWhenNavigationChangesLayout(t *testing.T) {
	a := New()
	a.Page = "forms"
	a.Component = "TextInput"
	tt := ui.NewTester(a.View, 1280, 1000)
	if err := tt.Click("名称"); err != nil {
		t.Fatal(err)
	}
	tt.Command("selectAll")
	tt.Type("原生")
	tt.SetSize(420, 1000)
	if !tt.Focused("名称") {
		t.Fatal("responsive search insertion lost component input identity")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Type("组件")
	if a.Name != "原生组件" {
		t.Fatalf("resized editing target changed: %s", a.Name)
	}
}

func TestGalleryControlledSheetFlowAndMobileNavigation(t *testing.T) {
	a := New()
	a.Page = "overlays"
	a.Component = "BottomSheet"
	tt := ui.NewTester(a.View, 1280, 1000)
	if err := tt.Click("分步底部面板"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("流程项目名称"); err != nil {
		t.Fatal(err)
	}
	tt.Command("selectAll")
	tt.Type("分步验证")
	if err := tt.Click("下一步：核对"); err != nil {
		t.Fatal(err)
	}
	if a.Advanced.SheetStep != "review" || !tt.HasText("项目名称：分步验证") || tt.HasText("填写资料") {
		t.Fatal("sheet content did not replace the previous step")
	}
	tt.Key(0, ui.KeyEscape)
	if a.Advanced.SheetFlowOpen || !tt.Focused("分步底部面板") {
		t.Fatal("flow dismissal did not update the controlled state/focus")
	}
	a.Page = "navigation"
	a.Component = "MobileNavToggle"
	tt.Frame()
	if err := tt.Click("打开站点导航"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("移动项目"); err != nil {
		t.Fatal(err)
	}
	if a.Advanced.MobileOpen || a.Advanced.Navigation != "projects" {
		t.Fatal("mobile drawer selection did not navigate and dismiss")
	}
}
