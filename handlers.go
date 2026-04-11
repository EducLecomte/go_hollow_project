package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (e *EditorApp) setupHandlers() {
	// 1. Raccourcis Globaux (Application)
	e.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			e.showHelp()
			return nil
		case tcell.KeyF10:
			e.App.Stop()
			return nil
		case tcell.KeyCtrlC:
			// Intercepte Ctrl+C pour éviter de quitter l'application
			// (Comportement de sécurité pour éviter le SIGINT du terminal)
			return nil
		}
		return event
	})

	// 2. Raccourcis spécifiques à l'Exploreur (FileList)
	e.FileList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			e.App.SetFocus(e.Editor)
			return nil
		case tcell.KeyCtrlX:
			e.App.Stop()
			return nil
		case tcell.KeyCtrlS:
			e.saveFile()
			return nil
		}
		return event
	})

	// 3. Raccourcis spécifiques à l'Éditeur (TextArea)
	e.Editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlX:
			e.App.SetFocus(e.FileList)
			return nil
		case tcell.KeyCtrlS:
			e.saveFile()
			return nil
		case tcell.KeyCtrlK:
			e.cutLine()
			return nil
		case tcell.KeyCtrlU, tcell.KeyCtrlV:
			e.pasteText()
			return nil
		}
		return event
	})
}

func (e *EditorApp) cutLine() {
	row, _, _, _ := e.Editor.GetCursor()
	text := e.Editor.GetText()
	lines := strings.Split(text, "\n")

	if row >= 0 && row < len(lines) {
		// Nano cut : on prend la ligne et on ajoute un retour à la ligne
		lineContent := lines[row]
		clipboard.WriteAll(lineContent + "\n")

		// Supprimer la ligne du tableau
		lines = append(lines[:row], lines[row+1:]...)

		// Si le fichier devient vide, on laisse une ligne vide
		if len(lines) == 0 {
			lines = []string{""}
		}

		e.Editor.SetText(strings.Join(lines, "\n"), true)
		e.updateStatusTemp("Ligne coupée (Ctrl+K)")
	}
}

func (e *EditorApp) pasteText() {
	text, err := clipboard.ReadAll()
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur Presse-papier: %v", err))
		return
	}
	if text == "" {
		return
	}

	text = strings.ReplaceAll(text, "\r", "")

	// Conversion 2D (row, col) -> 1D (index global) pour Replace
	row, col, _, _ := e.Editor.GetCursor()
	textInEditor := e.Editor.GetText()
	lines := strings.Split(textInEditor, "\n")

	pos := 0
	for i := 0; i < row && i < len(lines); i++ {
		pos += len(lines[i]) + 1 // +1 pour le \n
	}

	if row < len(lines) {
		if col > len(lines[row]) {
			col = len(lines[row])
		}
		pos += col
	} else {
		pos = len(textInEditor)
	}

	e.Editor.Replace(pos, 0, text)
	e.updateStatusTemp("Collé (Ctrl+U/V)")
}

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

// showHelp affiche une fenêtre modale avec la liste des raccourcis
func (e *EditorApp) showHelp() {
	helpText := tview.NewTextView().
		SetText(HelpContent).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	helpText.SetBorder(true).
		SetTitle(" Aide - Raccourcis (Echap pour fermer) ").
		SetTitleAlign(tview.AlignCenter)

	// Centre la modale à l'écran
	helpModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(helpText, 18, 1, true).
			AddItem(nil, 0, 1, false), 60, 1, true).
		AddItem(nil, 0, 1, false)

	e.Pages.AddPage("help", helpModal, true, true)
	e.App.SetFocus(helpText)

	helpText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyF1 {
			e.Pages.RemovePage("help")
			return nil
		}
		return event
	})
}
