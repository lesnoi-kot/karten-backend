package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store/models"
)

func (api *APIService) getTask(c echo.Context) error {
	id := c.Param("id")

	var task models.Task

	err := api.store.
		NewSelect().
		Model(&task).
		Where("id = ?", id).
		Relation("Comments").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}

		return err
	}

	return c.JSON(http.StatusOK, OK(task))
}

func (api *APIService) addTask(c echo.Context) error {
	taskListID := c.Param("id")

	var body struct {
		Name     string `json:"name" validate:"required,min=1,max=32"`
		Text     string `json:"text"`
		Position int64  `json:"position"`
		DueDate  string `json:"due_date" validate:"datetime="`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	dueDate, _ := time.Parse("", body.DueDate)

	task := models.Task{
		TaskListID: taskListID,
		Name:       body.Name,
		Text:       body.Text,
		Position:   body.Position,
		DueDate:    dueDate,
	}

	_, err := api.store.
		NewInsert().
		Model(&task).
		Column("task_list_id", "name", "text", "position", "due_date").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(task))
}

func (api *APIService) deleteTask(c echo.Context) error {
	id := c.Param("id")

	result, err := api.store.
		NewDelete().
		Model((*models.Task)(nil)).
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

func (api *APIService) editTask(c echo.Context) error {
	return nil
}
