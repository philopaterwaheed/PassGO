package pages

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/state"
	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type VaultDetailAction struct {
	Back bool
}

type VaultDetailPage struct {
	BackBtn         widget.Clickable
	ShowPasswordBtn widget.Clickable
	ShowPassword    bool
}

func NewVaultDetailPage() *VaultDetailPage { return &VaultDetailPage{} }

func (p *VaultDetailPage) Reset() {
	p.ShowPassword = false
}

func (p *VaultDetailPage) Layout(gtx layout.Context, th *material.Theme, v state.Vault) (layout.Dimensions, VaultDetailAction) {
	var action VaultDetailAction

	for p.BackBtn.Clicked(gtx) {
		action.Back = true
	}

	for p.ShowPasswordBtn.Clicked(gtx) {
		p.ShowPassword = !p.ShowPassword
	}

	passStr := ""
	if v.Password != "" {
		if p.ShowPassword {
			passStr = v.Password
		} else {
			passStr = strings.Repeat("•", 10)
		}
	}

	gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(520)))

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return ui.SurfacePanel(gtx, ui.SurfaceColor, ui.RadiusLarge, layout.UniformInset(unit.Dp(24)), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H5(th, v.Title)
					lbl.Color = th.Palette.ContrastBg
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return field(th, gtx, "Username", v.Username) }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if passStr == "" {
						return layout.Dimensions{}
					}
					return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(ui.MutedLabel(th, material.Caption(th, "Password")).Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Flexed(1, material.Body1(th, passStr).Layout),
									layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										btnTxt := "Show"
										if p.ShowPassword {
											btnTxt = "Hide"
										}
										return ui.GhostButton(gtx, th, &p.ShowPasswordBtn, btnTxt)
									}),
								)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return field(th, gtx, "URL", v.URL) }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions { return field(th, gtx, "Notes", v.Notes) }),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return ui.SecondaryButton(gtx, th, &p.BackBtn, "Back")
				}),
			)
		})
	}), action
}

func field(th *material.Theme, gtx layout.Context, label string, value string) layout.Dimensions {
	if strings.TrimSpace(value) == "" {
		return layout.Dimensions{}
	}
	return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(ui.MutedLabel(th, material.Caption(th, label)).Layout),
			layout.Rigid(material.Body1(th, value).Layout),
		)
	})
}
