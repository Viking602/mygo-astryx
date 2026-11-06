// Gallery is the native component explorer. -snapshot renders without a window.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"image/png"
	"log"
	"mygo-astryx/internal/catalog"
	"mygo-astryx/internal/gallery"
	"os"
	"path/filepath"
)

func main() {
	snapshot := flag.String("snapshot", "", "render a PNG without a window")
	appearance := flag.String("theme", "system", "system, light or dark")
	page := flag.String("page", "overview", "overview or a component category")
	component := flag.String("component", "", "show the live demo for one component")
	list := flag.Bool("list", false, "print the component and live demo index as JSON")
	suite := flag.String("snapshot-suite", "", "render every live demo in both themes and widths into a directory")
	width := flag.Int("width", 1280, "window/snapshot width")
	height := flag.Int("height", 900, "window/snapshot height")
	focus := flag.Bool("focus", false, "bring this Gallery window to the foreground after it is ready")
	flag.Parse()
	a := gallery.New()
	a.Page = *page
	if *list {
		type demoInfo struct {
			ID, Category, Title string
			Components          []string
		}
		demos := []demoInfo{}
		for _, d := range a.Demos(nil) {
			demos = append(demos, demoInfo{d.ID, d.Category, d.Title, d.Names})
		}
		type scenarioInfo struct{ ID, Component, Expected string }
		states := []scenarioInfo{}
		for _, s := range gallery.Scenarios() {
			states = append(states, scenarioInfo{s.ID, s.Component, s.Expected})
		}
		if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
			"Categories": catalog.Categories, "Components": catalog.Components, "Demos": demos, "Scenarios": states, "AstryxSource": catalog.Astryx.Source,
		}); err != nil {
			log.Fatal(err)
		}
		return
	}
	if *component != "" {
		found := false
		for _, item := range catalog.Components {
			if item.Name == *component {
				a.Page, a.Component, found = item.Category, item.Name, true
				break
			}
		}
		if !found {
			log.Fatal("unknown component: ", *component)
		}
	}
	switch *appearance {
	case "system":
	case "light":
		a.Mode = 1
	case "dark":
		a.Mode = 2
	default:
		log.Fatal("theme must be system, light or dark")
	}
	valid := a.Page == "overview"
	for _, cat := range catalog.Categories {
		if cat.ID == a.Page {
			valid = true
		}
	}
	if !valid {
		log.Fatal("unknown page: ", a.Page)
	}
	if *width < 320 || *height < 320 {
		log.Fatal("width and height must be at least 320")
	}
	if *suite != "" {
		if err := renderSuite(*suite); err != nil {
			log.Fatal(err)
		}
		return
	}
	if *snapshot != "" {
		tt := ui.NewTester(a.View, *width, *height)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		if err := os.MkdirAll(filepath.Dir(*snapshot), 0755); err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(*snapshot)
		if err != nil {
			log.Fatal(err)
		}
		if err = png.Encode(file, tt.Image()); err != nil {
			file.Close()
			log.Fatal(err)
		}
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Rendered %s (%d x %d)\n", *snapshot, *width, *height)
		return
	}
	mygo.App.WhenReady(func() {
		window := mygo.NewWindow(mygo.WindowOptions{
			Title:         "MyGo / Astryx — 原生组件 Gallery",
			Width:         *width,
			Height:        *height,
			TitleBarStyle: mygo.TitleBarHiddenInset,
			Content:       ui.View(a.View),
		})
		if *focus {
			window.Focus()
			mygo.App.Focus()
		}
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
