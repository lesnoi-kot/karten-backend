package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/entityservices"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type LabelDTO struct {
	ID      int    `json:"id"`
	BoardID string `json:"board_id"`
	UserID  int    `json:"user_id"`
	Name    string `json:"name"`
	Color   int    `json:"color"`
}

type TaskDTO struct {
	ID                  string     `json:"id"`
	UserID              int        `json:"user_id"`
	ShortID             string     `json:"short_id"`
	TaskListID          string     `json:"task_list_id"`
	Name                string     `json:"name"`
	Text                string     `json:"text"`
	HTML                string     `json:"html"`
	Position            int64      `json:"position"`
	SpentTime           int64      `json:"spent_time"`
	Archived            bool       `json:"archived"`
	DateCreated         time.Time  `json:"date_created"`
	DateStartedTracking *time.Time `json:"date_started_tracking"`
	DueDate             *time.Time `json:"due_date"`

	Comments    []*CommentDTO `json:"comments,omitempty"`
	Attachments []*FileDTO    `json:"attachments,omitempty"`
	Labels      []*LabelDTO   `json:"labels,omitempty"`
}

func (api *APIService) getTask(c echo.Context) error {
	taskID := c.Param("id")
	user := api.mustGetUserService(c)
	task, err := user.GetTask(&entityservices.GetTaskOptions{
		TaskID:             taskID,
		IncludeComments:    true,
		IncludeLabels:      true,
		IncludeAttachments: true,
		SkipTextRender:     false,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskToDTO(task)))
}

func (api *APIService) addTask(c echo.Context) error {
	var body struct {
		Name     string     `json:"name" validate:"required,min=1"`
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
	task, err := user.AddTask(&entityservices.AddTaskOptions{
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
	err := user.EditTask(&entityservices.EditTaskOptions{
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

	task, err := user.GetTask(&entityservices.GetTaskOptions{
		TaskID:                taskID,
		IncludeComments:       true,
		IncludeLabels:         true,
		IncludeAttachments:    true,
		SkipCommentTextRender: true,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(taskToDTO(task)))
}

func (api *APIService) deleteTask(c echo.Context) error {
	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.DeleteTask(&entityservices.DeleteTaskOptions{TaskID: taskID})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) addTaskAttachments(c echo.Context) error {
	var body struct {
		FilesID []string `json:"files_id" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.AttachFilesToTask(&entityservices.AttachFilesToTask{
		TaskID:  taskID,
		FilesID: body.FilesID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) deleteTaskAttachment(c echo.Context) error {
	var body struct {
		FileID string `json:"file_id" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskID := c.Param("id")
	_, err := api.store.ORM.NewDelete().
		Model((*store.AttachmentToTaskAssoc)(nil)).
		Where("task_id = ?", taskID).
		Where("file_id = ?", body.FileID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) startTaskTracking(c echo.Context) error {
	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	task, err := user.GetTask(&entityservices.GetTaskOptions{
		TaskID:                taskID,
		SkipTextRender:        true,
		SkipCommentTextRender: true,
	})
	if err != nil {
		return err
	}

	if task.DateStartedTracking != nil {
		return c.NoContent(http.StatusOK)
	}

	now := time.Now()
	err = user.EditTask(&entityservices.EditTaskOptions{
		TaskID:              taskID,
		DateStartedTracking: &now,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *APIService) stopTaskTracking(c echo.Context) error {
	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	task, err := user.GetTask(&entityservices.GetTaskOptions{
		TaskID:                taskID,
		SkipTextRender:        true,
		SkipCommentTextRender: true,
	})
	if err != nil {
		return err
	}

	if task.DateStartedTracking == nil {
		return c.NoContent(http.StatusOK)
	}

	now := time.Now()
	var nullTime time.Time
	timeSpentSec := task.SpentTime + now.Unix() - task.DateStartedTracking.Unix()

	err = user.EditTask(&entityservices.EditTaskOptions{
		TaskID:              taskID,
		DateStartedTracking: &nullTime,
		SpentTime:           &timeSpentSec,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *APIService) addLabelToTask(c echo.Context) error {
	var body struct {
		LabelID int `json:"label_id" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.AddLabelToTask(&entityservices.AddLabelToTaskOptions{
		TaskID:  taskID,
		LabelID: body.LabelID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *APIService) deleteLabelFromTask(c echo.Context) error {
	var body struct {
		LabelID int `json:"label_id" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	taskID := c.Param("id")
	user := api.mustGetUserService(c)

	err := user.DeleteLabelFromTask(&entityservices.AddLabelToTaskOptions{
		TaskID:  taskID,
		LabelID: body.LabelID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
