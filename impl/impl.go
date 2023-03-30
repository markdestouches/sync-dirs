package impl

import (
	"crypto/sha256"
	"example.com/sync-dir/decl"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileEntry struct {
	path  string
	isDir bool
}

type Directory struct {
	path        string
	fileEntries decl.FileEntryMap
}

type Source struct {
	dir *Directory
}

type Destination struct {
	dir *Directory
}

func (fe FileEntry) Name() string {
	return filepath.Base(fe.path)
}

func (fe FileEntry) Path() string {
	return fe.path
}

func (fe FileEntry) IsDir() bool {
	return fe.isDir
}

func NewDirectory(path string) *Directory {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fileEntries := make(decl.FileEntryMap)

	for _, entry := range dirEntries {
		fileEntry := FileEntry{path: filepath.Join(path, entry.Name()), isDir: entry.Type().IsDir()}
		hash := hashFile(fileEntry)
		fileEntries[hash] = fileEntry
	}

	return &Directory{path: path, fileEntries: fileEntries}
}

func (d *Directory) Path() string {
	return d.path
}

func (d *Directory) GetFileEntryMap() decl.FileEntryMap {
	return d.fileEntries
}

func (d *Directory) GetFileEntryByHash(fileContentHash decl.FileContentHash) (decl.FileEntry, bool) {
	fileEntry, exists := d.fileEntries[fileContentHash]
	return fileEntry, exists
}

func NewSource(path string) Source {
	return Source{dir: NewDirectory(path)}
}

func (s Source) Path() string {
	return s.dir.Path()
}

func (s Source) GetFileEntryMap() decl.FileEntryMap {
	return s.dir.GetFileEntryMap()
}

func (s Source) GetFileEntryByHash(fileContentHash decl.FileContentHash) (decl.FileEntry, bool) {
	fileEntry, exists := s.dir.GetFileEntryByHash(fileContentHash)
	return fileEntry, exists
}

func NewDestination(path string) Destination {
	return Destination{dir: NewDirectory(path)}
}

func (d Destination) Path() string {
	return d.dir.Path()
}

func (d Destination) GetFileEntryMap() decl.FileEntryMap {
	return d.dir.GetFileEntryMap()
}

func (d Destination) GetFileEntryByHash(fileContentHash decl.FileContentHash) (decl.FileEntry, bool) {
	fileEntry, exists := d.dir.GetFileEntryByHash(fileContentHash)
	return fileEntry, exists
}

func (dst Destination) CopyFile(srcFileEntry decl.FileEntry) {
	if srcFileEntry.IsDir() {
		log.Fatal("handling nested directories is not implemented")
	}

	srcFile, err := os.OpenFile(srcFileEntry.Path(), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	dstFilePath := filepath.Join(dst.Path(), srcFileEntry.Name())
	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
}

func (dst Destination) RenameFile(dstFileEntry decl.FileEntry, srcFileEntry decl.FileEntry) {
	oldFilePath := dstFileEntry.Path()
	newFilePath := filepath.Join(dst.dir.path, srcFileEntry.Name())
	err := os.Rename(oldFilePath, newFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func (dst Destination) DeleteFile(dstFileEntry decl.FileEntry) {
	err := os.RemoveAll(dstFileEntry.Path())
	if err != nil {
		log.Fatal(err)
	}
}

func hashFile(fileEntry FileEntry) [32]byte {
	if fileEntry.isDir {
		log.Fatal("hashing directories is not implemented")
	}

	file, err := os.OpenFile(fileEntry.path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	var hashArray [32]byte
	// hopefully this gets compiled away
	copy(hashArray[:], hash.Sum(nil))
	return hashArray
}
