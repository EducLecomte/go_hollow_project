package main

import "fmt"

const (
	HelpMsgDefault = "Prêt - TAB: Éditer | Ctrl+S: Sauver | Ctrl+X: Quitter"
	HelpMsgEdit    = "Focus: Éditeur (Ctrl+X pour revenir aux fichiers)"
	HelpMsgFiles   = "Focus: Fichiers (TAB pour éditer)"
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
