package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(newFysionTheme())
	w := a.NewWindow("Fysion App")

	w.SetContent(makeGUI())

	w.Show()
	a.Run()
}
