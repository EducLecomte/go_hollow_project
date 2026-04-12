# Hollow

![Banner](banner.png)

**Hollow** est un éditeur de texte TUI (Terminal User Interface) moderne écrit en Go. Il fusionne la simplicité d'utilisation de **Nano** avec la puissance de navigation et de gestion de fichiers distants de **mcedit** (Midnight Commander).

Ce projet est développé en collaboration avec Gemini Code Assist, dans un but récréatif et éducatif.


## Fonctionnalités

- **Explorateur de fichiers intégré** : Navigation fluide dans l'arborescence locale.
- **Éditeur de texte avec indicateur d'état** : Visualisation en temps réel des modifications non sauvegardées (symbole `*`).
- **Manipulation de lignes (Nano-style)** : Raccourcis `Ctrl+K` et `Ctrl+U` pour couper et coller des lignes entières.
- **Système de Fichiers Virtuel (VFS)** : Architecture prête pour le support FTP et Archives.
- **Presse-papier système** : Copier-coller intégré avec l'OS (X11/Wayland).
- **Aide contextuelle** : Barre de raccourcis dynamique et documentation interactive via `F1`.
- **Sauvegarde non-interruptive** : `Ctrl+S` enregistre votre travail sans fermer l'éditeur.

## Architecture Technique

Le projet est construit de manière modulaire :

- `main.go` : Point d'entrée, initialise l'application.
- `internal/app/` : Cœur de l'application (UI, Handlers, Actions).
- `internal/vfs/` : Abstraction et implémentations du système de fichiers (VFS).
- `internal/utils/` : Constantes d'aide et formatage utilitaire.

### Le coeur : L'interface VFS

L'extensibilité du projet repose sur l'interface `VFS`, permettant d'ajouter des protocoles sans modifier l'UI :

```go
type VFS interface {
    List(path string) ([]FileInfo, error)
    Read(path string) (io.ReadCloser, error)
    Write(path string, data io.Reader) error
}
```

## Raccourcis Clavier

| Touche | Action |
| :--- | :--- |
| `F1` | Afficher l'aide interactive |
| `TAB` | Basculer entre l'explorateur et le visualiseur |
| `Ctrl + X` | Quitter l'application (depuis la liste) |
| `F10` | Quitter l'application |
| `Ctrl + F` | Créer un nouveau fichier |
| `Ctrl + D` | Créer un nouveau dossier |
| `Ctrl + R` | Supprimer le fichier ou dossier sélectionné |
| `Ctrl + K` | Copier le fichier ou dossier sélectionné |
| `Ctrl + U` | Coller l'élément dans le dossier actuel |
| `Ctrl + S` | Sauvegarder le fichier (Éditeur) |
| `Ctrl + K` | Couper la ligne entière (concatène en bloc si répété) |
| `Ctrl + U` | Coller le bloc de texte coupé (Éditeur) |
| `Entrée` | Ouvrir un fichier ou entrer dans un dossier |

## Installation & Utilisation

### Prérequis

- Go 1.18+
- `xclip`, `xsel` ou `wl-clipboard` (pour le support du presse-papier sur Linux)

### Mise à jour des dépendances

Pour bénéficier des dernières améliorations de l'interface (comme la numérotation des lignes), assurez-vous de mettre à jour `tview` :
```bash
go get -u github.com/rivo/tview
```

### Lancement

```bash
go run .
```

## État du Projet

### Implémenté
- Navigation locale avec gestion des métadonnées (taille).
- Lecture et écriture réelle sur le disque.
- Architecture factorisée pour la maintenabilité.
- Indicateur visuel de modification dans l'éditeur.
- Mécanique de "Cut & Paste" de blocs de lignes sans perte du curseur (Ctrl+K / Ctrl+U).

### En cours / À venir
- [ ] **Client FTP** : Implémentation de `FTPFS` pour l'édition distante.
- [ ] **Explorateur d'archives** : Support des fichiers `.zip` et `.tar.gz`.
- [ ] **Recherche avancée** : Logique de recherche textuelle avec surlignage.
- [ ] **Numérotation des lignes** : Activation après mise à jour de la bibliothèque TUI.

---
*Dernière mise à jour : Dimanche 12 Avril 2026 - 11:30*