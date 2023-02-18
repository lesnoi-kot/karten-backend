package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func (api *APIService) getBoard(c echo.Context) error {
	id := c.Param("id")
	board, err := api.store.Boards.Get(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(board))
}

func (api *APIService) addBoard(c echo.Context) error {
	projectID := c.Param("id")

	var body struct {
		Name     string      `json:"name" validate:"required,min=1,max=32"`
		Color    store.Color `json:"color"`
		CoverURL string      `json:"cover_url"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	board := &store.Board{
		ProjectID: projectID,
		Name:      body.Name,
		Color:     body.Color,
	}

	if err := api.store.Boards.Add(context.Background(), board); err != nil {
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
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	board, err := api.store.Boards.Get(context.Background(), id)
	if err != nil {
		return err
	}

	board.Name = body.Name
	if body.Archived != nil {
		board.Archived = *body.Archived
	}
	if body.Color != nil {
		board.Color = *body.Color
	}

	if err = api.store.Boards.Update(context.Background(), board); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(board))
}

func (api *APIService) deleteBoard(c echo.Context) error {
	id := c.Param("id")

	if err := api.store.Boards.Delete(context.Background(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) setFavoriteBoard(c echo.Context, favorite bool) error {
	board := &store.Board{
		ID:       c.Param("id"),
		Favorite: favorite,
	}

	if err := api.store.Boards.UpdateColumns(context.Background(), board, "favorite"); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *APIService) favoriteBoard(c echo.Context) error {
	return api.setFavoriteBoard(c, true)
}

func (api *APIService) unfavoriteBoard(c echo.Context) error {
	return api.setFavoriteBoard(c, false)
}
