package vfs

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

// FtpFS implémente VFS pour le protocole FTP
type FtpFS struct {
	conn *ftp.ServerConn
	host string
	port int
}

// NewFtpFS crée une nouvelle connexion FTP et retourne une instance de FtpFS
func NewFtpFS(host string, port int, user, password string) (*FtpFS, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion: %v", err)
	}

	err = conn.Login(user, password)
	if err != nil {
		_ = conn.Quit()
		return nil, fmt.Errorf("erreur d'authentification: %v", err)
	}

	return &FtpFS{
		conn: conn,
		host: host,
		port: port,
	}, nil
}

func (f *FtpFS) List(path string) ([]FileInfo, error) {
	entries, err := f.conn.List(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		// Ignorer "." et ".." si présents
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		isDir := entry.Type == ftp.EntryTypeFolder
		files = append(files, FileInfo{
			Name:        entry.Name,
			IsDir:       isDir,
			Size:        int64(entry.Size),
			ModTime:     entry.Time,
			Permissions: "rwxr-xr-x", // FTP ne donne pas toujours les droits, valeur par défaut
			Owner:       "ftp",
		})
	}
	return files, nil
}

func (f *FtpFS) Read(path string) (io.ReadCloser, error) {
	resp, err := f.conn.Retr(path)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FtpFS) Write(path string, data io.Reader) error {
	return f.conn.Stor(path, data)
}

func (f *FtpFS) Mkdir(path string) error {
	return f.conn.MakeDir(path)
}

func (f *FtpFS) Copy(src, dst string) error {
	// Le FTP standard ne supporte pas la copie directe. 
	// On utilise la fonction générique CopyRecursiveBetweenVFS définie dans vfs.go
	// qui gère le transfert via un buffer local.
	return fmt.Errorf("la copie directe n'est pas supportée en FTP, utilisez CopyRecursiveBetweenVFS")
}

func (f *FtpFS) Remove(path string) error {
	// On essaie d'abord de supprimer comme un fichier
	err := f.conn.Delete(path)
	if err != nil {
		// Si ça échoue, on essaie comme un dossier (récursif)
		err = f.conn.RemoveDirRecur(path)
	}
	return err
}

func (f *FtpFS) Stat(path string) (FileInfo, error) {
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	
	entries, err := f.conn.List(dir)
	if err != nil {
		return FileInfo{}, err
	}

	for _, entry := range entries {
		if entry.Name == name {
			isDir := entry.Type == ftp.EntryTypeFolder
			return FileInfo{
				Name:        entry.Name,
				IsDir:       isDir,
				Size:        int64(entry.Size),
				ModTime:     entry.Time,
				Permissions: "rwxr-xr-x",
				Owner:       "ftp",
			}, nil
		}
	}

	return FileInfo{}, fmt.Errorf("fichier non trouvé")
}

func (f *FtpFS) Close() error {
	return f.conn.Quit()
}
