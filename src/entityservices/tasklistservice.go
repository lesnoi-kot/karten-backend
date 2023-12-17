package entityservices

import (
	"context"

	"github.com/lesnoi-kot/karten-backend/src/store"
)

type ITaskListService interface {
	GetTaskList(args *GetTaskListOptions) (*store.TaskList, error)
	AddTaskList(args *AddTaskListOptions) (*store.TaskList, error)
	EditTaskList(args *EditTaskListOptions) error
	ClearTaskList(args *ClearTaskListOptions) error
	DeleteTaskList(args *DeleteTaskListOptions) error
}

type TaskListServiceRequirements interface {
	StoreInjector
	ActorInjector
	UserServiceInjector
}

type TaskListService struct {
	TaskListServiceRequirements
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

func (taskListService TaskListService) GetTaskList(args *GetTaskListOptions) (*store.TaskList, error) {
	taskList := new(store.TaskList)
	q := taskListService.GetStore().ORM.NewSelect().
		Model(taskList).
		Where("task_list.id = ?", args.TaskListID).
		Where("task_list.user_id = ?", taskListService.GetActor().UserID).
		Where("task_list.archived = ?", false)

	if args.IncludeTasks {
		q = q.Relation("Tasks")
	}

	if err := q.Scan(context.TODO()); err != nil {
		return nil, err
	}

	return taskList, nil
}

func (taskListService TaskListService) AddTaskList(args *AddTaskListOptions) (*store.TaskList, error) {
	owns, err := taskListService.GetUserService().OwnsBoard(args.BoardID)
	if err != nil {
		return nil, err
	}

	if !owns {
		return nil, ErrPermissionDenied
	}

	taskList := &store.TaskList{
		BoardID:  args.BoardID,
		UserID:   taskListService.GetActor().UserID,
		Name:     args.Name,
		Color:    args.Color,
		Position: args.Position,
	}

	_, err = taskListService.GetStore().ORM.NewInsert().
		Model(taskList).
		Column("board_id", "user_id", "name", "color", "position").
		Returning("*").
		Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func (taskListService TaskListService) EditTaskList(args *EditTaskListOptions) error {
	q := taskListService.GetStore().ORM.NewUpdate().
		Model((*store.TaskList)(nil)).
		Where("id = ?", args.TaskListID).
		Where("user_id = ?", taskListService.GetActor().UserID)

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

	updateResult, err := q.Exec(context.TODO())
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (taskListService TaskListService) ClearTaskList(args *ClearTaskListOptions) error {
	_, err := taskListService.GetStore().ORM.NewDelete().
		Model((*store.Task)(nil)).
		Where("task_list_id = ?", args.TaskListID).
		Where("user_id = ?", taskListService.GetActor().UserID).
		Exec(context.TODO())

	return err
}

func (taskListService TaskListService) DeleteTaskList(args *DeleteTaskListOptions) error {
	deleteResult, err := taskListService.GetStore().ORM.NewDelete().
		Model((*store.TaskList)(nil)).
		Where("id = ?", args.TaskListID).
		Where("user_id = ?", taskListService.GetActor().UserID).
		Exec(context.TODO())

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}
