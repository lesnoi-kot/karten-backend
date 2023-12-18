package filestorage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileSystemStorage struct {
	RootPath string // Absolute path

	GenerateID func() FileID
}

func (s FileSystemStorage) Add(data io.Reader) (FileID, int64, error) {
	fileID := s.GenerateID()
	size, err := s.Set(fileID, data)
	return fileID, size, err
}

func (s FileSystemStorage) Get(key FileID) ([]byte, error) {
	path := filepath.Join(s.RootPath, key)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s FileSystemStorage) Set(key FileID, data io.Reader) (int64, error) {
	path := filepath.Join(s.RootPath, key)

	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, data)
}

func NewFileSystemStorage(rootPath string) (*FileSystemStorage, error) {
	if !filepath.IsAbs(rootPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		rootPath = filepath.Join(cwd, rootPath)
	}

	f, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("FileSystemStorage requires directory path: %w", err)
	}

	if !f.IsDir() {
		return nil, errors.New("FileSystemStorage requires directory path")
	}

	service := &FileSystemStorage{
		RootPath:   rootPath,
		GenerateID: RandomID,
	}

	return service, nil
}
