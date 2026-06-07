package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type WelcomePage struct {
	LoginBtn     widget.Clickable
	RegisterBtn  widget.Clickable
	GitHubBtn    widget.Clickable
	LinkedInBtn  widget.Clickable
	PortfolioBtn widget.Clickable
	CVBtn        widget.Clickable
}

func NewWelcomePage() *WelcomePage {
	return &WelcomePage{}
}

func (p *WelcomePage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for p.GitHubBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.GitHubURL)
	}
	for p.LinkedInBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.LinkedInURL)
	}
	for p.PortfolioBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.PortfolioURL)
	}
	for p.CVBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.CVURL)
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return ui.ConstrainWidth(gtx, unit.Dp(440), func(gtx layout.Context) layout.Dimensions {
			return ui.SurfacePanel(gtx, ui.SurfaceColor, ui.RadiusLarge, layout.UniformInset(unit.Dp(24)), func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.H4(th, "PassGO")
						lbl.Color = th.Palette.ContrastBg
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
					layout.Rigid(ui.MutedLabel(th, material.Body1(th, "Minimal, fast password manager")).Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(22)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return ui.PrimaryButton(gtx, th, &p.LoginBtn, "Log in")
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return ui.SecondaryButton(gtx, th, &p.RegisterBtn, "Create account")
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(22)}.Layout),
					layout.Rigid(p.creditBlock(th)),
				)
			})
		})
	})
}

func (p *WelcomePage) creditBlock(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(th, "Philopater Waheed | philopaterwaheed9@gmail.com")
				lbl.Color = ui.MutedColor
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return creditLinks(gtx, th, &p.GitHubBtn, &p.LinkedInBtn, &p.PortfolioBtn, &p.CVBtn)
			}),
		)
	}
}

func creditLinks(gtx layout.Context, th *material.Theme, githubBtn, linkedInBtn, portfolioBtn, cvBtn *widget.Clickable) layout.Dimensions {
	gtx.Constraints.Min.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(440)))

	if ui.IsCompact(gtx) {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(300)))
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return creditLinkRow(gtx, th,
						creditLink{click: githubBtn, icon: "GH", label: "GitHub"},
						creditLink{click: linkedInBtn, icon: "IN", label: "LinkedIn"},
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return creditLinkRow(gtx, th,
						creditLink{click: portfolioBtn, icon: "PF", label: "Portfolio"},
						creditLink{click: cvBtn, icon: "CV", label: "Resume"},
					)
				}),
			)
		})
	}

	row := func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(360)))
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.LinkButton(gtx, th, githubBtn, "GH", "GitHub")
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.LinkButton(gtx, th, linkedInBtn, "IN", "LinkedIn")
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.LinkButton(gtx, th, portfolioBtn, "PF", "Portfolio")
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.LinkButton(gtx, th, cvBtn, "CV", "Resume")
				}),
			)
		})
	}

	return layout.Center.Layout(gtx, row)
}

type creditLink struct {
	click *widget.Clickable
	icon  string
	label string
}

func creditLinkRow(gtx layout.Context, th *material.Theme, left, right creditLink) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.CompactLinkButton(gtx, th, left.click, left.icon, left.label)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.CompactLinkButton(gtx, th, right.click, right.icon, right.label)
			}),
		)
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
