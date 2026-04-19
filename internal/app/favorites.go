package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Favorite représente un marque-page vers un dossier.
type Favorite struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// getFavoritesFilePath retourne le chemin d'accès au fichier de configuration des favoris.
func (e *EditorApp) getFavoritesFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".config")
	}
	appConfigDir := filepath.Join(configDir, "hollow")
	if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
		os.MkdirAll(appConfigDir, 0755)
	}
	return filepath.Join(appConfigDir, "favorites.json")
}

// loadFavorites lit le fichier de favoris s'il existe et l'injecte dans la structure de l'application.
func (e *EditorApp) loadFavorites() {
	filePath := e.getFavoritesFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		e.Favorites = []Favorite{}
		return
	}
	var favs []Favorite
	if err := json.Unmarshal(data, &favs); err != nil {
		e.Favorites = []Favorite{}
		return
	}
	e.Favorites = favs
}

// saveFavorites enregistre la liste actuelle de favoris dans le fichier de configuration.
func (e *EditorApp) saveFavorites() {
	filePath := e.getFavoritesFilePath()
	data, err := json.MarshalIndent(e.Favorites, "", "  ")
	if err == nil {
		os.WriteFile(filePath, data, 0644)
	}
}

// showFavoritesDialog affiche la modale listant les dossiers favoris.
func (e *EditorApp) showFavoritesDialog() {
	list := tview.NewList().ShowSecondaryText(true)
	list.SetBorder(true).SetTitle(" Dossiers Favoris (Ctrl+B) ").SetTitleAlign(tview.AlignCenter)
	list.SetSelectedBackgroundColor(tcell.ColorWhite).SetSelectedTextColor(tcell.ColorBlack)

	refreshList := func() {
		list.Clear()
		list.AddItem("[+] Ajouter le dossier courant", e.CurrentDir, 'a', nil)
		for i, fav := range e.Favorites {
			shortcut := rune(0)
			if i < 9 {
				shortcut = rune('1' + i)
			}
			list.AddItem(fav.Name, fav.Path, shortcut, nil)
		}
	}

	refreshList()

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index == 0 {
			// Ajouter le dossier courant
			fav := Favorite{
				Name: filepath.Base(e.CurrentDir),
				Path: e.CurrentDir,
			}
			if fav.Name == "." || fav.Name == "/" || fav.Name == "" {
				fav.Name = "Racine"
			}
			e.Favorites = append(e.Favorites, fav)
			e.saveFavorites()
			refreshList()
			list.SetCurrentItem(len(e.Favorites)) // Sélectionner le nouvel élément rajouté
		} else {
			// Naviguer vers le favori
			if index-1 < len(e.Favorites) {
				targetPath := e.Favorites[index-1].Path
				
				// Revenir au système de fichiers local si on était dans un sous-système (ex: Archive)
				if e.PreviousFS != nil {
					e.FileSystem = e.PreviousFS
					e.PreviousFS = nil
				}
				
				e.CurrentDir = targetPath
				e.refreshFileList()
				e.Pages.RemovePage("favorites")
				e.App.SetFocus(e.FileList)
			}
		}
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyCtrlB || event.Key() == tcell.KeyLeft {
			e.Pages.RemovePage("favorites")
			e.App.SetFocus(e.FileList)
			return nil
		}
		if event.Key() == tcell.KeyDelete {
			index := list.GetCurrentItem()
			if index > 0 && index-1 < len(e.Favorites) {
				e.Favorites = append(e.Favorites[:index-1], e.Favorites[index:]...)
				e.saveFavorites()
				refreshList()
				// Ajuster la sélection
				if index < list.GetItemCount() {
					list.SetCurrentItem(index)
				} else {
					list.SetCurrentItem(list.GetItemCount() - 1)
				}
			}
			return nil
		}
		return event
	})

	e.showCenteredDialog("favorites", list, 60, 20)
}
