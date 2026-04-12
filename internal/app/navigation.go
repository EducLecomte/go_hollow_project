package app

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"context"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/EducLecomte/go_hollow_project/internal/vfs"
	"github.com/rivo/tview"
)

// refreshFileList met à jour la liste des fichiers du dossier actuel
func (e *EditorApp) refreshFileList() {
	e.FileList.Clear()
	e.FileList.AddItem("..", "Retour au parent", '.', nil)

	files, err := e.FileSystem.List(e.CurrentDir)
	if err != nil {
		e.CurrentFiles = nil
		e.updateStatus(fmt.Sprintf("[red]Erreur listage: %v", err))
		return
	}

	e.PathBar.SetText(fmt.Sprintf(" [yellow]Path: [white]%s", e.CurrentDir))
	e.CurrentFiles = files
	for _, f := range files {
		var displayName string
		if f.IsDir {
			// Dossier : Orange pour l'identification
			displayName = "[#ff8c00]" + f.Name + "/"
		} else {
			// Fichier : Nom simple
			displayName = f.Name
		}
		e.FileList.AddItem(displayName, "", 0, nil)
	}
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
				e.updateStatus(utils.HelpMsgFiles)
				return
			}
		}

		e.CurrentDir = filepath.Dir(e.CurrentDir)
		e.refreshFileList()
		e.updateStatus(utils.HelpMsgFiles)
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
						e.updateStatusTemp("[yellow]Ouverture annulée.")
					} else {
						e.updateStatusTemp(fmt.Sprintf("[red]Erreur d'ouverture d'archive: %v", err))
					}
					return
				}
				e.PreviousFS = e.FileSystem
				e.PreviousDir = e.CurrentDir
				e.FileSystem = archiveFS
				e.CurrentDir = "/"
				e.refreshFileList()
				e.updateStatusTemp(fmt.Sprintf("[green]Exploration de l'archive: %s", file.Name))
			})
		}()
	} else {
		e.openFile(targetPath)
	}
}

// openFile charge le contenu d'un fichier et ouvre l'éditeur
func (e *EditorApp) openFile(path string) {
	reader, err := e.FileSystem.Read(path)
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur lecture: %v", err))
		return
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur de buffer: %v", err))
		return
	}

	content := strings.ReplaceAll(buf.String(), "\r", "")
	e.FilePath = path

	// On délègue l'affichage à la nouvelle fenêtre d'édition
	e.showFullEditor(content)
}

// previewFile affiche un aperçu du contenu d'un fichier
func (e *EditorApp) previewFile(path string) {
	reader, err := e.FileSystem.Read(path)
	if err != nil {
		e.Viewer.SetText(fmt.Sprintf("[red]Erreur lecture: %v", err))
		return
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	// On limite la prélecture pour les gros fichiers par performance
	_, _ = io.CopyN(buf, reader, 10000)

	content := strings.ReplaceAll(buf.String(), "\r", "")

	// Application de la coloration syntaxique
	highlighted := utils.Highlight(content, path)
	// tview.TranslateANSI convertit les codes couleurs de chroma pour le TextView
	e.Viewer.SetText(tview.TranslateANSI(highlighted))
	e.Viewer.ScrollToBeginning()
	e.Viewer.SetTitle(fmt.Sprintf(" Visualiseur: %s ", filepath.Base(path)))
}

// previewDirectory affiche une arborescence simplifiée du contenu d'un dossier
func (e *EditorApp) previewDirectory(path string) {
	files, err := e.FileSystem.List(path)
	if err != nil {
		e.Viewer.SetText(fmt.Sprintf("[red]Erreur lecture dossier: %v", err))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Contenu de [yellow]%s[white] :\n\n", filepath.Base(path)))

	if len(files) == 0 {
		sb.WriteString("  [gray](Dossier vide)")
	} else {
		for i, f := range files {
			connector := "├── "
			if i == len(files)-1 {
				connector = "└── "
			}
			if f.IsDir {
				sb.WriteString(fmt.Sprintf("%s[darkorange]%s/[white]\n", connector, f.Name))
			} else {
				sb.WriteString(fmt.Sprintf("%s%s\n", connector, f.Name))
			}
		}
	}

	e.Viewer.SetText(sb.String())
	e.Viewer.ScrollToBeginning()
	e.Viewer.SetTitle(fmt.Sprintf(" Visualiseur: %s ", filepath.Base(path)))
}
