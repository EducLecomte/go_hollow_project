#!/bin/bash

# Arrêter le script en cas d'erreur
set -e

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY_NAME="hollow"

echo "--- Installation de Hollow ---"

# Vérification de Go
if ! command -v go &> /dev/null; then
    echo "Erreur : Go n'est pas installé sur votre système."
    exit 1
fi

echo "Compilation dans : $PROJECT_DIR"
cd "$PROJECT_DIR"

go build -o "$BINARY_NAME" ./cmd/hollow

echo "Compilation terminée."
echo ""
echo "Configuration recommandée :"
echo "Pour utiliser Hollow partout et activer le changement de répertoire"
echo "automatique à la sortie, ajoutez ceci à votre ~/.bashrc ou ~/.zshrc :"
echo ""
echo "-------------------------------------------------------------------"
echo "function hollow() {"
echo "    $PROJECT_DIR/$BINARY_NAME \"\$@\""
echo "    local tmp_file=\"/tmp/hollow_cwd_\$USER\""
echo "    if [ -f \"\$tmp_file\" ]; then"
echo "        cd \"\$(cat \"\$tmp_file\")\" && rm \"\$tmp_file\""
echo "    fi"
echo "}"
echo "-------------------------------------------------------------------"
echo ""
echo "Une fois ajouté, relancez votre terminal ou tapez : source ~/.bashrc"
