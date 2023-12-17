package entityservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/store"
)

type IBoardService interface {
	GetBoard(args *GetBoardOptions) (*store.Board, error)
	AddBoard(args *AddBoardOptions) (*store.Board, error)
	EditBoard(args *EditBoardOptions) error
	DeleteBoard(args *DeleteBoardOptions) error
	AddLabel(args *AddLabelOptions) (*store.Label, error)
	DeleteLabel(args *DeleteLabelOptions) error
	EditLabel(args *EditLabelOptions) error
	GetLabel(labelID store.LabelID) (*store.Label, error)
}

type BoardServiceRequirements interface {
	StoreInjector
	ActorInjector
	UserServiceInjector
	FileServiceInjector
}

type BoardService struct {
	BoardServiceRequirements
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

type DeleteBoardOptions struct {
	BoardID store.EntityID
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

func (boardService BoardService) GetBoard(args *GetBoardOptions) (*store.Board, error) {
	board := &store.Board{
		ID:     args.BoardID,
		UserID: boardService.GetActor().UserID,
	}

	q := boardService.GetStore().ORM.NewSelect().
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

	if err := q.Scan(context.TODO()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	// if service.UserID != board.UserID {
	// 	return nil, errors.New("Forbidden")
	// }

	if !args.SkipDateLastViewedUpdate {
		board.DateLastViewed = time.Now().UTC()

		updateResult, err := boardService.GetStore().ORM.NewUpdate().
			Model(board).
			Column("date_last_viewed").
			WherePK().
			Exec(context.TODO())
		if err != nil {
			return nil, err
		}

		if store.NoRowsAffected(updateResult) {
			return nil, store.ErrNotFound
		}
	}

	return board, nil
}

func (boardService BoardService) AddBoard(args *AddBoardOptions) (*store.Board, error) {
	board := &store.Board{
		ProjectID: args.ProjectID,
		UserID:    boardService.GetActor().UserID,
		Name:      args.Name,
		Color:     args.Color,
		CoverID:   nil,
	}

	if args.CoverID != nil {
		coverFile, _ := boardService.GetFileService().Get(context.TODO(), *args.CoverID)

		if coverFile.IsImage() {
			board.CoverID = args.CoverID
			board.Cover = coverFile
		}
	}

	_, err := boardService.GetStore().ORM.NewInsert().
		Model(board).
		Column("project_id", "user_id", "name", "color", "cover_id").
		Returning("*").
		Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (boardService BoardService) EditBoard(args *EditBoardOptions) error {
	q := boardService.GetStore().ORM.NewUpdate().
		Model((*store.Board)(nil)).
		Where("id = ?", args.BoardID).
		Where("user_id = ?", boardService.GetActor().UserID)

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

	updateResult, err := q.Exec(context.TODO())
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (boardService BoardService) DeleteBoard(args *DeleteBoardOptions) error {
	deleteResult, err := boardService.GetStore().ORM.NewDelete().
		Model((*store.Board)(nil)).
		Where("id = ?", args.BoardID).
		Where("user_id = ?", boardService.GetActor().UserID).
		Exec(context.TODO())

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (boardService BoardService) AddLabel(args *AddLabelOptions) (*store.Label, error) {
	owns, err := boardService.GetUserService().OwnsBoard(args.BoardID)
	if err != nil {
		return nil, err
	}

	if !owns {
		return nil, ErrPermissionDenied
	}

	label := &store.Label{
		BoardID: args.BoardID,
		UserID:  boardService.GetActor().UserID,
		Name:    args.Name,
		Color:   args.Color,
	}

	_, err = boardService.GetStore().ORM.NewInsert().
		Model(label).
		Column("board_id", "user_id", "name", "color").
		Returning("*").
		Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return label, nil
}

func (boardService BoardService) DeleteLabel(args *DeleteLabelOptions) error {
	deleteResult, err := boardService.GetStore().ORM.NewDelete().
		Model((*store.Label)(nil)).
		Where("id = ?", args.LabelID).
		Where("user_id = ?", boardService.GetActor().UserID).
		Exec(context.TODO())

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (boardService BoardService) EditLabel(args *EditLabelOptions) error {
	q := boardService.GetStore().ORM.NewUpdate().
		Model((*store.Label)(nil)).
		Where("id = ?", args.LabelID).
		Where("user_id = ?", boardService.GetActor().UserID)

	changedFields := 0

	if args.Name != nil {
		q = q.Set("name = ?", *args.Name)
		changedFields++
	}
	if args.Color != nil {
		q = q.Set("color = ?", *args.Color)
	}

	updateResult, err := q.Exec(context.TODO())
	if err != nil {
		return err
	} else if store.NoRowsAffected(updateResult) {
		return store.ErrNotFound
	}

	return nil
}

func (boardService BoardService) GetLabel(labelID store.LabelID) (*store.Label, error) {
	label := new(store.Label)
	err := boardService.GetStore().ORM.NewSelect().Model(label).Scan(context.TODO())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	return label, nil
}
