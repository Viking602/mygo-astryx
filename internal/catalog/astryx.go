package catalog

import (
	_ "embed"
	"encoding/json"
	"sort"
)

//go:embed astryx.json
var sourceManifest []byte

type SourcePart struct {
	Export, Status, Representation, Source string
	Native                                 []string
	Documented                             bool
}
type SourceFamily struct {
	Family, Mapping, Notes string
	Native                 []string
	Parts                  []SourcePart
}
type SourceIndex struct {
	Source   struct{ Repository, Revision, Entry string }
	Families []SourceFamily
}

var Astryx = func() SourceIndex {
	var index SourceIndex
	if err := json.Unmarshal(sourceManifest, &index); err != nil {
		panic("catalog: invalid Astryx manifest: " + err.Error())
	}
	return index
}()

var aliasesByAPI, familiesByAPI = sourceLookups()

func sourceLookups() (map[string][]string, map[string][]int) {
	aliases, indices := map[string][]string{}, map[string][]int{}
	for i, family := range Astryx.Families {
		for _, method := range family.Native {
			aliases[method] = append(aliases[method], family.Family)
			indices[method] = append(indices[method], i)
		}
		for _, part := range family.Parts {
			if len(part.Native) > 0 {
				method := part.Native[0]
				aliases[method] = append(aliases[method], part.Export)
			}
		}
	}
	return aliases, indices
}

// AliasesFor maps upstream names to the primary native API rather than treating
// every composition helper (Text, Stack, etc.) as another component alias.
func AliasesFor(names []string) []string {
	wanted := map[string]bool{}
	for _, name := range names {
		wanted[name] = true
	}
	seen := map[string]bool{}
	var result []string
	add := func(name string) {
		if !seen[name] && !wanted[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	for _, method := range names {
		for _, alias := range aliasesByAPI[method] {
			add(alias)
		}
	}
	return result
}

func FamiliesFor(names []string) []SourceFamily {
	wanted := map[int]bool{}
	for _, name := range names {
		for _, i := range familiesByAPI[name] {
			wanted[i] = true
		}
	}
	indices := make([]int, 0, len(wanted))
	for i := range wanted {
		indices = append(indices, i)
	}
	sort.Ints(indices)
	var result []SourceFamily
	for _, i := range indices {
		result = append(result, Astryx.Families[i])
	}
	return result
}

func SourceURL(family SourceFamily) string {
	return Astryx.Source.Repository + "/tree/" + Astryx.Source.Revision + "/packages/core/src/" + family.Family
}
