package entityservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lesnoi-kot/karten-backend/src/modules/markdown"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/samber/lo"
)

type ICommentService interface {
	GetComment(args *GetCommentOptions) (*store.Comment, error)
	AddComment(args *AddCommentOptions) (*store.Comment, error)
	EditComment(args *EditCommentOptions) error
	DeleteComment(args *DeleteCommentOptions) error
	AttachFilesToComment(args *AttachFilesToComment) error
}

type CommentServiceRequirements interface {
	StoreInjector
	ActorInjector
	UserServiceInjector
}

type CommentService struct {
	CommentServiceRequirements
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

type AttachFilesToComment struct {
	CommentID store.EntityID
	FilesID   []store.FileID
}

func (commentService CommentService) GetComment(args *GetCommentOptions) (*store.Comment, error) {
	comment := new(store.Comment)
	q := commentService.GetStore().ORM.
		NewSelect().
		Model(comment).
		Relation("Attachments").
		Relation("Author").
		Where("comment.id = ?", args.CommentID).
		Where("comment.user_id = ?", commentService.GetActor().UserID)

	err := q.Scan(context.TODO())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}

		return nil, err
	}

	comment.HTML = markdown.Render(comment.Text)
	return comment, nil
}

func (commentService CommentService) AddComment(args *AddCommentOptions) (*store.Comment, error) {
	comment := &store.Comment{
		TaskID: args.TaskID,
		UserID: commentService.GetActor().UserID,
		Text:   args.Text,
	}

	_, err := commentService.GetStore().ORM.NewInsert().
		Model(comment).
		Column("task_id", "user_id", "text").
		Returning("id").
		Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	comment, err = commentService.GetComment(&GetCommentOptions{
		CommentID: comment.ID,
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (commentService CommentService) EditComment(args *EditCommentOptions) error {
	q := commentService.GetStore().ORM.NewUpdate().
		Model((*store.Comment)(nil)).
		Where("id = ?", args.CommentID).
		Where("user_id = ?", commentService.GetActor().UserID)

	changedFields := 0

	if args.Text != nil {
		q = q.Set("text = ?", *args.Text)
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

func (commentService CommentService) DeleteComment(args *DeleteCommentOptions) error {
	deleteResult, err := commentService.GetStore().ORM.NewDelete().
		Model((*store.Comment)(nil)).
		Where("id = ?", args.CommentID).
		Where("user_id = ?", commentService.GetActor().UserID).
		Exec(context.TODO())

	if store.NoRowsAffected(deleteResult) {
		return store.ErrNotFound
	}

	return err
}

func (commentService CommentService) AttachFilesToComment(args *AttachFilesToComment) error {
	if len(args.FilesID) == 0 {
		return nil
	}

	if owns, err := commentService.GetUserService().OwnsComment(args.CommentID); err != nil {
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

	_, err := commentService.GetStore().ORM.NewInsert().Model(&assocs).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
