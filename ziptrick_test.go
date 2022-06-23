package ziptrick

import (
	"archive/zip"
	"testing"
)

func TestRenameFile(t *testing.T) {
	archive, err := NewArchive("archive.zip")
	if err != nil {
		t.Error(err)
	}
	err = archive.RenameFile("test.txt", "new_test.txt")
	if err != nil {
		t.Error(err)
	}
	err = archive.Write("new_archive.zip")
	if err != nil {
		t.Error(err)
	}
	reader, err := zip.OpenReader("new_archive.zip")
	if err != nil {
		t.Error(err)
	}
	_ = reader.Close()
}
