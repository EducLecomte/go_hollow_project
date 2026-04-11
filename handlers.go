package main

import (
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (e *EditorApp) setupHandlers() {
	e.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			e.showHelp()
			return nil
		case tcell.KeyF10:
			e.App.Stop()
			return nil
		case tcell.KeyTab:
			if e.App.GetFocus() == e.FileList {
				e.App.SetFocus(e.Editor)
				e.updateStatus(HelpMsgEdit)
				return nil
			}
			return event
		case tcell.KeyCtrlS:
			e.saveFile()
			return nil
		case tcell.KeyCtrlF:
			e.updateStatus("[blue]Recherche... (Fonctionnalité en cours de développement)")
			return nil
		case tcell.KeyCtrlX:
			if e.App.GetFocus() == e.Editor {
				e.App.SetFocus(e.FileList)
				e.updateStatus(HelpMsgFiles)
			} else {
				e.App.Stop()
			}
			return nil
		case tcell.KeyCtrlC:
			clipboard.WriteAll(e.Editor.GetText())
			e.updateStatus("Texte copié !")
			return nil
		case tcell.KeyCtrlK:
			// Logique de coupe (Kill line)
			clipboard.WriteAll(e.Editor.GetText())
			e.Editor.SetText("", true) // Simulation de coupe du contenu
			e.updateStatus("Texte coupé (Ctrl+K) !")
			return nil
		case tcell.KeyCtrlV:
			text, _ := clipboard.ReadAll()
			e.Editor.SetText(e.Editor.GetText()+text, true)
			e.updateStatus("Texte collé !")
			return nil
		}
		return event
	})
}

func (e *EditorApp) handleFileSelection(index int) {
	if index == 0 {
		e.CurrentDir = filepath.Dir(e.CurrentDir)
		e.refreshFileList()
		return
	}

	file := e.CurrentFiles[index-1]
	targetPath := filepath.Join(e.CurrentDir, file.Name)

	if file.IsDir {
		e.CurrentDir = targetPath
		e.refreshFileList()
	} else {
		e.openFile(targetPath)
	}
}

// showHelp affiche une fenêtre modale avec la liste des raccourcis
func (e *EditorApp) showHelp() {
	helpText := tview.NewTextView().
		SetText(HelpContent).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	helpText.SetBorder(true).
		SetTitle(" Aide - Raccourcis (Echap pour fermer) ").
		SetTitleAlign(tview.AlignCenter)

	// Centre la modale à l'écran
	helpModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(helpText, 18, 1, true).
			AddItem(nil, 0, 1, false), 60, 1, true).
		AddItem(nil, 0, 1, false)

	e.Pages.AddPage("help", helpModal, true, true)
	e.App.SetFocus(helpText)

	helpText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyF1 {
			e.Pages.RemovePage("help")
			return nil
		}
		return event
	})
}
