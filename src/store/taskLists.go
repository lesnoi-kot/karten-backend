package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type TaskListsStore struct {
	db bun.IDB
}

func (s TaskListsStore) Get(ctx context.Context, id string) (*TaskList, error) {
	taskList := new(TaskList)

	err := s.db.
		NewSelect().
		Model(taskList).
		Where("id = ?", id).
		Relation("Tasks").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return taskList, nil
}

func (s TaskListsStore) Add(ctx context.Context, item *TaskList) error {
	_, err := s.db.
		NewInsert().
		Model(item).
		Column("board_id", "name", "color", "position").
		Returning("*").
		Exec(ctx)

	return err
}

func (s TaskListsStore) Update(ctx context.Context, item *TaskList) error {
	result, err := s.db.NewUpdate().
		Model(item).
		Column("name", "position").
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

func (s TaskListsStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.
		NewDelete().
		Model((*TaskList)(nil)).
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
