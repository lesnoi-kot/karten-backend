package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/modules/images"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

type ProjectDTO struct {
	ID                 string      `json:"id"`
	ShortID            string      `json:"short_id"`
	Name               string      `json:"name"`
	AvatarURL          string      `json:"avatar_url,omitempty"`
	AvatarThumbnailURL string      `json:"avatar_thumbnail_url,omitempty"`
	Boards             []*BoardDTO `json:"boards,omitempty"`
}

func (api *APIService) getProjects(c echo.Context) error {
	projects, err := api.store.Projects.GetAll(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projectsToDTO(projects)))
}

func (api *APIService) getProject(c echo.Context) error {
	id := c.Param("id")

	project, err := api.store.Projects.Get(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projectToDTO(project)))
}

func (api *APIService) addProject(c echo.Context) error {
	var body struct {
		Name string `json:"name" form:"name" validate:"required,min=1,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	body.Name = strings.TrimSpace(body.Name)
	if err := c.Validate(&body); err != nil {
		return err
	}

	project := &store.Project{Name: body.Name}

	if avatarFileHeader, err := c.FormFile("avatar"); err == nil {
		avatarData, err := avatarFileHeader.Open()
		if err != nil {
			return err
		}

		image, err := images.ParseImage(avatarData)
		if err != nil {
			return echo.
				NewHTTPError(http.StatusBadRequest, "Image can't be decoded").
				SetInternal(err)
		}

		avatarData, err = avatarFileHeader.Open()
		if err != nil {
			return err
		}

		avatarFile, err := api.store.Files.Add(context.Background(), store.AddFileOptions{
			Name:     avatarFileHeader.Filename,
			Data:     avatarData,
			MIMEType: image.MIMEType,
		})
		if err != nil {
			return fmt.Errorf("Save image for new project error: %w", err)
		}

		api.logger.Infof("Added avatar for project %s: %s", project.Name, avatarFile.ID)

		avatarData, err = avatarFileHeader.Open()
		if err != nil {
			return err
		}

		thumbnailImageData, err := images.MakeProjectAvatarThumbnail(avatarData)
		if err != nil {
			return err
		}

		thumbnailFile, err := api.store.Files.AddImageThumbnail(context.Background(), store.AddImageThumbnailOptions{
			AddFileOptions: store.AddFileOptions{
				Name:     avatarFileHeader.Filename,
				Data:     thumbnailImageData,
				MIMEType: image.MIMEType,
			},
			OriginalImageFileID: avatarFile.ID,
		})
		if err != nil {
			return fmt.Errorf("Save image for new project error: %w", err)
		}

		api.logger.Infof("Added avatar thumbnail for project %s: %s", project.Name, thumbnailFile.ID)

		project.AvatarID = avatarFile.ID
		project.Avatar = &store.ImageFile{
			File:       *avatarFile,
			Thumbnails: []store.File{*thumbnailFile},
		}
	}

	if err := api.store.Projects.Add(context.Background(), project); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(projectToDTO(project)))
}

func (api *APIService) deleteProject(c echo.Context) error {
	id := c.Param("id")

	if err := api.store.Projects.Delete(context.Background(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) deleteProjects(c echo.Context) error {
	if err := api.store.Projects.DeleteAll(context.Background()); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) clearProject(c echo.Context) error {
	id := c.Param("id")

	if err := api.store.Projects.Clear(context.Background(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *APIService) editProject(c echo.Context) error {
	var body struct {
		Name string `json:"name" validate:"required,min=1,max=32"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	id := c.Param("id")
	project, err := api.store.Projects.Get(context.Background(), id)
	if err != nil {
		return err
	}

	project.Name = body.Name
	if err := api.store.Projects.Update(context.Background(), project); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, OK(project))
}
