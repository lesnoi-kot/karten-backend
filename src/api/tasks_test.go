package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type tasksSuite struct {
	baseAPITestSuite
}

func TestTasks(t *testing.T) {
	suite.Run(t, new(tasksSuite))
}

func (s *tasksSuite) TestGetTask() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodGet, "/tasks/123", nil)
		rec := httptest.NewRecorder()
		s.tasksMock.On("Get", "123").Return(
			&store.Task{
				ID:          "123",
				TaskListID:  "111",
				Name:        "Test",
				Text:        "",
				Position:    333,
				DateCreated: testDate,
				DueDate:     &testDate,
			},
			nil,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": {
				"id": "123",
				"task_list_id": "111",
				"name": "Test",
				"text": "",
				"archived": false,
				"position": 333,
				"date_created": "1970-01-01T00:00:00Z",
				"due_date": "1970-01-01T00:00:00Z"
			}
		}`, rec.Body.String())
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodGet, "/tasks/123", nil)
		rec := httptest.NewRecorder()
		s.tasksMock.On("Get", "123").Return((*store.Task)(nil), store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func (s *tasksSuite) TestAddTask() {
	s.Run("200", func() {
		req := httptest.NewRequest(
			http.MethodPost, "/task-lists/123/tasks",
			strings.NewReader(`{"name": "qqq", "text": "looool", "position": 111, "due_date": "1970-01-01T00:00:00Z"}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.tasksMock.
			On("Add", mock.Anything).
			Run(func(args mock.Arguments) {
				task := args.Get(0).(*store.Task)
				task.ID = "555"
				task.DateCreated = testDate
			}).
			Return(nil).
			Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": {
				"id": "555",
				"task_list_id": "123",
				"name": "qqq",
				"text": "looool",
				"archived": false,
				"position": 111,
				"date_created": "1970-01-01T00:00:00Z",
				"due_date": "1970-01-01T00:00:00Z"
			}
		}`, rec.Body.String())
	})

	s.Run("400 - Invalid name", func() {
		req := httptest.NewRequest(
			http.MethodPost, "/task-lists/123/tasks",
			strings.NewReader(`{"name": "  ", "text": "looool", "position": 111, "due_date": "1970-01-01T00:00:00Z"}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})

	s.Run("400 - Invalid due_date", func() {
		req := httptest.NewRequest(
			http.MethodPost, "/task-lists/123/tasks",
			strings.NewReader(`{"name": "q", "text": "looool", "position": 111, "due_date": "1970-01-01Thrrhsrhy"}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})

	s.Run("500 - Store error", func() {
		req := httptest.NewRequest(
			http.MethodPost, "/task-lists/123/tasks",
			strings.NewReader(`{"name": "q", "text": "looool", "position": 111}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.tasksMock.
			On("Add", mock.Anything).
			Return(errors.New("some error")).
			Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusInternalServerError, rec.Code)
	})
}

func (s *tasksSuite) TestEditTask() {
	s.Run("200", func() {
		req := httptest.NewRequest(
			http.MethodPatch, "/tasks/345",
			strings.NewReader(`{
				"name": "qqq",
				"text": "looool",
				"position": 111,
				"due_date": "1970-01-01T00:00:00Z",
				"task_list_id": "777"
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.tasksMock.
			On("Update", mock.Anything).
			Return(nil).
			Once()
		s.tasksMock.On("Get", "345").Return(
			&store.Task{
				ID:          "345",
				TaskListID:  "222",
				Name:        "Test",
				Text:        "",
				Position:    333,
				DateCreated: testDate,
				DueDate:     &testDate,
			},
			nil,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": {
				"id": "345",
				"archived": false,
				"task_list_id": "777",
				"name": "qqq",
				"text": "looool",
				"position": 111,
				"date_created": "1970-01-01T00:00:00Z",
				"due_date": "1970-01-01T00:00:00Z"
			}
		}`, rec.Body.String())
	})
}

func (s *tasksSuite) TestDeleteTask() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodDelete, "/tasks/123", nil)
		rec := httptest.NewRecorder()
		s.tasksMock.On("Delete", "123").Return(nil).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNoContent, rec.Code)
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodDelete, "/tasks/123", nil)
		rec := httptest.NewRecorder()
		s.tasksMock.On("Delete", "123").Return(store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}
