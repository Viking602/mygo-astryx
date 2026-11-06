package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo/ui"
	"mygo-astryx/internal/catalog"
	"mygo-astryx/internal/gallery"
)

// Every component maps to an actual live demo. Compound pieces share a demo;
// rendering it once per configuration keeps the suite useful and small.
func renderSuite(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	type capture struct {
		Demo, Theme, File string
		Width, Height     int
		Components        []string
	}
	report := struct {
		Passed           bool
		ComponentEntries int
		Captures         []capture
		ScenarioCaptures []capture
		SourceRevision   string
	}{ComponentEntries: len(catalog.Components), SourceRevision: catalog.Astryx.Source.Revision}
	if err := os.WriteFile(filepath.Join(dir, "report.json"), []byte("{\"Passed\":false,\"Status\":\"running\"}\n"), 0644); err != nil {
		return err
	}
	for _, mode := range []int{1, 2} {
		theme := []string{"", "light", "dark"}[mode]
		for _, size := range [][2]int{{1280, 900}, {420, 780}} {
			for _, demo := range gallery.New().Demos(nil) {
				a := gallery.New()
				a.Page, a.Component, a.Mode = demo.Category, demo.Names[0], mode
				tt := ui.NewTester(a.View, size[0], size[1])
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				if !tt.HasText(demo.Title) || !tt.HasText(strings.Join(demo.Names, " · ")) {
					return fmt.Errorf("demo %s did not render its component content", demo.ID)
				}
				name := fmt.Sprintf("%s-%s-%d.png", demo.ID, theme, size[0])
				file, err := os.Create(filepath.Join(dir, name))
				if err != nil {
					return err
				}
				err = png.Encode(file, tt.Image())
				closeErr := file.Close()
				if err != nil {
					return err
				}
				if closeErr != nil {
					return closeErr
				}
				report.Captures = append(report.Captures, capture{demo.ID, theme, name, size[0], size[1], demo.Names})
			}
			for _, scenario := range gallery.Scenarios() {
				a := gallery.New()
				a.Component = scenario.Component
				a.Mode = mode
				for _, component := range catalog.Components {
					if component.Name == scenario.Component {
						a.Page = component.Category
						break
					}
				}
				if scenario.Prepare != nil {
					scenario.Prepare(a)
				}
				tt := ui.NewTester(a.View, size[0], size[1])
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				if scenario.Open != nil {
					if err := scenario.Open(tt); err != nil {
						return fmt.Errorf("scenario %s: %w", scenario.ID, err)
					}
				}
				if _, ok := tt.Find(scenario.Expected); !ok {
					return fmt.Errorf("scenario %s did not render %s", scenario.ID, scenario.Expected)
				}
				name := fmt.Sprintf("state-%s-%s-%d.png", scenario.ID, theme, size[0])
				file, err := os.Create(filepath.Join(dir, name))
				if err != nil {
					return err
				}
				err = png.Encode(file, tt.Image())
				closeErr := file.Close()
				if err != nil {
					return err
				}
				if closeErr != nil {
					return closeErr
				}
				report.ScenarioCaptures = append(report.ScenarioCaptures, capture{scenario.ID, theme, name, size[0], size[1], []string{scenario.Component}})
			}
		}
	}
	report.Passed = true
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(dir, "report.json"), append(data, '\n'), 0644); err != nil {
		return err
	}
	fmt.Printf("Rendered %d demo captures and %d state captures covering %d component entries into %s\n", len(report.Captures), len(report.ScenarioCaptures), report.ComponentEntries, dir)
	return nil
}
