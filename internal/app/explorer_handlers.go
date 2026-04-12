package app

import (
	"path/filepath"

	"github.com/gdamore/tcell/v2"
)

// setupExplorerHandlers gère les entrées clavier pour la liste de fichiers
func (e *EditorApp) setupExplorerHandlers() {
	e.FileList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			e.App.SetFocus(e.Viewer)
			return nil
		case tcell.KeyCtrlX:
			e.showQuitConfirmation()
			return nil
		case tcell.KeyCtrlF:
			e.showNewFileDialog()
			return nil
		case tcell.KeyCtrlD:
			e.showNewDirDialog()
			return nil
		case tcell.KeyCtrlK:
			index := e.FileList.GetCurrentItem()
			if index > 0 && index-1 < len(e.CurrentFiles) {
				file := e.CurrentFiles[index-1]
				e.prepareCopyFile(filepath.Join(e.CurrentDir, file.Name))
			}
			return nil
		case tcell.KeyCtrlU:
			e.pasteFile()
			return nil
		case tcell.KeyCtrlR:
			e.showDeleteConfirmation()
			return nil
		}
		return event
	})
}

// handleFileSelection gère l'ouverture des fichiers et la navigation dans les dossiers
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
