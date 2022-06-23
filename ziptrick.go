package ziptrick

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

type Archive struct {
	Data []byte
}

func NewArchive(path string) (*Archive, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	archive := new(Archive)
	archive.Data = data
	return archive, nil
}

func (archive *Archive) Write(path string) error {
	err := os.WriteFile(path, archive.Data, 0644)
	return err
}

func (archive *Archive) RenameFile(oldName string, newName string) error {
	if bytes.Count(archive.Data, []byte(oldName)) != 2 {
		return errors.New("ziptrick: filename not found")
	}
	archive.patchSize16([]byte{80, 75, 3, 4}, 26, len(newName))
	archive.patchSize16([]byte{80, 75, 1, 2}, 28, len(newName))
	size := archive.readSize32([]byte{80, 75, 5, 6}, 12)
	archive.patchSize32([]byte{80, 75, 5, 6}, 12, int(size)-len(oldName)+len(newName))
	offset := archive.readSize32([]byte{80, 75, 5, 6}, 16)
	archive.patchSize32([]byte{80, 75, 5, 6}, 16, int(offset)-len(oldName)+len(newName))
	archive.Data = bytes.ReplaceAll(archive.Data, []byte(oldName), []byte(newName))
	return nil
}

func (archive *Archive) readSize32(signature []byte, offset int) uint32 {
	signatureOffset := bytes.Index(archive.Data, signature)
	if signatureOffset >= 0 {
		return binary.LittleEndian.Uint32(archive.Data[signatureOffset+offset : signatureOffset+offset+4])
	}
	return 0
}

func (archive *Archive) patchSize16(signature []byte, offset int, size int) {
	signatureOffset := bytes.Index(archive.Data, signature)
	if signatureOffset >= 0 {
		binary.LittleEndian.PutUint16(archive.Data[signatureOffset+offset:signatureOffset+offset+2], uint16(size))
	}
}

func (archive *Archive) patchSize32(signature []byte, offset int, size int) {
	signatureOffset := bytes.Index(archive.Data, signature)
	if signatureOffset >= 0 {
		binary.LittleEndian.PutUint32(archive.Data[signatureOffset+offset:signatureOffset+offset+4], uint32(size))
	}
}
