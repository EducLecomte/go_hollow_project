package vfs

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ArchiveNode struct {
	Info     FileInfo
	Children map[string]*ArchiveNode
	ZipFile  *zip.File // nil si c'est un tar
}

type ArchiveFS struct {
	ArchivePath string
	Root        *ArchiveNode
	ZipReader   *zip.ReadCloser
	IsTar       bool
	IsGzip      bool
}

func NewArchiveFS(ctx context.Context, path string) (*ArchiveFS, error) {
	fs := &ArchiveFS{
		ArchivePath: path,
		Root: &ArchiveNode{
			Info: FileInfo{Name: "/", IsDir: true},
			Children: make(map[string]*ArchiveNode),
		},
	}

	ext := strings.ToLower(filepath.Ext(path))
	isGz := ext == ".gz" || ext == ".tgz"
	isTar := ext == ".tar" || isGz

	if isTar {
		fs.IsTar = true
		fs.IsGzip = isGz
		err := fs.scanTar(ctx)
		if err != nil {
			return nil, err
		}
	} else if ext == ".zip" {
		err := fs.scanZip(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("format d'archive non supporté")
	}

	return fs, nil
}

func (a *ArchiveFS) addNode(path string, info FileInfo, zipFile *zip.File) {
	// nettoyage path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "./")
	if path == "" {
		return
	}

	parts := strings.Split(path, "/")
	current := a.Root

	for i, part := range parts {
		if part == "" {
			continue
		}
		
		isLast := i == len(parts)-1
		
		if _, exists := current.Children[part]; !exists {
			// Si c'est un dossier intermédiaire (qui n'était pas explicite) ou si c'est le dernier élément
			isDir := true
			if isLast && !info.IsDir {
				isDir = false
			}

			nodeInfo := FileInfo{
				Name: part,
				IsDir: isDir,
				Permissions: "r--r--r--", // Read Only
				Owner: "archive",
			}
			
			if isLast {
				nodeInfo.Size = info.Size
				nodeInfo.ModTime = info.ModTime
				if info.Permissions != "" {
				    nodeInfo.Permissions = info.Permissions
				}
				nodeInfo.Owner = info.Owner
			}

			current.Children[part] = &ArchiveNode{
				Info: nodeInfo,
				Children: make(map[string]*ArchiveNode),
				ZipFile: nil,
			}
			
			if isLast && zipFile != nil {
				current.Children[part].ZipFile = zipFile
			}
		}
		current = current.Children[part]
	}
}

func (a *ArchiveFS) scanZip(ctx context.Context) error {
	r, err := zip.OpenReader(a.ArchivePath)
	if err != nil {
		return err
	}
	a.ZipReader = r
	for _, f := range r.File {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		a.addNode(f.Name, FileInfo{
			Name: filepath.Base(f.Name),
			IsDir: f.FileInfo().IsDir(),
			Size: f.FileInfo().Size(),
			ModTime: f.FileInfo().ModTime(),
		}, f)
	}
	return nil
}

func (a *ArchiveFS) scanTar(ctx context.Context) error {
	f, err := os.Open(a.ArchivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var tr *tar.Reader
	if a.IsGzip {
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gzr.Close()
		tr = tar.NewReader(gzr)
	} else {
		tr = tar.NewReader(f)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		a.addNode(hdr.Name, FileInfo{
			Name: filepath.Base(hdr.Name),
			IsDir: hdr.FileInfo().IsDir(),
			Size: hdr.Size,
			ModTime: hdr.ModTime,
			Permissions: hdr.FileInfo().Mode().String(),
			Owner: hdr.Uname,
		}, nil)
	}
	return nil
}

func (a *ArchiveFS) getNode(path string) *ArchiveNode {
	path = strings.TrimPrefix(path, "/")
	if path == "" || path == "." {
		return a.Root
	}
	parts := strings.Split(path, "/")
	current := a.Root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if next, ok := current.Children[part]; ok {
			current = next
		} else {
			return nil
		}
	}
	return current
}

func (a *ArchiveFS) List(path string) ([]FileInfo, error) {
	node := a.getNode(path)
	if node == nil || !node.Info.IsDir {
		return nil, fmt.Errorf("dossier non trouvé")
	}

	var files []FileInfo
	for _, child := range node.Children {
		files = append(files, child.Info)
	}
	return files, nil
}

type tarReadCloser struct {
	f   *os.File
	gzr *gzip.Reader
	tr  *tar.Reader
}

func (t *tarReadCloser) Read(p []byte) (n int, err error) {
	return t.tr.Read(p)
}

func (t *tarReadCloser) Close() error {
	if t.gzr != nil {
		t.gzr.Close()
	}
	return t.f.Close()
}

func (a *ArchiveFS) Read(path string) (io.ReadCloser, error) {
	if a.ZipReader != nil {
		node := a.getNode(path)
		if node != nil && node.ZipFile != nil {
			return node.ZipFile.Open()
		}
		return nil, fmt.Errorf("fichier non trouvé dans le zip")
	}

	if a.IsTar {
		f, err := os.Open(a.ArchivePath)
		if err != nil {
			return nil, err
		}

		var gzr *gzip.Reader
		var tr *tar.Reader

		if a.IsGzip {
			gzr, err = gzip.NewReader(f)
			if err != nil {
				f.Close()
				return nil, err
			}
			tr = tar.NewReader(gzr)
		} else {
			tr = tar.NewReader(f)
		}

		cleanPath := strings.TrimPrefix(path, "/")
		cleanPath = strings.TrimPrefix(cleanPath, "./")
		
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			hdrClean := strings.TrimPrefix(hdr.Name, "./")
			hdrClean = strings.TrimPrefix(hdrClean, "/")
			if hdrClean == cleanPath {
				return &tarReadCloser{f: f, gzr: gzr, tr: tr}, nil
			}
		}

		if gzr != nil {
			gzr.Close()
		}
		f.Close()
		return nil, fmt.Errorf("fichier introuvable dans le tar")
	}

	return nil, fmt.Errorf("lecture impossible")
}

func (a *ArchiveFS) Write(path string, data io.Reader) error {
	return fmt.Errorf("les archives sont montées en lecture seule")
}

func (a *ArchiveFS) Mkdir(path string) error {
	return fmt.Errorf("les archives sont montées en lecture seule")
}

func (a *ArchiveFS) Copy(src, dst string) error {
	return fmt.Errorf("les archives sont montées en lecture seule")
}

func (a *ArchiveFS) Remove(path string) error {
	return fmt.Errorf("les archives sont montées en lecture seule")
}
