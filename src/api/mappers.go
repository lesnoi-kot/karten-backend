package api

import (
	"github.com/samber/lo"

	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/lesnoi-kot/karten-backend/src/urlprovider"
)

func projectToDTO(project *store.Project) *ProjectDTO {
	dto := &ProjectDTO{
		ID:      project.ID,
		ShortID: project.ShortID,
		Name:    project.Name,
		Boards:  project.Boards,
	}

	if project.Avatar != nil {
		dto.AvatarURL = urlprovider.GetFileURL(&project.Avatar.File)

		if len(project.Avatar.Thumbnails) > 0 {
			dto.AvatarThumbnailURL = urlprovider.GetFileURL(&project.Avatar.Thumbnails[0])
		}
	}

	return dto
}

func projectsToDTO(projects []*store.Project) []*ProjectDTO {
	dtos := lo.Map(projects, func(project *store.Project, _ int) *ProjectDTO {
		return projectToDTO(project)
	})

	return dtos
}
