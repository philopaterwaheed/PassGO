package pages

import (
	"image"
	"image/color"
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
	Add      bool
	OpenID   string
	EditID   string
	DeleteID string
}

type VaultListPage struct {
	AddBtn widget.Clickable
	List   widget.List

	cardClicks   []widget.Clickable
	editClicks   []widget.Clickable
	deleteClicks []widget.Clickable
	showClicks   []widget.Clickable
	showPassword []bool
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
	ensureClickables(&p.editClicks, len(vaults))
	ensureClickables(&p.deleteClicks, len(vaults))
	ensureClickables(&p.showClicks, len(vaults))
	ensureBools(&p.showPassword, len(vaults))

	for i := range vaults {
		for p.showClicks[i].Clicked(gtx) {
			p.showPassword[i] = !p.showPassword[i]
		}
		for p.deleteClicks[i].Clicked(gtx) {
			action.DeleteID = vaults[i].ID
		}
		for p.editClicks[i].Clicked(gtx) {
			action.EditID = vaults[i].ID
		}
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
				editClick := &p.editClicks[i]
				deleteClick := &p.deleteClicks[i]
				showClick := &p.showClicks[i]
				isShowing := p.showPassword[i]

				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return card(gtx, th, v, editClick, deleteClick, showClick, isShowing)
					})
				})
			})
		}),
	)

	return dims, action
}

func ensureBools(dst *[]bool, n int) {
	if len(*dst) >= n {
		return
	}
	for i := len(*dst); i < n; i++ {
		*dst = append(*dst, false)
	}
}

func ensureClickables(dst *[]widget.Clickable, n int) {
	if len(*dst) >= n {
		return
	}
	for i := len(*dst); i < n; i++ {
		*dst = append(*dst, widget.Clickable{})
	}
}

func card(gtx layout.Context, th *material.Theme, v state.Vault, editBtn, delBtn, showBtn *widget.Clickable, isShowing bool) layout.Dimensions {
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
							layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if v.Password == "" {
									return layout.Dimensions{}
								}
								
								txt := strings.Repeat("•", 10)
								if isShowing {
									txt = v.Password
								}

								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										lbl := material.Body2(th, txt)
										return lbl.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										btnTxt := "Show"
										if isShowing {
											btnTxt = "Hide"
										}
										btn := material.Button(th, showBtn, btnTxt)
										btn.TextSize = unit.Sp(10)
										btn.Inset = layout.UniformInset(unit.Dp(4))
										btn.Background = th.Palette.ContrastBg
										btn.Color = th.Palette.ContrastFg
										return btn.Layout(gtx)
									}),
								)
							}),
						)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, editBtn, "Edit")
						btn.Background = color.NRGBA{R: 50, G: 150, B: 200, A: 255}
						btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						btn.Inset = layout.UniformInset(unit.Dp(8))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, delBtn, "Delete")
						btn.Background = color.NRGBA{R: 200, G: 50, B: 50, A: 255}
						btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						btn.Inset = layout.UniformInset(unit.Dp(8))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
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
