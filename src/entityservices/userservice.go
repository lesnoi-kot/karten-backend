package entityservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/uptrace/bun"
)

type IUserService interface {
	IsValidUser() (bool, error)
	GetUser(args *GetUserOptions) (*store.User, error)
	EditUser(args *EditUserOptions) error
	DeleteUser() error

	DeleteAllProjects() error

	OwnsProject(projectID store.EntityID) (bool, error)
	OwnsBoard(boardID store.EntityID) (bool, error)
	OwnsTask(taskID store.EntityID) (bool, error)
	OwnsComment(commentID store.EntityID) (bool, error)
}

type UserServiceRequirements interface {
	StoreInjector
	ActorInjector
}

type UserService struct {
	UserServiceRequirements
}

type GetUserOptions struct {
	FullInfo      bool
	IncludeAvatar bool
}

type EditUserOptions struct {
	Name *string
}

func (userService *UserService) IsValidUser() (bool, error) {
	ok, err := userService.GetStore().ORM.NewSelect().
		Model((*store.User)(nil)).
		Where("id = ?", userService.GetActor().UserID).
		Exists(context.Background())
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (userService *UserService) GetUser(args *GetUserOptions) (*store.User, error) {
	user := new(store.User)
	q := userService.GetStore().ORM.NewSelect().
		Model(user).
		Column("id", "name", "date_created").
		Where("? = ?", bun.Ident("user.id"), userService.GetActor().UserID)

	if args.FullInfo {
		q = q.Column("social_id", "login", "email", "url")
	}
	if args.IncludeAvatar {
		q = q.Relation("Avatar")
	}

	err := q.Scan(context.Background())

	if errors.Is(err, sql.ErrNoRows) {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userService *UserService) EditUser(args *EditUserOptions) error {
	if args == nil {
		return nil
	}

	q := userService.GetStore().ORM.NewUpdate().
		Model((*store.User)(nil)).
		Where("id = ?", userService.GetActor().UserID)

	if args.Name != nil && strings.TrimSpace(*args.Name) != "" {
		q = q.Set("name = ?", *args.Name)
	}

	result, err := q.Exec(context.Background())
	if err != nil {
		return err
	}

	if store.NoRowsAffected(result) {
		return store.ErrNotFound
	}

	return nil
}

func (userService *UserService) DeleteUser() error {
	result, err := userService.GetStore().ORM.NewDelete().
		Model((*store.User)(nil)).
		Where("id = ?", userService.GetActor().UserID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	if store.NoRowsAffected(result) {
		return store.ErrNotFound
	}

	return nil
}

func (userService UserService) DeleteAllProjects() error {
	_, err := userService.GetStore().ORM.NewDelete().
		Model((*store.Project)(nil)).
		Where("user_id = ?", userService.GetActor().UserID).
		Exec(context.Background())
	return err
}

func (userService UserService) OwnsProject(projectID store.EntityID) (bool, error) {
	return userService.GetStore().ORM.NewSelect().
		Model((*store.Project)(nil)).
		Where("id = ?", projectID).
		Where("user_id = ?", userService.GetActor().UserID).
		Exists(context.Background())
}

func (userService UserService) OwnsBoard(boardID store.EntityID) (bool, error) {
	return userService.GetStore().ORM.NewSelect().
		Model((*store.Board)(nil)).
		Where("id = ?", boardID).
		Where("user_id = ?", userService.GetActor().UserID).
		Exists(context.Background())
}

func (userService UserService) OwnsTask(taskID store.EntityID) (bool, error) {
	return userService.GetStore().ORM.NewSelect().
		Model((*store.Task)(nil)).
		Where("id = ?", taskID).
		Where("user_id = ?", userService.GetActor().UserID).
		Exists(context.Background())
}

func (userService UserService) OwnsComment(commentID store.EntityID) (bool, error) {
	return userService.GetStore().ORM.NewSelect().
		Model((*store.Comment)(nil)).
		Where("id = ?", commentID).
		Where("user_id = ?", userService.GetActor().UserID).
		Exists(context.Background())
}
