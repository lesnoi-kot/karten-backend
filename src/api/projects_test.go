package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/lesnoi-kot/karten-backend/src/store"
)

type projectsSuite struct {
	baseAPITestSuite
}

func TestProjects(t *testing.T) {
	suite.Run(t, new(projectsSuite))
}

func (s *projectsSuite) TestGetProjects() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		rec := httptest.NewRecorder()

		s.projectsMock.On("GetAll").Return(
			[]*store.Project{
				{ID: "1", Name: "A"},
				{ID: "2", Name: "B"},
				{ID: "3", Name: "C"},
			},
			nil,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": [
				{"id": "1", "name": "A"},
				{"id": "2", "name": "B"},
				{"id": "3", "name": "C"}
			]
		}`, rec.Body.String())
	})

	s.Run("Error from the store", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		rec := httptest.NewRecorder()
		s.projectsMock.On("GetAll").Return([]*store.Project(nil), errors.New("BOOM")).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusInternalServerError, rec.Code)
	})
}

func (s *projectsSuite) TestGetProject() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects/id-1", nil)
		rec := httptest.NewRecorder()

		s.projectsMock.On("Get", "id-1").Return(
			&store.Project{ID: "id-1", Name: "First"},
			nil,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": {"id": "id-1", "name": "First"}
		}`, rec.Body.String())
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects/id-not-found", nil)
		rec := httptest.NewRecorder()
		s.projectsMock.On("Get", "id-not-found").Return(
			(*store.Project)(nil),
			store.ErrNotFound,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusNotFound, rec.Code)
	})

	s.Run("Invalid path", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects/gfgd/fdgfgdf/g", nil)
		rec := httptest.NewRecorder()
		s.projectsMock.On("Get").Return(
			(*store.Project)(nil),
			store.ErrNotFound,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func (s *projectsSuite) TestAddProject() {
	s.Run("Valid input", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/projects",
			strings.NewReader(`{"name":"Unit-test"}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.projectsMock.
			On("Add", mock.Anything).
			Run(func(args mock.Arguments) {
				project := args.Get(0).(*store.Project)
				project.ID = "new-id"
			}).
			Return(nil).
			Once()

		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": {"id": "new-id", "name": "Unit-test"}
		}`, rec.Body.String())
	})

	s.Run("Invalid input", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/projects",
			strings.NewReader(`{"name":""}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (s *projectsSuite) TestDeleteProject() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodDelete, "/projects/111", nil)
		rec := httptest.NewRecorder()
		s.projectsMock.On("Delete", "111").Return(nil).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNoContent, rec.Code)
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodDelete, "/projects/666", nil)
		rec := httptest.NewRecorder()
		s.projectsMock.On("Delete", "666").Return(store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func (s *projectsSuite) TestEditProject() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/111", strings.NewReader(`{"name":"New"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.projectsMock.
			On("Get", "111").
			Return(&store.Project{ID: "111", Name: "KEKW"}, nil).
			Once()

		s.projectsMock.
			On("Update", mock.Anything).
			Return(nil).
			Once()

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": {"id": "111", "name": "New"}
		}`, rec.Body.String())
	})

	s.Run("400", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/111", strings.NewReader(`{"name":""}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/111", strings.NewReader(`{"name":"x"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.projectsMock.
			On("Get", "111").
			Return((*store.Project)(nil), store.ErrNotFound).
			Once()

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}
