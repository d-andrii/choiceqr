package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(200, 200), app.Title("Choice Online"))
		if err := run(w); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	app.Main()
}

func getOnlineSum(c *Client, date string) float64 {
	checks, err := c.GetChecksForDay(date)
	if err != nil {
		log.Fatal(err)
	}

	online := 0
	for _, x := range checks {
		if x.Status == "closed" && x.PayBy == "online" {
			online += x.Total
		}
	}

	return float64(online) / 100
}

func run(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())

	date := time.Now()
	sum := -1.

	c, err := NewClient()
	if err != nil {
		return err
	}

	err = c.Auth()
	if err != nil {
		return err
	}

	reload := func() {
		sum = -1
		w.Invalidate()
		sum = getOnlineSum(c, date.Format("2006-01-02"))
		w.Invalidate()
	}

	prevButton := &widget.Clickable{}
	nextButton := &widget.Clickable{}

	go reload()

	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			if prevButton.Clicked() {
				date = date.Add(-24 * time.Hour)
				go reload()
			}

			if nextButton.Clicked() {
				date = date.Add(24 * time.Hour)
				go reload()
			}

			layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(layout.Spacer{Height: 30}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Start}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Max.Y = gtx.Dp(40)
								i, _ := widget.NewIcon(icons.NavigationArrowBack)
								b := material.IconButton(th, prevButton, i, "Previous Date")
								b.Size = unit.Dp(20)
								b.Inset = layout.UniformInset(unit.Dp(4))
								return b.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: 10}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									t := material.Body2(th, date.Format("2006-01-02"))
									t.Alignment = text.Middle
									return t.Layout(gtx)
								})
							}),
							layout.Rigid(layout.Spacer{Width: 10}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Max.Y = gtx.Dp(40)
								i, _ := widget.NewIcon(icons.NavigationArrowForward)
								b := material.IconButton(th, nextButton, i, "Next Date")
								b.Size = unit.Dp(20)
								b.Inset = layout.UniformInset(unit.Dp(4))
								return b.Layout(gtx)
							}),
						)
					})
				}),
				layout.Rigid(layout.Spacer{Height: 50}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if sum == -1 {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max = image.Point{X: gtx.Dp(56), Y: gtx.Dp(56)}
							l := material.Loader(th)
							return l.Layout(gtx)
						})
					} else {
						t := material.Body1(th, fmt.Sprintf("Online: %.2f", sum))
						t.Alignment = text.Middle
						return t.Layout(gtx)
					}
				}),
			)

			e.Frame(gtx.Ops)
		}
	}

	return nil
}
