package main

import (
	"example.com/sync-dir/impl"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestSyncDirsCopyFromSource(t *testing.T) {
	testDirPath := t.TempDir()
	sourceDirName := "mock_source"
	destinationDirName := "mock_destination"
	sourceDirPath := filepath.Join(testDirPath, sourceDirName)
	destinationDirPath := filepath.Join(testDirPath, destinationDirName)

	err := os.Mkdir(sourceDirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(destinationDirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	sourceFileName := "source_file_name"
	file, err := os.OpenFile(filepath.Join(sourceDirPath, sourceFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	sourceFileContent := "source file content"
	_, err = file.WriteString(sourceFileContent)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	src := impl.NewSource(sourceDirPath)
	dst := impl.NewDestination(destinationDirPath)

	syncDirs(src, dst)
	t.Log("------")

	dirEntries, err := os.ReadDir(destinationDirPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(dirEntries) != 1 {
		t.Fatal("destination dir contains wrong number of files")
	}
	if dirEntries[0].Name() != sourceFileName {
		t.Fatal("copied file name does not match source file name")
	}

	destinationFileContent, err := os.ReadFile(filepath.Join(destinationDirPath, sourceFileName))
	if err != nil {
		t.Fatal(err)
	}

	if string(destinationFileContent) != sourceFileContent {
		t.Fatal("destination file content does not match source file content")
	}
}

func TestSyncDirsRenameFileInDestination(t *testing.T) {
	testDirPath := t.TempDir()
	sourceDirName := "mock_source"
	destinationDirName := "mock_destination"
	sourceDirPath := filepath.Join(testDirPath, sourceDirName)
	destinationDirPath := filepath.Join(testDirPath, destinationDirName)

	err := os.Mkdir(sourceDirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(destinationDirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	sourceFileName := "source_file_name"
	srcFile, err := os.OpenFile(filepath.Join(sourceDirPath, sourceFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	sourceFileContent := "source file content"
	_, err = srcFile.WriteString(sourceFileContent)
	if err != nil {
		t.Fatal(err)
	}
	srcFile.Close()

	srcFile, err = os.OpenFile(filepath.Join(sourceDirPath, sourceFileName), os.O_APPEND|os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	destinationFileName := "destination_file_name"
	dstFile, err := os.OpenFile(filepath.Join(destinationDirPath, destinationFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	io.Copy(dstFile, srcFile)
	srcFile.Close()
	dstFile.Close()

	src := impl.NewSource(sourceDirPath)
	dst := impl.NewDestination(destinationDirPath)

	syncDirs(src, dst)
	destinationFileContent, err := os.ReadFile(filepath.Join(destinationDirPath, sourceFileName))
	if sourceFileContent != string(destinationFileContent) {
		t.Fatalf("destinations file content <%s> does not match source file content <%s>", destinationFileContent, sourceFileContent)
	}
}

func TestSyncDirsDeleteFromDestination(t *testing.T) {
	testDirPath := t.TempDir()
	sourceDirName := "mock_source"
	destinationDirName := "mock_destination"
	sourceDirPath := filepath.Join(testDirPath, sourceDirName)
	destinationDirPath := filepath.Join(testDirPath, destinationDirName)

	err := os.Mkdir(sourceDirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(destinationDirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	destinationFileName := "destination_file_name"
	dstFile, err := os.OpenFile(filepath.Join(destinationDirPath, destinationFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	destinationFileContent := "destination file content"
	_, err = dstFile.WriteString(destinationFileContent)
	if err != nil {
		t.Fatal(err)
	}
	dstFile.Close()

	src := impl.NewSource(sourceDirPath)
	dst := impl.NewDestination(destinationDirPath)
	syncDirs(src, dst)

	destinationFilePath := filepath.Join(destinationDirPath, destinationFileName)
	dstFile, err = os.OpenFile(destinationFilePath, os.O_RDONLY, os.ModePerm)
	expectedErrMsg := "open " + destinationFilePath + ": no such file or directory"
	if err.Error() != expectedErrMsg {
		t.Fatal("Destination file must not exist")
	}
}
