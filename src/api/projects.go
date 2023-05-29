package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/lesnoi-kot/karten-backend/src/userservice"
	"github.com/samber/lo"
)

type ProjectDTO struct {
	ID                 string      `json:"id"`
	UserID             int         `json:"user_id"`
	ShortID            string      `json:"short_id"`
	Name               string      `json:"name"`
	AvatarURL          string      `json:"avatar_url,omitempty"`
	AvatarThumbnailURL string      `json:"avatar_thumbnail_url,omitempty"`
	Boards             []*BoardDTO `json:"boards,omitempty"`
}

func (api *APIService) getProjects(c echo.Context) error {
	includes := c.QueryParams()["include"]
	userService := api.mustGetUserService(c)
	projects, err := userService.GetProjects(&userservice.GetProjectsOptions{
		IncludeBoards: lo.Contains(includes, "boards"),
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projectsToDTO(projects)))
}

func (api *APIService) getProject(c echo.Context) error {
	projectID := c.Param("id")
	userService := api.mustGetUserService(c)

	project, err := userService.GetProject(&userservice.GetProjectOptions{
		ProjectID:     projectID,
		IncludeBoards: true,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projectToDTO(project)))
}

func (api *APIService) addProject(c echo.Context) error {
	var body struct {
		Name     string        `json:"name" validate:"required,min=1,max=32"`
		AvatarID *store.FileID `json:"avatar_id"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	userService := api.mustGetUserService(c)
	project, err := userService.AddProject(&userservice.AddProjectOptions{
		Name:     body.Name,
		AvatarID: body.AvatarID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projectToDTO(project)))
}

func (api *APIService) deleteProject(c echo.Context) error {
	projectID := c.Param("id")
	userService := api.mustGetUserService(c)

	err := userService.DeleteProject(&userservice.DeleteProjectOptions{
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) deleteProjects(c echo.Context) error {
	if err := api.mustGetUserService(c).DeleteAllProjects(); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) clearProject(c echo.Context) error {
	projectID := c.Param("id")
	if err := api.mustGetUserService(c).ClearProject(projectID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) editProject(c echo.Context) error {
	var body struct {
		Name     *string       `json:"name" validate:"omitempty,min=1,max=32"`
		AvatarID *store.FileID `json:"avatar_id"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	projectID := c.Param("id")
	userService := api.mustGetUserService(c)

	err := userService.EditProject(&userservice.EditProjectOptions{
		ProjectID: projectID,
		Name:      body.Name,
		AvatarID:  body.AvatarID,
	})
	if err != nil {
		return err
	}

	project, err := userService.GetProject(&userservice.GetProjectOptions{
		ProjectID:     projectID,
		IncludeBoards: false,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projectToDTO(project)))
}
