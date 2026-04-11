package main

import (
	"path/filepath"
	"strings"

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
			// Coupe la ligne actuelle (Style Nano)
			row, col := e.Editor.GetCursor()
			fullText := e.Editor.GetText()
			lines := strings.Split(fullText, "\n")

			if row < len(lines) {
				lineContent := lines[row]
				clipboard.WriteAll(lineContent)
				// On remplace la ligne (et le caractère newline) par rien
				e.Editor.Replace(row, 0, row+1, 0, "")
				e.Editor.SetCursor(row, col)
				e.updateStatus("Ligne coupée !")
			}
			return nil

		case tcell.KeyCtrlV:
			// Insère le texte à la position du curseur
			text, _ := clipboard.ReadAll()
			row, col := e.Editor.GetCursor()
			// Replace avec les mêmes coordonnées de début et de fin = Insertion
			e.Editor.Replace(row, col, row, col, text)
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
