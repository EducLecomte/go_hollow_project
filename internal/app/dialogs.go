package app

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/EducLecomte/go_hollow_project/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// showCenteredDialog est une fonction utilitaire pour positionner n'importe quel composant au centre de l'écran par-dessus l'interface principale.
func (e *EditorApp) showCenteredDialog(pageName string, item tview.Primitive, width, height int) {
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(item, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)

	e.Pages.AddPage(pageName, flex, true, true)
	e.App.SetFocus(item)
}

// showHelp affiche une fenêtre modale contenant la liste complète des raccourcis clavier et l'aide utilisateur.
func (e *EditorApp) showHelp() {
	previousFocus := e.App.GetFocus()
	helpText := tview.NewTextView().
		SetText(utils.HelpContent).
		SetDynamicColors(true).
		SetScrollable(true).
		SetTextAlign(tview.AlignLeft)

	helpText.SetBorder(true).
		SetTitle(" Aide Hollow - Raccourcis ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderPadding(1, 1, 2, 2)

	e.showCenteredDialog("help", helpText, 65, 20)

	helpText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyF1 || event.Rune() == 'q' {
			e.Pages.RemovePage("help")
			if previousFocus != nil {
				e.App.SetFocus(previousFocus)
			}
			return nil
		}
		return event
	})
}

// showQuitConfirmation affiche une boîte de dialogue demandant confirmation avant de quitter définitivement l'application.
func (e *EditorApp) showQuitConfirmation() {
	modal := tview.NewModal().
		SetText("Voulez-vous vraiment quitter Hollow ?").
		AddButtons([]string{"Oui", "Non"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Oui" {
				e.saveLastDir()
				e.App.Stop()
			}
			e.Pages.RemovePage("quit")
		})
	e.Pages.AddPage("quit", modal, true, true)
}

// showDeleteConfirmation affiche une confirmation avant de supprimer physiquement un fichier ou un dossier.
func (e *EditorApp) showDeleteConfirmation() {
	index := e.FileList.GetCurrentItem()
	if index <= 0 || index-1 >= len(e.CurrentFiles) {
		return
	}
	file := e.CurrentFiles[index-1]
	path := filepath.Join(e.CurrentDir, file.Name)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Voulez-vous vraiment supprimer %s ?", file.Name)).
		AddButtons([]string{"Supprimer", "Annuler"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Supprimer" {
				e.deleteElement(path)
			}
			e.Pages.RemovePage("delete")
			e.App.SetFocus(e.FileList)
		})
	e.Pages.AddPage("delete", modal, true, true)
}

// showNewFileDialog affiche une boîte de saisie pour nommer et créer un nouveau fichier dans le dossier courant.
func (e *EditorApp) showNewFileDialog() {
	inputField := tview.NewInputField().SetLabel(" Nom du nouveau fichier: ")
	inputField.SetBorder(true).SetTitle(" Nouveau Fichier ").SetTitleAlign(tview.AlignCenter)

	e.showCenteredDialog("newfile", inputField, 60, 3)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			name := inputField.GetText()
			if name != "" {
				e.createFile(name)
				e.App.SetFocus(e.Viewer)
			}
			e.Pages.RemovePage("newfile")
		} else if key == tcell.KeyEsc {
			e.Pages.RemovePage("newfile")
			e.App.SetFocus(e.FileList)
		}
	})
}

// showSaveConfirmation demande à l'utilisateur s'il souhaite sauvegarder ses modifications avant de fermer l'éditeur plein écran.
func (e *EditorApp) showSaveConfirmation(content string) {
	modal := tview.NewModal().
		SetText("Voulez-vous sauvegarder les modifications avant de quitter ?").
		AddButtons([]string{"Sauvegarder", "Ignorer", "Annuler"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Sauvegarder":
				e.saveFromFullEditor(content)
				e.Pages.RemovePage("edit_screen")
				e.App.SetFocus(e.FileList)
			case "Ignorer":
				e.Pages.RemovePage("edit_screen")
				e.App.SetFocus(e.FileList)
			case "Annuler":
				// On ferme juste la modale, le focus revient à l'éditeur
			}
			e.Pages.RemovePage("save_confirm")
		})

	e.Pages.AddPage("save_confirm", modal, true, true)
}

// showNewDirDialog affiche une boîte de saisie pour créer un nouveau répertoire.
func (e *EditorApp) showNewDirDialog() {
	inputField := tview.NewInputField().SetLabel(" Nom du nouveau dossier: ")
	inputField.SetBorder(true).SetTitle(" Nouveau Dossier ").SetTitleAlign(tview.AlignCenter)

	e.showCenteredDialog("newdir", inputField, 60, 3)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			name := inputField.GetText()
			if name != "" {
				e.createDir(name)
			}
			e.Pages.RemovePage("newdir")
			e.App.SetFocus(e.FileList)
		} else if key == tcell.KeyEsc {
			e.Pages.RemovePage("newdir")
			e.App.SetFocus(e.FileList)
		}
	})
}

// showFTPDialog affiche un message d'information concernant la future fonctionnalité de connexion distante.
func (e *EditorApp) showFTPDialog() {
	modal := tview.NewModal().
		SetText("Connexion Distante (FTP) : Fonctionnalité à configurer.").
		AddButtons([]string{"Fermer"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.Pages.RemovePage("ftp")
			e.App.SetFocus(e.FileList)
		})
	e.Pages.AddPage("ftp", modal, true, true)
}

// showLoadingDialog affiche une modale d'attente pour les opérations longues avec option d'annulation.
func (e *EditorApp) showLoadingDialog(title string, message string, cancelFunc context.CancelFunc) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Annuler"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Annuler" && cancelFunc != nil {
				cancelFunc()
			}
		})
	e.Pages.AddPage("loading", modal, true, true)
}
