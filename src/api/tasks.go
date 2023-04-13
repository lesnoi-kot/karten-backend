package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/userservice"
)

type TaskDTO struct {
	ID          string     `json:"id"`
	ShortID     string     `json:"short_id"`
	TaskListID  string     `json:"task_list_id"`
	Name        string     `json:"name"`
	Text        string     `json:"text"`
	Position    int64      `json:"position"`
	Archived    bool       `json:"archived"`
	DateCreated time.Time  `json:"date_created"`
	DueDate     *time.Time `json:"due_date,omitempty"`

	Comments []*CommentDTO `json:"comments,omitempty"`
}

func (api *APIService) getTask(c echo.Context) error {
	taskID := c.Param("id")
	user := api.mustGetUserService(c)
	task, err := user.GetTask(&userservice.GetTaskOptions{
		TaskID:          taskID,
		IncludeComments: true,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskToDTO(task)))
}

func (api *APIService) addTask(c echo.Context) error {
	var body struct {
		Name     string     `json:"name" validate:"required,min=1,max=32"`
		Text     string     `json:"text"`
		Position int64      `json:"position"`
		DueDate  *time.Time `json:"due_date"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskListID := c.Param("id")
	user := api.mustGetUserService(c)
	task, err := user.AddTask(&userservice.AddTaskOptions{
		TaskListID: taskListID,
		Name:       body.Name,
		Text:       body.Text,
		Position:   body.Position,
		DueDate:    body.DueDate,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskToDTO(task)))
}

func (api *APIService) editTask(c echo.Context) error {
	var body struct {
		TaskListID *string    `json:"task_list_id"`
		Name       *string    `json:"name" validate:"omitempty,min=1,max=32"`
		Text       *string    `json:"text"`
		Position   *int64     `json:"position"`
		DueDate    *time.Time `json:"due_date"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	if body.Name != nil {
		*body.Name = strings.TrimSpace(*body.Name)
	}
	if body.Text != nil {
		*body.Text = strings.TrimSpace(*body.Text)
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskID := c.Param("id")
	user := api.mustGetUserService(c)
	err := user.EditTask(&userservice.EditTaskOptions{
		TaskID:     taskID,
		TaskListID: body.TaskListID,
		Name:       body.Name,
		Text:       body.Text,
		Position:   body.Position,
		DueDate:    body.DueDate,
	})
	if err != nil {
		return err
	}

	task, err := user.GetTask(&userservice.GetTaskOptions{
		TaskID:          taskID,
		IncludeComments: false,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskToDTO(task)))
}

func (api *APIService) deleteTask(c echo.Context) error {
	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.DeleteTask(&userservice.DeleteTaskOptions{TaskID: taskID})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
