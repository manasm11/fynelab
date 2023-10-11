package main

import (
	"errors"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	dialogs "github.com/manasm11/fynelab/pkg/dialog"
)

const introText = `Here you can create new project!

Or open an existing one that you created earlier.`

type gui struct {
	win   fyne.Window
	title binding.String
}

func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)
	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillContain
	return container.NewStack(toolbar, container.NewPadded(logo))
}

func (g *gui) makeGUI() fyne.CanvasObject {
	var (
		top         = makeBanner()
		left        = widget.NewLabel("Left")
		right       = widget.NewLabel("Right")
		contentRect = canvas.NewRectangle(color.Gray{Y: 0xee})
		seps        = [3]fyne.CanvasObject{widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator()}
	)
	directory := widget.NewLabelWithData(g.title)
	content := container.NewStack(contentRect, directory)
	objs := []fyne.CanvasObject{top, left, right, content, seps[0], seps[1], seps[2]}

	return container.New(newFysionLayout(top, left, right, content, seps), objs...)
}

func (g *gui) openProjectDialog() {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, g.win)
			return
		}
		if dir == nil {
			return
		}

		g.openProject(dir)
	}, g.win)
}

func (g *gui) openProject(dir fyne.ListableURI) {
	name := dir.Name()
	g.title.Set(name)
}

func (g *gui) showCreate(w fyne.Window) {
	var (
		wizard *dialogs.Wizard
		intro  = widget.NewLabel(introText)

		open   = widget.NewButton("Open Project", func() { wizard.Hide(); g.openProjectDialog() })
		create = widget.NewButton("Create Project", func() { wizard.Push("Project Details", g.makeCreateDetail(wizard)) })

		buttons = container.NewGridWithColumns(2, open, create)
		home    = container.NewVBox(intro, buttons)
	)
	create.Importance = widget.HighImportance
	wizard = dialogs.NewWizard("Create Project", home)
	wizard.Show(w)
	wizard.Resize(home.MinSize().AddWidthHeight(40, 80))
}

func (g *gui) makeCreateDetail(wizard *dialogs.Wizard) fyne.CanvasObject {
	var (
		homeDir, _ = os.UserHomeDir()
		parent     = storage.NewFileURI(homeDir)
		chosen, _  = storage.ListerForURI(parent)
		dir        *widget.Button
		name       = widget.NewEntry()
		form       *widget.Form
	)

	name.Validator = func(s string) error {
		if s == "" {
			return errors.New("project name is required")
		}
		return nil
	}

	dir = widget.NewButton(chosen.Name(), func() {
		d := dialog.NewFolderOpen(func(l fyne.ListableURI, err error) {
			if err != nil || l == nil {
				return
			}
			chosen = l
			dir.SetText(l.Name())
		}, g.win)
		d.SetLocation(chosen)
		d.Show()
	})
	form = widget.NewForm(
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Parent Directory", dir),
	)
	form.OnSubmit = func() {
		if name.Text == "" {
			return
		}

		project, err := createProject(name.Text, chosen)
		if err != nil {
			dialog.ShowError(err, g.win)
			return
		}

		wizard.Hide()
		g.openProject(project)
	}
	return form
}
