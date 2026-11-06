package astryx

import "github.com/egoist/mygo/ui"

// These original outline icons use the current text color. They are provided
// for examples and common actions; callers can supply any MyGo SVG.
var (
	IconPlus         = mustIcon(`<path d="M12 5v14M5 12h14"/>`)
	IconMinus        = mustIcon(`<path d="M5 12h14"/>`)
	IconClose        = mustIcon(`<path d="m6 6 12 12M18 6 6 18"/>`)
	IconCheck        = mustIcon(`<path d="m5 12.5 4.5 4.5L19 7.5"/>`)
	IconChevronDown  = mustIcon(`<path d="m6 9 6 6 6-6"/>`)
	IconChevronRight = mustIcon(`<path d="m9 6 6 6-6 6"/>`)
	IconChevronLeft  = mustIcon(`<path d="m15 6-6 6 6 6"/>`)
	IconMore         = mustIcon(`<circle cx="5" cy="12" r="1.2" fill="currentColor"/><circle cx="12" cy="12" r="1.2" fill="currentColor"/><circle cx="19" cy="12" r="1.2" fill="currentColor"/>`)
	IconCopy         = mustIcon(`<rect x="9" y="9" width="10" height="10" rx="2"/><path d="M15 9V7a2 2 0 0 0-2-2H7a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2"/>`)
	IconInfo         = mustIcon(`<circle cx="12" cy="12" r="9"/><path d="M12 16v-4M12 8h.01"/>`)
	IconUpload       = mustIcon(`<path d="M12 15V4m0 0-4 4m4-4 4 4M4 17v1a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-1"/>`)
	IconSearch       = mustIcon(`<circle cx="10.5" cy="10.5" r="6.5"/><path d="m16 16 5 5"/>`)
	IconArrow        = mustIcon(`<path d="M5 12h14m-6-6 6 6-6 6"/>`)
	IconSettings     = mustIcon(`<path d="M4 7h16M4 17h16"/><circle cx="8" cy="7" r="3"/><circle cx="16" cy="17" r="3"/>`)
)

func mustIcon(body string) *ui.SVG {
	icon, err := ui.ParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
	if err != nil {
		panic(err)
	}
	return icon
}
