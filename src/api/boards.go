package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type BoardDTO struct {
	ID             string      `json:"id"`
	ShortID        string      `json:"short_id"`
	UserID         int         `json:"user_id"`
	Name           string      `json:"name"`
	ProjectID      string      `json:"project_id"`
	Archived       bool        `json:"archived"`
	Favorite       bool        `json:"favorite"`
	DateCreated    time.Time   `json:"date_created"`
	DateLastViewed time.Time   `json:"date_last_viewed"`
	Color          store.Color `json:"color"`
	CoverURL       string      `json:"cover_url,omitempty"`

	ProjectName string         `json:"project_name"`
	TaskLists   []*TaskListDTO `json:"task_lists,omitempty"`
	Labels      []*LabelDTO    `json:"labels,omitempty"`
}

func (api *APIService) getBoard(c echo.Context) error {
	boardID := c.Param("id")
	mode := c.QueryParam("mode")
	userService := api.mustGetUserService(c)
	board, err := userService.BoardService.GetBoard(&entityservices.GetBoardOptions{
		BoardID:          boardID,
		IncludeTaskLists: mode != "shallow",
		IncludeTasks:     mode != "shallow",
		IncludeProject:   true,
	})
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
	userService := api.mustGetUserService(c)
	board, err := userService.BoardService.AddBoard(&entityservices.AddBoardOptions{
		ProjectID: projectID,
		Name:      body.Name,
		Color:     body.Color,
		CoverID:   body.CoverID,
	})
	if err != nil {
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

	if body.Name != nil {
		*body.Name = strings.TrimSpace(*body.Name)
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	userService := api.mustGetUserService(c)
	boardID := c.Param("id")
	err := userService.BoardService.EditBoard(&entityservices.EditBoardOptions{
		BoardID:  boardID,
		Name:     body.Name,
		Archived: body.Archived,
		Color:    body.Color,
		CoverID:  body.CoverID,
	})
	if err != nil {
		return err
	}

	board, err := userService.BoardService.GetBoard(&entityservices.GetBoardOptions{
		BoardID:                  boardID,
		SkipDateLastViewedUpdate: true,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(boardToDTO(board)))
}

func (api *APIService) deleteBoard(c echo.Context) error {
	boardID := c.Param("id")
	userService := api.mustGetUserService(c)

	err := userService.BoardService.DeleteBoard(&entityservices.DeleteBoardOptions{BoardID: boardID})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) setFavoriteBoard(c echo.Context, favorite bool) error {
	boardID := c.Param("id")
	userService := api.mustGetUserService(c)

	err := userService.BoardService.EditBoard(&entityservices.EditBoardOptions{
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

func (api *APIService) addLabel(c echo.Context) error {
	var body struct {
		Name  string `json:"name" validate:"required,min=1,max=32"`
		Color int    `json:"color"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	boardID := c.Param("id")
	userService := api.mustGetUserService(c)
	label, err := userService.BoardService.AddLabel(&entityservices.AddLabelOptions{
		BoardID: boardID,
		Name:    body.Name,
		Color:   body.Color,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(labelToDTO(label)))
}

func (api *APIService) deleteLabel(c echo.Context) error {
	labelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	userService := api.mustGetUserService(c)
	err = userService.BoardService.DeleteLabel(&entityservices.DeleteLabelOptions{
		LabelID: labelID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *APIService) editLabel(c echo.Context) error {
	var body struct {
		Name  *string      `json:"name" validate:"omitempty,min=1,max=32"`
		Color *store.Color `json:"color"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	if body.Name != nil {
		*body.Name = strings.TrimSpace(*body.Name)
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	labelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	userService := api.mustGetUserService(c)
	err = userService.BoardService.EditLabel(&entityservices.EditLabelOptions{
		LabelID: labelID,
		Name:    body.Name,
		Color:   body.Color,
	})
	if err != nil {
		return err
	}

	label, err := userService.BoardService.GetLabel(labelID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(labelToDTO(label)))
}
