package userservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/uptrace/bun"
)

var ErrPermissionDenied error = errors.New("Permission denied")

type UserService struct {
	Context context.Context
	UserID  store.UserID
	Store   *store.Store
}

type GetUserOptions struct {
	FullInfo      bool
	IncludeAvatar bool
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

type GetBoardOptions struct {
	BoardID                  store.EntityID
	SkipDateLastViewedUpdate bool
	IncludeTaskLists         bool
	IncludeTasks             bool
}

type AddBoardOptions struct {
	ProjectID store.EntityID
	Name      string
	Color     store.Color
	CoverID   *store.FileID
}

type EditBoardOptions struct {
	BoardID  store.EntityID
	Name     *string
	Archived *bool
	Color    *store.Color
	Favorite *bool
	CoverID  *store.FileID
}

type GetTaskListOptions struct {
	TaskListID   store.EntityID
	IncludeTasks bool
}

type EditTaskListOptions struct {
	TaskListID store.EntityID
	Name       *string
	Archived   *bool
	Color      *store.Color
	Position   *int64
}

type DeleteBoardOptions struct {
	BoardID store.EntityID
}

type AddTaskListOptions struct {
	BoardID  store.EntityID
	Name     string
	Color    store.Color
	Position int64
}

type DeleteTaskListOptions struct {
	TaskListID store.EntityID
}

type GetTaskOptions struct {
	TaskID          store.EntityID
	IncludeComments bool
}

type AddTaskOptions struct {
	TaskListID store.EntityID
	Name       string
	Text       string
	Position   int64
	DueDate    *time.Time
}

type EditTaskOptions struct {
	TaskID     store.EntityID
	TaskListID *store.EntityID
	Name       *string
	Text       *string
	Position   *int64
	DueDate    *time.Time
}

type DeleteTaskOptions struct {
	TaskID store.EntityID
}

func (user *UserService) SetContext(ctx context.Context) {
	user.Context = ctx
}

func (user *UserService) IsValidUser() (bool, error) {
	ok, err := user.Store.ORM.NewSelect().
		Model((*store.User)(nil)).
		Where("id = ?", user.UserID).
		Exists(user.Context)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (userService *UserService) GetUser(args *GetUserOptions) (*store.User, error) {
	user := new(store.User)
	q := userService.Store.ORM.NewSelect().
		Model(user).
		Column("id", "name", "date_created").
		Where("? = ?", bun.Ident("user.id"), userService.UserID)

	if args.FullInfo {
		q = q.Column("social_id", "login", "email", "url")
	}
	if args.IncludeAvatar {
		q = q.Relation("Avatar")
	}

	err := q.Scan(userService.Context)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (user *UserService) Delete() error {
	return user.Store.Users.Delete(user.Context, user.UserID)
}

func (user UserService) GetProjects() ([]*store.Project, error) {
	var projects []*store.Project

	err := user.Store.ORM.NewSelect().
		Model(&projects).
		Where("project.user_id = ?", user.UserID).
		Relation("Boards").
		Relation("Boards.Cover").
		Relation("Avatar").
		Relation("Avatar.Thumbnails").
		Scan(user.Context)

	return projects, err
}

func (user UserService) GetProject(args *GetProjectOptions) (*store.Project, error) {
	project := new(store.Project)
	q := user.Store.ORM.NewSelect().
		Model(project).
		Where("project.id = ?", args.ProjectID).
		Where("project.user_id = ?", user.UserID).
		Relation("Avatar").
		Relation("Avatar.Thumbnails")

	if args.IncludeBoards {
		q = q.Relation("Boards").Relation("Boards.Cover")
	}

	err := q.Scan(user.Context)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return project, nil
}

func (user UserService) AddProject(args *AddProjectOptions) (*store.Project, error) {
	project := &store.Project{
		UserID: user.UserID,
		Name:   args.Name,
	}

	if args.AvatarID != nil {
		avatarFile, err := user.Store.Files.GetImage(user.Context, *args.AvatarID)
		if err != nil {
			return nil, err
		}

		project.AvatarID = avatarFile.ID
		project.Avatar = avatarFile
	}

	_, err := user.Store.ORM.NewInsert().
		Model(project).
		Column("user_id", "name", "avatar_id").
		Returning("*").
		Exec(user.Context)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (user UserService) EditProject(args *EditProjectOptions) error {
	if args.AvatarID == nil && args.Name == nil {
		return nil
	}

	q := user.Store.ORM.NewUpdate().
		Model((*store.Project)(nil)).
		Where("id = ?", args.ProjectID).
		Where("user_id = ?", user.UserID)

	if args.Name != nil && *args.Name != "" {
		q = q.Set("name = ?", *args.Name)
	}
	if args.AvatarID != nil {
		q = q.Set("avatar_id = ?", *args.AvatarID)
	}

	updateResult, err := q.Exec(user.Context)
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (user UserService) ClearProject(projectID store.EntityID) error {
	owns, err := user.OwnsProject(projectID)
	if err != nil {
		return err
	}
	if !owns {
		return ErrPermissionDenied
	}

	_, err = user.Store.ORM.NewDelete().
		Model((*store.Board)(nil)).
		Where("project_id = ?", projectID).
		Exec(user.Context)

	return err
}

func (user UserService) DeleteProject(args *DeleteProjectOptions) error {
	result, err := user.Store.ORM.NewDelete().
		Model((*store.Project)(nil)).
		Where("id = ?", args.ProjectID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	if store.NoRowsAffected(result) {
		return store.ErrNotFound
	}

	return err
}

func (user UserService) DeleteAllProjects() error {
	_, err := user.Store.ORM.NewDelete().
		Model((*store.Project)(nil)).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)
	return err
}

func (user UserService) OwnsProject(projectID store.EntityID) (bool, error) {
	return user.Store.ORM.NewSelect().
		Model((*store.Project)(nil)).
		Where("id = ?", projectID).
		Where("user_id = ?", user.UserID).
		Exists(user.Context)
}

func (user UserService) GetBoard(args *GetBoardOptions) (*store.Board, error) {
	board := &store.Board{ID: args.BoardID, UserID: user.UserID}

	q := user.Store.ORM.NewSelect().
		Model(board).
		WherePK("id", "user_id").
		Where("archived = ?", false).
		Relation("Cover")

	if args.IncludeTaskLists {
		q = q.Relation("TaskLists")

		if args.IncludeTasks {
			q = q.Relation("TaskLists.Tasks")
		}
	}

	if err := q.Scan(user.Context); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	if user.UserID != board.UserID {
		return nil, errors.New("Forbidden")
	}

	if !args.SkipDateLastViewedUpdate {
		board.DateLastViewed = time.Now().UTC()

		updateResult, err := user.Store.ORM.NewUpdate().
			Model(board).
			Column("date_last_viewed").
			WherePK().
			Exec(user.Context)
		if err != nil {
			return nil, err
		}

		if store.NoRowsAffected(updateResult) {
			return nil, store.ErrNotFound
		}
	}

	return board, nil
}

func (user UserService) AddBoard(args *AddBoardOptions) (*store.Board, error) {
	board := &store.Board{
		ProjectID: args.ProjectID,
		UserID:    user.UserID,
		Name:      args.Name,
		Color:     args.Color,
		CoverID:   nil,
	}

	if args.CoverID != nil {
		coverFile, _ := user.Store.Files.Get(user.Context, *args.CoverID)

		if coverFile.IsImage() {
			board.CoverID = args.CoverID
			board.Cover = coverFile
		}
	}

	_, err := user.Store.ORM.NewInsert().
		Model(board).
		Column("project_id", "user_id", "name", "color", "cover_id").
		Returning("*").
		Exec(user.Context)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (user UserService) EditBoard(args *EditBoardOptions) error {
	q := user.Store.ORM.NewUpdate().
		Model((*store.Board)(nil)).
		Where("id = ?", args.BoardID).
		Where("user_id = ?", user.UserID)

	if args.Name != nil && *args.Name != "" {
		q = q.Set("name = ?", *args.Name)
	}
	if args.Archived != nil {
		q = q.Set("archived = ?", *args.Archived)
	}
	if args.Color != nil {
		q = q.Set("color = ?", *args.Color)
	}
	if args.Favorite != nil {
		q = q.Set("favorite = ?", *args.Favorite)
	}
	if args.CoverID != nil {
		q = q.Set("cover_id = ?", *args.CoverID)
	}

	updateResult, err := q.Exec(user.Context)
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (user UserService) DeleteBoard(args *DeleteBoardOptions) error {
	deleteResult, err := user.Store.ORM.NewDelete().
		Model((*store.Board)(nil)).
		Where("id = ?", args.BoardID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (user UserService) OwnsBoard(boardID store.EntityID) (bool, error) {
	return user.Store.ORM.NewSelect().
		Model((*store.Board)(nil)).
		Where("id = ?", boardID).
		Where("user_id = ?", user.UserID).
		Exists(user.Context)
}

func (user UserService) GetTaskList(args *GetTaskListOptions) (*store.TaskList, error) {
	taskList := new(store.TaskList)
	q := user.Store.ORM.NewSelect().
		Model(taskList).
		Where("id = ?", args.TaskListID).
		Where("user_id = ?", user.UserID).
		Where("archived = ?", false)

	if args.IncludeTasks {
		q = q.Relation("Tasks")
	}

	if err := q.Scan(user.Context); err != nil {
		return nil, err
	}

	return taskList, nil
}

func (user UserService) AddTaskList(args *AddTaskListOptions) (*store.TaskList, error) {
	owns, err := user.OwnsBoard(args.BoardID)
	if err != nil {
		return nil, err
	}

	if !owns {
		return nil, ErrPermissionDenied
	}

	taskList := &store.TaskList{
		BoardID:  args.BoardID,
		UserID:   user.UserID,
		Name:     args.Name,
		Color:    args.Color,
		Position: args.Position,
	}

	_, err = user.Store.ORM.NewInsert().
		Model(taskList).
		Column("board_id", "user_id", "name", "color", "position").
		Returning("*").
		Exec(user.Context)
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func (user UserService) EditTaskList(args *EditTaskListOptions) error {
	q := user.Store.ORM.NewUpdate().
		Model((*store.TaskList)(nil)).
		Where("id = ?", args.TaskListID).
		Where("user_id = ?", user.UserID)

	if args.Name != nil && *args.Name != "" {
		q = q.Set("name = ?", *args.Name)
	}
	if args.Archived != nil {
		q = q.Set("archived = ?", *args.Archived)
	}
	if args.Color != nil {
		q = q.Set("color = ?", *args.Color)
	}
	if args.Position != nil {
		q = q.Set("position = ?", *args.Position)
	}

	updateResult, err := q.Exec(user.Context)
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (user UserService) DeleteTaskList(args *DeleteTaskListOptions) error {
	deleteResult, err := user.Store.ORM.NewDelete().
		Model((*store.TaskList)(nil)).
		Where("id = ?", args.TaskListID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (user UserService) GetTask(args *GetTaskOptions) (*store.Task, error) {
	task := new(store.Task)

	q := user.Store.ORM.
		NewSelect().
		Model(task).
		Where("id = ?", args.TaskID).
		Where("user_id = ?", user.UserID)

	if args.IncludeComments {
		q = q.Relation("Comments")
	}

	err := q.Scan(user.Context)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return task, nil
}

func (user UserService) AddTask(args *AddTaskOptions) (*store.Task, error) {
	task := &store.Task{
		UserID:     user.UserID,
		TaskListID: args.TaskListID,
		Name:       args.Name,
		Text:       args.Text,
		Position:   args.Position,
		DueDate:    args.DueDate,
	}

	_, err := user.Store.ORM.NewInsert().
		Model(task).
		Column("task_list_id", "user_id", "name", "text", "position", "due_date").
		Returning("*").
		Exec(user.Context)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (user UserService) EditTask(args *EditTaskOptions) error {
	q := user.Store.ORM.NewUpdate().
		Model((*store.Task)(nil)).
		Where("id = ?", args.TaskID).
		Where("user_id = ?", user.UserID)

	if args.TaskListID != nil {
		q = q.Set("task_list_id = ?", *args.TaskListID)
	}
	if args.Name != nil && *args.Name != "" {
		q = q.Set("name = ?", *args.Name)
	}
	if args.Text != nil {
		q = q.Set("text = ?", *args.Text)
	}
	if args.Position != nil {
		q = q.Set("position = ?", *args.Position)
	}
	if args.DueDate != nil {
		q = q.Set("due_date = ?", *args.DueDate)
	}

	updateResult, err := q.Exec(user.Context)
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (user UserService) DeleteTask(args *DeleteTaskOptions) error {
	deleteResult, err := user.Store.ORM.NewDelete().
		Model((*store.Task)(nil)).
		Where("id = ?", args.TaskID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}
