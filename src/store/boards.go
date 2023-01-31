package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type BoardsStore struct {
	db bun.IDB
}

func (s BoardsStore) Get(ctx context.Context, id string) (*Board, error) {
	updateResult, err := s.db.
		NewUpdate().
		Model((*Board)(nil)).
		Where("id = ?", id).
		Set("date_last_viewed = ?", time.Now().UTC()).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	if noRowsAffected(updateResult) {
		return nil, ErrNotFound
	}

	board := new(Board)

	err = s.db.
		NewSelect().
		Model(board).
		Where("id = ?", id).
		Relation("TaskLists").
		Relation("TaskLists.Tasks").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return board, nil
}

func (s BoardsStore) Filter(ctx context.Context, id string) (*Board, error) {
	board := new(Board)

	err := s.db.
		NewSelect().
		Model(board).
		Where("id = ?", id).
		Relation("TaskLists").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return board, nil
}

func (s BoardsStore) Add(ctx context.Context, board *Board) error {
	_, err := s.db.
		NewInsert().
		Model(board).
		Column("project_id", "name", "color", "cover_id").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s BoardsStore) Update(ctx context.Context, board *Board) error {
	result, err := s.db.NewUpdate().
		Model(board).
		Column("name", "archived", "color", "cover_id").
		Where("id = ?", board.ID).
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

func (s BoardsStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.
		NewDelete().
		Model((*Board)(nil)).
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
