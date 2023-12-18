package entityservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lesnoi-kot/karten-backend/src/store"
)

type IProjectService interface {
	GetProject(args GetProjectOptions) (*store.Project, error)
	GetProjects(args GetProjectsOptions) ([]*store.Project, error)
	AddProject(args AddProjectOptions) (*store.Project, error)
	EditProject(args EditProjectOptions) error
	ClearProject(projectID store.EntityID) error
	DeleteProject(args DeleteProjectOptions) error
}

type ProjectServiceRequirements interface {
	StoreInjector
	ActorInjector
	FileServiceInjector
}

type ProjectService struct {
	ProjectServiceRequirements
}

type GetProjectsOptions struct {
	IncludeBoards bool
}

type GetProjectOptions struct {
	ProjectID     store.EntityID
	IncludeBoards bool
}

type AddProjectOptions struct {
	Name     string
	AvatarID *store.FileID
}

type EditProjectOptions struct {
	ProjectID store.EntityID
	Name      *string
	AvatarID  *store.FileID
}

type DeleteProjectOptions struct {
	ProjectID store.EntityID
}

func (projectService ProjectService) GetProject(args GetProjectOptions) (*store.Project, error) {
	project := new(store.Project)
	q := projectService.GetStore().ORM.NewSelect().
		Model(project).
		Where("project.id = ?", args.ProjectID).
		Where("project.user_id = ?", projectService.GetActor().UserID).
		Relation("Avatar").
		Relation("Avatar.Thumbnails")

	if args.IncludeBoards {
		q = q.Relation("Boards").Relation("Boards.Cover")
	}

	err := q.Scan(context.TODO())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return project, nil
}

func (projectService ProjectService) GetProjects(args GetProjectsOptions) ([]*store.Project, error) {
	var projects []*store.Project

	q := projectService.GetStore().ORM.NewSelect().
		Model(&projects).
		Where("project.user_id = ?", projectService.GetActor().UserID).
		Relation("Avatar").
		Relation("Avatar.Thumbnails")

	if args.IncludeBoards {
		q = q.Relation("Boards").Relation("Boards.Cover")
	}

	err := q.Scan(context.TODO())
	if err != nil {
		return nil, err
	}

	return projects, err
}

func (projectService ProjectService) AddProject(args AddProjectOptions) (*store.Project, error) {
	project := &store.Project{
		UserID: projectService.GetActor().UserID,
		Name:   args.Name,
	}

	if args.AvatarID != nil {
		avatarFile, err := projectService.GetFileService().GetImage(context.TODO(), *args.AvatarID)
		if err != nil {
			return nil, err
		}

		project.AvatarID = avatarFile.ID
		project.Avatar = avatarFile
	}

	_, err := projectService.GetStore().ORM.NewInsert().
		Model(project).
		Column("user_id", "name", "avatar_id").
		Returning("*").
		Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (projectService ProjectService) EditProject(args EditProjectOptions) error {
	if args.AvatarID == nil && args.Name == nil {
		return nil
	}

	q := projectService.GetStore().ORM.NewUpdate().
		Model((*store.Project)(nil)).
		Where("id = ?", args.ProjectID).
		Where("user_id = ?", projectService.GetActor().UserID)

	if args.Name != nil && *args.Name != "" {
		q = q.Set("name = ?", *args.Name)
	}
	if args.AvatarID != nil {
		q = q.Set("avatar_id = ?", *args.AvatarID)
	}

	updateResult, err := q.Exec(context.TODO())
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (projectService ProjectService) DeleteProject(args DeleteProjectOptions) error {
	result, err := projectService.GetStore().ORM.NewDelete().
		Model((*store.Project)(nil)).
		Where("id = ?", args.ProjectID).
		Where("user_id = ?", projectService.GetActor().UserID).
		Exec(context.TODO())

	if store.NoRowsAffected(result) {
		return store.ErrNotFound
	}

	return err
}

func (projectService ProjectService) ClearProject(projectID store.EntityID) error {
	_, err := projectService.GetStore().ORM.NewDelete().
		Model((*store.Board)(nil)).
		Where("project_id = ?", projectID).
		Where("user_id = ?", projectService.GetActor().UserID).
		Exec(context.TODO())

	return err
}
