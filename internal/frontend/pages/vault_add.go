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

	ErrorMsg string
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

func (p *VaultAddPage) Reset() {
	p.TitleEd.SetText("")
	p.UsernameEd.SetText("")
	p.PasswordEd.SetText("")
	p.URLEd.SetText("")
	p.NotesEd.SetText("")
	p.ErrorMsg = ""
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

	gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(520)))

	content := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.ErrorMsg == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, material.Body2(th, p.ErrorMsg).Layout)
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
			e := material.Editor(th, &p.PasswordEd, "Password")
			e.TextSize = unit.Sp(16)
			return e.Layout(gtx)
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
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(96))
			return e.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, &p.SaveBtn, "Save")
			btn.TextSize = unit.Sp(16)
			return btn.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
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
