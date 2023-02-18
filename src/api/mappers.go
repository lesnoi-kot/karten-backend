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
	}

	if project.Avatar != nil {
		dto.AvatarURL = urlprovider.GetFileURL(&project.Avatar.File)

		if len(project.Avatar.Thumbnails) > 0 {
			dto.AvatarThumbnailURL = urlprovider.GetFileURL(&project.Avatar.Thumbnails[0])
		}
	}

	if len(project.Boards) > 0 {
		dto.Boards = lo.Map(project.Boards, func(board *store.Board, index int) *BoardDTO {
			return boardToDTO(board)
		})
	}

	return dto
}

func projectsToDTO(projects []*store.Project) []*ProjectDTO {
	dtos := lo.Map(projects, func(project *store.Project, _ int) *ProjectDTO {
		return projectToDTO(project)
	})

	return dtos
}

func fileToDTO(file *store.File) *FileDTO {
	dto := &FileDTO{
		ID:       file.ID,
		URL:      urlprovider.GetFileURL(file),
		Name:     file.Name,
		MimeType: file.MimeType,
		Size:     file.Size,
	}

	return dto
}

func filesToDTO(files []*store.File) []*FileDTO {
	dtos := lo.Map(files, func(file *store.File, _ int) *FileDTO {
		return fileToDTO(file)
	})

	return dtos
}

func boardToDTO(board *store.Board) *BoardDTO {
	dto := &BoardDTO{
		ID:             board.ID,
		ShortID:        board.ShortID,
		Name:           board.Name,
		ProjectID:      board.ProjectID,
		Archived:       board.Archived,
		Favorite:       board.Favorite,
		DateCreated:    board.DateCreated,
		DateLastViewed: board.DateLastViewed,
		Color:          board.Color,
	}

	if board.Cover != nil {
		dto.CoverURL = urlprovider.GetFileURL(board.Cover)
	}

	if len(board.TaskLists) > 0 {
		dto.TaskLists = lo.Map(board.TaskLists, func(taskList *store.TaskList, index int) *TaskListDTO {
			return taskListToDTO(taskList)
		})
	}

	return dto
}

func taskListToDTO(taskList *store.TaskList) *TaskListDTO {
	dto := &TaskListDTO{
		ID:          taskList.ID,
		BoardID:     taskList.BoardID,
		Position:    taskList.Position,
		Name:        taskList.Name,
		Archived:    taskList.Archived,
		DateCreated: taskList.DateCreated,
		Color:       taskList.Color,
	}

	if len(taskList.Tasks) > 0 {
		dto.Tasks = lo.Map(taskList.Tasks, func(task *store.Task, index int) *TaskDTO {
			return taskToDTO(task)
		})
	}

	return dto
}

func taskToDTO(task *store.Task) *TaskDTO {
	dto := &TaskDTO{
		ID:          task.ID,
		ShortID:     task.ShortID,
		TaskListID:  task.TaskListID,
		Position:    task.Position,
		Name:        task.Name,
		Text:        task.Text,
		Archived:    task.Archived,
		DateCreated: task.DateCreated,
		DueDate:     task.DueDate,
	}

	if len(task.Comments) > 0 {
		dto.Comments = lo.Map(task.Comments, func(comment *store.Comment, index int) *CommentDTO {
			return commentToDTO(comment)
		})
	}

	return dto
}

func commentToDTO(comment *store.Comment) *CommentDTO {
	dto := &CommentDTO{
		ID:          comment.ID,
		TaskID:      comment.TaskID,
		Author:      comment.Author,
		Text:        comment.Text,
		DateCreated: comment.DateCreated,
	}

	return dto
}
