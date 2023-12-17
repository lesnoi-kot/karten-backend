package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type CommentDTO struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"task_id"`
	UserID      int       `json:"user_id"`
	Text        string    `json:"text"`
	HTML        string    `json:"html"`
	DateCreated time.Time `json:"date_created"`

	Author      *PublicUserDTO `json:"author"`
	Attachments []*FileDTO     `json:"attachments"`
}

func (api *APIService) getComment(c echo.Context) error {
	commentID := c.Param("id")
	user := api.mustGetUserService(c)
	comment, err := user.CommentService.GetComment(&entityservices.GetCommentOptions{
		CommentID: commentID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(commentToDTO(comment)))
}

func (api *APIService) addComment(c echo.Context) error {
	var body struct {
		Text        string   `json:"text" validate:"required,min=1"`
		Attachments []string `json:"attachments"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Text = strings.TrimSpace(body.Text)
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	comment, err := user.CommentService.AddComment(&entityservices.AddCommentOptions{
		TaskID: taskID,
		Text:   body.Text,
	})
	if err != nil {
		return err
	}

	err = user.CommentService.AttachFilesToComment(&entityservices.AttachFilesToComment{
		CommentID: comment.ID,
		FilesID:   body.Attachments,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(commentToDTO(comment)))
}

func (api *APIService) editComment(c echo.Context) error {
	var body struct {
		Text *string `json:"text" validate:"omitempty,min=1"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	if body.Text != nil {
		*body.Text = strings.TrimSpace(*body.Text)
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	commentID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.CommentService.EditComment(&entityservices.EditCommentOptions{
		CommentID: commentID,
		Text:      body.Text,
	})
	if err != nil {
		return err
	}

	comment, err := user.CommentService.GetComment(&entityservices.GetCommentOptions{CommentID: commentID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(commentToDTO(comment)))
}

func (api *APIService) deleteComment(c echo.Context) error {
	commentID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.CommentService.DeleteComment(&entityservices.DeleteCommentOptions{CommentID: commentID})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) addCommentAttachments(c echo.Context) error {
	var body struct {
		FilesID []string `json:"files_id" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	commentID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.CommentService.AttachFilesToComment(&entityservices.AttachFilesToComment{
		CommentID: commentID,
		FilesID:   body.FilesID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) deleteCommentAttachment(c echo.Context) error {
	var body struct {
		FileID string `json:"file_id" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	commentID := c.Param("id")
	_, err := api.store.ORM.NewDelete().
		Model((*store.AttachmentToCommentAssoc)(nil)).
		Where("comment_id = ?", commentID).
		Where("file_id = ?", body.FileID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
