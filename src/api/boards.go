package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/lesnoi-kot/karten-backend/src/userservice"
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
	boardID := c.Param("id")
	userService := api.getUserService(c)

	board, err := userService.GetBoard(&userservice.GetBoardOptions{BoardID: boardID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(boardToDTO(board)))
}

func (api *APIService) addBoard(c echo.Context) error {
	var body struct {
		Name    string        `json:"name" validate:"required,min=1,max=32"`
		Color   store.Color   `json:"color"`
		CoverID *store.FileID `json:"cover_id"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	projectID := c.Param("id")
	userID, _ := getUserID(c)
	board := &store.Board{
		ProjectID: projectID,
		UserID:    userID,
		Name:      body.Name,
		Color:     body.Color,
		CoverID:   nil,
	}

	if body.CoverID != nil {
		cover, _ := api.store.Files.Get(context.Background(), *body.CoverID)

		if cover.IsImage() {
			board.CoverID = body.CoverID
			board.Cover = cover
		}
	}

	if err := api.store.Boards.Add(context.Background(), board); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(boardToDTO(board)))
}

func (api *APIService) editBoard(c echo.Context) error {
	var body struct {
		Name     *string       `json:"name" validate:"omitempty,min=1,max=32"`
		Archived *bool         `json:"archived"`
		Color    *store.Color  `json:"color"`
		CoverID  *store.FileID `json:"cover_id"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	*body.Name = strings.TrimSpace(*body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	userService := api.getUserService(c)
	boardID := c.Param("id")
	board, err := userService.GetBoard(&userservice.GetBoardOptions{
		BoardID:                  boardID,
		SkipDateLastViewedUpdate: true,
	})
	if err != nil {
		return err
	}

	err = userService.UpdateBoard(&userservice.EditBoardOptions{
		BoardID:  board.ID,
		Name:     body.Name,
		Archived: body.Archived,
		Color:    body.Color,
		CoverID:  body.CoverID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(boardToDTO(board)))
}

func (api *APIService) deleteBoard(c echo.Context) error {
	boardID := c.Param("id")
	userService := api.getUserService(c)

	err := userService.DeleteBoard(&userservice.DeleteBoardOptions{BoardID: boardID})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) setFavoriteBoard(c echo.Context, favorite bool) error {
	boardID := c.Param("id")
	userService := api.getUserService(c)

	err := userService.UpdateBoard(&userservice.EditBoardOptions{
		BoardID:  boardID,
		Favorite: &favorite,
	})
	if err != nil {
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
