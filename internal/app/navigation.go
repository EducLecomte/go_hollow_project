package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/EducLecomte/go_hollow_project/internal/vfs"
	"github.com/rivo/tview"
)

// refreshFileList recharge la liste des fichiers du répertoire courant et met à jour l'affichage de l'explorateur.
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

// handleFileSelection traite l'action de validation sur un élément de la liste (navigation, ouverture de fichier ou d'archive).
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

// openFile lit le contenu d'un fichier via le VFS et lance l'interface d'édition plein écran.
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

// previewFile lit les premiers octets d'un fichier pour en afficher un aperçu colorisé dans le visualiseur.
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

// previewDirectory génère une représentation textuelle arborescente du contenu d'un dossier pour le visualiseur.
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
