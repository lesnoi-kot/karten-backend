package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/api"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var testDate = time.Unix(0, 0).UTC()

type mockBoardsStore struct {
	mock.Mock
}

type boardsSuite struct {
	suite.Suite
	store      *store.Store
	api        *api.APIService
	boardsMock *mockBoardsStore
}

func (suite *boardsSuite) SetupTest() {
	suite.boardsMock = new(mockBoardsStore)
	suite.store = &store.Store{
		DB:       nil,
		Projects: nil,
		Boards:   suite.boardsMock,
	}
	suite.api = api.NewAPI(api.APIConfig{
		Store:     suite.store,
		Logger:    zap.NewNop().Sugar(),
		APIPrefix: "",
	})
}

func (s *boardsSuite) TestGetBoard() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodGet, "/boards/123", nil)
		rec := httptest.NewRecorder()
		s.boardsMock.On("Get", "123").Return(
			&store.Board{
				ID:             "123",
				Name:           "Test",
				ProjectID:      "111",
				Archived:       false,
				DateCreated:    testDate,
				DateLastViewed: testDate,
				Color:          0,
				CoverURL:       "",
			},
			nil,
		).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "123",
				"name": "Test",
				"project_id": "111",
				"archived": false,
				"date_created": "1970-01-01T00:00:00Z",
				"date_last_viewed": "1970-01-01T00:00:00Z",
				"color": 0,
				"cover_url": ""
			}
		}`, rec.Body.String())
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodGet, "/boards/123", nil)
		rec := httptest.NewRecorder()
		s.boardsMock.On("Get", "123").Return((*store.Board)(nil), store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func (s *boardsSuite) TestAddBoard() {
	s.Run("200", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/projects/project-id/boards",
			strings.NewReader(`{
				"name": "Apple",
				"color": 1,
				"cover_url": "www"
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.boardsMock.
			On("Add", mock.Anything).
			Run(func(args mock.Arguments) {
				board := args.Get(0).(*store.Board)
				board.ID = "123"
				board.DateCreated = testDate
				board.DateLastViewed = testDate
			}).
			Return(nil).
			Once()

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "123",
				"name": "Apple",
				"project_id": "project-id",
				"archived": false,
				"date_created": "1970-01-01T00:00:00Z",
				"date_last_viewed": "1970-01-01T00:00:00Z",
				"color": 1,
				"cover_url": "www"
			}
		}`, rec.Body.String())
	})

	s.Run("Invalid json", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/projects/project-id/boards",
			strings.NewReader(`{
				"name": "Apple",
				"color": 1,
				"cover_url": "www",
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})

	s.Run("Invalid name", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/projects/project-id/boards",
			strings.NewReader(`{
				"name": "",
				"color": 1,
				"cover_url": "www"
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (s *boardsSuite) TestEditBoard() {
	s.Run("Edit name", func() {
		req := httptest.NewRequest(
			http.MethodPatch,
			"/boards/777",
			strings.NewReader(`{
				"name": "Pineapple"
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.boardsMock.On("Get", "777").Return(
			&store.Board{
				ID:             "777",
				Name:           "Test",
				ProjectID:      "111",
				Archived:       false,
				DateCreated:    testDate,
				DateLastViewed: testDate,
				Color:          3,
				CoverURL:       "ru",
			},
			nil,
		).Once()

		s.boardsMock.On("Edit", mock.Anything).Return(nil).Once()

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "777",
				"name": "Pineapple",
				"project_id": "111",
				"archived": false,
				"date_created": "1970-01-01T00:00:00Z",
				"date_last_viewed": "1970-01-01T00:00:00Z",
				"color": 3,
				"cover_url": "ru"
			}
		}`, rec.Body.String())
	})

	s.Run("Edit full", func() {
		req := httptest.NewRequest(
			http.MethodPatch,
			"/boards/777",
			strings.NewReader(`{
				"name": "qwerty",
				"archived": true,
				"color": 123,
				"cover_url": "com"
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		s.boardsMock.On("Get", "777").Return(
			&store.Board{
				ID:             "777",
				Name:           "Test",
				ProjectID:      "111",
				Archived:       false,
				DateCreated:    testDate,
				DateLastViewed: testDate,
				Color:          3,
				CoverURL:       "ru",
			},
			nil,
		).Once()

		s.boardsMock.On("Edit", mock.Anything).Return(nil).Once()

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"error": null,
			"data": {
				"id": "777",
				"name": "qwerty",
				"project_id": "111",
				"archived": true,
				"date_created": "1970-01-01T00:00:00Z",
				"date_last_viewed": "1970-01-01T00:00:00Z",
				"color": 123,
				"cover_url": "com"
			}
		}`, rec.Body.String())
	})

	s.Run("Invalid name", func() {
		req := httptest.NewRequest(
			http.MethodPatch,
			"/boards/777",
			strings.NewReader(`{
				"name": ""
			}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (s *boardsSuite) TestDeleteBoard() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodDelete, "/boards/777", nil)
		rec := httptest.NewRecorder()
		s.boardsMock.On("Delete", "777").Return(nil).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNoContent, rec.Code)
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodDelete, "/boards/777", nil)
		rec := httptest.NewRecorder()
		s.boardsMock.On("Delete", "777").Return(store.ErrNotFound).Once()
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func TestBoards(t *testing.T) {
	suite.Run(t, new(boardsSuite))
}

func (m mockBoardsStore) Get(id string) (*store.Board, error) {
	args := m.Called(id)
	return args.Get(0).(*store.Board), args.Error(1)
}

func (m mockBoardsStore) Add(a *store.Board) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m mockBoardsStore) Edit(a *store.Board) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m mockBoardsStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
