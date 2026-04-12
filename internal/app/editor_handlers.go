package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// showFullEditor affiche l'interface d'édition plein écran
func (e *EditorApp) showFullEditor(content string) {
	initialContent := content
	textArea := tview.NewTextArea()
	textArea.SetText(content, false)
	// textArea.SetShowLineNumbers(true) // Nécessite une mise à jour : go get -u github.com/rivo/tview
	textArea.SetBorder(true)

	// Fonction pour mettre à jour le titre avec un indicateur de modification
	updateTitle := func(modified bool) {
		status := ""
		if modified {
			status = "[red]*[white] "
		}
		textArea.SetTitle(fmt.Sprintf(" %sÉdition: %s ", status, filepath.Base(e.FilePath)))
	}

	updateTitle(false)

	textArea.SetChangedFunc(func() {
		updateTitle(textArea.GetText() != initialContent)
	})

	// Instructions en bas de l'éditeur
	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetText(" [yellow]Ctrl+S:[white] Sauver | [yellow]Ctrl+X:[white] Quitter")

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textArea, 0, 1, true).
		AddItem(footer, 1, 0, false)

	e.Pages.AddPage("edit_screen", layout, true, true)
	e.App.SetFocus(textArea)

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlX || event.Key() == tcell.KeyEsc {
			if textArea.GetText() != initialContent {
				e.showSaveConfirmation(textArea.GetText())
			} else {
				e.Pages.RemovePage("edit_screen")
				e.App.SetFocus(e.FileList)
			}
			return nil
		}
		if event.Key() == tcell.KeyCtrlS {
			e.saveFromFullEditor(textArea.GetText())
			initialContent = textArea.GetText()
			updateTitle(false)
			return nil
		}
		if event.Key() == tcell.KeyCtrlK { // Couper la ligne (Nano style)
			text := textArea.GetText()
			lines := strings.Split(text, "\n")
			row, _, _, _ := textArea.GetCursor()
			if row < len(lines) {
				e.Clipboard = lines[row]
				lines = append(lines[:row], lines[row+1:]...)
				textArea.SetText(strings.Join(lines, "\n"), true)
			}
			return nil
		}
		if event.Key() == tcell.KeyCtrlU { // Coller la ligne
			if e.Clipboard != "" {
				text := textArea.GetText()
				lines := strings.Split(text, "\n")
				row, _, _, _ := textArea.GetCursor()
				newLines := append(lines[:row], append([]string{e.Clipboard}, lines[row:]...)...)
				textArea.SetText(strings.Join(newLines, "\n"), true)
			}
			return nil
		}
		return event
	})
}

// saveFromFullEditor gère la persistance des données depuis le mode édition
func (e *EditorApp) saveFromFullEditor(content string) {
	reader := strings.NewReader(content)
	err := e.FileSystem.Write(e.FilePath, reader)
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur de sauvegarde: %v", err))
	} else {
		e.refreshFileList()
		e.previewFile(e.FilePath)
		e.updateStatus(fmt.Sprintf("[green]Enregistré: %s", e.FilePath))
	}
}
