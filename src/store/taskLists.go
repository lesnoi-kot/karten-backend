package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type TaskListsStore struct {
	db *bun.DB
}

func (s TaskListsStore) Get(id string) (*TaskList, error) {
	taskList := new(TaskList)

	err := s.db.
		NewSelect().
		Model(taskList).
		Where("id = ?", id).
		Relation("Tasks").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return taskList, nil
}

func (s TaskListsStore) Add(item *TaskList) error {
	_, err := s.db.
		NewInsert().
		Model(item).
		Column("board_id", "name", "color", "position").
		Returning("*").
		Exec(context.Background())

	return err
}

func (s TaskListsStore) Update(item *TaskList) error {
	result, err := s.db.NewUpdate().
		Model(item).
		Column("name").
		Where("id = ?", item.ID).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	if noRowsAffected(result) {
		return ErrNotFound
	}

	return nil
}

func (s TaskListsStore) Delete(id string) error {
	result, err := s.db.
		NewDelete().
		Model((*TaskList)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	if noRowsAffected(result) {
		return ErrNotFound
	}

	return nil
}
