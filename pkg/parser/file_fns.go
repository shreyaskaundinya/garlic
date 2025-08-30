package parser

import (
	"os"
	"path"

	"github.com/shreyaskaundinya/garlic/pkg/utils"
)

func NewFile(path string, fileType string) *File {
	f := &File{
		Type: fileType,
		Path: path,
		Body: make([]byte, 0),
	}

	return f
}

// Read the file and populate the Body field
func (f *File) ReadFile() error {
	b, err := os.ReadFile(f.Path)

	if err != nil {
		return err
	}

	f.Body = b
	return nil
}

func (f *File) WriteToDest(DestFolder string, FileName string, body []byte) error {
	log := utils.NewLogger()

	p := path.Join(DestFolder, FileName)
	log.Infow("Writing to : ", "path", p)

	return os.WriteFile(p, body, 0644)
}
