package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/modules/images"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type BoardDTO struct {
	ID             string      `json:"id"`
	ShortID        string      `json:"short_id"`
	Name           string      `json:"name"`
	ProjectID      string      `json:"project_id"`
	Archived       bool        `json:"archived"`
	Favorite       bool        `json:"favorite"`
	DateCreated    time.Time   `json:"date_created"`
	DateLastViewed time.Time   `json:"date_last_viewed"`
	Color          store.Color `json:"color"`
	CoverURL       string      `json:"cover_url,omitempty"`

	TaskLists []*TaskListDTO `json:"task_lists,omitempty"`
}

func (api *APIService) getBoard(c echo.Context) error {
	id := c.Param("id")
	board, err := api.store.Boards.Get(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(boardToDTO(board)))
}

func (api *APIService) addBoard(c echo.Context) error {
	projectID := c.Param("id")

	var body struct {
		Name    string      `form:"name" json:"name" validate:"required,min=1,max=32"`
		Color   store.Color `form:"color" json:"color"`
		CoverID *string     `form:"cover_id" json:"cover_id"`
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
		CoverID:   body.CoverID,
	}

	if body.CoverID != nil && api.store.Files.IsDefaultCover(context.Background(), *body.CoverID) {
		board.CoverID = body.CoverID
	}

	if coverFileHeader, err := c.FormFile("cover"); board.CoverID == nil && err == nil {
		coverData, err := coverFileHeader.Open()
		if err != nil {
			return err
		}

		image, err := images.ParseImage(coverData)
		if err != nil {
			return echo.
				NewHTTPError(http.StatusBadRequest, "Image can't be decoded").
				SetInternal(err)
		}

		coverData, err = coverFileHeader.Open()
		if err != nil {
			return err
		}

		coverFile, err := api.store.Files.Add(context.Background(), store.AddFileOptions{
			Name:     coverFileHeader.Filename,
			Data:     coverData,
			MIMEType: image.MIMEType,
		})
		if err != nil {
			return fmt.Errorf("Save image for new board error: %w", err)
		}

		api.logger.Infof("Added cover image for board %s: %s", board.Name, coverFile.ID)

		board.CoverID = &coverFile.ID
		board.Cover = coverFile
	}

	if err := api.store.Boards.Add(context.Background(), board); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(boardToDTO(board)))
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
