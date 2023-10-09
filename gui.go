package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)
	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillContain
	return container.NewStack(toolbar, container.NewPadded(logo))
}

func makeGUI() fyne.CanvasObject {
	var (
		top     = makeBanner()
		left    = widget.NewLabel("Left")
		right   = widget.NewLabel("Right")
		content = canvas.NewRectangle(color.Gray{Y: 0xee})
		seps    = [3]fyne.CanvasObject{widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator()}
		objs    = []fyne.CanvasObject{top, left, right, content, seps[0], seps[1], seps[2]}
	)

	return container.New(newFysionLayout(top, left, right, content, seps), objs...)
}
