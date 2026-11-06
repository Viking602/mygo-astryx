package astryx

import "github.com/egoist/mygo/ui"

// TopNavMegaMenu composes rich menu items with an optional featured panel.
// Pointer and keyboard users get the same content; narrow windows use a drawer.
func (k *UI) TopNavMegaMenu(label string, items, featured func(), options ...HoverCardOptions) ui.Element {
	o := option(options)
	if o.Width <= 0 {
		o.Width = 640
	}
	o.Label = label
	o.OpenOnClick = true
	trigger := k.Button(label, ButtonOptions{Variant: Outline, EndIcon: IconChevronDown, Disabled: o.Disabled}).AlignSelf(ui.Start).Label(label).Role(ui.RoleMenuButton)
	returnFocus := ui.Local(trigger, "mega-menu-focus", func() ui.Handle { return ui.Handle{} })
	trigger.Bind(returnFocus)
	var closeMenu func()
	build := func() {
		previous := k.closeMenu
		k.closeMenu = closeMenu
		defer func() { k.closeMenu = previous }()
		w, _ := k.c.Size()
		row := ui.Row(k.c).Gap(24).AlignItems(ui.Stretch)
		if w < 700 {
			row.Column()
		}
		row.Children(func() {
			ui.Column(k.c).Grow(1).Gap(8).Children(items)
			if featured != nil {
				ui.Column(k.c).Grow(1).Gap(8).Padding(14).Background(k.tokens.Subtle).Radius(k.tokens.Radius).Children(featured)
			}
		})
	}
	w, _ := k.c.Size()
	if w < 700 {
		open := ui.Local(trigger, "mega-menu-mobile", func() bool { return false })
		if trigger.Clicked() && !o.Disabled {
			*open = !*open
		}
		if o.Disabled {
			*open = false
		}
		trigger.Expanded(*open)
		closeMenu = func() { *open = false; returnFocus.Focus() }
		k.Drawer(label, open, build, DialogOptions{Width: 380})
	} else {
		state := hoverState(trigger)
		closeMenu = func() { state.Open = false; state.Suppressed = true; returnFocus.Focus() }
		o.FocusOnOpen = trigger.FocusVisible()
		panel := k.HoverCard(trigger, build, o)
		trigger.Expanded(panel.Valid())
	}
	return trigger
}

// TopNavMegaMenuItem closes its containing mega menu before running the host
// action. Ordinary Item remains suitable for controls that should keep it open.
func (k *UI) TopNavMegaMenuItem(label string, options ...ItemOptions) ui.Element {
	o := option(options)
	onClick, closeMenu := o.OnClick, k.closeMenu
	o.OnClick = func() {
		if closeMenu != nil {
			closeMenu()
		}
		if onClick != nil {
			onClick()
		}
	}
	return k.Item(label, o)
}
