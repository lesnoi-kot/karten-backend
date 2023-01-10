package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type CommentsStore struct {
	db bun.IDB
}

func (s CommentsStore) Get(ctx context.Context, id string) (*Comment, error) {
	comment := new(Comment)

	err := s.db.
		NewSelect().
		Model(comment).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return comment, nil
}

func (s CommentsStore) Add(ctx context.Context, item *Comment) error {
	_, err := s.db.NewInsert().
		Model(item).
		Column("task_id", "text", "author").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s CommentsStore) Update(ctx context.Context, item *Comment) error {
	result, err := s.db.NewUpdate().
		Model(item).
		Column("text").
		Where("id = ?", item.ID).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	}

	if noRowsAffected(result) {
		return ErrNotFound
	}

	return nil
}

func (s CommentsStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.NewDelete().
		Model((*Comment)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	if noRowsAffected(result) {
		return ErrNotFound
	}

	return nil
}
