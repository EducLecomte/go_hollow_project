package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/rivo/tview"
)

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

func (e *EditorApp) createFile(name string) {
	path := filepath.Join(e.CurrentDir, name)
	// On écrit un fichier vide
	err := e.FileSystem.Write(path, strings.NewReader(""))
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur création: %v", err))
		return
	}
	e.refreshFileList()
	e.openFile(path)
	e.updateStatus(fmt.Sprintf("[green]Fichier créé: %s", name))
}

func (e *EditorApp) createDir(name string) {
	path := filepath.Join(e.CurrentDir, name)
	err := e.FileSystem.Mkdir(path)
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur dossier: %v", err))
		return
	}
	e.refreshFileList()
	e.updateStatus(fmt.Sprintf("[green]Dossier créé: %s", name))
}

func (e *EditorApp) prepareCopyFile(path string) {
	e.CopiedPath = path
	e.updateStatusTemp(fmt.Sprintf("Élément prêt à copier: %s", filepath.Base(path)))
}

func (e *EditorApp) pasteFile() {
	if e.CopiedPath == "" {
		e.updateStatusTemp("[red]Rien à coller")
		return
	}

	baseName := filepath.Base(e.CopiedPath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)

	finalName := baseName
	counter := 1

	for {
		exists := false
		for _, f := range e.CurrentFiles {
			if f.Name == finalName {
				exists = true
				break
			}
		}
		if !exists {
			break
		}
		if counter == 1 {
			finalName = fmt.Sprintf("%s_copy%s", nameWithoutExt, ext)
		} else {
			finalName = fmt.Sprintf("%s_copy%d%s", nameWithoutExt, counter, ext)
		}
		counter++
	}

	dst := filepath.Join(e.CurrentDir, finalName)

	err := e.FileSystem.Copy(e.CopiedPath, dst)
	if err != nil {
		e.updateStatusTemp(fmt.Sprintf("[red]Erreur collage: %v", err))
		return
	}

	e.refreshFileList()
	e.updateStatusTemp(fmt.Sprintf("[green]Élément collé: %s", finalName))
}

func (e *EditorApp) deleteElement(path string) {
	err := e.FileSystem.Remove(path)
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur suppression: %v", err))
		return
	}
	e.refreshFileList()
	e.updateStatus(fmt.Sprintf("[green]Supprimé: %s", filepath.Base(path)))
}

// saveLastDir enregistre le répertoire actuel dans un fichier temporaire pour le shell
func (e *EditorApp) saveLastDir() {
	path := fmt.Sprintf("/tmp/hollow_cwd_%s", os.Getenv("USER"))
	_ = os.WriteFile(path, []byte(e.CurrentDir), 0644)
}
