package vfs

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

// FileInfo représente un fichier ou dossier de manière universelle (Local, FTP, Zip)
type FileInfo struct {
	Name        string
	IsDir       bool
	Size        int64
	ModTime     time.Time
	Permissions string
	Owner       string
}

// VFS est l'interface que devront implémenter tes différents modules
type VFS interface {
	List(path string) ([]FileInfo, error)
	Read(path string) (io.ReadCloser, error)
	Write(path string, data io.Reader) error
	Mkdir(path string) error
	Copy(src, dst string) error
	Remove(path string) error
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

		// Extraction des infos Unix (Propriétaire)
		owner := "unknown"
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			u, err := user.LookupId(fmt.Sprint(stat.Uid))
			if err == nil {
				owner = u.Username
			}
		}

		files = append(files, FileInfo{
			Name:        entry.Name(),
			IsDir:       entry.IsDir(),
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			Permissions: info.Mode().String(),
			Owner:       owner,
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

func (l *LocalFS) Copy(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if err := os.MkdirAll(dst, info.Mode()); err != nil {
			return err
		}
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := l.Copy(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
				return err
			}
		}
		return nil
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (l *LocalFS) Remove(path string) error {
	return os.RemoveAll(path)
}
