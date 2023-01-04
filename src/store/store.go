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

var ErrNotFound = errors.New("item not found")

type EditProjectArgs struct {
	ID   string
	Name string
}

type AddBoardArgs struct {
	Name string
}

type Repo[T any] interface {
	Get(id string) (*T, error)
	GetAll() ([]*T, error)
	Add(item *T) (*T, error)
	Edit(item *T) (*T, error)
	Delete(id string) error
}

type Store struct {
	*bun.DB

	Projects interface {
		Get(id string) (*Project, error)
		GetAll() ([]*Project, error)
		Add(name string) (*Project, error)
		Edit(args EditProjectArgs) (*Project, error)
		Delete(id string) error
	}
	Boards interface {
		Get(id string) (*Board, error)
		GetAll() ([]*Board, error)
		Add(args AddBoardArgs) (*Board, error)
		Edit(args *Board) (*Board, error)
		Delete(id string) error
	}
}

type StoreConfig struct {
	DSN    string
	Logger *zap.SugaredLogger
	Debug  bool
}

func NewStore(c StoreConfig) (*Store, error) {
	stdDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.DSN)))
	db := bun.NewDB(stdDB, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}

	if c.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.WithWriter(
				&zapio.Writer{Log: c.Logger.Desugar(), Level: zap.DebugLevel},
			),
		))
	}

	store := &Store{
		DB:       db,
		Projects: &ProjectsStore{db},
		Boards:   nil,
	}

	return store, nil
}
