package api

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lesnoi-kot/karten-backend/src/modules/images"
	"github.com/lesnoi-kot/karten-backend/src/store"
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

func (api *APIService) uploadImage(c echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	defer file.Close()

	image, err := images.ParseImage(bytes.NewReader(data))
	if err != nil {
		return echo.
			NewHTTPError(http.StatusBadRequest, "Image can't be decoded").
			SetInternal(err)
	}

	dbFile, err := api.store.Files.Add(context.Background(), store.AddFileOptions{
		Name:     fileHeader.Filename,
		Data:     bytes.NewReader(data),
		MIMEType: image.MIMEType,
	})
	if err != nil {
		return fmt.Errorf("Save image error: %w", err)
	}

	api.logger.Debugf("Added image %s", dbFile.ID)
	return c.JSON(http.StatusOK, OK(fileToDTO(dbFile)))
}
