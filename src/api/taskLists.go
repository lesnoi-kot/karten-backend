package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type TaskListDTO struct {
	ID          string      `json:"id"`
	BoardID     string      `json:"board_id"`
	UserID      int         `json:"user_id"`
	Name        string      `json:"name"`
	Archived    bool        `json:"archived"`
	Position    int64       `json:"position"`
	DateCreated time.Time   `json:"date_created"`
	Color       store.Color `json:"color"`

	Tasks []*TaskDTO `json:"tasks,omitempty"`
}

func (api *APIService) getTaskList(c echo.Context) error {
	taskListID := c.Param("id")
	userService := api.mustGetUserService(c)

	taskList, err := userService.TaskListService.GetTaskList(&entityservices.GetTaskListOptions{
		TaskListID:   taskListID,
		IncludeTasks: true,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskListToDTO(taskList)))
}

func (api *APIService) addTaskList(c echo.Context) error {
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

	boardID := c.Param("id")
	userService := api.mustGetUserService(c)
	taskList, err := userService.TaskListService.AddTaskList(&entityservices.AddTaskListOptions{
		BoardID:  boardID,
		Name:     body.Name,
		Color:    body.Color,
		Position: body.Position,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskListToDTO(taskList)))
}

func (api *APIService) editTaskList(c echo.Context) error {
	var body struct {
		Name     *string      `json:"name" validate:"required,min=1,max=32"`
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

	taskListID := c.Param("id")
	userService := api.mustGetUserService(c)

	err := userService.TaskListService.EditTaskList(&entityservices.EditTaskListOptions{
		TaskListID: taskListID,
		Name:       body.Name,
		Archived:   body.Archived,
		Color:      body.Color,
		Position:   body.Position,
	})
	if err != nil {
		return err
	}

	taskList, err := userService.TaskListService.GetTaskList(&entityservices.GetTaskListOptions{
		TaskListID:   taskListID,
		IncludeTasks: false,
	})

	return c.JSON(http.StatusOK, OK(taskListToDTO(taskList)))
}

func (api *APIService) deleteTaskList(c echo.Context) error {
	taskListID := c.Param("id")
	userService := api.mustGetUserService(c)

	err := userService.TaskListService.DeleteTaskList(&entityservices.DeleteTaskListOptions{
		TaskListID: taskListID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) clearTaskList(c echo.Context) error {
	taskListID := c.Param("id")
	userService := api.mustGetUserService(c)

	err := userService.TaskListService.ClearTaskList(&entityservices.ClearTaskListOptions{
		TaskListID: taskListID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
