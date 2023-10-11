package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const sideWidth = 220

type fysionLayout struct {
	top, left, right, content fyne.CanvasObject
	seps                      [3]fyne.CanvasObject
}

func newFysionLayout(top, left, right, content fyne.CanvasObject, seps [3]fyne.CanvasObject) fyne.Layout {
	return &fysionLayout{top: top, left: left, right: right, content: content, seps: seps}
}

func (l *fysionLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	topHeight := l.top.MinSize().Height
	l.top.Resize(fyne.NewSize(size.Width, topHeight))

	l.left.Move(fyne.NewPos(0, topHeight))
	l.left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.right.Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.right.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.content.Move(fyne.NewPos(sideWidth, topHeight))
	l.content.Resize(fyne.NewSize(size.Width-2*sideWidth, size.Height-topHeight))

	sepThickness := theme.SeparatorThicknessSize()

	l.seps[0].Move(fyne.NewPos(0, topHeight))
	l.seps[0].Resize(fyne.NewSize(size.Width, sepThickness))

	l.seps[1].Move(fyne.NewPos(sideWidth, topHeight))
	l.seps[1].Resize(fyne.NewSize(sepThickness, size.Height-topHeight))

	l.seps[1].Move(fyne.NewPos(sideWidth, topHeight))
	l.seps[1].Resize(fyne.NewSize(sepThickness, size.Height-topHeight))
}

func (l *fysionLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	borders := fyne.NewSize(2*sideWidth, l.top.MinSize().Height)
	return borders.AddWidthHeight(100, 100)
}
