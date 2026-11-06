package astryx

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"time"
)

type TimestampOptions struct {
	Format   string
	Location *time.Location
	Static   bool
	Copy     bool
}

func relativeTime(delta time.Duration) string {
	future := delta < 0
	if future {
		delta = -delta
	}
	n, unit := int64(delta/time.Second), "秒"
	switch {
	case delta < 5*time.Second:
		return "刚刚"
	case delta >= 365*24*time.Hour:
		n, unit = int64(delta/(365*24*time.Hour)), "年"
	case delta >= 30*24*time.Hour:
		n, unit = int64(delta/(30*24*time.Hour)), "个月"
	case delta >= 24*time.Hour:
		n, unit = int64(delta/(24*time.Hour)), "天"
	case delta >= time.Hour:
		n, unit = int64(delta/time.Hour), "小时"
	case delta >= time.Minute:
		n, unit = int64(delta/time.Minute), "分钟"
	}
	if future {
		return fmt.Sprintf("%d%s后", n, unit)
	}
	return fmt.Sprintf("%d%s前", n, unit)
}
func (k *UI) Timestamp(value time.Time, options ...TimestampOptions) ui.Element {
	o := option(options)
	if o.Location != nil {
		value = value.In(o.Location)
	}
	text := "未设置时间"
	if !value.IsZero() {
		switch o.Format {
		case "date":
			text = value.Format("2006-01-02")
		case "time":
			text = value.Format("15:04:05")
		case "absolute", "date_time":
			text = value.Format("2006-01-02 15:04:05 MST")
		default:
			text = relativeTime(k.c.Now().Sub(value))
			if !o.Static {
				quantum := time.Minute
				if d := k.c.Now().Sub(value); d < time.Minute && d > -time.Minute {
					quantum = time.Second
				}
				k.c.After(quantum)
			}
		}
	}
	root := ui.Row(k.c).AlignItems(ui.Center).Gap(6).Children(func() {
		e := k.Caption(text)
		if !value.IsZero() {
			k.Tooltip(e, value.Format(time.RFC3339))
			if o.Copy {
				k.IconButton("复制时间", IconCopy, ButtonOptions{Variant: Ghost, Size: Small}).OnClick(func() { k.c.WriteClipboard(value.Format(time.RFC3339)) })
			}
		}
	})
	return root
}

// TimerState is an application-owned stopwatch model. Start/Pause use the
// caller's clock, which also makes tests deterministic.
type TimerState struct {
	Elapsed time.Duration
	Running bool
	started time.Time
}

func (s *TimerState) Start(now time.Time) {
	if !s.Running {
		s.started = now
		s.Running = true
	}
}
func (s *TimerState) Pause(now time.Time) {
	if s.Running {
		if !s.started.IsZero() {
			s.Elapsed += max(now.Sub(s.started), 0)
		}
		s.Running = false
		s.started = time.Time{}
	}
}
func (s *TimerState) Reset() { *s = TimerState{} }
func (s *TimerState) Duration(now time.Time) time.Duration {
	d := s.Elapsed
	if s.Running && !s.started.IsZero() {
		d += max(now.Sub(s.started), 0)
	}
	return d
}

type TimerOptions struct {
	Clock bool
	Label string
}

func (k *UI) Timer(state *TimerState, options ...TimerOptions) ui.Element {
	o := option(options)
	if state.Running && state.started.IsZero() {
		state.started = k.c.Now()
	}
	d := state.Duration(k.c.Now())
	secs := int64(d / time.Second)
	text := ""
	if o.Clock {
		if secs >= 3600 {
			text = fmt.Sprintf("%d:%02d:%02d", secs/3600, secs/60%60, secs%60)
		} else {
			text = fmt.Sprintf("%d:%02d", secs/60, secs%60)
		}
	} else {
		switch {
		case secs >= 3600:
			text = fmt.Sprintf("%dh %dm", secs/3600, secs/60%60)
		case secs >= 60:
			text = fmt.Sprintf("%dm %ds", secs/60, secs%60)
		default:
			text = fmt.Sprintf("%ds", secs)
		}
	}
	if state.Running {
		k.c.After(time.Second - d%time.Second)
	}
	e := k.Text(text).Font("monospace")
	if o.Label != "" {
		e.Label(o.Label).Value(text)
	}
	return e
}
