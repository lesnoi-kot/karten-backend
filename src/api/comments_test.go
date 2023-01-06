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

type commentsSuite struct {
	baseAPITestSuite
}

func TestComments(t *testing.T) {
	suite.Run(t, new(commentsSuite))
}

func (s *commentsSuite) TestAddComment() {
	s.Run("200", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/tasks/123/comments",
			strings.NewReader(`{"text":"KEKW"}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.commentsMock.
			On("Add", mock.Anything).
			Run(func(args mock.Arguments) {
				comment := args.Get(0).(*store.Comment)
				comment.ID = "qqq"
				comment.DateCreated = testDate
			}).
			Return(nil).
			Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "qqq",
				"text": "KEKW",
				"task_id": "123",
				"date_created": "1970-01-01T00:00:00Z",
				"author": "Author"
			}
		}`, rec.Body.String())
	})

	s.Run("Invalid text", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/tasks/123/comments",
			strings.NewReader(`{"text":""}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (s *commentsSuite) TestEditComment() {
	s.Run("Edit full", func() {
		req := httptest.NewRequest(
			http.MethodPatch,
			"/comments/777",
			strings.NewReader(`{
				"text": "qwerty!!!"
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.commentsMock.On("Get", "777").Return(
			&store.Comment{
				ID:          "777",
				Text:        "Test",
				TaskID:      "111",
				Author:      "Admin",
				DateCreated: testDate,
			},
			nil,
		).Once()

		s.commentsMock.On("Update", mock.Anything).Return(nil).Once()

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "777",
				"text": "qwerty!!!",
				"author": "Admin",
				"task_id": "111",
				"date_created": "1970-01-01T00:00:00Z"
			}
		}`, rec.Body.String())
	})

	s.Run("400 - Empty text", func() {
		req := httptest.NewRequest(
			http.MethodPatch,
			"/comments/777",
			strings.NewReader(`{"text": ""}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (s *commentsSuite) TestDeleteComment() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodDelete, "/comments/777", nil)
		rec := httptest.NewRecorder()
		s.commentsMock.On("Delete", "777").Return(nil).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNoContent, rec.Code)
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodDelete, "/comments/777", nil)
		rec := httptest.NewRecorder()
		s.commentsMock.On("Delete", "777").Return(store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}
