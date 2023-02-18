package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FileDTO struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
}

func (api *APIService) getCoverImages(c echo.Context) error {
	files, err := api.store.Files.GetDefaultCovers(context.Background())
	if err != nil {
		return err
	}

	filesDTO := make([]FileDTO, len(files))

	for i := 0; i < len(files); i++ {
		filesDTO[i] = *fileToDTO(&files[i].File)
	}

	return c.JSON(http.StatusOK, OK(filesDTO))
}
