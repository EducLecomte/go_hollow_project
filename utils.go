package main

import "fmt"

const (
	HelpMsgDefault = "F1: Aide | TAB: Éditer | Ctrl+F: Chercher | Ctrl+S: Sauver | Ctrl+X: Quitter"
	HelpMsgEdit    = "F1: Aide | Ctrl+F: Chercher | Ctrl+S: Sauver | Ctrl+K: Couper | Ctrl+U: Coller | Ctrl+C/V: Copier/Coller | Ctrl+X: Retour"
	HelpMsgFiles   = "F1: Aide | TAB: Éditer | Ctrl+F: Chercher | Ctrl+S: Sauver | Ctrl+X: Quitter"

	HelpContent = `
 [yellow]Navigation & Système[white]
 --------------------
 F1          : Afficher cette aide
 F10         : Quitter l'application
 TAB         : Basculer vers l'Éditeur
 Ctrl + X    : Quitter l'app (depuis l'explorateur)
 Ctrl + X    : Revenir à l'explorateur (depuis l'éditeur)
 Entrée      : Ouvrir un fichier / Entrer dans un dossier

 [yellow]Édition[white]
 -------
 Ctrl + S    : Sauvegarder le fichier courant
 Ctrl + F    : Rechercher dans le fichier ou la liste
 Ctrl + K    : Couper la ligne de texte (Nano style)
 Ctrl + U    : Coller la ligne coupée (Uncut)
 Ctrl + C    : Copier le texte vers le presse-papier
 Ctrl + V    : Coller le texte depuis le presse-papier
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
