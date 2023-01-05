package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type taskListsSuite struct {
	baseAPITestSuite
}

func TestTaskLists(t *testing.T) {
	suite.Run(t, new(taskListsSuite))
}

func (s *taskListsSuite) TestGetTaskList() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodGet, "/task-lists/111", nil)
		rec := httptest.NewRecorder()
		s.taskListsMock.On("Get", "111").Return(
			&store.TaskList{
				ID:          "111",
				BoardID:     "board-id",
				Name:        "Test",
				Archived:    false,
				Position:    666,
				DateCreated: testDate,
				Color:       3,
			},
			nil,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "111",
				"board_id": "board-id",
				"name": "Test",
				"archived": false,
				"position": 666,
				"date_created": "1970-01-01T00:00:00Z",
				"color": 3
			}
		}`, rec.Body.String())
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodGet, "/task-lists/111", nil)
		rec := httptest.NewRecorder()
		s.taskListsMock.On("Get", "111").Return((*store.TaskList)(nil), store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func (s *taskListsSuite) TestAddTaskList() {
	s.Run("200", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/boards/board-id/task-lists",
			strings.NewReader(`{"name": "Test", "position": 666, "color": 3}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.taskListsMock.
			On("Add", mock.Anything).
			Run(func(args mock.Arguments) {
				taskList := args.Get(0).(*store.TaskList)
				taskList.ID = "777"
				taskList.DateCreated = testDate
			}).
			Return(nil).
			Once()

		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "777",
				"board_id": "board-id",
				"name": "Test",
				"archived": false,
				"position": 666,
				"date_created": "1970-01-01T00:00:00Z",
				"color": 3
			}
		}`, rec.Body.String())
	})

	s.Run("Invalid name", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/boards/board-id/task-lists",
			strings.NewReader(`{"name": ""}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (s *taskListsSuite) TestEditTaskList() {
	s.Run("Edit name", func() {
		req := httptest.NewRequest(
			http.MethodPatch,
			"/task-lists/888",
			strings.NewReader(`{
				"name": "Pineapple",
				"position": 123,
				"archived": true,
				"color": 5
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.taskListsMock.On("Get", "888").Return(
			&store.TaskList{
				ID:          "888",
				BoardID:     "board-id",
				Name:        "Test",
				Archived:    false,
				Position:    666,
				DateCreated: testDate,
				Color:       3,
			},
			nil,
		).Once()
		s.taskListsMock.On("Update", mock.Anything).Return(nil).Once()

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "888",
				"board_id": "board-id",
				"name": "Pineapple",
				"archived": true,
				"position": 123,
				"date_created": "1970-01-01T00:00:00Z",
				"color": 5
			}
		}`, rec.Body.String())
	})
}

func (s *taskListsSuite) TestDeleteTaskList() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodDelete, "/task-lists/777", nil)
		rec := httptest.NewRecorder()
		s.taskListsMock.On("Delete", "777").Return(nil).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNoContent, rec.Code)
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodDelete, "/task-lists/777", nil)
		rec := httptest.NewRecorder()
		s.taskListsMock.On("Delete", "777").Return(store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}
