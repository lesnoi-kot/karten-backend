package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func (api *APIService) getTask(c echo.Context) error {
	id := c.Param("id")
	task, err := api.store.Tasks.Get(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(task))
}

func (api *APIService) addTask(c echo.Context) error {
	taskListID := c.Param("id")

	var body struct {
		Name     string     `json:"name" validate:"required,min=1,max=32"`
		Text     string     `json:"text"`
		Position int64      `json:"position"`
		DueDate  *time.Time `json:"due_date"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	task := &store.Task{
		TaskListID: taskListID,
		Name:       body.Name,
		Text:       body.Text,
		Position:   body.Position,
		DueDate:    body.DueDate,
	}

	if err := api.store.Tasks.Add(context.Background(), task); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(task))
}

func (api *APIService) editTask(c echo.Context) error {
	var body struct {
		TaskListID *string    `json:"task_list_id"`
		Name       *string    `json:"name" validate:"min=1,max=32"`
		Text       *string    `json:"text"`
		Position   *int64     `json:"position"`
		DueDate    *time.Time `json:"due_date"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	if body.Name != nil {
		*body.Name = strings.TrimSpace(*body.Name)
	}
	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	id := c.Param("id")
	task, err := api.store.Tasks.Get(context.Background(), id)
	if err != nil {
		return err
	}

	if body.Name != nil {
		task.Name = *body.Name
	}
	if body.Text != nil {
		task.Text = *body.Text
	}
	if body.Position != nil {
		task.Position = *body.Position
	}
	if body.DueDate != nil {
		task.DueDate = body.DueDate
	}
	if body.TaskListID != nil {
		task.TaskListID = *body.TaskListID
	}

	if err := api.store.Tasks.Update(context.Background(), task); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(task))
}

func (api *APIService) deleteTask(c echo.Context) error {
	id := c.Param("id")

	if err := api.store.Tasks.Delete(context.Background(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
