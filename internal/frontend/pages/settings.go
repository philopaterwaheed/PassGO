package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SettingsPage struct {
	DarkMode widget.Bool

	initialized bool
}

func NewSettingsPage() *SettingsPage { return &SettingsPage{} }

func (p *SettingsPage) Layout(gtx layout.Context, th *material.Theme, email string, apiBaseURL string, darkMode *bool) (layout.Dimensions, bool) {
	if darkMode != nil && !p.initialized {
		p.DarkMode.Value = *darkMode
		p.initialized = true
	}

	changed := false

	dims := layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(material.Body1(th, "Account").Layout),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(material.Body2(th, "Signed in as: "+email).Layout),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
			layout.Rigid(material.Body1(th, "Appearance").Layout),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				prev := p.DarkMode.Value
				sw := material.Switch(th, &p.DarkMode, "Dark theme")
				d := sw.Layout(gtx)
				if darkMode != nil && prev != p.DarkMode.Value {
					*darkMode = p.DarkMode.Value
					changed = true
				}
				return d
			}),
		)
	})

	return dims, changed
}
