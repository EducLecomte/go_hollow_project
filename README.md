# Hollow

![Banner](banner.png)

**Hollow** est un éditeur de texte TUI (Terminal User Interface) moderne et ultra-fluide écrit en Go. Il fusionne la simplicité d'utilisation de **Nano** avec la puissance de navigation et de gestion de fichiers distants inspirée de **mcedit** (Midnight Commander).

Ce projet est développé avec une IA, dans un but récréatif et pédagogique.

## Aperçu

![Explorateur](screenshot-explorer.png)
*L'explorateur de fichiers avec prévisualisation intelligente et navigation asynchrone.*

## Fonctionnalités Clés

- **Explorateur de fichiers multi-protocoles** : Navigation fluide dans l'arborescence locale et distante (FTP).
- **Architecture Asynchrone** : Chargement des fichiers et prévisualisations en arrière-plan avec système d'annulation intelligent (Context). L'interface ne "gèle" jamais, même sur des connexions réseaux (FTP) lentes, et les opérations longues peuvent être interrompues instantanément.
- **Client FTP Intégré** : Connectez-vous à des serveurs distants via `F3` et éditez vos fichiers comme s'ils étaient sur votre disque.
  !FTP
- **Sécurité et Robustesse** : Détection automatique des fichiers binaires (images, exécutables) avec avertissements pour éviter les affichages illisibles ou les plantages.
- **Explorateur d'archives** : Navigation transparente et extraction à la volée du contenu des fichiers `.zip`, `.tar` et `.tar.gz`.
- **Prévisualisation intelligente** : Coloration syntaxique temps réel (via Chroma) et vue arborescente pour les dossiers.
- **Éditeur de texte riche** : Mode plein écran, numérotation des lignes, recherche textuelle (`Ctrl+F`), et gestion des fins de ligne.
  !Éditeur
- **Barre latérale des Favoris** : Enregistrez vos dossiers fréquents et accédez-y instantanément via une barre latérale rétractable (`Ctrl+B`).
- **Aide Contextuelle Dynamique** : Appuyez sur `F1` à tout moment pour voir les raccourcis spécifiques au mode actuel.
  !Aide

## Architecture Technique

Le projet repose sur une abstraction puissante du système de fichiers (**VFS**) située dans `internal/vfs/`, permettant d'ajouter facilement de nouveaux protocoles (SFTP, S3, etc.) sans toucher à la logique de l'interface utilisateur.

## Raccourcis Clavier

### Navigation (Explorateur / Visualiseur / Favoris)
| Touche | Action |
| :--- | :--- |
| `Ctrl + G` | Aide contextuelle (adaptée au panneau actif) |
| `Ctrl + T` | Ouvrir le dialogue de connexion FTP |
| `Ctrl + B` | Afficher / Masquer la barre latérale des Favoris |
| `TAB` | Passer au panneau suivant (Favoris → Explorateur → Visualiseur) |
| `Shift + TAB` | Passer au panneau précédent (cycle inverse) |
| `Entrée` | Ouvrir un fichier ou entrer dans un dossier / archive |
| `Ctrl + A` | Ajouter / Retirer le dossier courant des favoris |
| `Ctrl + P` | Recherche Globale (Fuzzy Finder sur tout le disque) |
| `1-9` | Accès rapide direct aux favoris (Home & Racine par défaut) |
| `Ctrl + N` | Renommer le favori sélectionné |
| `Ctrl + F` | Créer un nouveau fichier |
| `Ctrl + D` | Créer un nouveau dossier |
| `Ctrl + O` | Modifier les permissions (Chmod - local uniquement) |
| `Ctrl + R` / `Suppr` | Supprimer l'élément sélectionné dans l'explorateur ou les favoris |
| `Ctrl + E` | Extraire une archive (ou un fichier d'une archive) |
| `Ctrl + K` / `Ctrl + U` | Copier / Coller un élément |
| `Ctrl + X` | Quitter Hollow (demande confirmation) |

### Édition (Éditeur Plein Écran)
| Touche | Action |
| :--- | :--- |
| `Ctrl + G` | Aide contextuelle (Édition) |
| `Ctrl + S` | Sauvegarder les modifications |
| `Ctrl + F` | Rechercher dans le texte (Suivant avec Entrée) |
| `Ctrl + K` | Couper la ligne actuelle (Nano-style, concatène si répété) |
| `Ctrl + U` | Coller le bloc de lignes coupé |
| `Esc` / `Ctrl + X` | Fermer l'éditeur (confirmation si non sauvegardé) |

## Installation & Utilisation

### Prérequis
- `curl` et `wget` (pour l'installation rapide)

### Installation (Utilisateurs)
Pour installer la version native pré-compilée sur Linux (Debian, Ubuntu, Kali, etc.) sans avoir besoin de Go :

```bash
curl -sL https://raw.githubusercontent.com/EducLecomte/go_hollow_project/main/install.sh | bash
```

Ou via le script local si vous avez déjà cloné le projet :
```bash
chmod +x install.sh
./install.sh
```

---
*Dernière mise à jour majeure : Dimanche 19 Avril 2026 - 20:08*