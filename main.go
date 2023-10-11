package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(newFysionTheme())
	w := a.NewWindow("Fysion App")
	w.Resize(fyne.NewSize(1000, 700))

	ui := &gui{win: w, title: binding.NewString()}
	ui.title.AddListener(binding.NewDataListener(func() {
		title, _ := ui.title.Get()
		w.SetTitle("Fysion App: " + title)
	}))

	w.SetMainMenu(ui.makeMenu())
	w.SetContent(ui.makeGUI())

	flag.Usage = func() {
		fmt.Println("Usage: fysion [project directory]")
	}
	flag.Parse()
	if len(flag.Args()) > 0 {
		dirPath := flag.Args()[0]
		dirPath, err := filepath.Abs(dirPath)
		if err != nil {
			fmt.Println("Error opening project", err)
			return
		}
		uri := storage.NewFileURI(dirPath)
		dir, err := storage.ListerForURI(uri)
		if err != nil {
			fmt.Println("Error opening project", err)
			return
		}
		ui.openProject(dir)
	} else {
		ui.showCreate(w)
	}
	w.Show()
	a.Run()
}

func (g *gui) makeMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("Open Project", g.openProjectDialog)),
	)
}
