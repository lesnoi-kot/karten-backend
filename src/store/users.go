package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type UsersStore struct {
	db bun.IDB
}

func (s UsersStore) Get(ctx context.Context, id int) (*User, error) {
	user := new(User)

	err := s.db.
		NewSelect().
		Model(user).
		Column("id", "name", "date_created").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s UsersStore) GetBySocialID(ctx context.Context, socialID string) (*User, error) {
	user := new(User)

	err := s.db.
		NewSelect().
		Model(user).
		Where("social_id = ?", socialID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s UsersStore) Add(ctx context.Context, item *User) error {
	_, err := s.db.NewInsert().
		Model(item).
		Column("social_id", "name", "login", "email", "url").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s UsersStore) Update(ctx context.Context, item *User) error {
	result, err := s.db.NewUpdate().
		Model(item).
		Column("name").
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

func (s UsersStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.NewDelete().
		Model((*User)(nil)).
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
