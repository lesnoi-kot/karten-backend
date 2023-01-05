package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func (api *APIService) getProjects(c echo.Context) error {
	projects, err := api.store.Projects.GetAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projects))
}

func (api *APIService) getProject(c echo.Context) error {
	id := c.Param("id")

	project, err := api.store.Projects.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}

		return err
	}

	return c.JSON(http.StatusOK, OK(project))
}

func (api *APIService) addProject(c echo.Context) error {
	var body struct {
		Name string `json:"name" validate:"required,min=1,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	project := &store.Project{Name: body.Name}

	if err := api.store.Projects.Add(project); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(project))
}

func (api *APIService) deleteProject(c echo.Context) error {
	id := c.Param("id")
	err := api.store.Projects.Delete(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) editProject(c echo.Context) error {
	var body struct {
		Name string `json:"name" validate:"required,min=1,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}
	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	id := c.Param("id")
	project, err := api.store.Projects.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}

		return err
	}

	project.Name = body.Name
	if err := api.store.Projects.Update(project); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(project))
}
