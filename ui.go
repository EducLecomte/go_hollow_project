package main

import (
	"fmt"
	"path/filepath"
	"time"

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
	CopiedPath   string
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
		e.updateStatus(HelpMsgFiles)
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
			return
		}
		file := e.CurrentFiles[index-1]

		if file.IsDir {
			e.FileSizeBox.SetText(fmt.Sprintf("[green]Type: [white]Dossier\n[green]Droits: [white]%s\n[green]Owner: [white]%s", file.Permissions, file.Owner))
		} else {
			e.FileSizeBox.SetText(fmt.Sprintf("[green]Taille: [white]%s\n[green]Droits: [white]%s\n[green]Owner: [white]%s", formatSize(file.Size), file.Permissions, file.Owner))
		}
	})

	e.Editor.SetBorder(true).SetTitle(" Éditeur ").SetBorderColor(tcell.ColorWhite)
	e.Editor.SetPlaceholder("Entrez votre texte ici... (Ctrl+S pour sauver, Ctrl+K/U pour couper/coller)")
	e.Editor.SetFocusFunc(func() {
		e.Editor.SetBorderColor(tcell.ColorYellow)
		// Force le style de sélection en noir sur blanc pour l'éditeur
		e.Editor.SetSelectedStyle(tcell.StyleDefault.
			Background(tcell.ColorWhite).
			Foreground(tcell.ColorBlack))
		e.updateStatus(HelpMsgEdit)
	})
	e.Editor.SetBlurFunc(func() {
		e.Editor.SetBorderColor(tcell.ColorWhite)
	})

	e.Status.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	e.updateStatus(HelpMsgDefault)

	// Layout de la colonne de gauche (Explorateur + Info taille)
	leftColumn := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(e.FileList, 0, 1, true).
		AddItem(e.FileSizeBox, 5, 0, false)

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
	header := fmt.Sprintf("%-25s [green]Date", "..")
	e.FileList.AddItem(header, "Retour au parent", '.', nil)

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
		modTimeStr := f.ModTime.Format("2006-01-02 15:04")
		if f.IsDir {
			// Dossier : On accole le slash au nom pour une identification immédiate, puis on aligne la colonne
			displayName = fmt.Sprintf("[darkorange]%-25s [gray]%s", f.Name+"/", modTimeStr)
		} else {
			// Fichier : Nom sans tag, date en gris pour l'uniformité
			displayName = fmt.Sprintf("%-25s [gray]%s", f.Name, modTimeStr)
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
			if e.App.GetFocus() == e.Editor {
				e.updateStatus(HelpMsgEdit)
			} else {
				e.updateStatus(HelpMsgFiles)
			}
		})
	}()
}
