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

	fileTree binding.URITree
}

func (g *gui) makeBanner() fyne.CanvasObject {
	title := canvas.NewText("App Creator", theme.ForegroundColor())
	title.TextSize = 14
	title.TextStyle = fyne.TextStyle{Bold: true}

	g.title.AddListener(binding.NewDataListener(func() {
		name, _ := g.title.Get()
		if name == "" {
			name = "App Creator"
		}
		title.Text = name
		title.Refresh()
	}))

	home := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {})
	left := container.NewHBox(home, title)

	logo := canvas.NewImageFromResource(resourceLogoPng)
	logo.FillMode = canvas.ImageFillContain

	return container.NewStack(container.NewPadded(left), container.NewPadded(logo))
}

func (g *gui) makeGUI() fyne.CanvasObject {
	g.fileTree = binding.NewURITree()
	var (
		top   = g.makeBanner()
		files = widget.NewTreeWithData(g.fileTree, func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("filename.jpg")
		}, func(data binding.DataItem, branch bool, obj fyne.CanvasObject) {
			l := obj.(*widget.Label)
			data.(binding.URI).Get()
			u, _ := data.(binding.URI).Get()

			name := u.Name()
			l.SetText(name)
		})
		left = widget.NewAccordion(
			widget.NewAccordionItem("Files", files),
			widget.NewAccordionItem("Screens", widget.NewLabel("TODO Screens")),
		)
		right       = widget.NewRichTextFromMarkdown("## Settings")
		contentRect = canvas.NewRectangle(color.Gray{Y: 0xee})
		seps        = [3]fyne.CanvasObject{widget.NewSeparator(), widget.NewSeparator(), widget.NewSeparator()}
		name, _     = g.title.Get()
		window      = container.NewInnerWindow(name, widget.NewLabel("App Preview Here"))
		picker      = widget.NewSelect([]string{"Desktop", "iPhone 15 Max"}, func(s string) {})
		preview     = container.NewBorder((container.NewHBox(picker)), nil, nil, nil, container.NewCenter(window))
	)

	window.CloseIntercept = func() {}

	picker.Selected = "Desktop"

	left.Open(0)
	left.MultiOpen = true

	content := container.NewStack(contentRect, container.NewPadded(preview))
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
