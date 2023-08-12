package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/app"
)

type App struct {
	window *app.Window
	client *Client
	theme  *material.Theme
	date   time.Time
	sum    float64
}

func (a *App) reload() {
	log.Printf("Loading data for %s", a.date.String())
	a.showLoader()
	a.showOnlineSum()
}

func (a *App) showLoader() {
	a.sum = -1
	a.window.Invalidate()
}

func (a *App) showOnlineSum() {
	checks, err := a.client.GetChecksForDay(a.date.Format("2006-01-02"))
	if err != nil {
		log.Fatal(err)
	}

	online := 0
	for _, x := range checks {
		if (x.Status == "closed" || x.Status == "pre_closed") && x.PayBy == "online" {
			online += x.Total
		}
	}

	a.sum = float64(online) / 100
	a.window.Invalidate()
}

func (a *App) setDate(d time.Time) {
	a.date = d
	go a.reload()
}

func (a *App) Loop() error {
	go a.reload()

	var ops op.Ops
	prevButton := &widget.Clickable{}
	nextButton := &widget.Clickable{}

	for e := range a.window.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			if prevButton.Clicked() {
				a.setDate(a.date.Add(-24 * time.Hour))
			}

			if nextButton.Clicked() {
				a.setDate(a.date.Add(24 * time.Hour))
			}

			layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(layout.Spacer{Height: 30}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Start}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return IconButton(gtx, prevButton, a.theme, "<-")
							}),
							layout.Rigid(layout.Spacer{Width: 10}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									t := material.Body2(a.theme, a.date.Format("2006-01-02"))
									t.Alignment = text.Middle
									t.TextSize = unit.Sp(18)
									return t.Layout(gtx)
								})
							}),
							layout.Rigid(layout.Spacer{Width: 10}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return IconButton(gtx, nextButton, a.theme, "->")
							}),
						)
					})
				}),
				layout.Rigid(layout.Spacer{Height: 50}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if a.sum == -1 {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max = image.Point{X: gtx.Dp(56), Y: gtx.Dp(56)}
							l := material.Loader(a.theme)
							return l.Layout(gtx)
						})
					} else {
						t := material.Body1(a.theme, fmt.Sprintf("Online: %.2f", a.sum))
						t.Alignment = text.Middle
						t.TextSize = unit.Sp(32)
						return t.Layout(gtx)
					}
				}),
			)

			e.Frame(gtx.Ops)
		}
	}

	return nil
}

func NewApp(w *app.Window) (*App, error) {
	log.Println("Initialising UI")

	th := material.NewTheme()
	th.Face = "monospace"
	log.Println("Creating a client")

	c, err := NewClient()
	if err != nil {
		return nil, err
	}

	log.Println("Authenticating")

	if err = c.Auth(); err != nil {
		return nil, err
	}

	return &App{
		window: w,
		client: c,
		theme:  th,
		date:   time.Now(),
		sum:    -1,
	}, nil
}
