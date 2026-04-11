package main

import (
	"fmt"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type EditorApp struct {
	App          *tview.Application
	PathBar      *tview.TextView
	FileList     *tview.List
	FileSizeBox  *tview.TextView
	Editor       *tview.TextArea
	Status       *tview.TextView
	Pages        *tview.Pages
	FilePath     string
	CurrentDir   string
	FileSystem   VFS
	CurrentFiles []FileInfo
}

func NewEditorApp() *EditorApp {
	// Initialisation par défaut en local
	localFS := &LocalFS{}
	wd, err := filepath.Abs(".")
	if err != nil {
		wd = "/"
	}

	e := &EditorApp{
		App:         tview.NewApplication(),
		PathBar:     tview.NewTextView(),
		FileList:    tview.NewList(),
		FileSizeBox: tview.NewTextView(),
		Editor:      tview.NewTextArea(),
		Status:      tview.NewTextView(),
		Pages:       tview.NewPages(),
		CurrentDir:  wd,
		FileSystem:  localFS,
	}

	e.setupUI()
	e.setupHandlers()
	e.refreshFileList()
	return e
}

func (e *EditorApp) setupUI() {
	e.PathBar.SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorGreen)

	e.FileList.SetBorder(true).SetTitle(" Explorer ")
	e.FileList.ShowSecondaryText(false) // Rend la liste compacte (une seule ligne)
	e.FileList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		e.handleFileSelection(index)
	})

	// Encart pour le poids du fichier
	e.FileSizeBox.SetBorder(true).SetTitle(" Info ")
	e.FileSizeBox.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)

	// Mise à jour de l'encart quand on navigue dans la liste
	e.FileList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index == 0 {
			e.FileSizeBox.SetText("[gray]Parent Directory")
			return
		}
		file := e.CurrentFiles[index-1]
		if file.IsDir {
			e.FileSizeBox.SetText("[yellow]Directory")
		} else {
			e.FileSizeBox.SetText(fmt.Sprintf("[blue]Size: [white]%s", formatSize(file.Size)))
		}
	})

	e.Editor.SetBorder(true).SetTitle(" Éditeur ")
	e.Editor.SetPlaceholder("Entrez votre texte ici... (Ctrl+S pour sauver, Ctrl+C/V pour copier/coller)")

	e.Status.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	e.updateStatus(HelpMsgDefault)

	// Layout de la colonne de gauche (Explorateur + Info taille)
	leftColumn := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(e.FileList, 0, 1, true).
		AddItem(e.FileSizeBox, 3, 0, false)

	// Layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(e.PathBar, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(leftColumn, 0, 1, true).
			AddItem(e.Editor, 0, 2, false), 0, 1, true).
		AddItem(e.Status, 1, 0, false)

	e.Pages.AddPage("main", mainFlex, true, true)
}

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
		suffix := ""
		if f.IsDir {
			suffix = "/"
		}
		// On n'affiche plus que le nom, le poids est dans l'encart du bas
		e.FileList.AddItem(f.Name+suffix, "", 0, nil)
	}
}
func (e *EditorApp) updateStatus(msg string) {
	// On s'assure que les messages s'affichent sur une seule ligne
	e.Status.SetText(fmt.Sprintf("[yellow]%s", msg))
}
