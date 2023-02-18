package filestorage

import (
	"io"

	"github.com/google/uuid"
)

type FileID = string

type FileStorage interface {
	Get(key FileID) ([]byte, error)
	Set(key FileID, data io.Reader) (int64, error)
	Add(data io.Reader) (FileID, int64, error)
}

func RandomID() FileID {
	uuid4, _ := uuid.NewRandom()
	return uuid4.String()
}
