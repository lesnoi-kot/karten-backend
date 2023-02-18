package store

import (
	"context"
	"io"
	"path"

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
	ext := path.Ext(opts.Name)
	storageID := filestorage.RandomID() + ext
	bytesCount, err := s.fileStorage.Set(storageID, opts.Data)
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

func (s FilesInfoStore) GetDefaultCovers(ctx context.Context) ([]ImageFile, error) {
	var covers []ImageFile
	subquery := s.db.NewSelect().Model((*CoverImageToFileAssoc)(nil)).Limit(4)

	err := s.db.NewSelect().
		Model(&covers).
		Relation("Thumbnails").
		Where("id in (?)", subquery).
		Scan(ctx)

	return covers, err
}

func (s FilesInfoStore) IsDefaultCover(ctx context.Context, fileID string) bool {
	exists, err := s.db.NewSelect().
		Model((*CoverImageToFileAssoc)(nil)).
		Where("id = ?", fileID).
		Exists(ctx)
	if err != nil {
		return false
	}

	return exists
}
