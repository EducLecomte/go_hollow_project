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
		}

		// Le TextView gère nativement les flèches pour le scroll
		// On laisse passer les touches par défaut
		return event
	})
}
