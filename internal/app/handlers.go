package app

import (
	"fmt"
	"path/filepath"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (e *EditorApp) setupHandlers() {
	// 1. Raccourcis Globaux (Application)
	e.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		allowedCtrlKeys := map[tcell.Key]bool{
			tcell.KeyCtrlS: true, tcell.KeyCtrlF: true, tcell.KeyCtrlD: true,
			tcell.KeyCtrlR: true, tcell.KeyCtrlK: true, tcell.KeyCtrlU: true,
			tcell.KeyCtrlV: true, tcell.KeyCtrlX: true,
			tcell.KeyTab: true, tcell.KeyEnter: true,
			tcell.KeyBackspace: true, tcell.KeyBackspace2: true,
		}

		switch event.Key() {
		case tcell.KeyF1:
			if e.Pages.HasPage("help") || e.Pages.HasPage("quit") ||
				e.Pages.HasPage("newfile") || e.Pages.HasPage("newdir") ||
				e.Pages.HasPage("delete") {
				return event
			}
			e.showHelp()
			return nil
		case tcell.KeyF10:
			if e.Pages.HasPage("help") || e.Pages.HasPage("quit") ||
				e.Pages.HasPage("newfile") || e.Pages.HasPage("newdir") ||
				e.Pages.HasPage("delete") {
				return event
			}
			e.showQuitConfirmation()
			return nil
		case tcell.KeyCtrlC:
			return nil
		}

		if event.Modifiers()&tcell.ModCtrl != 0 {
			if !allowedCtrlKeys[event.Key()] {
				return nil
			}
		}
		if event.Modifiers()&tcell.ModAlt != 0 {
			return nil
		}
		return event
	})

	e.setupExplorerHandlers()
	e.setupEditorHandlers()
}

func (e *EditorApp) showHelp() {
	previousFocus := e.App.GetFocus()
	helpText := tview.NewTextView().
		SetText(utils.HelpContent).
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextAlign(tview.AlignLeft)

	helpText.SetBorder(true).
		SetTitle(" Aide Hollow - Raccourcis ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderPadding(1, 1, 2, 2)

	helpModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(helpText, 20, 1, true).
			AddItem(nil, 0, 1, false), 65, 1, true).
		AddItem(nil, 0, 1, false)

	e.Pages.AddPage("help", helpModal, true, true)
	e.App.SetFocus(helpText)

	helpText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyF1 || event.Rune() == 'q' {
			e.Pages.RemovePage("help")
			if previousFocus != nil {
				e.App.SetFocus(previousFocus)
			}
			return nil
		}
		return event
	})
}

func (e *EditorApp) showQuitConfirmation() {
	modal := tview.NewModal().
		SetText("Voulez-vous vraiment quitter Hollow ?").
		AddButtons([]string{"Oui", "Non"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Oui" {
				e.saveLastDir()
				e.App.Stop()
			}
			e.Pages.RemovePage("quit")
		})
	e.Pages.AddPage("quit", modal, true, true)
}

func (e *EditorApp) showDeleteConfirmation() {
	index := e.FileList.GetCurrentItem()
	if index <= 0 || index-1 >= len(e.CurrentFiles) {
		return
	}
	file := e.CurrentFiles[index-1]
	path := filepath.Join(e.CurrentDir, file.Name)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Voulez-vous vraiment supprimer %s ?", file.Name)).
		AddButtons([]string{"Supprimer", "Annuler"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Supprimer" {
				e.deleteElement(path)
			}
			e.Pages.RemovePage("delete")
			e.App.SetFocus(e.FileList)
		})
	e.Pages.AddPage("delete", modal, true, true)
}

func (e *EditorApp) showNewFileDialog() {
	inputField := tview.NewInputField().SetLabel(" Nom du nouveau fichier: ")
	inputField.SetBorder(true).SetTitle(" Nouveau Fichier ").SetTitleAlign(tview.AlignCenter)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(inputField, 3, 1, true).
			AddItem(nil, 0, 1, false), 60, 1, true).
		AddItem(nil, 0, 1, false)

	e.Pages.AddPage("newfile", modal, true, true)
	e.App.SetFocus(inputField)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			name := inputField.GetText()
			if name != "" {
				e.createFile(name)
				e.App.SetFocus(e.Editor)
			}
			e.Pages.RemovePage("newfile")
		} else if key == tcell.KeyEsc {
			e.Pages.RemovePage("newfile")
			e.App.SetFocus(e.FileList)
		}
	})
}

// showNewDirDialog affiche un champ de saisie pour créer un nouveau dossier
func (e *EditorApp) showNewDirDialog() {
	inputField := tview.NewInputField().SetLabel(" Nom du nouveau dossier: ")
	inputField.SetBorder(true).SetTitle(" Nouveau Dossier ").SetTitleAlign(tview.AlignCenter)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(inputField, 3, 1, true).
			AddItem(nil, 0, 1, false), 60, 1, true).
		AddItem(nil, 0, 1, false)

	e.Pages.AddPage("newdir", modal, true, true)
	e.App.SetFocus(inputField)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			name := inputField.GetText()
			if name != "" {
				e.createDir(name)
			}
			e.Pages.RemovePage("newdir")
			e.App.SetFocus(e.FileList)
		} else if key == tcell.KeyEsc {
			e.Pages.RemovePage("newdir")
			e.App.SetFocus(e.FileList)
		}
	})
}
