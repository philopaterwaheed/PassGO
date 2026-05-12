package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/state"
)

type VaultListAction struct {
	Add    bool
	OpenID string
}

type VaultListPage struct {
	AddBtn widget.Clickable
	List   widget.List

	cardClicks []widget.Clickable
}

func NewVaultListPage() *VaultListPage {
	return &VaultListPage{
		List: widget.List{List: layout.List{Axis: layout.Vertical}},
	}
}

func (p *VaultListPage) Layout(gtx layout.Context, th *material.Theme, vaults []state.Vault) (layout.Dimensions, VaultListAction) {
	var action VaultListAction

	if p.AddBtn.Clicked(gtx) {
		action.Add = true
	}

	ensureClickables(&p.cardClicks, len(vaults))

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &p.AddBtn, "Add")
					btn.TextSize = unit.Sp(14)
					return btn.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if len(vaults) == 0 {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(520)))
					return material.Body1(th, "No vaults yet. Tap Add to create one.").Layout(gtx)
				})
			}

			return p.List.Layout(gtx, len(vaults), func(gtx layout.Context, i int) layout.Dimensions {
				v := vaults[i]
				click := &p.cardClicks[i]

				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					d := click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return card(gtx, th, v)
					})
					if click.Clicked(gtx) {
						action.OpenID = v.ID
					}
					return d
				})
			})
		}),
	)

	return dims, action
}

func ensureClickables(dst *[]widget.Clickable, n int) {
	if len(*dst) >= n {
		return
	}
	for i := len(*dst); i < n; i++ {
		*dst = append(*dst, widget.Clickable{})
	}
}

func card(gtx layout.Context, th *material.Theme, v state.Vault) layout.Dimensions {
	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		border := widget.Border{Color: th.Palette.Fg, CornerRadius: unit.Dp(14), Width: unit.Dp(1)}
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.H6(th, v.Title)
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if v.Username == "" {
							return layout.Dimensions{}
						}
						return material.Body2(th, v.Username).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if v.URL == "" {
							return layout.Dimensions{}
						}
						return material.Caption(th, v.URL).Layout(gtx)
					}),
				)
			})
		})
	})
}
