package api

import (
	"github.com/samber/lo"

	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/lesnoi-kot/karten-backend/src/urlprovider"
)

func projectToDTO(project *store.Project) *ProjectDTO {
	dto := &ProjectDTO{
		ID:      project.ID,
		UserID:  int(project.UserID),
		ShortID: project.ShortID,
		Name:    project.Name,
	}

	if project.Avatar != nil {
		dto.AvatarURL = urlprovider.GetFileURL(&project.Avatar.File)

		if len(project.Avatar.Thumbnails) > 0 {
			dto.AvatarThumbnailURL = urlprovider.GetFileURL(project.Avatar.Thumbnails[0])
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

func imageFileToDTO(file *store.ImageFile) *ImageFileDTO {
	dto := &ImageFileDTO{
		FileDTO: FileDTO{
			ID:       file.ID,
			URL:      urlprovider.GetFileURL(&file.File),
			Name:     file.Name,
			MimeType: file.MimeType,
			Size:     file.Size,
		},
		Thumbnails: filesToDTO(file.Thumbnails),
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
		UserID:         int(board.UserID),
		ShortID:        board.ShortID,
		Name:           board.Name,
		ProjectID:      board.ProjectID,
		Archived:       board.Archived,
		Favorite:       board.Favorite,
		DateCreated:    board.DateCreated,
		DateLastViewed: board.DateLastViewed,
		Color:          board.Color,
	}

	if board.Project != nil {
		dto.ProjectName = board.Project.Name
	}

	if board.Cover != nil {
		dto.CoverURL = urlprovider.GetFileURL(board.Cover)
	}

	if len(board.TaskLists) > 0 {
		dto.TaskLists = lo.Map(board.TaskLists, func(taskList *store.TaskList, index int) *TaskListDTO {
			return taskListToDTO(taskList)
		})
	}

	if len(board.Labels) > 0 {
		dto.Labels = lo.Map(board.Labels, func(label *store.Label, index int) *LabelDTO {
			return labelToDTO(label)
		})
	}

	return dto
}

func taskListToDTO(taskList *store.TaskList) *TaskListDTO {
	dto := &TaskListDTO{
		ID:          taskList.ID,
		BoardID:     taskList.BoardID,
		UserID:      int(taskList.UserID),
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

func labelToDTO(label *store.Label) *LabelDTO {
	return &LabelDTO{
		ID:      label.ID,
		BoardID: label.BoardID,
		UserID:  int(label.UserID),
		Name:    label.Name,
		Color:   label.Color,
	}
}

func taskToDTO(task *store.Task) *TaskDTO {
	dto := &TaskDTO{
		ID:                  task.ID,
		UserID:              int(task.UserID),
		ShortID:             task.ShortID,
		TaskListID:          task.TaskListID,
		Position:            task.Position,
		SpentTime:           task.SpentTime,
		Name:                task.Name,
		Text:                task.Text,
		HTML:                task.HTML,
		Archived:            task.Archived,
		DateCreated:         task.DateCreated,
		DateStartedTracking: task.DateStartedTracking,
		DueDate:             task.DueDate,
	}

	if len(task.Comments) > 0 {
		dto.Comments = lo.Map(task.Comments, func(comment *store.Comment, index int) *CommentDTO {
			return commentToDTO(comment)
		})
	}

	if len(task.Attachments) > 0 {
		dto.Attachments = lo.Map(task.Attachments, func(file *store.File, index int) *FileDTO {
			return fileToDTO(file)
		})
	}

	if len(task.Labels) > 0 {
		dto.Labels = lo.Map(task.Labels, func(file *store.Label, index int) *LabelDTO {
			return labelToDTO(file)
		})
	}

	return dto
}

func commentToDTO(comment *store.Comment) *CommentDTO {
	dto := &CommentDTO{
		ID:          comment.ID,
		TaskID:      comment.TaskID,
		UserID:      int(comment.UserID),
		Text:        comment.Text,
		HTML:        comment.HTML,
		DateCreated: comment.DateCreated,
	}

	if len(comment.Attachments) > 0 {
		dto.Attachments = lo.Map(comment.Attachments, func(file *store.File, index int) *FileDTO {
			return fileToDTO(file)
		})
	}

	if comment.Author != nil {
		dto.Author = publicUserToDTO(comment.Author)
	}

	return dto
}

func userToDTO(user *store.User) *UserDTO {
	dto := &UserDTO{
		ID:          int(user.ID),
		SocialID:    user.SocialID,
		Name:        user.Name,
		Login:       user.Login,
		Email:       user.Email,
		URL:         user.URL,
		AvatarURL:   urlprovider.GetFileURL(user.Avatar),
		DateCreated: user.DateCreated,
	}

	return dto
}

func publicUserToDTO(user *store.User) *PublicUserDTO {
	dto := &PublicUserDTO{
		ID:          int(user.ID),
		Name:        user.Name,
		AvatarURL:   urlprovider.GetFileURL(user.Avatar),
		DateCreated: user.DateCreated,
	}

	return dto
}
