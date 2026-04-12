package app

import (
	"github.com/gdamore/tcell/v2"
)

// setupViewerHandlers gère les entrées clavier pour la zone de visualisation (lecture seule)
func (e *EditorApp) setupViewerHandlers() {
	e.Viewer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			e.App.SetFocus(e.FileList)
			return nil
		}

		// Le TextView (Viewer) gère nativement les flèches pour le défilement
		return event
	})
}
