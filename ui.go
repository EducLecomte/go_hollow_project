package main

import (
	"fmt"
	"path/filepath"

	"github.com/rivo/tview"
)

type EditorApp struct {
	App          *tview.Application
	FileList     *tview.List
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
	wd, _ := filepath.Abs(".")

	e := &EditorApp{
		App:        tview.NewApplication(),
		FileList:   tview.NewList(),
		Editor:     tview.NewTextArea(),
		Status:     tview.NewTextView(),
		Pages:      tview.NewPages(),
		CurrentDir: wd,
		FileSystem: localFS,
	}

	e.setupUI()
	e.setupHandlers()
	e.refreshFileList()
	return e
}

func (e *EditorApp) setupUI() {
	e.FileList.SetBorder(true).SetTitle(" Fichiers (Alt+F) ")
	e.FileList.ShowSecondaryText(false) // Rend la liste compacte (une seule ligne)
	e.FileList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		e.handleFileSelection(index)
	})

	e.Editor.SetBorder(true).SetTitle(" Éditeur ")
	e.Editor.SetPlaceholder("Entrez votre texte ici... (Ctrl+S pour sauver, Ctrl+C/V pour copier/coller)")

	e.Status.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	e.updateStatus(HelpMsgDefault)

	// Layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(e.FileList, 30, 1, true).
			AddItem(e.Editor, 0, 4, false), 0, 1, true).
		AddItem(e.Status, 1, 0, false)

	e.Pages.AddPage("main", mainFlex, true, true)
}

func (e *EditorApp) refreshFileList() {
	e.FileList.Clear()
	e.FileList.AddItem("..", "Retour au parent", '.', nil)

	files, err := e.FileSystem.List(e.CurrentDir)
	if err != nil {
		e.updateStatus(fmt.Sprintf("[red]Erreur listage: %v", err))
		return
	}

	e.CurrentFiles = files
	for _, f := range files {
		suffix := ""
		if f.IsDir {
			suffix = "/"
		}
		// Affiche le nom et la taille sur la même ligne
		display := fmt.Sprintf("%-30s [blue]%10s", f.Name+suffix, formatSize(f.Size))
		e.FileList.AddItem(display, "", 0, nil)
	}
	e.FileList.SetTitle(fmt.Sprintf(" Fichiers: %s ", e.CurrentDir))
}
func (e *EditorApp) updateStatus(msg string) {
	// On s'assure que les messages s'affichent sur une seule ligne
	e.Status.SetText(fmt.Sprintf("[yellow]%s", msg))
}
