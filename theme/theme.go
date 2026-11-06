// Package theme defines semantic design tokens and their MyGo theme mapping.
package theme

import "github.com/egoist/mygo/ui"

// Theme is a value. Copy it before customization; do not mutate a theme while
// another window is building or painting with it.
type Theme struct {
	Native                                                   ui.Theme
	Muted, Subtle, Elevated                                  ui.Color
	Info, InfoMuted, SuccessMuted, WarningMuted, DangerMuted ui.Color
	Shadow                                                   ui.Color
	RadiusSmall, Radius, RadiusLarge                         float32
	Space, FontSmall, FontBody, FontTitle                    float32
	// Indicator can replace decorative checkbox/radio/check visuals. The
	// owning control retains all input and accessibility semantics.
	Indicator func(c *ui.Context, kind, state string, size float32) (ui.Element, bool)
}

// FromNative preserves the host's current colors, typography and preferences.
func FromNative(base *ui.Theme) Theme {
	if base == nil {
		base = ui.LightTheme()
	}
	t := Theme{Native: *base}
	t.Space = base.Space(1)
	t.Radius = base.Radius
	t.RadiusSmall = max(base.Radius-2, 2)
	t.RadiusLarge = base.Radius + 4
	t.FontBody = base.Rem(1)
	t.FontSmall = base.Rem(0.9)
	t.FontTitle = base.Rem(1.8)
	t.Muted = base.TextMuted
	t.Subtle = base.Surface
	t.Elevated = base.Background
	t.Info = ui.Hex("#2457b2")
	if base.Dark {
		t.Info = ui.Hex("#93b4ff")
	}
	t.InfoMuted = base.Background.Mix(t.Info, 0.12)
	t.SuccessMuted = base.Background.Mix(base.Success, 0.12)
	t.WarningMuted = base.Background.Mix(base.Warning, 0.12)
	t.DangerMuted = base.Background.Mix(base.Danger, 0.12)
	t.Shadow = ui.RGBA(0, 0, 0, 0.08)
	if base.Dark {
		t.Shadow = ui.RGBA(0, 0, 0, 0.3)
	}
	return t
}

// Neutral creates an Astryx-inspired neutral theme, respecting desktop contrast
// and text scaling. Pass preferences on every build to follow desktop changes.
func Neutral(dark bool, p ui.Preferences) Theme {
	b := ui.LightTheme()
	if dark {
		b = ui.DarkTheme()
	}
	if dark {
		b.Background, b.Surface = ui.Hex("#131316"), ui.Hex("#25252a")
		b.SurfaceHover, b.SurfacePressed = ui.Hex("#303036"), ui.Hex("#3b3b42")
		b.Text, b.TextMuted, b.Border = ui.Hex("#f4f4f6"), ui.Hex("#a4a4b0"), ui.Hex("#34343c")
		b.Accent, b.AccentText = ui.Hex("#eeeef0"), ui.Hex("#18181b")
		b.AccentHover, b.AccentPressed = ui.Hex("#dcdce0"), ui.Hex("#c6c6cc")
		b.Success, b.Warning, b.Danger = ui.Hex("#91d5a0"), ui.Hex("#f0cf85"), ui.Hex("#ffa1a9")
	} else {
		b.Background, b.Surface = ui.Hex("#f6f6f4"), ui.Hex("#ededeb")
		b.SurfaceHover, b.SurfacePressed = ui.Hex("#e3e3e0"), ui.Hex("#d6d6d2")
		b.Text, b.TextMuted, b.Border = ui.Hex("#1d1d20"), ui.Hex("#66666e"), ui.Hex("#dcdcd7")
		b.Accent, b.AccentText = ui.Hex("#202024"), ui.Hex("#ffffff")
		b.AccentHover, b.AccentPressed = ui.Hex("#36363b"), ui.Hex("#101012")
		b.Success, b.Warning, b.Danger = ui.Hex("#237238"), ui.Hex("#855a12"), ui.Hex("#b42332")
	}
	b.Radius, b.Spacing = 8, 4
	b.Focus = b.Accent.Alpha(0.75)
	b.Selection = b.Accent.Alpha(0.18)
	if p.TextScale > 0 {
		b.FontSize *= p.TextScale
	}
	if p.HighContrast {
		b.Border = b.Border.Mix(b.Text, 0.5)
		b.TextMuted = b.TextMuted.Mix(b.Text, 0.4)
		b.Focus = b.Text
	}
	t := FromNative(b)
	if dark {
		t.Elevated = ui.Hex("#1b1b1f")
		t.InfoMuted = t.Elevated.Mix(t.Info, 0.14)
		t.SuccessMuted = t.Elevated.Mix(b.Success, 0.14)
		t.WarningMuted = t.Elevated.Mix(b.Warning, 0.14)
		t.DangerMuted = t.Elevated.Mix(b.Danger, 0.14)
	} else {
		t.Elevated = ui.Hex("#ffffff")
		t.InfoMuted = t.Elevated.Mix(t.Info, 0.11)
		t.SuccessMuted = t.Elevated.Mix(b.Success, 0.11)
		t.WarningMuted = t.Elevated.Mix(b.Warning, 0.12)
		t.DangerMuted = t.Elevated.Mix(b.Danger, 0.11)
	}
	return t
}

// System resolves the Astryx neutral theme for the window's current light/dark
// appearance and accessibility preferences.
func System(c *ui.Context) Theme {
	dark := false
	if base := c.Theme(); base != nil {
		dark = base.Dark
	}
	return Neutral(dark, c.Preferences())
}
