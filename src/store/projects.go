package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type ProjectsStore struct {
	db *bun.DB
}

func (s ProjectsStore) Get(id string) (*Project, error) {
	project := new(Project)

	err := s.db.
		NewSelect().
		Model(&project).
		Where("id = ?", id).
		Relation("Boards").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return project, nil
}

func (s ProjectsStore) GetAll() ([]*Project, error) {
	var projects []*Project

	err := s.db.
		NewSelect().
		Model(&projects).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s ProjectsStore) Add(name string) (*Project, error) {
	project := &Project{Name: name}

	_, err := s.db.NewInsert().
		Model(project).
		Column("name").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s ProjectsStore) Edit(args EditProjectArgs) (*Project, error) {
	project := &Project{Name: args.Name}

	result, err := s.db.
		NewUpdate().
		Model(&project).
		Column("name").
		Where("id = ?", args.ID).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return nil, ErrNotFound
	}

	return project, nil
}

func (s ProjectsStore) Delete(id string) error {
	result, err := s.db.
		NewDelete().
		Model((*Project)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return ErrNotFound
	}

	return nil
}
