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
		SetTextAlign(tview.AlignCenter).
		SetText(" [yellow]F1:[white] Aide | [yellow]Ctrl+S:[white] Sauver | [yellow]Ctrl+K:[white] Couper bloc | [yellow]Ctrl+U:[white] Coller | [yellow]Ctrl+X:[white] Quitter")

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textArea, 0, 1, true).
		AddItem(footer, 1, 0, false)

	e.Pages.AddPage("edit_screen", layout, true, true)
	e.App.SetFocus(textArea)

	lastActionWasCut := false

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key != tcell.KeyCtrlK {
			lastActionWasCut = false
		}

		if key == tcell.KeyCtrlX || key == tcell.KeyEsc {
			if textArea.GetText() != initialContent {
				e.showSaveConfirmation(textArea.GetText())
			} else {
				e.Pages.RemovePage("edit_screen")
				e.App.SetFocus(e.FileList)
			}
			return nil
		}
		if key == tcell.KeyCtrlS {
			e.saveFromFullEditor(textArea.GetText())
			initialContent = textArea.GetText()
			updateTitle(false)
			return nil
		}
		if key == tcell.KeyCtrlK { // Couper la ligne (Nano style)
			text := textArea.GetText()
			lines := strings.Split(text, "\n")
			row, _, _, _ := textArea.GetCursor()
			if row < len(lines) {
				if lastActionWasCut {
					e.Clipboard += "\n" + lines[row]
				} else {
					e.Clipboard = lines[row]
				}
				lastActionWasCut = true

				// Calculate byte offsets for the current line
				start := 0
				for i := 0; i < row; i++ {
					start += len(lines[i]) + 1
				}
				
				end := start + len(lines[row])
				if row < len(lines)-1 {
					end += 1 // include the trailing newline
				} else if row > 0 {
					start -= 1 // include the preceding newline if last line
				}

				textArea.Replace(start, end, "")
			}
			return nil
		}
		if key == tcell.KeyCtrlU { // Coller la ligne
			if e.Clipboard != "" {
				text := textArea.GetText()
				lines := strings.Split(text, "\n")
				row, _, _, _ := textArea.GetCursor()
				
				start := 0
				for i := 0; i < row; i++ {
					start += len(lines[i]) + 1
				}
				
				textArea.Replace(start, start, e.Clipboard+"\n")
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
