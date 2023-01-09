package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func (api *APIService) addComment(c echo.Context) error {
	var body struct {
		Text string `json:"text" validate:"required,min=1,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	body.Text = strings.TrimSpace(body.Text)
	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
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
		return echo.ErrBadRequest
	}

	body.Text = strings.TrimSpace(body.Text)
	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	id := c.Param("id")
	comment, err := api.store.Comments.Get(context.Background(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}

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
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}
