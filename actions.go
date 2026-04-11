package main

import (
	"bytes"
	"fmt"
	"io"
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
