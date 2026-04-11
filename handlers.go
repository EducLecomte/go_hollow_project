package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (e *EditorApp) setupHandlers() {
	// 1. Raccourcis Globaux (Application)
	e.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			e.showHelp()
			return nil
		case tcell.KeyF10:
			e.App.Stop()
			return nil
		case tcell.KeyCtrlC:
			// Si l'éditeur n'a pas le focus, on consomme l'événement (nil) pour éviter
			// que le terminal n'intercepte le signal SIGINT et ne ferme l'application.
			// Si l'éditeur a le focus, on laisse l'événement passer pour qu'il soit géré par l'éditeur.
			if e.App.GetFocus() != e.Editor {
				return nil
			}
		}
		return event
	})

	// 2. Raccourcis spécifiques à l'Exploreur (FileList)
	e.FileList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			e.App.SetFocus(e.Editor)
			return nil
		case tcell.KeyCtrlX:
			e.App.Stop()
			return nil
		case tcell.KeyCtrlS:
			e.saveFile()
			return nil
		case tcell.KeyCtrlF:
			e.updateStatus("[blue]Recherche dans l'exploreur... (Bientôt disponible)")
			return nil
		}
		return event
	})

	// 3. Raccourcis spécifiques à l'Éditeur (TextArea)
	e.Editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlX:
			e.App.SetFocus(e.FileList)
			return nil
		case tcell.KeyCtrlS:
			e.saveFile()
			return nil
		case tcell.KeyCtrlF:
			e.updateStatus("[blue]Recherche dans le texte... (Bientôt disponible)")
			return nil
		case tcell.KeyCtrlC:
			row, _, _, _ := e.Editor.GetCursor()
			lines := strings.Split(e.Editor.GetText(), "\n")
			if row >= 0 && row < len(lines) {
				if err := clipboard.WriteAll(lines[row]); err != nil {
					e.updateStatus(fmt.Sprintf("[red]Erreur Clipboard: %v", err))
				} else {
					e.updateStatusTemp("Ligne copiée !")
				}
			}
			return nil
		case tcell.KeyCtrlK:
			row, _, _, _ := e.Editor.GetCursor()
			fullText := e.Editor.GetText()
			lines := strings.Split(fullText, "\n")
			if row < 0 || row >= len(lines) {
				return nil
			}
			clipboard.WriteAll(lines[row])
			newLines := make([]string, 0, len(lines)-1)
			newLines = append(newLines, lines[:row]...)
			newLines = append(newLines, lines[row+1:]...)
			e.Editor.SetText(strings.Join(newLines, "\n"), true)
			e.updateStatusTemp("Ligne coupée !")
			return nil
		case tcell.KeyCtrlU, tcell.KeyCtrlV:
			text, err := clipboard.ReadAll()
			if err != nil {
				e.updateStatus(fmt.Sprintf("[red]Erreur Presse-papier: %v", err))
				return nil
			}
			if text != "" {
				text = strings.ReplaceAll(text, "\r", "")
				row, col, _, _ := e.Editor.GetCursor()
				e.Editor.Replace(row, col, text)
				msg := "Texte collé !"
				if event.Key() == tcell.KeyCtrlU {
					msg = "Ligne collée (Ctrl+U) !"
				}
				e.updateStatusTemp(msg)
			}
			return nil
		}
		return event
	})
}

// copyLineFromEditor extrait la ligne actuelle et l'envoie au presse-papier
func (e *EditorApp) copyLineFromEditor() {
	row, _, _, _ := e.Editor.GetCursor()
	lines := strings.Split(e.Editor.GetText(), "\n")
	if row >= 0 && row < len(lines) {
		if err := clipboard.WriteAll(lines[row]); err != nil {
			e.updateStatus(fmt.Sprintf("[red]Erreur Clipboard: %v", err))
		} else {
			e.updateStatusTemp("Ligne copiée !")
		}
	}
}

func (e *EditorApp) handleFileSelection(index int) {
	if index == 0 {
		e.CurrentDir = filepath.Dir(e.CurrentDir)
		e.refreshFileList()
		return
	}

	if e.CurrentFiles == nil || index-1 >= len(e.CurrentFiles) {
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
