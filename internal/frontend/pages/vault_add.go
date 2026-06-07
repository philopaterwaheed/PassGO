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

type VaultAddAction struct {
	Back  bool
	Save  bool
	Vault state.Vault
}

type VaultAddPage struct {
	TitleEd    widget.Editor
	UsernameEd widget.Editor
	PasswordEd widget.Editor
	URLEd      widget.Editor
	NotesEd    widget.Editor

	SaveBtn widget.Clickable
	BackBtn widget.Clickable

	ShowPasswordBtn widget.Clickable
	ShowPassword    bool

	ErrorMsg  string
	IsLoading bool
	EditingID string
}

func NewVaultAddPage() *VaultAddPage {
	return &VaultAddPage{
		TitleEd:    widget.Editor{SingleLine: true, Submit: true},
		UsernameEd: widget.Editor{SingleLine: true, Submit: true},
		PasswordEd: widget.Editor{SingleLine: true, Submit: true, Mask: '*'},
		URLEd:      widget.Editor{SingleLine: true, Submit: true},
		NotesEd:    widget.Editor{SingleLine: false, Submit: false},
	}
}

func (p *VaultAddPage) Reset(v *state.Vault) {
	if v != nil {
		p.TitleEd.SetText(v.Title)
		p.UsernameEd.SetText(v.Username)
		p.PasswordEd.SetText(v.Password)
		p.URLEd.SetText(v.URL)
		p.NotesEd.SetText(v.Notes)
		p.EditingID = v.ID
	} else {
		p.TitleEd.SetText("")
		p.UsernameEd.SetText("")
		p.PasswordEd.SetText("")
		p.URLEd.SetText("")
		p.NotesEd.SetText("")
		p.EditingID = ""
	}
	p.ErrorMsg = ""
	p.IsLoading = false
	p.ShowPassword = false
	p.PasswordEd.Mask = '*'
}

func (p *VaultAddPage) TryBuildVault() (state.Vault, bool) {
	title := strings.TrimSpace(p.TitleEd.Text())
	if title == "" {
		p.ErrorMsg = "Title is required"
		return state.Vault{}, false
	}

	p.ErrorMsg = ""
	return state.Vault{
		Title:    title,
		Username: strings.TrimSpace(p.UsernameEd.Text()),
		Password: p.PasswordEd.Text(),
		URL:      strings.TrimSpace(p.URLEd.Text()),
		Notes:    strings.TrimSpace(p.NotesEd.Text()),
	}, true
}

func (p *VaultAddPage) Layout(gtx layout.Context, th *material.Theme) (layout.Dimensions, VaultAddAction) {
	var action VaultAddAction

	for p.ShowPasswordBtn.Clicked(gtx) {
		p.ShowPassword = !p.ShowPassword
		if p.ShowPassword {
			p.PasswordEd.Mask = 0
		} else {
			p.PasswordEd.Mask = '*'
		}
	}

	for p.SaveBtn.Clicked(gtx) {
		v, ok := p.TryBuildVault()
		if ok {
			action.Save = true
			action.Vault = v
		}
	}

	for p.BackBtn.Clicked(gtx) {
		action.Back = true
	}

	dims := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(520)))
		return ui.SurfacePanel(gtx, ui.SurfaceColor, ui.RadiusLarge, layout.UniformInset(unit.Dp(24)), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.ErrorMsg == "" {
						return layout.Dimensions{}
					}
					lbl := material.Body2(th, p.ErrorMsg)
					lbl.Color = ui.DangerColor
					return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, lbl.Layout)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, &p.TitleEd, "Title")
					e.TextSize = unit.Sp(16)
					return e.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, &p.UsernameEd, "Username")
					e.TextSize = unit.Sp(16)
					return e.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							e := material.Editor(th, &p.PasswordEd, "Password")
							e.TextSize = unit.Sp(16)
							return e.Layout(gtx)
						}),
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
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, &p.URLEd, "URL")
					e.TextSize = unit.Sp(16)
					return e.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, &p.NotesEd, "Notes")
					e.TextSize = unit.Sp(16)
					height := gtx.Dp(unit.Dp(96))
					gtx.Constraints.Min.Y = height
					gtx.Constraints.Max.Y = height
					return e.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btnText := "Add vault"
					if p.EditingID != "" {
						btnText = "Save changes"
					}
					if p.IsLoading {
						btnText = "Saving..."
					}
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return ui.PrimaryButton(gtx, th, &p.SaveBtn, btnText)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return ui.SecondaryButton(gtx, th, &p.BackBtn, "Back")
				}),
			)
		})
	})

	return dims, action
}
