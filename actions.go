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

	e.FilePath = path
	e.Editor.SetText(buf.String(), true)
	e.updateStatus(fmt.Sprintf("Ouvert: %s", filepath.Base(path)))
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
		e.updateStatus(fmt.Sprintf("[green]Enregistré: %s", e.FilePath))
	}
}
