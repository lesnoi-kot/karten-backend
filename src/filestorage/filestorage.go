package filestorage

import "io"

type FileID = string

type FileStorage interface {
	Get(key FileID) ([]byte, error)
	Set(key FileID, data io.Reader) (int64, error)
	Add(data io.Reader) (FileID, int64, error)
}
