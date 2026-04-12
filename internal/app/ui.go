package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/EducLecomte/go_hollow_project/internal/vfs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type EditorApp struct {
	// Infrastructure Tview
	App   *tview.Application
	Pages *tview.Pages

	// Composants de l'interface principale
	PathBar     *tview.TextView
	FileList    *tview.List
	FileSizeBox *tview.TextView
	Viewer      *tview.TextView
	Status      *tview.TextView

	// Système de fichiers et Navigation
	FileSystem   vfs.VFS
	CurrentDir   string
	CurrentFiles []vfs.FileInfo
	PreviousFS   vfs.VFS
	PreviousDir  string

	// État de l'éditeur et Presse-papiers
	FilePath      string
	CopiedPath    string
	Clipboard     string
	LastSearch    string
	LastSearchPos int

	// Gestion de l'asynchronisme
	previewCancel context.CancelFunc
}

// NewEditorApp initialise une nouvelle instance de l'application Hollow.
// Par défaut, elle démarre sur le système de fichiers local dans le répertoire courant.
func NewEditorApp() *EditorApp {
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

// setupUI configure la disposition des widgets, les styles et les comportements de base de l'interface.
func (e *EditorApp) setupUI() {
	e.PathBar.SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDarkGreen)

	e.FileList.SetBorder(true).SetTitle(" Exploreur ").SetBorderColor(tcell.ColorYellow)
	e.FileList.ShowSecondaryText(false)
	e.FileList.SetSelectedBackgroundColor(tcell.ColorWhite).
		SetSelectedTextColor(tcell.ColorBlack)

	e.FileList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		e.handleFileSelection(index)
	})

	// Gestion dynamique de la couleur des dossiers et mise à jour de l'encart d'info
	isUpdatingList := false
	// Mise à jour asynchrone du visualiseur pour éviter les blocages (surtout en FTP)
	e.FileList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// 1. Annulation de la prévisualisation précédente
		if e.previewCancel != nil {
			e.previewCancel()
		}

		if index == 0 {
			e.FileSizeBox.SetText("[gray]Parent Directory")
			e.Viewer.SetText("").SetTitle(" Visualiseur ")
			return
		}

		if e.CurrentFiles == nil || index-1 >= len(e.CurrentFiles) {
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		e.previewCancel = cancel

		file := e.CurrentFiles[index-1]
		modTimeStr := file.ModTime.Format("2006-01-02 15:04")
		path := filepath.Join(e.CurrentDir, file.Name)

		// Mise à jour immédiate des infos basiques (synchrone)
		if file.IsDir {
			e.FileSizeBox.SetText(fmt.Sprintf("[green]Type: [white]Dossier\n[green]Date: [white]%s\n[green]Droits: [white]%s\n[green]Owner: [white]%s", modTimeStr, file.Permissions, file.Owner))
		} else {
			e.FileSizeBox.SetText(fmt.Sprintf("[green]Taille: [white]%s\n[green]Date: [white]%s\n[green]Droits: [white]%s\n[green]Owner: [white]%s", utils.FormatSize(file.Size), modTimeStr, file.Permissions, file.Owner))
		}

		// Prévisualisation asynchrone (E/S et Coloration)
		go func() {
			// Petite pause pour éviter de charger inutilement lors d'un défilement rapide
			time.Sleep(100 * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			default:
			}

			if file.IsDir {
				e.previewDirectory(ctx, path)
			} else {
				e.previewFile(ctx, path)
			}
		}()

		// 2. Gestion dynamique de la couleur des dossiers pour le contraste
		if isUpdatingList {
			return
		}
		isUpdatingList = true
		defer func() { isUpdatingList = false }()

		for i := 0; i < e.FileList.GetItemCount(); i++ {
			m, s := e.FileList.GetItemText(i)
			if !strings.HasSuffix(m, "/") && !strings.HasPrefix(m, "[#ff8c00]") {
				continue
			}

			// Nettoyage du nom
			name := strings.TrimPrefix(m, "[#ff8c00]")
			if strings.HasSuffix(name, "/") {
				if i == index {
					// Sélectionné : pas de tag pour être noir sur blanc
					if m != name {
						e.FileList.SetItemText(i, name, s)
					}
				} else {
					// Non sélectionné : orange
					if !strings.HasPrefix(m, "[#ff8c00]") {
						e.FileList.SetItemText(i, "[#ff8c00]"+name, s)
					}
				}
			}
		}
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

	e.Viewer.SetBorder(true).SetTitle(" Visualiseur ").SetBorderColor(tcell.ColorWhite)
	e.Viewer.SetDynamicColors(true).SetRegions(true) // Active le support des couleurs ANSI/Tags
	e.Viewer.SetFocusFunc(func() {
		e.Viewer.SetBorderColor(tcell.ColorYellow)
		e.updateStatus(utils.HelpMsgView)
	})
	e.Viewer.SetBlurFunc(func() {
		e.Viewer.SetBorderColor(tcell.ColorWhite)
	})
	e.Viewer.SetWrap(true)    // Rétablit le retour à la ligne automatique
	e.Viewer.SetDrawFunc(nil) // Supprime la fonction de synchronisation obsolète

	// Capture de F1 et autres touches dans le visualiseur
	e.Viewer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			e.showHelp(utils.HelpContentExplorer)
			return nil
		}
		if event.Key() == tcell.KeyTab {
			e.App.SetFocus(e.FileList)
			return nil
		}
		return event
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

// updateStatus met à jour le texte de la barre d'état en bas de l'écran.
func (e *EditorApp) updateStatus(msg string) {
	// On s'assure que les messages s'affichent sur une seule ligne
	e.Status.SetText(fmt.Sprintf("[yellow]%s", msg))
}

// updateStatusTemp affiche un message temporaire dans la barre d'état et le restaure après un délai de 5 secondes.
func (e *EditorApp) updateStatusTemp(msg string) {
	e.updateStatus(msg)

	go func() {
		time.Sleep(5 * time.Second)
		// tview n'est pas thread-safe, on utilise QueueUpdateDraw pour mettre à jour l'UI
		e.App.QueueUpdateDraw(func() {
			// Restauration du message d'aide selon le focus actuel
			focus := e.App.GetFocus()
			if focus == e.Viewer {
				e.updateStatus(utils.HelpMsgView)
			} else if _, ok := focus.(*tview.TextArea); ok {
				e.updateStatus(utils.HelpMsgEdit)
			} else {
				e.updateStatus(utils.HelpMsgFiles)
			}
		})
	}()
}

// connectFTP initialise une connexion à un serveur distant et bascule le système de fichiers de l'application.
func (e *EditorApp) connectFTP(host string, port int, user, pass string) error {
	ftpFS, err := vfs.NewFtpFS(host, port, user, pass)
	if err != nil {
		return err
	}

	// Sauvegarde du système actuel pour permettre le retour
	e.PreviousFS = e.FileSystem
	e.PreviousDir = e.CurrentDir

	e.FileSystem = ftpFS
	e.CurrentDir = "/"
	e.refreshFileList()
	return nil
}
