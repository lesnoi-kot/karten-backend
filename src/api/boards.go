package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func (api *APIService) getBoard(c echo.Context) error {
	id := c.Param("id")
	var board store.Board

	err := api.store.
		NewSelect().
		Model(&board).
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}

		return err
	}

	return c.JSON(http.StatusOK, OK(board))
}

func (api *APIService) addBoard(c echo.Context) error {
	projectID := c.Param("id")

	var body struct {
		Name     string      `json:"name" validate:"required,max=32"`
		Color    store.Color `json:"color"`
		CoverURL string      `json:"cover_url"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	board := store.Board{
		ProjectID: projectID,
		Name:      body.Name,
		Color:     body.Color,
		CoverURL:  body.CoverURL,
	}

	_, err := api.store.
		NewInsert().
		Model(&board).
		Column("project_id", "name", "color", "cover_url").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(board))
}

func (api *APIService) editBoard(c echo.Context) error {
	id := c.Param("id")

	var body struct {
		Name     string       `json:"name" validate:"required,min=1,max=32"`
		Archived *bool        `json:"archived"`
		Color    *store.Color `json:"color"`
		CoverURL *string      `json:"cover_url"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	body.Name = strings.TrimSpace(body.Name)

	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	board := store.Board{
		Name:     body.Name,
		Archived: *body.Archived,
		Color:    *body.Color,
		CoverURL: *body.CoverURL,
	}

	query := api.store.NewUpdate().Model(&board)

	if body.Archived != nil {
		query = query.Column("archived")
	}
	if body.Color != nil {
		query = query.Column("color")
	}
	if body.CoverURL != nil {
		query = query.Column("cover_url")
	}

	result, err := query.
		Column("name").
		Where("id = ?", id).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, OK(board))
}

func (api *APIService) deleteBoard(c echo.Context) error {
	id := c.Param("id")

	result, err := api.store.
		NewDelete().
		Model((*store.Board)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, OK(nil))
}
