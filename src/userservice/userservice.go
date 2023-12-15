package userservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/samber/lo"

	"github.com/lesnoi-kot/karten-backend/src/modules/markdown"
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

type GetBoardOptions struct {
	BoardID                  store.EntityID
	SkipDateLastViewedUpdate bool
	IncludeTaskLists         bool
	IncludeTasks             bool
	IncludeProject           bool
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

type ClearTaskListOptions struct {
	TaskListID store.EntityID
}

type DeleteTaskListOptions struct {
	TaskListID store.EntityID
}

type GetTaskOptions struct {
	TaskID                store.EntityID
	IncludeComments       bool
	IncludeLabels         bool
	IncludeAttachments    bool
	SkipTextRender        bool
	SkipCommentTextRender bool
}

type AddTaskOptions struct {
	TaskListID store.EntityID
	Name       string
	Text       string
	Position   int64
	DueDate    *time.Time
}

type EditTaskOptions struct {
	TaskID              store.EntityID
	TaskListID          *store.EntityID
	Name                *string
	Text                *string
	Position            *int64
	SpentTime           *int64
	DueDate             *time.Time
	DateStartedTracking *time.Time
}

type DeleteTaskOptions struct {
	TaskID store.EntityID
}

type GetCommentOptions struct {
	CommentID store.EntityID
}

type AddCommentOptions struct {
	TaskID store.EntityID
	Text   string
}

type EditCommentOptions struct {
	CommentID store.EntityID
	Text      *string
}

type DeleteCommentOptions struct {
	CommentID store.EntityID
}

type AddLabelOptions struct {
	BoardID store.EntityID
	Name    string
	Color   store.Color
}

type DeleteLabelOptions struct {
	LabelID store.LabelID
}

type AddLabelToTaskOptions struct {
	TaskID  store.EntityID
	LabelID store.LabelID
}

type EditLabelOptions struct {
	LabelID store.LabelID
	Name    *string
	Color   *store.Color
}

type AttachFilesToComment struct {
	CommentID store.EntityID
	FilesID   []store.FileID
}

type AttachFilesToTask struct {
	TaskID  store.EntityID
	FilesID []store.FileID
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

func (user UserService) GetProjects(args *GetProjectsOptions) ([]*store.Project, error) {
	var projects []*store.Project

	q := user.Store.ORM.NewSelect().
		Model(&projects).
		Where("project.user_id = ?", user.UserID).
		Relation("Avatar").
		Relation("Avatar.Thumbnails")

	if args != nil && args.IncludeBoards {
		q = q.Relation("Boards").Relation("Boards.Cover")
	}

	err := q.Scan(user.Context)
	if err != nil {
		return nil, err
	}

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
		Relation("Cover").
		Relation("Labels")

	if args.IncludeProject {
		q = q.Relation("Project")
	}

	if args.IncludeTaskLists {
		q = q.Relation("TaskLists")

		if args.IncludeTasks {
			q = q.Relation("TaskLists.Tasks").
				Relation("TaskLists.Tasks.Comments").
				Relation("TaskLists.Tasks.Attachments").
				Relation("TaskLists.Tasks.Labels")
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

	changedFields := 0

	if args.Name != nil && *args.Name != "" {
		q = q.Set("name = ?", *args.Name)
		changedFields++
	}
	if args.Archived != nil {
		q = q.Set("archived = ?", *args.Archived)
		changedFields++
	}
	if args.Color != nil {
		q = q.Set("color = ?", *args.Color)
		q = q.Set("cover_id = ?", nil)
		changedFields++
	}
	if args.Favorite != nil {
		q = q.Set("favorite = ?", *args.Favorite)
		changedFields++
	}
	if args.CoverID != nil {
		q = q.Set("cover_id = ?", *args.CoverID)
		changedFields++
	}

	if changedFields == 0 {
		return nil
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

func (user UserService) OwnsTask(taskID store.EntityID) (bool, error) {
	return user.Store.ORM.NewSelect().
		Model((*store.Task)(nil)).
		Where("id = ?", taskID).
		Where("user_id = ?", user.UserID).
		Exists(user.Context)
}

func (user UserService) GetTaskList(args *GetTaskListOptions) (*store.TaskList, error) {
	taskList := new(store.TaskList)
	q := user.Store.ORM.NewSelect().
		Model(taskList).
		Where("task_list.id = ?", args.TaskListID).
		Where("task_list.user_id = ?", user.UserID).
		Where("task_list.archived = ?", false)

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

func (user UserService) ClearTaskList(args *ClearTaskListOptions) error {
	_, err := user.Store.ORM.NewDelete().
		Model((*store.Task)(nil)).
		Where("task_list_id = ?", args.TaskListID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	return err
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
		Where("task.id = ?", args.TaskID).
		Where("task.user_id = ?", user.UserID)

	if args.IncludeAttachments {
		q = q.Relation("Attachments")
	}

	if args.IncludeLabels {
		q = q.Relation("Labels")
	}

	if args.IncludeComments {
		q = q.Relation("Comments").
			Relation("Comments.Attachments").
			Relation("Comments.Author")
	}

	err := q.Scan(user.Context)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	if !args.SkipCommentTextRender {
		for _, comment := range task.Comments {
			comment.HTML = markdown.Render(comment.Text)
		}
	}

	if !args.SkipTextRender {
		task.HTML = markdown.Render(task.Text)
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
	if args.DateStartedTracking != nil {
		if args.DateStartedTracking.IsZero() {
			q = q.Set("date_started_tracking = ?", nil)
		} else {
			q = q.Set("date_started_tracking = ?", *args.DateStartedTracking)
		}
	}
	if args.SpentTime != nil {
		q = q.Set("spent_time = ?", *args.SpentTime)
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

func (user UserService) AddLabelToTask(args *AddLabelToTaskOptions) error {
	if owns, err := user.OwnsTask(args.TaskID); err != nil {
		return err
	} else if !owns {
		return ErrPermissionDenied
	}

	assoc := &store.LabelToTaskAssoc{
		TaskID:  args.TaskID,
		LabelID: args.LabelID,
	}

	_, err := user.Store.ORM.NewInsert().Model(assoc).Exec(context.Background())
	return err
}

func (user UserService) DeleteLabelFromTask(args *AddLabelToTaskOptions) error {
	if owns, err := user.OwnsTask(args.TaskID); err != nil {
		return err
	} else if !owns {
		return ErrPermissionDenied
	}

	_, err := user.Store.ORM.NewDelete().
		Model((*store.LabelToTaskAssoc)(nil)).
		Where("task_id = ?", args.TaskID).
		Where("label_id = ?", args.LabelID).
		Exec(context.Background())
	return err
}

func (user UserService) GetComment(args *GetCommentOptions) (*store.Comment, error) {
	comment := new(store.Comment)
	q := user.Store.ORM.
		NewSelect().
		Model(comment).
		Relation("Attachments").
		Relation("Author").
		Where("comment.id = ?", args.CommentID).
		Where("comment.user_id = ?", user.UserID)

	err := q.Scan(user.Context)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	comment.HTML = markdown.Render(comment.Text)
	return comment, nil
}

func (user UserService) AddComment(args *AddCommentOptions) (*store.Comment, error) {
	comment := &store.Comment{
		TaskID: args.TaskID,
		UserID: user.UserID,
		Text:   args.Text,
	}

	_, err := user.Store.ORM.NewInsert().
		Model(comment).
		Column("task_id", "user_id", "text").
		Returning("id").
		Exec(user.Context)
	if err != nil {
		return nil, err
	}

	comment, err = user.GetComment(&GetCommentOptions{
		CommentID: comment.ID,
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (user UserService) EditComment(args *EditCommentOptions) error {
	q := user.Store.ORM.NewUpdate().
		Model((*store.Comment)(nil)).
		Where("id = ?", args.CommentID).
		Where("user_id = ?", user.UserID)

	changedFields := 0

	if args.Text != nil {
		q = q.Set("text = ?", *args.Text)
		changedFields++
	}

	if changedFields == 0 {
		return nil
	}

	updateResult, err := q.Exec(user.Context)
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (user UserService) DeleteComment(args *DeleteCommentOptions) error {
	deleteResult, err := user.Store.ORM.NewDelete().
		Model((*store.Comment)(nil)).
		Where("id = ?", args.CommentID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (user UserService) OwnsComment(commentID store.EntityID) (bool, error) {
	return user.Store.ORM.NewSelect().
		Model((*store.Comment)(nil)).
		Where("id = ?", commentID).
		Where("user_id = ?", user.UserID).
		Exists(user.Context)
}

func (user UserService) AddLabel(args *AddLabelOptions) (*store.Label, error) {
	owns, err := user.OwnsBoard(args.BoardID)
	if err != nil {
		return nil, err
	}

	if !owns {
		return nil, ErrPermissionDenied
	}

	label := &store.Label{
		BoardID: args.BoardID,
		UserID:  user.UserID,
		Name:    args.Name,
		Color:   args.Color,
	}

	_, err = user.Store.ORM.NewInsert().
		Model(label).
		Column("board_id", "user_id", "name", "color").
		Returning("*").
		Exec(user.Context)
	if err != nil {
		return nil, err
	}

	return label, nil
}

func (user UserService) DeleteLabel(args *DeleteLabelOptions) error {
	deleteResult, err := user.Store.ORM.NewDelete().
		Model((*store.Label)(nil)).
		Where("id = ?", args.LabelID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (user UserService) EditLabel(args *EditLabelOptions) error {
	q := user.Store.ORM.NewUpdate().
		Model((*store.Label)(nil)).
		Where("id = ?", args.LabelID).
		Where("user_id = ?", user.UserID)

	changedFields := 0

	if args.Name != nil {
		q = q.Set("name = ?", *args.Name)
		changedFields++
	}
	if args.Color != nil {
		q = q.Set("color = ?", *args.Color)
	}

	updateResult, err := q.Exec(user.Context)
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (user UserService) GetLabel(labelID store.LabelID) (*store.Label, error) {
	label := new(store.Label)
	err := user.Store.ORM.NewSelect().Model(label).Scan(user.Context)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return label, nil
}

func (user UserService) AttachFilesToTask(args *AttachFilesToTask) error {
	if len(args.FilesID) == 0 {
		return nil
	}

	if owns, err := user.OwnsTask(args.TaskID); err != nil {
		return err
	} else if !owns {
		return ErrPermissionDenied
	}

	assocs := lo.Map(args.FilesID, func(fileID store.FileID, _ int) *store.AttachmentToTaskAssoc {
		return &store.AttachmentToTaskAssoc{
			TaskID: args.TaskID,
			FileID: fileID,
		}
	})

	_, err := user.Store.ORM.NewInsert().Model(&assocs).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (user UserService) AttachFilesToComment(args *AttachFilesToComment) error {
	if len(args.FilesID) == 0 {
		return nil
	}

	if owns, err := user.OwnsComment(args.CommentID); err != nil {
		return err
	} else if !owns {
		return ErrPermissionDenied
	}

	assocs := lo.Map(args.FilesID, func(fileID store.FileID, _ int) *store.AttachmentToCommentAssoc {
		return &store.AttachmentToCommentAssoc{
			CommentID: args.CommentID,
			FileID:    fileID,
		}
	})

	_, err := user.Store.ORM.NewInsert().Model(&assocs).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
