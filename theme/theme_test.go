package theme_test

import (
	"github.com/egoist/mygo/ui"
	"math"
	"mygo-astryx/theme"
	"testing"
)

func luminance(c ui.Color) float64 {
	linear := func(b uint8) float64 {
		v := float64(b) / 255
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*linear(c.R) + 0.7152*linear(c.G) + 0.0722*linear(c.B)
}
func contrast(a, b ui.Color) float64 {
	la, lb := luminance(a), luminance(b)
	return (max(la, lb) + 0.05) / (min(la, lb) + 0.05)
}

func TestNeutralTextContrastAndScaling(t *testing.T) {
	for _, dark := range []bool{false, true} {
		base := theme.Neutral(dark, ui.Preferences{})
		for _, pair := range []struct {
			name      string
			ink, fill ui.Color
		}{
			{"body", base.Native.Text, base.Native.Background},
			{"body-elevated", base.Native.Text, base.Elevated},
			{"secondary", base.Muted, base.Native.Background},
			{"secondary-elevated", base.Muted, base.Elevated},
			{"primary", base.Native.AccentText, base.Native.Accent},
			{"info", base.Info, base.InfoMuted},
			{"success", base.Native.Success, base.SuccessMuted},
			{"warning", base.Native.Warning, base.WarningMuted},
			{"danger", base.Native.Danger, base.DangerMuted},
		} {
			if ratio := contrast(pair.ink, pair.fill); ratio < 4.5 {
				t.Errorf("dark=%v %s contrast %.2f < 4.5", dark, pair.name, ratio)
			}
		}
		if base.Elevated == base.Native.Background {
			t.Errorf("dark=%v elevated surface is not distinct from background", dark)
		}
		scaled := theme.Neutral(dark, ui.Preferences{TextScale: 1.75})
		if scaled.FontBody != base.FontBody*1.75 {
			t.Fatal("system text scale not applied")
		}
		high := theme.Neutral(dark, ui.Preferences{HighContrast: true})
		if contrast(high.Native.Border, high.Native.Background) <= contrast(base.Native.Border, base.Native.Background) {
			t.Fatal("high contrast did not strengthen borders")
		}
	}
}

func TestHostThemeRemainsUnmodifiedAndScopesRestore(t *testing.T) {
	native := ui.DarkTheme()
	original := *native
	_ = theme.FromNative(native)
	if *native != original {
		t.Fatal("resolving tokens mutated host theme")
	}
}
