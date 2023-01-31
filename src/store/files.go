package store

import (
	"context"
	"io"

	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/uptrace/bun"
)

func init() {
}

type FilesInfoStore struct {
	db          bun.IDB
	fileStorage filestorage.FileStorage
}

type AddFileOptions struct {
	Name     string
	MIMEType string
	Data     io.Reader
}

func (s FilesInfoStore) Add(ctx context.Context, opts AddFileOptions) (*File, error) {
	storageID, bytesCount, err := s.fileStorage.Add(opts.Data)
	if err != nil {
		return nil, err
	}

	file := &File{
		StorageObjectID: storageID,
		Name:            opts.Name,
		MimeType:        opts.MIMEType,
		Size:            int(bytesCount),
	}

	_, err = s.db.NewInsert().
		Model(file).
		Column("storage_object_id", "name", "size", "mime_type").
		Returning("id").
		Exec(ctx)

	return file, err
}

type AddImageThumbnailOptions struct {
	AddFileOptions
	OriginalImageFileID string
}

func (s FilesInfoStore) AddImageThumbnail(ctx context.Context, opts AddImageThumbnailOptions) (*File, error) {
	thumbnail, err := s.Add(ctx, opts.AddFileOptions)
	if err != nil {
		return nil, err
	}

	link := &ImageThumbnailAssoc{
		ID:      thumbnail.ID,
		ImageID: opts.OriginalImageFileID,
	}

	_, err = s.db.NewInsert().Model(link).Exec(ctx)
	return thumbnail, err
}
