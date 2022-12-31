package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store/models"
)

func (api *APIService) getProjects(c echo.Context) error {
	var projects []models.Project

	if err := api.store.NewSelect().
		Model(&projects).
		Scan(context.Background()); err != nil {
		return err
	}

	res := make([]struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}, len(projects))

	for i, project := range projects {
		res[i].ID = project.ID
		res[i].Name = project.Name
	}

	return c.JSON(http.StatusOK, OK(res))
}

func (api *APIService) getProject(c echo.Context) error {
	id := c.Param("id")
	var project models.Project

	err := api.store.
		NewSelect().
		Model(&project).
		Where("id = ?", id).
		Relation("Boards").
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}

		return err
	}

	return c.JSON(http.StatusOK, OK(project))
}

func (api *APIService) addProject(c echo.Context) error {
	var body struct {
		Name string `json:"name" validate:"required,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}
	if err := c.Validate(&body); err != nil {
		return echo.ErrBadRequest
	}

	project := models.Project{Name: body.Name}

	_, err := api.store.
		NewInsert().
		Model(&project).
		Column("name").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(project))
}

func (api *APIService) deleteProject(c echo.Context) error {
	id := c.Param("id")

	result, err := api.store.
		NewDelete().
		Model((*models.Project)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return echo.ErrNotFound
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) editProject(c echo.Context) error {
	id := c.Param("id")

	var body struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	}

	project := models.Project{Name: body.Name}

	result, err := api.store.
		NewUpdate().
		Model(&project).
		Column("name").
		Where("id = ?", id).
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, OK(project))
}
