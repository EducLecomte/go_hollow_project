package main

import (
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

func (e *EditorApp) setupHandlers() {
	e.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if e.App.GetFocus() == e.FileList {
				e.App.SetFocus(e.Editor)
				e.updateStatus(HelpMsgEdit)
				return nil
			}
			return event
		case tcell.KeyCtrlS:
			e.saveFile()
			return nil
		case tcell.KeyCtrlX:
			if e.App.GetFocus() == e.Editor {
				e.App.SetFocus(e.FileList)
				e.updateStatus(HelpMsgFiles)
			} else {
				e.App.Stop()
			}
			return nil
		case tcell.KeyCtrlC:
			clipboard.WriteAll(e.Editor.GetText())
			e.updateStatus("Texte copié !")
			return nil
		case tcell.KeyCtrlV:
			text, _ := clipboard.ReadAll()
			e.Editor.SetText(e.Editor.GetText()+text, true)
			e.updateStatus("Texte collé !")
			return nil
		}
		return event
	})
}

func (e *EditorApp) handleFileSelection(index int) {
	if index == 0 {
		e.CurrentDir = filepath.Dir(e.CurrentDir)
		e.refreshFileList()
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
