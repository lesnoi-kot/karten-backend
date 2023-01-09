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

func (s ProjectsStore) Get(ctx context.Context, id string) (*Project, error) {
	project := new(Project)

	err := s.db.
		NewSelect().
		Model(project).
		Where("id = ?", id).
		Relation("Boards").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return project, nil
}

func (s ProjectsStore) GetAll(ctx context.Context) ([]*Project, error) {
	var projects []*Project

	err := s.db.NewSelect().
		Model(&projects).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s ProjectsStore) Add(ctx context.Context, project *Project) error {
	_, err := s.db.NewInsert().
		Model(project).
		Column("name").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s ProjectsStore) Update(ctx context.Context, project *Project) error {
	result, err := s.db.NewUpdate().
		Model(project).
		Column("name").
		Where("id = ?", project.ID).
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

func (s ProjectsStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.NewDelete().
		Model((*Project)(nil)).
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
