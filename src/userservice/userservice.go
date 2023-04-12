package userservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/store"
)

type UserService struct {
	Context context.Context
	UserID  store.UserID
	Store   *store.Store
}

type GetBoardOptions struct {
	BoardID                  store.EntityID
	SkipDateLastViewedUpdate bool
	IncludeTaskLists         bool
	IncludeTasks             bool
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

func (user *UserService) SetContext(ctx context.Context) {
	user.Context = ctx
}

func (user UserService) GetBoard(args *GetBoardOptions) (*store.Board, error) {
	board := &store.Board{ID: args.BoardID, UserID: user.UserID}

	q := user.Store.ORM.NewSelect().
		Model(board).
		WherePK("id", "user_id").
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

func (user UserService) UpdateBoard(args *EditBoardOptions) error {
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
	_, err := user.Store.ORM.NewDelete().
		Model((*store.Board)(nil)).
		Where("id = ?", args.BoardID).
		Where("user_id = ?", user.UserID).
		Exec(user.Context)

	return err
}
