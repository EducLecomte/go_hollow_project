package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/EducLecomte/go_hollow_project/internal/vfs"
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
		if e.CurrentDir == "/" || e.CurrentDir == "." || e.CurrentDir == "" {
			if e.PreviousFS != nil {
				// Sortie du système de fichiers virtuel
				e.FileSystem = e.PreviousFS
				e.CurrentDir = e.PreviousDir
				e.PreviousFS = nil
				e.refreshFileList()
				return
			}
		}

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
	} else if utils.IsArchive(file.Name) {
		ctx, cancel := context.WithCancel(context.Background())
		e.showLoadingDialog("Chargement", fmt.Sprintf("Ouverture de %s en cours...", file.Name), cancel)

		go func() {
			archiveFS, err := vfs.NewArchiveFS(ctx, targetPath)
			
			e.App.QueueUpdateDraw(func() {
				e.Pages.RemovePage("loading")
				
				if err != nil {
					if err == context.Canceled {
						e.updateStatus("[yellow]Ouverture annulée.")
					} else {
						e.updateStatus(fmt.Sprintf("[red]Erreur d'ouverture d'archive: %v", err))
					}
					return
				}
				e.PreviousFS = e.FileSystem
				e.PreviousDir = e.CurrentDir
				e.FileSystem = archiveFS
				e.CurrentDir = "/"
				e.refreshFileList()
				e.updateStatus(fmt.Sprintf("[green]Exploration de l'archive: %s", file.Name))
			})
		}()
	} else {
		e.openFile(targetPath)
	}
}
