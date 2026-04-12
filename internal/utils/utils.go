package utils

import "fmt"

const (
	HelpMsgDefault = "[yellow]F1:[white] Aide | [yellow]TAB:[white] Visualiser | [yellow]Ctrl+F:[white] Fich. | [yellow]Ctrl+D:[white] Doss. | [yellow]Ctrl+R:[white] Suppr. | [yellow]Ctrl+K/U:[white] C/V | [yellow]Ctrl+X:[white] Quitter"
	HelpMsgEdit    = "[yellow]F1:[white] Aide | [yellow]TAB:[white] Retour | [yellow]Ctrl+K/U:[white] Couper/Coller bloc | [yellow]Flèches:[white] Naviguer"
	HelpMsgFiles   = "[yellow]F1:[white] Aide | [yellow]TAB:[white] Visualiser | [yellow]Ctrl+F:[white] Fich. | [yellow]Ctrl+D:[white] Doss. | [yellow]Ctrl+R:[white] Suppr. | [yellow]Ctrl+K/U:[white] C/V | [yellow]Ctrl+X:[white] Quitter"

	HelpContent = `
 [yellow]Navigation & Système[white]
 --------------------
 F1          : Afficher cette aide
 F10         : Quitter l'application
 TAB         : Basculer entre l'Explorateur et le Visualiseur
 Ctrl + X    : Quitter l'app (depuis l'explorateur)
 Ctrl + F    : Créer un nouveau fichier (depuis l'explorateur)
 Ctrl + D    : Créer un nouveau dossier (depuis l'explorateur)
 Ctrl + R    : Supprimer un fichier/dossier (depuis l'explorateur)
 Ctrl + K    : Copier un fichier/dossier (depuis l'explorateur)
 Ctrl + U    : Coller un fichier/dossier (depuis l'explorateur)
 Entrée      : Ouvrir un fichier / Entrer dans un dossier

 [yellow]Édition[white]
 -------
 Ctrl + S    : Sauvegarder le fichier courant
 Ctrl + K    : Couper la ligne (concatène en bloc si répété)
 Ctrl + U    : Coller le texte ou le bloc coupé
 `
)

func FormatSize(b int64) string {
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
