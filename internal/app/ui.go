package app

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/EducLecomte/go_hollow_project/internal/vfs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type EditorApp struct {
	App          *tview.Application
	PathBar      *tview.TextView
	FileList     *tview.List
	FileSizeBox  *tview.TextView
	Viewer       *tview.TextView
	Status       *tview.TextView
	Pages        *tview.Pages
	FilePath     string
	CopiedPath   string
	CurrentDir   string
	FileSystem   vfs.VFS
	CurrentFiles []vfs.FileInfo
}

func NewEditorApp() *EditorApp {
	// Initialisation par défaut en local
	localFS := &vfs.LocalFS{}
	wd, err := filepath.Abs(".")
	if err != nil {
		wd = "/"
	}

	e := &EditorApp{
		App:         tview.NewApplication(),
		PathBar:     tview.NewTextView(),
		FileList:    tview.NewList(),
		FileSizeBox: tview.NewTextView(),
		Viewer:      tview.NewTextView(),
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
		SetBackgroundColor(tcell.ColorDarkGreen)

	e.FileList.SetBorder(true).SetTitle(" Exploreur ").SetBorderColor(tcell.ColorYellow)
	e.FileList.ShowSecondaryText(false) // Rend la liste compacte (une seule ligne)
	e.FileList.SetSelectedBackgroundColor(tcell.ColorWhite).
		SetSelectedTextColor(tcell.ColorBlack)
	e.FileList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		e.handleFileSelection(index)
	})
	e.FileList.SetFocusFunc(func() {
		e.FileList.SetBorderColor(tcell.ColorYellow)
		e.updateStatus(utils.HelpMsgFiles)
	})
	e.FileList.SetBlurFunc(func() {
		e.FileList.SetBorderColor(tcell.ColorWhite)
	})

	// Encart pour le poids du fichier
	e.FileSizeBox.SetBorder(true).SetTitle(" Info ")
	e.FileSizeBox.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)

	// Mise à jour de l'encart quand on navigue dans la liste
	e.FileList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index == 0 {
			e.FileSizeBox.SetText("[gray]Parent Directory")
			e.Viewer.SetText("").SetTitle(" Visualiseur ")
			return
		}
		file := e.CurrentFiles[index-1]
		modTimeStr := file.ModTime.Format("2006-01-02 15:04")
		if file.IsDir {
			e.FileSizeBox.SetText(fmt.Sprintf("[green]Type: [white]Dossier\n[green]Date: [white]%s\n[green]Droits: [white]%s\n[green]Owner: [white]%s", modTimeStr, file.Permissions, file.Owner))
			e.previewDirectory(filepath.Join(e.CurrentDir, file.Name))
		} else {
			e.FileSizeBox.SetText(fmt.Sprintf("[green]Taille: [white]%s\n[green]Date: [white]%s\n[green]Droits: [white]%s\n[green]Owner: [white]%s", utils.FormatSize(file.Size), modTimeStr, file.Permissions, file.Owner))
			e.previewFile(filepath.Join(e.CurrentDir, file.Name))
		}
	})

	e.Viewer.SetBorder(true).SetTitle(" Visualiseur ").SetBorderColor(tcell.ColorWhite)
	e.Viewer.SetDynamicColors(true).SetRegions(true) // Active le support des couleurs ANSI/Tags
	e.Viewer.SetFocusFunc(func() {
		e.Viewer.SetBorderColor(tcell.ColorYellow)
		e.updateStatus(utils.HelpMsgEdit)
	})
	e.Viewer.SetBlurFunc(func() {
		e.Viewer.SetBorderColor(tcell.ColorWhite)
	})

	e.Status.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	e.updateStatus(utils.HelpMsgDefault)

	// Layout de la colonne de gauche (Explorateur + Info taille)
	leftColumn := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(e.FileList, 0, 1, true).
		AddItem(e.FileSizeBox, 6, 0, false)

	// Layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(e.PathBar, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(leftColumn, 0, 1, true).
			AddItem(e.Viewer, 0, 2, false), 0, 1, true).
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
		var displayName string
		if f.IsDir {
			// Dossier : On accole le slash au nom pour une identification immédiate
			displayName = "[darkorange]" + f.Name + "/"
		} else {
			// Fichier : Nom simple
			displayName = f.Name
		}
		e.FileList.AddItem(displayName, "", 0, nil)
	}
}
func (e *EditorApp) updateStatus(msg string) {
	// On s'assure que les messages s'affichent sur une seule ligne
	e.Status.SetText(fmt.Sprintf("[yellow]%s", msg))
}

// updateStatusTemp affiche un message puis restaure les raccourcis après 5 secondes
func (e *EditorApp) updateStatusTemp(msg string) {
	e.updateStatus(msg)

	go func() {
		time.Sleep(5 * time.Second)
		// tview n'est pas thread-safe, on utilise QueueUpdateDraw pour mettre à jour l'UI
		e.App.QueueUpdateDraw(func() {
			// Restauration du message d'aide selon le focus actuel
			if e.App.GetFocus() == e.Viewer {
				e.updateStatus(utils.HelpMsgEdit)
			} else {
				e.updateStatus(utils.HelpMsgFiles)
			}
		})
	}()
}
