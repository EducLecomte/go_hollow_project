package main

import (
	"io"
	"os"
	"time"
)

// FileInfo représente un fichier ou dossier de manière universelle (Local, FTP, Zip)
type FileInfo struct {
	Name    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

// VFS est l'interface que devront implémenter tes différents modules
type VFS interface {
	List(path string) ([]FileInfo, error)
	Read(path string) (io.ReadCloser, error)
	Write(path string, data io.Reader) error
	Mkdir(path string) error
}

// LocalFS implémente VFS pour le système de fichiers local
type LocalFS struct{}

func (l *LocalFS) List(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			// Gérer l'erreur ou ignorer le fichier si les infos ne sont pas accessibles
			continue
		}
		files = append(files, FileInfo{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}
	return files, nil
}

func (l *LocalFS) Read(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (l *LocalFS) Write(path string, data io.Reader) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, data)
	return err
}

func (l *LocalFS) Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}
