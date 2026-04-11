package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

// setupEditorHandlers gère les entrées clavier pour la zone d'édition
func (e *EditorApp) setupEditorHandlers() {
	e.Editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlX:
			e.App.SetFocus(e.FileList)
			return nil
		case tcell.KeyCtrlS:
			e.saveFile()
			return nil
		case tcell.KeyCtrlK:
			e.cutLine()
			return nil
		case tcell.KeyCtrlU, tcell.KeyCtrlV:
			e.pasteText(event.Key() == tcell.KeyCtrlU)
			return nil
		}
		return event
	})
}

// cutLine : Supprime la ligne et la met dans le presse-papier
func (e *EditorApp) cutLine() {
	row, _, _, _ := e.Editor.GetCursor()
	text := e.Editor.GetText()
	lines := strings.Split(text, "\n")

	if row >= 0 && row < len(lines) {
		content := lines[row]
		clipboard.WriteAll(content)

		// Suppression de la ligne
		var newLines []string
		for i, l := range lines {
			if i != row {
				newLines = append(newLines, l)
			}
		}
		if len(newLines) == 0 {
			newLines = []string{""}
		}
		e.Editor.SetText(strings.Join(newLines, "\n"), true)
		e.updateStatusTemp("Ligne coupée !")
	}
}

// pasteText : Colle le texte à la position exacte du curseur
func (e *EditorApp) pasteText(isUncut bool) {
	text, err := clipboard.ReadAll()
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur Presse-papier: %v", err))
		return
	}
	if text == "" {
		return
	}

	text = strings.ReplaceAll(text, "\r", "")

	// Conversion des coordonnées 2D en index 1D pour coller au bon endroit
	row, col, _, _ := e.Editor.GetCursor()
	content := e.Editor.GetText()
	lines := strings.Split(content, "\n")

	pos := 0
	for i := 0; i < row && i < len(lines); i++ {
		pos += len(lines[i]) + 1
	}

	if row < len(lines) {
		if col > len(lines[row]) {
			col = len(lines[row])
		}
		pos += col
	} else {
		pos = len(content)
	}

	e.Editor.Replace(pos, 0, text)

	msg := "Texte collé !"
	if isUncut {
		msg = "Ligne collée (Ctrl+U) !"
	}
	e.updateStatusTemp(msg)
}
