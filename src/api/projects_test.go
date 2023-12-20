package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/api"
)

func (s *APITestSuite) TestGetProjects() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects?order_by=date", nil)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": [
				{
					"id": "1f894df2-f233-4885-81ef-e21aee62e2cd",
					"name": "Science pals",
					"user_id": 33,
					"short_id": "e21aee62e2cd"
				},
				{
					"id": "2d2712eb-266d-4626-b017-697a67907e28",
					"name": "Just Do It vol. 2023",
					"user_id": 33,
					"short_id": "697a67907e28"
				},
				{
					"id": "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					"name": "TO BE DELETED",
					"user_id": 33,
					"short_id": "34adb39a6701"
				},
				{
					"id": "2f146153-ee2f-4968-a241-11a4f00bf212",
					"name": "Business",
					"user_id": 33,
					"short_id": "11a4f00bf212"
				}
			]
		}`, rec.Body.String())
	})
}

func (s *APITestSuite) TestGetProject() {
	s.Run("200", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects/2f146153-ee2f-4968-a241-11a4f00bf212", nil)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)
		s.JSONEq(`{
			"data": {
				"id": "2f146153-ee2f-4968-a241-11a4f00bf212",
				"name": "Business",
				"user_id": 33,
				"short_id": "11a4f00bf212",
				"boards": [
					{
						"id": "36461fd5-2eb4-42b9-a921-b4428a448cfa",
						"user_id": 33,
						"short_id": "b4428a448cfa",
						"name": "Red board",
						"project_id": "2f146153-ee2f-4968-a241-11a4f00bf212",
						"archived": false,
						"favorite": true,
						"color": 3359829,
						"cover_url": "http://nginx:80/9878f395-9987-42f4-8a49-282812612d5e",
						"project_name": "Business",
						"date_created": "0001-01-01T00:00:00Z",
						"date_last_viewed": "0001-01-01T00:00:00Z"
					}
				]
			}
		}`, rec.Body.String())
	})
	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects/ffffffff-ee2f-4968-a241-11a4f00bf212", nil)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
	s.Run("Invalid path", func() {
		req := httptest.NewRequest(http.MethodGet, "/projects/gfgd/fdgfgdf/g", nil)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func (s *APITestSuite) TestAddProjectValidation() {
	s.Run("Invalid input", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/projects",
			strings.NewReader(`{"name":""}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (s *APITestSuite) TestAddAndDeleteProject() {
	projectID := ""

	s.Run("Should add new project", func() {
		req := httptest.NewRequest(
			http.MethodPost,
			"/projects",
			strings.NewReader(`{"name":"Integration-test"}`),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)
		s.api.Server().Handler.ServeHTTP(rec, req)

		s.Equal(http.StatusOK, rec.Code)

		var response api.Response[api.ProjectDTO]
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		s.Require().NoError(err)
		s.Equal("Integration-test", response.Data.Name)
		projectID = response.Data.ID
	})

	s.Run("Should delete this new project", func() {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/projects/%s", projectID), nil)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNoContent, rec.Code)
	})

	s.Run("Forbidden", func() {
		req := httptest.NewRequest(http.MethodDelete, "/projects/6d8065f4-032f-494d-b124-b7bdf3c6480b", nil)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}

func (s *APITestSuite) TestEditProject() {
	s.Run("Change some project name", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/1f894df2-f233-4885-81ef-e21aee62e2cd", strings.NewReader(`{"name":"New"}`))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)

		var response api.Response[api.ProjectDTO]
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		s.Require().NoError(err)
		s.Equal("New", response.Data.Name)
		s.Equal("1f894df2-f233-4885-81ef-e21aee62e2cd", response.Data.ID)
	})

	s.Run("Unchange it back", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/1f894df2-f233-4885-81ef-e21aee62e2cd", strings.NewReader(`{"name":"Science pals"}`))
		rec := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)

		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusOK, rec.Code)

		var response api.Response[api.ProjectDTO]
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		s.Require().NoError(err)
		s.Equal("Science pals", response.Data.Name)
		s.Equal("1f894df2-f233-4885-81ef-e21aee62e2cd", response.Data.ID)
	})

	s.Run("400", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/2f146153-ee2f-4968-a241-11a4f00bf212", strings.NewReader(`{"name":""}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})

	s.Run("400", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/111", strings.NewReader(`{"name":"XXX"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusBadRequest, rec.Code)
	})

	s.Run("404", func() {
		req := httptest.NewRequest(http.MethodPatch, "/projects/6d8065f4-032f-494d-b124-b7bdf3c6480b", strings.NewReader(`{"name":"x"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		s.authorizeRequest(req, testUser.ID)
		s.setCSRF(req)
		s.api.Server().Handler.ServeHTTP(rec, req)
		s.Equal(http.StatusNotFound, rec.Code)
	})
}
