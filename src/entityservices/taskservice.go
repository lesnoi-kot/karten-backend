package entityservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/modules/markdown"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/samber/lo"
)

type ITaskService interface {
	GetTask(args *GetTaskOptions) (*store.Task, error)
	AddTask(args *AddTaskOptions) (*store.Task, error)
	EditTask(args *EditTaskOptions) error
	DeleteTask(args *DeleteTaskOptions) error

	AddLabelToTask(args *AddLabelToTaskOptions) error
	DeleteLabelFromTask(args *AddLabelToTaskOptions) error

	AttachFilesToTask(args *AttachFilesToTask) error
}

type TaskServiceRequirements interface {
	StoreInjector
	ActorInjector
	UserServiceInjector
}

type TaskService struct {
	TaskServiceRequirements
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

type AttachFilesToTask struct {
	TaskID  store.EntityID
	FilesID []store.FileID
}

func (taskService TaskService) GetTask(args *GetTaskOptions) (*store.Task, error) {
	task := new(store.Task)

	q := taskService.GetStore().ORM.
		NewSelect().
		Model(task).
		Where("task.id = ?", args.TaskID).
		Where("task.user_id = ?", taskService.GetActor().UserID)

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

	err := q.Scan(context.TODO())
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

func (taskService TaskService) AddTask(args *AddTaskOptions) (*store.Task, error) {
	task := &store.Task{
		UserID:     taskService.GetActor().UserID,
		TaskListID: args.TaskListID,
		Name:       args.Name,
		Text:       args.Text,
		Position:   args.Position,
		DueDate:    args.DueDate,
	}

	_, err := taskService.GetStore().ORM.NewInsert().
		Model(task).
		Column("task_list_id", "user_id", "name", "text", "position", "due_date").
		Returning("*").
		Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (taskService TaskService) EditTask(args *EditTaskOptions) error {
	q := taskService.GetStore().ORM.NewUpdate().
		Model((*store.Task)(nil)).
		Where("id = ?", args.TaskID).
		Where("user_id = ?", taskService.GetActor().UserID)

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

	updateResult, err := q.Exec(context.TODO())
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (taskService TaskService) DeleteTask(args *DeleteTaskOptions) error {
	deleteResult, err := taskService.GetStore().ORM.NewDelete().
		Model((*store.Task)(nil)).
		Where("id = ?", args.TaskID).
		Where("user_id = ?", taskService.GetActor().UserID).
		Exec(context.TODO())

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (taskService TaskService) AddLabelToTask(args *AddLabelToTaskOptions) error {
	if owns, err := taskService.GetUserService().OwnsTask(args.TaskID); err != nil {
		return err
	} else if !owns {
		return ErrPermissionDenied
	}

	assoc := &store.LabelToTaskAssoc{
		TaskID:  args.TaskID,
		LabelID: args.LabelID,
	}

	_, err := taskService.GetStore().ORM.NewInsert().Model(assoc).Exec(context.Background())
	return err
}

func (taskService TaskService) DeleteLabelFromTask(args *AddLabelToTaskOptions) error {
	if owns, err := taskService.GetUserService().OwnsTask(args.TaskID); err != nil {
		return err
	} else if !owns {
		return ErrPermissionDenied
	}

	_, err := taskService.GetStore().ORM.NewDelete().
		Model((*store.LabelToTaskAssoc)(nil)).
		Where("task_id = ?", args.TaskID).
		Where("label_id = ?", args.LabelID).
		Exec(context.Background())
	return err
}

func (taskService TaskService) AttachFilesToTask(args *AttachFilesToTask) error {
	if len(args.FilesID) == 0 {
		return nil
	}

	if owns, err := taskService.GetUserService().OwnsTask(args.TaskID); err != nil {
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

	_, err := taskService.GetStore().ORM.NewInsert().Model(&assocs).Exec(context.TODO())
	if err != nil {
		return err
	}

	return nil
}
