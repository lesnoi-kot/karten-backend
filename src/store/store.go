package store

import (
	"context"
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

type Repo[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	Add(ctx context.Context, item *T) error
	Update(ctx context.Context, item *T) error
	Delete(ctx context.Context, id string) error
}

type Entities struct {
	Projects interface {
		Repo[Project]

		GetAll(ctx context.Context) ([]*Project, error)
	}
	Boards interface {
		Repo[Board]
	}
	TaskLists interface {
		Repo[TaskList]
	}
	Tasks interface {
		Repo[Task]
	}
	Comments interface {
		Repo[Comment]
	}
}

type Store struct {
	Entities
	db *bun.DB
}

type TxStore struct {
	Entities
	tx bun.Tx
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
			// Adapt zap Logger to io.Writer
			bundebug.WithWriter(
				&zapio.Writer{Log: c.Logger.Desugar(), Level: zap.DebugLevel},
			),
		))
	}

	store := &Store{
		db: db,
		Entities: Entities{
			Projects:  ProjectsStore{db},
			Boards:    BoardsStore{db},
			TaskLists: TaskListsStore{db},
			Tasks:     TasksStore{db},
			Comments:  CommentsStore{db},
		},
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) BeginTx(ctx context.Context) (*TxStore, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	return newTxStore(tx), nil
}

func (s *Store) RunInTx(ctx context.Context, fn func(ctx context.Context, s *TxStore) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	txStore := newTxStore(tx)

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

func newTxStore(tx bun.Tx) *TxStore {
	return &TxStore{
		tx: tx,
		Entities: Entities{
			Projects:  ProjectsStore{tx},
			Boards:    BoardsStore{tx},
			TaskLists: TaskListsStore{tx},
			Tasks:     TasksStore{tx},
			Comments:  CommentsStore{tx},
		},
	}
}
