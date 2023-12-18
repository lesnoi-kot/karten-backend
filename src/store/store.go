package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapio"
)

var ErrNotFound = errors.New("resource not found")

type Store struct {
	ORM *bun.DB
}

type StoreConfig struct {
	DSN    string
	Logger *zap.SugaredLogger
	Debug  bool
}

func NewStore(cfg StoreConfig) *Store {
	stdDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN)))
	db := bun.NewDB(stdDB, pgdialect.New())

	db.RegisterModel((*User)(nil))
	db.RegisterModel((*Project)(nil))

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

	store := &Store{db}
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
