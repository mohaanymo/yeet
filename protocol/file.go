package protocol

import (
	"fmt"
	"os"
)

type Metadata struct {
	Filename string
	FilenameLen int
	FileSize 	int64
}

type File struct {
	Name string
	Path	string
	Metadata Metadata
	Reader *os.File
}

// NewFile check existence of the file then open it and parse it's metadata
func NewFile(name, filePath string) (*File, error) {
	ok, err := IsPathExist(filePath)
	if !ok{
		return nil, fmt.Errorf("error while checking the file path: %v", err)
	}

	// Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }

    // Get file info
    fileInfo, err := file.Stat()
    if err != nil {
        return nil, fmt.Errorf("failed to get file info: %w", err)
    }
    filename := fileInfo.Name()
    fileSize := fileInfo.Size()
	
	// Creating metadata
	metadata := Metadata{
		Filename: filename,
		FilenameLen: len(filename),
		FileSize: fileSize,
	}

	return &File{
		Name: name,
		Path: filePath,
		Metadata: metadata,
		Reader: file,
	}, nil
}
