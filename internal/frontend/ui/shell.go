package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/state"
)

type Shell struct {
	VaultBtn    widget.Clickable
	SettingsBtn widget.Clickable
	LogoutBtn   widget.Clickable

	// internal lists for nav
	railList  widget.List
	bottomRow [2]widget.Clickable
}

func NewShell() *Shell {
	return &Shell{
		railList: widget.List{List: layout.List{Axis: layout.Vertical}},
	}
}

func (s *Shell) Layout(gtx layout.Context, th *material.Theme, st *state.AppState, title string, body layout.Widget) layout.Dimensions {
	compact := IsCompact(gtx)
	if compact {
		return s.layoutCompact(gtx, th, st, title, body)
	}
	return s.layoutDesktop(gtx, th, st, title, body)
}

func (s *Shell) layoutDesktop(gtx layout.Context, th *material.Theme, st *state.AppState, title string, body layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(unit.Dp(220))
			gtx.Constraints.Max.X = gtx.Dp(unit.Dp(220))
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.H4(th, "PassGO")
						lbl.Color = th.Palette.ContrastBg
						return layout.Inset{Bottom: unit.Dp(24), Left: unit.Dp(12)}.Layout(gtx, lbl.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return s.navButton(gtx, th, &s.VaultBtn, "Vault", st.Nav == state.NavVault)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return s.navButton(gtx, th, &s.SettingsBtn, "Settings", st.Nav == state.NavSettings)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: gtx.Constraints.Min}
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !st.IsAuthed() {
							return layout.Dimensions{}
						}
						btn := material.Button(th, &s.LogoutBtn, "Log out")
						btn.TextSize = unit.Sp(14)
						btn.Background = color.NRGBA{0, 0, 0, 0}
						btn.Color = th.Palette.Fg
						btn.Inset = layout.UniformInset(unit.Dp(12))
						return btn.Layout(gtx)
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								lbl := material.H5(th, title)
								return lbl.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
					layout.Flexed(1, body),
				)
			})
		}),
	)
}

func (s *Shell) layoutCompact(gtx layout.Context, th *material.Theme, st *state.AppState, title string, body layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.H5(th, title)
						lbl.Color = th.Palette.ContrastBg
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !st.IsAuthed() {
							return layout.Dimensions{}
						}
						btn := material.Button(th, &s.LogoutBtn, "Log out")
						btn.TextSize = unit.Sp(14)
						btn.Background = color.NRGBA{0, 0, 0, 0}
						btn.Color = th.Palette.Fg
						return btn.Layout(gtx)
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, body)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Simple bottom nav when authed, otherwise no chrome.
			if !st.IsAuthed() {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return s.navButton(gtx, th, &s.VaultBtn, "Vault", st.Nav == state.NavVault)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return s.navButton(gtx, th, &s.SettingsBtn, "Settings", st.Nav == state.NavSettings)
					}),
				)
			})
		}),
	)
}

func (s *Shell) navButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, label string, selected bool) layout.Dimensions {
	btn := material.Button(th, click, label)
	btn.TextSize = unit.Sp(16)
	btn.Inset = layout.UniformInset(unit.Dp(12))
	if selected {
		btn.CornerRadius = unit.Dp(8)
	} else {
		btn.Background = color.NRGBA{0, 0, 0, 0}
		btn.Color = th.Palette.Fg
	}
	return btn.Layout(gtx)
}
