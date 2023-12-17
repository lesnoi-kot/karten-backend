package entityservices

import (
	"context"
	"errors"
	"io"
	"path"
	"strings"

	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

const DefaultCoverLimit = 4

type IFileService interface {
	Get(ctx context.Context, fileID store.FileID) (*store.File, error)
	GetImage(ctx context.Context, fileID store.FileID) (*store.ImageFile, error)
	Add(ctx context.Context, opts AddFileOptions) (*store.File, error)
	AddImage(ctx context.Context, opts AddFileOptions) (*store.ImageFile, error)
	AddImageThumbnail(ctx context.Context, opts AddImageThumbnailOptions) (*store.File, error)
	GetDefaultCovers(ctx context.Context) ([]store.ImageFile, error)
	IsDefaultCover(ctx context.Context, fileID store.FileID) bool
	IsImage(ctx context.Context, fileID store.FileID) bool
	Delete(ctx context.Context, fileID store.FileID) error
}

type FileServiceRequirements interface {
	StoreInjector
}

type FileService struct {
	FileServiceRequirements
	FileStorage filestorage.FileStorage
}

type AddFileOptions struct {
	Name     string
	MIMEType string
	Data     io.Reader
}

type AddImageThumbnailOptions struct {
	AddFileOptions
	OriginalImageFileID string
}

func (s FileService) Add(ctx context.Context, opts AddFileOptions) (*store.File, error) {
	ext := path.Ext(opts.Name)
	storageID := filestorage.RandomID() + ext
	bytesCount, err := s.FileStorage.Set(storageID, opts.Data)
	if err != nil {
		return nil, err
	}

	file := &store.File{
		StorageObjectID: storageID,
		Name:            opts.Name,
		MimeType:        opts.MIMEType,
		Size:            int(bytesCount),
	}

	_, err = s.GetStore().ORM.NewInsert().
		Model(file).
		Column("storage_object_id", "name", "size", "mime_type").
		Returning("id").
		Exec(ctx)

	return file, err
}

func (s FileService) AddImageThumbnail(ctx context.Context, opts AddImageThumbnailOptions) (*store.File, error) {
	thumbnail, err := s.Add(ctx, opts.AddFileOptions)
	if err != nil {
		return nil, err
	}

	link := &store.ImageThumbnailAssoc{
		ID:      thumbnail.ID,
		ImageID: opts.OriginalImageFileID,
	}

	_, err = s.GetStore().ORM.NewInsert().Model(link).Exec(ctx)
	return thumbnail, err
}

func (s FileService) GetDefaultCovers(ctx context.Context) ([]store.ImageFile, error) {
	var covers []store.ImageFile
	subquery := s.GetStore().ORM.NewSelect().
		Model((*store.CoverImageToFileAssoc)(nil)).
		Limit(DefaultCoverLimit)

	err := s.GetStore().ORM.NewSelect().
		Model(&covers).
		Relation("Thumbnails").
		Where("id in (?)", subquery).
		Scan(ctx)

	return covers, err
}

func (s FileService) IsDefaultCover(ctx context.Context, fileID store.FileID) bool {
	exists, err := s.GetStore().ORM.NewSelect().
		Model((*store.CoverImageToFileAssoc)(nil)).
		Where("id = ?", fileID).
		Exists(ctx)
	if err != nil {
		return false
	}

	return exists
}

func (s FileService) IsImage(ctx context.Context, fileID store.FileID) bool {
	exists, err := s.GetStore().ORM.NewSelect().
		Model((*store.File)(nil)).
		Where("id = ?", fileID).
		Where("mime_type LIKE ?", "image/%").
		Exists(ctx)
	if err != nil {
		return false
	}

	return exists
}

func (s FileService) Get(ctx context.Context, fileID store.FileID) (*store.File, error) {
	file := new(store.File)
	err := s.GetStore().ORM.NewSelect().
		Model(file).
		Where("id = ?", fileID).
		Scan(ctx)

	return file, err
}

func (s FileService) GetImage(ctx context.Context, fileID store.FileID) (*store.ImageFile, error) {
	img := new(store.ImageFile)
	err := s.GetStore().ORM.NewSelect().
		Model(img).
		Relation("Thumbnails").
		Where("id = ?", fileID).
		Scan(ctx)

	if !img.IsImage() {
		return nil, store.ErrNotFound
	}

	return img, err
}

func (s FileService) AddImage(ctx context.Context, opts AddFileOptions) (*store.ImageFile, error) {
	if !strings.HasPrefix(opts.MIMEType, "image/") {
		return nil, errors.New("Input file is not an image")
	}

	file, err := s.Add(ctx, opts)
	if err != nil {
		return nil, err
	}

	image := &store.ImageFile{
		File:       *file,
		Thumbnails: make([]*store.File, 0),
	}

	return image, err
}

func (s FileService) Delete(ctx context.Context, fileID store.FileID) error {
	result, err := s.GetStore().ORM.NewDelete().
		Model((*store.File)(nil)).
		Where("id = ?", fileID).
		Exec(ctx)
	if err != nil {
		return err
	}

	if store.NoRowsAffected(result) {
		return store.ErrNotFound
	}

	return nil
}
