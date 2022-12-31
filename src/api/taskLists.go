package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store/models"
)

func (api *APIService) getTaskList(c echo.Context) error {
	id := c.Param("id")

	var taskList models.TaskList

	err := api.store.
		NewSelect().
		Model(&taskList).
		Where("id = ?", id).
		Relation("Tasks").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}

		return err
	}

	return c.JSON(http.StatusOK, OK(taskList))
}

func (api *APIService) addTaskList(c echo.Context) error {
	boardID := c.Param("id")

	var body struct {
		Name  string `json:"name"`
		Color int    `json:"color"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	taskList := models.TaskList{
		BoardID: boardID,
		Name:    body.Name,
		Color:   body.Color,
	}

	_, err := api.store.
		NewInsert().
		Model(&taskList).
		Column("board_id", "name", "color").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskList))
}

func (api *APIService) deleteTaskList(c echo.Context) error {
	id := c.Param("id")

	result, err := api.store.
		NewDelete().
		Model((*models.TaskList)(nil)).
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

func (api *APIService) editTaskList(c echo.Context) error {
	id := c.Param("id")

	var body struct {
		Name     string        `json:"name" validate:"required,min=1,max=32"`
		Archived *bool         `json:"archived"`
		Position *int64        `json:"position"`
		Color    *models.Color `json:"color"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	body.Name = strings.TrimSpace(body.Name)

	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	board := models.TaskList{
		Name:     body.Name,
		Archived: *body.Archived,
		Position: *body.Position,
		Color:    *body.Color,
	}

	query := api.store.NewUpdate().Model(&board)

	if body.Archived != nil {
		query = query.Column("archived")
	}
	if body.Position != nil {
		query = query.Column("position")
	}
	if body.Color != nil {
		query = query.Column("color")
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
