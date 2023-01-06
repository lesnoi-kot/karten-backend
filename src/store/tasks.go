package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type TasksStore struct {
	db *bun.DB
}

func (s *TasksStore) Get(id string) (*Task, error) {
	task := new(Task)

	err := s.db.
		NewSelect().
		Model(task).
		Where("id = ?", id).
		Relation("Comments").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return task, nil
}

func (s *TasksStore) Add(task *Task) error {
	_, err := s.db.
		NewInsert().
		Model(task).
		Column("task_list_id", "name", "text", "position", "due_date").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *TasksStore) Update(task *Task) error {
	result, err := s.db.NewUpdate().
		Model(task).
		Column("task_list_id", "name", "text", "position", "due_date").
		Where("id = ?", task.ID).
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

func (s *TasksStore) Delete(id string) error {
	result, err := s.db.
		NewDelete().
		Model((*Task)(nil)).
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
