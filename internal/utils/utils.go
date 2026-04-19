package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	HelpMsgDefault = "[yellow]F1:[white] Aide | [yellow]F3:[white] FTP | [yellow]Ctrl+F/D:[white] Fich/Doss | [yellow]Ctrl+E:[white] Extr | [yellow]Suppr:[white] Suppr | [yellow]Ctrl+X:[white] Quitter"
	HelpMsgEdit    = "[yellow]F1:[white] Aide | [yellow]Ctrl+S:[white] Sauver | [yellow]Ctrl+F:[white] Chercher | [yellow]Ctrl+K/U:[white] C/V | [yellow]Esc:[white] Quitter"
	HelpMsgView    = "[yellow]F1:[white] Aide | [yellow]TAB/Ctrl+X:[white] Explorer | [yellow]Flèches/Molette:[white] Défiler"
	HelpMsgArchive = "[yellow]F1:[white] Aide | [yellow]Entrée:[white] Aperçu | [yellow]Ctrl+E:[white] Extraire ficher | [yellow]..:[white] Sortir"
	HelpMsgFiles   = HelpMsgDefault

	HelpContentExplorer = `
 [yellow]Navigation & Système[white]
 --------------------
 F1          : Afficher cette aide
 F3          : Connexion FTP / Distante
 TAB / Ctrl+X: Basculer entre l'Explorateur et le Visualiseur
 Ctrl + X    : Quitter l'application (quand l'Explorateur a le focus)
 Entrée      : Ouvrir un fichier / Entrer dans un dossier
 ..          : Remonter au dossier parent
 
 [yellow]Opérations Fichiers[white]
 --------------------
 Ctrl + F    : Créer un nouveau fichier
 Ctrl + D    : Créer un nouveau dossier
 Suppr       : Supprimer l'élément sélectionné
 Ctrl + E    : Extraire l'archive sélectionnée
 Ctrl + K    : Copier (mémoriser le chemin)
 Ctrl + U    : Coller (copie physique)
 `

	HelpContentArchive = `
 [yellow]Exploration d'Archive[white]
 ----------------------
 F1          : Afficher cette aide
 Entrée      : Visualiser un fichier dans l'archive
 ..          : Remonter dans l'arborescence (ou sortir de l'archive si à la racine)
 TAB         : Basculer vers le visualiseur
 
 [yellow]Actions Spécifiques[white]
 --------------------
 Ctrl + E    : Extraire le fichier ou dossier sélectionné vers le dossier hôte
 `

	HelpContentEditor = `
 [yellow]Édition de Texte[white]
 ----------------
 F1          : Afficher cette aide
 Ctrl + S    : Sauvegarder les modifications
 Ctrl + F    : Rechercher (Entrée pour suivant)
 Ctrl + K    : Couper la ligne (Nano style, concatène si répété)
 Ctrl + U    : Coller le texte ou le bloc coupé
 Esc / Ctrl+X : Fermer l'éditeur (demande confirmation si modifié)
 
 [yellow]Déplacement[white]
 -----------
 Flèches     : Se déplacer dans le texte
 Page Up/Down: Défilement rapide
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

func IsArchive(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".zip" || ext == ".tar" || ext == ".gz" || ext == ".tgz" {
		return true
	}
	// Check for double extension like .tar.gz
	if strings.HasSuffix(strings.ToLower(name), ".tar.gz") {
		return true
	}
	return false
}

// IsBinary détecte si un contenu est binaire en cherchant des octets nuls (null bytes) dans les premiers octets.
func IsBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	// On analyse un échantillon du début du fichier
	limit := 1024
	if len(data) < limit {
		limit = len(data)
	}

	for i := 0; i < limit; i++ {
		if data[i] == 0 { // Un octet nul est le signe quasi certain d'un fichier binaire
			return true
		}
	}
	return false
}

// GetBinaryFileDescription renvoie une description amicale du type de fichier basé sur l'extension.
func GetBinaryFileDescription(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	// Images
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff", ".svg", ".ico":
		return "une image"
	// Archives
	case ".zip", ".tar", ".gz", ".tgz", ".rar", ".7z", ".xz", ".bz2":
		return "une archive compressée"
	// Executables/Librairies
	case ".exe", ".dll", ".so", ".dylib", ".bin", ".elf", ".out":
		return "un exécutable ou une bibliothèque système"
	// Vidéos
	case ".mp4", ".mkv", ".avi", ".mov", ".flv", ".webm":
		return "une vidéo"
	// Audios
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return "un fichier audio"
	// Documents/Base de données
	case ".pdf":
		return "un document PDF"
	case ".epub", ".mobi":
		return "un livre numérique"
	case ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods":
		return "un document bureautique"
	case ".sqlite", ".sqlite3", ".db":
		return "une base de données"
	case ".iso", ".img":
		return "une image disque"
	default:
		return "un fichier binaire ou non textuel"
	}
}
