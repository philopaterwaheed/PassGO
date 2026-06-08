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
	CopyUsernameBtn widget.Clickable
	CopyPasswordBtn widget.Clickable
	CopyURLBtn      widget.Clickable
	CopyNotesBtn    widget.Clickable
	ShowPassword    bool
	CopiedField     string
}

func NewVaultDetailPage() *VaultDetailPage { return &VaultDetailPage{} }

func (p *VaultDetailPage) Reset() {
	p.ShowPassword = false
	p.CopiedField = ""
}

func (p *VaultDetailPage) Layout(gtx layout.Context, th *material.Theme, v state.Vault, loading bool, loadError string) (layout.Dimensions, VaultDetailAction) {
	var action VaultDetailAction

	for p.BackBtn.Clicked(gtx) {
		action.Back = true
	}

	for p.ShowPasswordBtn.Clicked(gtx) {
		p.ShowPassword = !p.ShowPassword
	}
	for p.CopyUsernameBtn.Clicked(gtx) {
		ui.CopyText(gtx, v.Username)
		p.CopiedField = "username"
	}
	for p.CopyPasswordBtn.Clicked(gtx) {
		ui.CopyText(gtx, v.Password)
		p.CopiedField = "password"
	}
	for p.CopyURLBtn.Clicked(gtx) {
		ui.CopyText(gtx, v.URL)
		p.CopiedField = "url"
	}
	for p.CopyNotesBtn.Clicked(gtx) {
		ui.CopyText(gtx, v.Notes)
		p.CopiedField = "notes"
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
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if loading {
						return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, ui.MutedLabel(th, material.Body2(th, "Decrypting vault...")).Layout)
					}
					if loadError != "" {
						errLbl := material.Body2(th, loadError)
						errLbl.Color = ui.DangerColor
						return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, errLbl.Layout)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.copyableField(gtx, th, "Username", v.Username, &p.CopyUsernameBtn, "username")
				}),
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
									layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return ui.CopyButton(gtx, th, &p.CopyPasswordBtn, p.CopiedField == "password")
									}),
								)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.copyableField(gtx, th, "URL", v.URL, &p.CopyURLBtn, "url")
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.copyableField(gtx, th, "Notes", v.Notes, &p.CopyNotesBtn, "notes")
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return ui.SecondaryButton(gtx, th, &p.BackBtn, "Back")
				}),
			)
		})
	}), action
}

func (p *VaultDetailPage) copyableField(gtx layout.Context, th *material.Theme, label string, value string, click *widget.Clickable, fieldKey string) layout.Dimensions {
	if strings.TrimSpace(value) == "" {
		return layout.Dimensions{}
	}
	return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(ui.MutedLabel(th, material.Caption(th, label)).Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, material.Body1(th, value).Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return ui.CopyButton(gtx, th, click, p.CopiedField == fieldKey)
					}),
				)
			}),
		)
	})
}
