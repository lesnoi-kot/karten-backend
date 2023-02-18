package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type TaskListDTO struct {
	ID          string      `json:"id"`
	BoardID     string      `json:"board_id"`
	Name        string      `json:"name"`
	Archived    bool        `json:"archived"`
	Position    int64       `json:"position"`
	DateCreated time.Time   `json:"date_created"`
	Color       store.Color `json:"color"`

	Tasks []*TaskDTO `json:"tasks,omitempty"`
}

func (api *APIService) getTaskList(c echo.Context) error {
	id := c.Param("id")

	taskList, err := api.store.TaskLists.Get(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskList))
}

func (api *APIService) addTaskList(c echo.Context) error {
	boardID := c.Param("id")

	var body struct {
		Name     string `json:"name" validate:"required,min=1,max=32"`
		Color    int    `json:"color"`
		Position int64  `json:"position"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskList := &store.TaskList{
		BoardID:  boardID,
		Name:     body.Name,
		Color:    body.Color,
		Position: body.Position,
	}

	if err := api.store.TaskLists.Add(context.Background(), taskList); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskList))
}

func (api *APIService) editTaskList(c echo.Context) error {
	var body struct {
		Name     *string      `json:"name" validate:"min=1,max=32"`
		Archived *bool        `json:"archived"`
		Position *int64       `json:"position"`
		Color    *store.Color `json:"color"`
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

	id := c.Param("id")
	taskList, err := api.store.TaskLists.Get(context.Background(), id)
	if err != nil {
		return err
	}

	if body.Name != nil {
		taskList.Name = *body.Name
	}
	if body.Archived != nil {
		taskList.Archived = *body.Archived
	}
	if body.Position != nil {
		taskList.Position = *body.Position
	}
	if body.Color != nil {
		taskList.Color = *body.Color
	}

	if err := api.store.TaskLists.Update(context.Background(), taskList); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskList))
}

func (api *APIService) deleteTaskList(c echo.Context) error {
	id := c.Param("id")

	if err := api.store.TaskLists.Delete(context.Background(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
