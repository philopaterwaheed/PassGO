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
	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type VaultListAction struct {
	Add            bool
	Retry          bool
	Unlock         bool
	UnlockPassword string
	OpenID         string
	EditID         string
	DeleteID       string
}

type VaultListPage struct {
	AddBtn              widget.Clickable
	RetryBtn            widget.Clickable
	UnlockBtn           widget.Clickable
	MasterPasswordInput widget.Editor
	List                widget.List
	UnlockError         string

	rows map[string]*vaultListRow
}

type vaultListRow struct {
	cardClick    widget.Clickable
	editClick    widget.Clickable
	deleteClick  widget.Clickable
	showClick    widget.Clickable
	urlClick     widget.Clickable
	urlOpenClick widget.Clickable
	showPassword bool
}

func NewVaultListPage() *VaultListPage {
	return &VaultListPage{
		List: widget.List{List: layout.List{Axis: layout.Vertical}},
		MasterPasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
	}
}

func (p *VaultListPage) Action(gtx layout.Context, vaults []state.Vault, locked bool) VaultListAction {
	var action VaultListAction

	for p.AddBtn.Clicked(gtx) {
		action.Add = true
	}
	for p.RetryBtn.Clicked(gtx) {
		action.Retry = true
	}
	for p.UnlockBtn.Clicked(gtx) {
		action.Unlock = true
		action.UnlockPassword = p.MasterPasswordInput.Text()
	}

	if locked {
		return action
	}

	p.pruneRows(vaults)

	for i := range vaults {
		row := p.row(vaults[i].ID)

		for row.showClick.Clicked(gtx) {
			if !vaults[i].Decrypted {
				action.OpenID = vaults[i].ID
				continue
			}
			row.showPassword = !row.showPassword
		}
		for row.deleteClick.Clicked(gtx) {
			action.DeleteID = vaults[i].ID
		}
		for row.editClick.Clicked(gtx) {
			action.EditID = vaults[i].ID
		}
		for row.cardClick.Clicked(gtx) {
			action.OpenID = vaults[i].ID
		}
		for row.urlClick.Clicked(gtx) {
			if url := ui.NormalizeURL(vaults[i].URL); url != "" {
				ui.OpenURLLogged(url)
			}
		}
		for row.urlOpenClick.Clicked(gtx) {
			if url := ui.NormalizeURL(vaults[i].URL); url != "" {
				ui.OpenURLLogged(url)
			}
		}
	}

	return action
}

func (p *VaultListPage) Layout(gtx layout.Context, th *material.Theme, vaults []state.Vault, loading bool, loadError string, locked bool) (layout.Dimensions, VaultListAction) {
	action := p.Action(gtx, vaults, locked)

	if locked {
		return p.unlockLayout(gtx, th), action
	}

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{} }),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.PrimaryButton(gtx, th, &p.AddBtn, "Add")
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if loading && len(vaults) == 0 {
				return loadingState(gtx, th)
			}

			if loadError != "" && len(vaults) == 0 {
				return loadErrorState(gtx, th, &p.RetryBtn, loadError)
			}

			if len(vaults) == 0 {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(520)))
					return ui.MutedLabel(th, material.Body1(th, "No vaults yet. Tap Add to create one.")).Layout(gtx)
				})
			}

			return p.List.Layout(gtx, len(vaults), func(gtx layout.Context, i int) layout.Dimensions {
				v := vaults[i]
				row := p.row(v.ID)

				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return card(gtx, th, v, &row.cardClick, &row.editClick, &row.deleteClick, &row.showClick, &row.urlClick, &row.urlOpenClick, row.showPassword)
				})
			})
		}),
	)

	return dims, action
}

func (p *VaultListPage) row(id string) *vaultListRow {
	if p.rows == nil {
		p.rows = make(map[string]*vaultListRow)
	}
	row := p.rows[id]
	if row == nil {
		row = &vaultListRow{}
		p.rows[id] = row
	}
	return row
}

func (p *VaultListPage) pruneRows(vaults []state.Vault) {
	if len(p.rows) == 0 {
		return
	}

	seen := make(map[string]struct{}, len(vaults))
	for _, v := range vaults {
		seen[v.ID] = struct{}{}
	}
	for id := range p.rows {
		if _, ok := seen[id]; !ok {
			delete(p.rows, id)
		}
	}
}

func (p *VaultListPage) unlockLayout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(420)))

		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.H6(th, "Unlock vault")
				lbl.Font.Weight = font.Bold
				return layout.Center.Layout(gtx, lbl.Layout)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		}

		if p.UnlockError != "" {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, p.UnlockError)
					lbl.Color = ui.DangerColor
					return layout.Center.Layout(gtx, lbl.Layout)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			)
		}

		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, &p.MasterPasswordInput, "Master password")
				e.TextSize = unit.Sp(16)
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return ui.PrimaryButton(gtx, th, &p.UnlockBtn, "Unlock")
			}),
		)

		return ui.SurfacePanel(gtx, ui.SurfaceColor, ui.RadiusLarge, layout.UniformInset(unit.Dp(24)), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, children...)
		})
	})
}

func loadErrorState(gtx layout.Context, th *material.Theme, retryBtn *widget.Clickable, message string) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(520)))
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.H6(th, "Could not load vaults")
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, message)
				lbl.Color = ui.MutedColor
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.PrimaryButton(gtx, th, retryBtn, "Retry")
			}),
		)
	})
}

func loadingState(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				size := gtx.Dp(unit.Dp(36))
				gtx.Constraints.Min = image.Pt(size, size)
				gtx.Constraints.Max = image.Pt(size, size)
				return material.Loader(th).Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, "Loading vaults...")
				lbl.Color = ui.MutedColor
				return lbl.Layout(gtx)
			}),
		)
	})
}

func card(gtx layout.Context, th *material.Theme, v state.Vault, openBtn, editBtn, delBtn, showBtn, urlBtn, urlOpenBtn *widget.Clickable, isShowing bool) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		radius := gtx.Dp(ui.RadiusLarge)

		// Background color (very transparent contrast background)
		bg := ui.SurfaceColor

		content := func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				mainRow := func(gtx layout.Context) layout.Dimensions {
					return vaultSummary(gtx, th, v, openBtn, showBtn, urlBtn, isShowing)
				}

				actions := func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if v.URL == "" {
								return layout.Dimensions{}
							}
							return ui.SecondaryButton(gtx, th, urlOpenBtn, "Open")
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if v.URL == "" {
								return layout.Dimensions{}
							}
							return layout.Spacer{Width: unit.Dp(6)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if v.Decrypted && v.Password == "" {
								return layout.Dimensions{}
							}
							btnTxt := "Show"
							if v.Decrypted && isShowing {
								btnTxt = "Hide"
							}
							return ui.SecondaryButton(gtx, th, showBtn, btnTxt)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if v.Decrypted && v.Password == "" {
								return layout.Dimensions{}
							}
							return layout.Spacer{Width: unit.Dp(6)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return ui.SecondaryButton(gtx, th, editBtn, "Edit")
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, delBtn, "Delete")
							btn.Background = ui.DangerColor
							btn.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
							btn.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}
							btn.TextSize = unit.Sp(13)
							btn.CornerRadius = ui.RadiusSmall
							return btn.Layout(gtx)
						}),
					)
				}

				if ui.IsCompact(gtx) {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(mainRow),
						layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
						layout.Rigid(actions),
					)
				}

				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, mainRow),
					layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
					layout.Rigid(actions),
				)
			})
		}

		mac := op.Record(gtx.Ops)
		dims := content(gtx)
		call := mac.Stop()

		defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, radius).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, bg)

		border := widget.Border{
			Color:        ui.BorderColor,
			CornerRadius: ui.RadiusLarge,
			Width:        unit.Dp(1),
		}
		border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: dims.Size}
		})

		call.Add(gtx.Ops)
		return dims
	})
}

func vaultAvatar(gtx layout.Context, th *material.Theme, title string) layout.Dimensions {
	size := gtx.Dp(unit.Dp(48))

	initial := "?"
	if len(title) > 0 {
		initial = strings.ToUpper(string(title[0]))
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
}

func vaultSummary(gtx layout.Context, th *material.Theme, v state.Vault, openBtn, showBtn, urlBtn *widget.Clickable, isShowing bool) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return openBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = max(gtx.Constraints.Min.Y, gtx.Dp(unit.Dp(56)))
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return vaultAvatar(gtx, th, v.Title)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return vaultIdentity(gtx, th, v)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if v.URL == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return urlBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					gtx.Constraints.Min.Y = max(gtx.Constraints.Min.Y, gtx.Dp(unit.Dp(44)))
					return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(th, v.URL)
						lbl.Color = th.Palette.ContrastBg
						return lbl.Layout(gtx)
					})
				})
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !v.Decrypted {
				return layout.Dimensions{}
			}
			if v.Password == "" && v.Decrypted {
				return layout.Dimensions{}
			}

			txt := strings.Repeat("•", 10)
			if isShowing && v.Decrypted {
				txt = v.Password
			}

			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, txt)
					return lbl.Layout(gtx)
				}),
			)
		}),
	)
}

func vaultIdentity(gtx layout.Context, th *material.Theme, v state.Vault) layout.Dimensions {
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
			lbl.Color = ui.MutedColor
			return lbl.Layout(gtx)
		}),
	)
}
