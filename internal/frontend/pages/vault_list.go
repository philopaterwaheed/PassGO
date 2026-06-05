package pages

import (
	"image"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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

	for p.AddBtn.Clicked(gtx) {
		action.Add = true
	}

	ensureClickables(&p.cardClicks, len(vaults))

	for i := range vaults {
		for p.cardClicks[i].Clicked(gtx) {
			action.OpenID = vaults[i].ID
		}
	}

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
					return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return card(gtx, th, v)
					})
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
	return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		radius := gtx.Dp(unit.Dp(16))

		// Background color (very transparent contrast background)
		bg := th.Palette.ContrastBg
		bg.A = 15

		content := func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// Avatar circle
						size := gtx.Dp(unit.Dp(48))
						
						initial := "?"
						if len(v.Title) > 0 {
							initial = strings.ToUpper(string(v.Title[0]))
						}
						
						return layout.Stack{}.Layout(gtx,
							layout.Expanded(func(gtx layout.Context) layout.Dimensions {
								bounds := image.Rectangle{Max: image.Point{X: size, Y: size}}
								defer clip.UniformRRect(bounds, size/2).Push(gtx.Ops).Pop()
								
								avatarBg := th.Palette.ContrastBg
								avatarBg.A = 180
								paint.Fill(gtx.Ops, avatarBg)
								return layout.Dimensions{Size: bounds.Max}
							}),
							layout.Stacked(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min = image.Point{X: size, Y: size}
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.H6(th, initial)
									lbl.Color = th.Palette.ContrastFg
									lbl.Font.Weight = font.Bold
									return lbl.Layout(gtx)
								})
							}),
						)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body1(th, v.Title)
								lbl.Font.Weight = font.Bold
								return lbl.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if v.Username == "" {
									return layout.Dimensions{}
								}
								lbl := material.Body2(th, v.Username)
								c := th.Palette.Fg
								c.A = 180
								lbl.Color = c
								return lbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if v.URL == "" {
									return layout.Dimensions{}
								}
								lbl := material.Caption(th, v.URL)
								c := th.Palette.Fg
								c.A = 150
								lbl.Color = c
								return lbl.Layout(gtx)
							}),
						)
					}),
				)
			})
		}

		mac := op.Record(gtx.Ops)
		dims := content(gtx)
		call := mac.Stop()

		defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, radius).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, bg)

		border := widget.Border{
			Color:        th.Palette.ContrastBg,
			CornerRadius: unit.Dp(16),
			Width:        unit.Dp(1.5),
		}
		border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: dims.Size}
		})

		call.Add(gtx.Ops)
		return dims
	})
}
