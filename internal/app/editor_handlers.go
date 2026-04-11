package app

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

// setupEditorHandlers gère les entrées clavier pour la zone d'édition
func (e *EditorApp) setupEditorHandlers() {
	// Réinitialise le curseur et l'affichage en haut de page lors de la prise de focus
	// Cela couvre le passage par "Tab" et la sélection depuis l'explorateur.
	e.Editor.SetFocusFunc(func() {
		// SetText avec keepCursor=false force le curseur à (0,0)
		// et assure que le haut du fichier est visible.
		e.Editor.SetText(e.Editor.GetText(), false)
	})

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
	// GetCursor renvoie 4 valeurs : row, col, width, height
	row, _, _, _ := e.Editor.GetCursor()
	lines := strings.Split(e.Editor.GetText(), "\n")

	if row < 0 || row >= len(lines) {
		return
	}

	lineText := lines[row]
	// On ajoute un retour à la ligne pour forcer l'insertion de ligne lors du collage
	clipboard.WriteAll(lineText + "\n")

	// Calcul de la position 1D du début de la ligne (colonne 0)
	pos := 0
	for i := 0; i < row; i++ {
		pos += len(lines[i]) + 1
	}

	// Longueur à supprimer : le texte de la ligne plus le séparateur
	fullText := e.Editor.GetText()
	count := len(lineText) + 1
	if pos+count > len(fullText) {
		count = len(fullText) - pos
	}

	e.Editor.Replace(pos, count, "")
	e.updateStatusTemp("Ligne coupée")
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

	// Pour que la ligne actuelle "descende", on colle systématiquement au début de la ligne
	// On s'assure d'ignorer les 4 valeurs de retour pour éviter l'erreur de compilation
	row, _, _, _ := e.Editor.GetCursor()
	lines := strings.Split(e.Editor.GetText(), "\n")

	pos := 0
	for i := 0; i < row; i++ {
		if i < len(lines) {
			pos += len(lines[i]) + 1
		}
	}

	// On s'assure que le texte se termine par un saut de ligne pour l'effet d'insertion
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	e.Editor.Replace(pos, 0, text)

	msg := "Texte collé !"
	if isUncut {
		msg = "Ligne rétablie (Ctrl+U) !"
	}
	e.updateStatusTemp(msg)
}
