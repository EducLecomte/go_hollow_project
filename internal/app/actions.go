package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

	// Normalisation des fins de ligne (CRLF -> LF) pour éviter les bugs de curseur
	content := strings.ReplaceAll(buf.String(), "\r", "")

	e.FilePath = path
	e.Editor.SetText(content, true)
	e.Editor.SetTitle(fmt.Sprintf(" Éditeur: %s ", filepath.Base(path)))
}

func (e *EditorApp) saveFile() {
	if e.FilePath == "" {
		e.updateStatus("[red]Erreur: Aucun fichier ouvert")
		return
	}

	content := e.Editor.GetText()
	reader := strings.NewReader(content)

	err := e.FileSystem.Write(e.FilePath, reader)
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur de sauvegarde: %v", err))
	} else {
		// Rafraîchir la liste pour mettre à jour la taille affichée
		e.refreshFileList()
		e.updateStatus(fmt.Sprintf("[green]Enregistré: %s", e.FilePath))
	}
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
	dst := filepath.Join(e.CurrentDir, baseName)

	// Si on colle dans le même dossier, on ajoute un suffixe pour éviter l'écrasement
	if e.CopiedPath == dst {
		ext := filepath.Ext(baseName)
		name := strings.TrimSuffix(baseName, ext)
		dst = filepath.Join(e.CurrentDir, name+"_copy"+ext)
	}

	err := e.FileSystem.Copy(e.CopiedPath, dst)
	if err != nil {
		e.updateStatusTemp(fmt.Sprintf("[red]Erreur collage: %v", err))
		return
	}

	e.refreshFileList()
	e.updateStatusTemp(fmt.Sprintf("[green]Élément collé: %s", filepath.Base(dst)))
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
