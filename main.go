package main

import (
	"example.com/sync-dir/decl"
	"example.com/sync-dir/impl"
	"flag"
	"fmt"
)

func syncDirs(src decl.Source, dst decl.Destination) {
	for srcFileHash, srcFileEntry := range src.GetFileEntryMap() {
		dstFileEntry, exists := dst.GetFileEntryByHash(srcFileHash)
		if !exists {
			dst.CopyFile(srcFileEntry)
		} else if srcFileEntry.Name() != dstFileEntry.Name() {
			dst.RenameFile(dstFileEntry, srcFileEntry)
		}
		delete(dst.GetFileEntryMap(), srcFileHash)
	}

	for _, dstFileEntry := range dst.GetFileEntryMap() {
		dst.DeleteFile(dstFileEntry)
	}
}

func main() {
	srcPath := flag.String("src", "", "path to source directory")
	dstPath := flag.String("dst", "", "path to destination directory")
	flag.Parse()
	src := impl.NewSource(*srcPath)
	dst := impl.NewDestination(*dstPath)
	fmt.Println("syncing: ", *srcPath, *dstPath)
	syncDirs(src, dst)
}