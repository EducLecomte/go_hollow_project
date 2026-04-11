#!/bin/bash

# Script d'installation pour Hollow
PROJECT_DIR="/home/gus/Documents/Dev/go_hollow_project"

echo "Compilation de Hollow..."
cd "$PROJECT_DIR" && go build -o hollow .

echo "Compilation terminée."
echo ""
echo "Pour activer le changement de répertoire automatique à la sortie,"
echo "ajoutez les lignes suivantes à votre ~/.bashrc ou ~/.zshrc :"
echo ""
echo "----------------------------------------------------------------"
echo "hollow() {"
echo "    $PROJECT_DIR/hollow \"\$@\""
echo "    if [ -f /tmp/hollow_cwd_\$USER ]; then"
echo "        cd \"\$(cat /tmp/hollow_cwd_\$USER)\" && rm /tmp/hollow_cwd_\$USER"
echo "    fi"
echo "}"
echo "----------------------------------------------------------------"