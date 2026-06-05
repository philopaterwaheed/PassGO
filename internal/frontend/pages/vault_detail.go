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

	content := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.H5(th, v.Title).Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions { return field(th, gtx, "Username", v.Username) }),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if passStr == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.Caption(th, "Password").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(1, material.Body1(th, passStr).Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btnTxt := "Show"
								if p.ShowPassword {
									btnTxt = "Hide"
								}
								btn := material.Button(th, &p.ShowPasswordBtn, btnTxt)
								btn.TextSize = unit.Sp(12)
								btn.Inset = layout.UniformInset(unit.Dp(6))
								btn.Background = th.Palette.ContrastBg
								btn.Color = th.Palette.ContrastFg
								return btn.Layout(gtx)
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
			gtx2 := gtx
			gtx2.Constraints.Min.X = gtx2.Constraints.Max.X
			return ui.OutlinedButton(gtx2, th, &p.BackBtn, "Back", th.Palette.Fg, th.Palette.ContrastBg)
		}),
	)

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return content
		})
	}), action
}

func field(th *material.Theme, gtx layout.Context, label string, value string) layout.Dimensions {
	if strings.TrimSpace(value) == "" {
		return layout.Dimensions{}
	}
	return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(material.Caption(th, label).Layout),
			layout.Rigid(material.Body1(th, value).Layout),
		)
	})
}
