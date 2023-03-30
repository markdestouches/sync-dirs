package decl

type FileContentHash interface{}

type FileEntry interface {
	Name() string
	IsDir() bool
	Path() string
}

type FileEntryMap map[FileContentHash]FileEntry

type Directory interface {
	GetFileEntryMap() FileEntryMap
	GetFileEntryByHash(FileContentHash) (FileEntry, bool)
}

type Source interface {
	Directory
}

type Destination interface {
	Directory
	CopyFile(FileEntry)
	RenameFile(FileEntry, FileEntry)
	DeleteFile(FileEntry)
}
