package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapio"
)

var ErrNotFound = errors.New("resource not found")

type Repo[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	Add(ctx context.Context, item *T) error
	Update(ctx context.Context, item *T) error
	Delete(ctx context.Context, id string) error
}

type Entities struct {
	Users interface {
		Get(ctx context.Context, id UserID) (*User, error)
		GetBySocialID(ctx context.Context, socialID string) (*User, error)
		Add(ctx context.Context, item *User) error
		Update(ctx context.Context, item *User) error
		Delete(ctx context.Context, id UserID) error
	}
	Files interface {
		Get(ctx context.Context, fileID FileID) (*File, error)
		GetImage(ctx context.Context, fileID FileID) (*ImageFile, error)
		Add(ctx context.Context, opts AddFileOptions) (*File, error)
		AddImage(ctx context.Context, opts AddFileOptions) (*ImageFile, error)
		AddImageThumbnail(ctx context.Context, opts AddImageThumbnailOptions) (*File, error)
		GetDefaultCovers(ctx context.Context) ([]ImageFile, error)
		IsDefaultCover(ctx context.Context, fileID FileID) bool
		IsImage(ctx context.Context, fileID FileID) bool
		Delete(ctx context.Context, fileID FileID) error
	}
}

type Store struct {
	Entities
	fileStorage filestorage.FileStorage
	ORM         *bun.DB
}

type TxStore struct {
	Entities
	fileStorage filestorage.FileStorage
	tx          bun.Tx
}

type StoreConfig struct {
	DSN         string
	Logger      *zap.SugaredLogger
	Debug       bool
	FileStorage filestorage.FileStorage
}

func NewStore(cfg StoreConfig) *Store {
	stdDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN)))
	db := bun.NewDB(stdDB, pgdialect.New())

	db.RegisterModel((*ImageThumbnailAssoc)(nil))
	db.RegisterModel((*AttachmentToTaskAssoc)(nil))
	db.RegisterModel((*AttachmentToCommentAssoc)(nil))
	db.RegisterModel((*LabelToTaskAssoc)(nil))

	if cfg.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			// Adapt zap Logger to io.Writer
			bundebug.WithWriter(
				&zapio.Writer{Log: cfg.Logger.Desugar(), Level: zap.DebugLevel},
			),
		))
	}

	store := &Store{
		ORM:         db,
		fileStorage: cfg.FileStorage,
		Entities: Entities{
			Users: UsersStore{db},
			Files: FilesInfoStore{db, cfg.FileStorage},
		},
	}

	return store
}

func (s *Store) Ping() error {
	if err := s.ORM.Ping(); err != nil {
		return fmt.Errorf("DB connection error: %w", err)
	}

	return nil
}

func (s *Store) Close() error {
	return s.ORM.Close()
}

func (s *Store) BeginTx(ctx context.Context) (*TxStore, error) {
	tx, err := s.ORM.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	return newTxStore(tx, s.fileStorage), nil
}

func (s *Store) RunInTx(ctx context.Context, fn func(ctx context.Context, s *TxStore) error) error {
	tx, err := s.ORM.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	txStore := newTxStore(tx, s.fileStorage)

	defer func() {
		if p := recover(); p != nil {
			txStore.Rollback()
			panic(p)
		}
	}()

	if err = fn(ctx, txStore); err != nil {
		if txErr := txStore.Rollback(); txErr != nil {
			return fmt.Errorf("%w; %s", txErr, err)
		}

		return err
	}

	return txStore.Commit()
}

func (s *TxStore) Commit() error {
	return s.tx.Commit()
}

func (s *TxStore) Rollback() error {
	return s.tx.Rollback()
}

func newTxStore(tx bun.Tx, fileStorage filestorage.FileStorage) *TxStore {
	return &TxStore{
		tx: tx,
		Entities: Entities{
			Users: UsersStore{tx},
			Files: FilesInfoStore{tx, fileStorage},
		},
	}
}
