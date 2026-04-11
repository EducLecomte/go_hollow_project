package main

import "fmt"

const (
	HelpMsgDefault = "F1: Aide | TAB: Editer | Ctrl+F: Fichier | Ctrl+D: Dossier | Ctrl+S: Sauver | Ctrl+X: Quitter"
	HelpMsgEdit    = "F1: Aide | Ctrl+S: Sauver | Ctrl+K: Couper | Ctrl+U/V: Coller | Ctrl+X: Retour"
	HelpMsgFiles   = "F1: Aide | TAB: Editer | Ctrl+F: Fichier | Ctrl+D: Dossier | Ctrl+S: Sauver | Ctrl+X: Quitter"

	HelpContent = `
 [yellow]Navigation & Système[white]
 --------------------
 F1          : Afficher cette aide
 F10         : Quitter l'application
 TAB         : Basculer vers l'Éditeur
 Ctrl + X    : Quitter l'app (depuis l'explorateur)
 Ctrl + F    : Créer un nouveau fichier (depuis l'explorateur)
 Ctrl + D    : Créer un nouveau dossier (depuis l'explorateur)
 Ctrl + X    : Revenir à l'explorateur (depuis l'éditeur)
 Entrée      : Ouvrir un fichier / Entrer dans un dossier

 [yellow]Édition[white]
 -------
 Ctrl + S    : Sauvegarder le fichier courant
 Ctrl + K    : Couper la ligne de texte (Nano style)
 Ctrl + U/V  : Coller le texte (Uncut/Paste)
 `
)

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
