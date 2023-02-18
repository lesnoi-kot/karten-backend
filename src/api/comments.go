package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type CommentDTO struct{
	ID          string    `json:"id"`
	TaskID      string    `json:"task_id"`
	Author      string    `json:"author"`
	Text        string    `json:"text"`
	DateCreated time.Time `json:"date_created"`
}

func (api *APIService) addComment(c echo.Context) error {
	var body struct {
		Text string `json:"text" validate:"required,min=1,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Text = strings.TrimSpace(body.Text)
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskID := c.Param("id")
	comment := &store.Comment{
		TaskID: taskID,
		Author: "Author",
		Text:   body.Text,
	}

	if err := api.store.Comments.Add(context.Background(), comment); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(comment))
}

func (api *APIService) editComment(c echo.Context) error {
	var body struct {
		Text string `json:"text" validate:"required,min=1,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Text = strings.TrimSpace(body.Text)
	if err := c.Validate(&body); err != nil {
		return err
	}

	id := c.Param("id")
	comment, err := api.store.Comments.Get(context.Background(), id)
	if err != nil {
		return err
	}

	comment.Text = body.Text
	if err := api.store.Comments.Update(context.Background(), comment); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(comment))
}

func (api *APIService) deleteComment(c echo.Context) error {
	id := c.Param("id")

	if err := api.store.Comments.Delete(context.Background(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
