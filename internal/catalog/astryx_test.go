package catalog

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

// The upstream export manifest prevents a plausible-looking component list
// from silently losing a family or mapping it to a nonexistent native API.
func TestPinnedAstryxFamiliesHaveNativeGalleryAPIs(t *testing.T) {
	data, err := os.ReadFile("astryx.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Source   struct{ Repository, Revision, Entry string }
		Families []struct {
			Family string
			Native []string
			Notes  string
			Parts  []SourcePart
		}
		PartStatusCounts map[string]int `json:"part_status_counts"`
	}
	if err = json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Source.Repository != "https://github.com/meta/astryx" || len(manifest.Source.Revision) != 40 || len(manifest.Families) != 99 {
		t.Fatal("missing pinned source or changed upstream family scope")
	}
	apis := map[string]bool{}
	for _, c := range Components {
		apis[c.Name] = true
	}
	seen := map[string]bool{}
	counts := map[string]int{}
	statuses := map[string]bool{"native-api": true, "native-equivalent": true, "native-model": true, "native-slot": true, "native-composition": true, "platform-adaptation": true, "deprecated-migration": true, "partial": true}
	for _, family := range manifest.Families {
		if family.Family == "" || seen[family.Family] || family.Notes == "" || len(family.Native) == 0 {
			t.Fatalf("incomplete/duplicate family: %s", family.Family)
		}
		seen[family.Family] = true
		for _, method := range family.Native {
			if !apis[method] {
				t.Errorf("%s points to missing native Gallery API %s", family.Family, method)
			}
		}
		for _, part := range family.Parts {
			if part.Export == "" || part.Representation == "" || !statuses[part.Status] || !strings.HasPrefix(part.Source, "packages/core/src/") || len(part.Native) == 0 {
				t.Errorf("unmapped export in %s", family.Family)
			}
			counts[part.Status]++
			for _, method := range part.Native {
				if !apis[method] {
					t.Errorf("%s/%s points to missing API %s", family.Family, part.Export, method)
				}
			}
			if strings.Contains(part.Representation, "父组件的原生实现或内容回调") {
				t.Errorf("generic unverified coverage claim for %s", part.Export)
			}
			if part.Export == "BottomSheetSwitcher" && part.Status != "deprecated-migration" {
				t.Fatal("deprecated switcher lost its explicit migration disposition")
			}
		}
	}
	if !reflect.DeepEqual(counts, manifest.PartStatusCounts) {
		t.Fatal("export status totals do not match individual dispositions")
	}
}

func TestSourceAliasesResolveToTheirPrimaryNativeAPI(t *testing.T) {
	for method, alias := range map[string]string{"Checkbox": "CheckboxInput", "Select": "Selector", "Tabs": "TabList", "Scroll": "ScrollableArea", "Progress": "ProgressBar"} {
		found := false
		for _, name := range AliasesFor([]string{method}) {
			if name == alias {
				found = true
			}
		}
		if !found {
			t.Errorf("%s is not searchable through %s", alias, method)
		}
	}
}
