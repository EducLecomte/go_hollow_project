# Hollow

![Banner](banner.png)

**Hollow** est un éditeur de texte TUI (Terminal User Interface) moderne écrit en Go. Il fusionne la simplicité d'utilisation de **Nano** avec la puissance de navigation et de gestion de fichiers distants de **mcedit** (Midnight Commander).

Ce projet est développé avec une IA, dans un but récréatif et éducatif.


## Fonctionnalités

- **Explorateur de fichiers intégré** : Navigation fluide dans l'arborescence locale.
- **Prévisualisation intelligente** : Coloration syntaxique pour les fichiers (via Chroma) et vue arborescente pour les dossiers.
- **Éditeur de texte avancé** : Mode plein écran avec indicateur d'état (`*`), support de la coloration et gestion des fins de ligne.
- **Manipulation de lignes (Nano-style)** : Raccourcis `Ctrl+K` et `Ctrl+U` pour couper et coller des lignes entières.
- **Système de Fichiers Virtuel (VFS)** : Support natif du système local et des archives (Lecture seule).
- **Aide contextuelle** : Barre de raccourcis dynamique et documentation interactive via `F1`.
- **Sauvegarde non-interruptive** : `Ctrl+S` enregistre votre travail sans fermer l'éditeur.
- **Intégration Shell** : Script d'installation permettant de synchroniser le répertoire de travail du terminal avec celui de l'éditeur à la fermeture.

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
| `Ctrl + F` | Créer un nouveau fichier |
| `Ctrl + D` | Créer un nouveau dossier |
| `Ctrl + R` | Supprimer le fichier ou dossier sélectionné |
| `Ctrl + K` | Copier le fichier (Explorateur) / Couper la ligne (Éditeur) |
| `Ctrl + U` | Coller le fichier (Explorateur) / Coller le bloc (Éditeur) |
| `Ctrl + S` | Sauvegarder le fichier (Éditeur) |
| `Entrée` | Ouvrir un fichier ou entrer dans un dossier |
| `Echap` | Fermer un dialogue ou quitter l'éditeur |

## Installation & Utilisation

### Prérequis

- Go 1.18+

### Installation rapide

Utilisez le script d'installation fourni pour compiler le projet et configurer votre shell :
```bash
chmod +x install.sh
./install.sh
```

## État du Projet

### Implémenté
- Navigation locale avec gestion des métadonnées (taille).
- Lecture et écriture réelle sur le disque.
- Architecture factorisée pour la maintenabilité.
- Indicateur visuel de modification dans l'éditeur.
- Mécanique de "Cut & Paste" de blocs de lignes sans perte du curseur (Ctrl+K / Ctrl+U).
- Explorateur d'archives: Navigation transparente dans les fichiers `.zip`, `.tar` et `.tar.gz`.

### En cours / À venir
- [ ] **Client FTP** : Implémentation de `FTPFS` pour l'édition distante.
- [ ] **Recherche avancée** : Logique de recherche textuelle avec surlignage.
- [ ] **Numérotation des lignes** : Activation après mise à jour de la bibliothèque TUI.

---
*Dernière mise à jour : Dimanche 12 Avril 2026 - 12:45*