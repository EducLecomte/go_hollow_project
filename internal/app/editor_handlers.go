package app

import (
	"github.com/gdamore/tcell/v2"
)

// setupEditorHandlers gère les entrées clavier pour la zone d'édition
func (e *EditorApp) setupEditorHandlers() {
	e.Viewer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			e.App.SetFocus(e.FileList)
			return nil
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight,
			tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
			// On laisse passer les touches de navigation pour permettre de parcourir le fichier
			return event
		}

		// On bloque tout le reste (caractères, Entrée, Backspace, etc.) pour simuler le ReadOnly
		return nil
	})
}
