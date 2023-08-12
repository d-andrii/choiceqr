package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func IconButton(gtx layout.Context, btn *widget.Clickable, th *material.Theme, text string) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Dp(40)
	b := material.Button(th, btn, text)
	return b.Layout(gtx)
}
