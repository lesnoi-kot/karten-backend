package store_test

import (
	"os"
	"testing"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var testDate = time.Unix(0, 0).UTC()

type storeSuite struct {
	suite.Suite

	store *store.Store
}

func TestStore(t *testing.T) {
	if os.Getenv("STORE_DSN") == "" {
		t.Skip("Store intergation tests are skipped!")
	}

	suite.Run(t, new(storeSuite))
}

func (s *storeSuite) SetupTest() {
	var err error

	s.store, err = store.NewStore(store.StoreConfig{
		DSN:    os.Getenv("STORE_DSN"),
		Logger: zap.NewNop().Sugar(),
		Debug:  false,
	})

	s.Require().NoError(err)
}

func (s *storeSuite) TestProjectsStore() {
	s.Run("GetAll", func() {
		all, err := s.store.Projects.GetAll()

		s.Require().NoError(err)

		expected := []*store.Project{
			{
				ID:   "fd5f451d-fac6-4bc7-a677-34adb39a6701",
				Name: "Personal",
			},
			{
				ID:   "2f146153-ee2f-4968-a241-11a4f00bf212",
				Name: "Business",
			},
			{
				ID:     "2d2712eb-266d-4626-b017-697a67907e28",
				Name:   "Just Do It vol. 2023",
				Boards: nil,
			},
			{
				ID:     "1f894df2-f233-4885-81ef-e21aee62e2cd",
				Name:   "Science pals",
				Boards: nil,
			},
		}

		s.Equal(expected, all)
	})

	s.Run("Get", func() {
		project, err := s.store.Projects.Get("fd5f451d-fac6-4bc7-a677-34adb39a6701")
		s.Require().NoError(err)

		expected := &store.Project{
			ID:   "fd5f451d-fac6-4bc7-a677-34adb39a6701",
			Name: "Personal",
			Boards: []*store.Board{
				{
					ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					ID:             "29e247c3-69f1-4397-8bab-b1dd10ae28b2",
					Name:           "Pet projects",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          0,
					CoverURL:       "",
				},
				{
					ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					ID:             "f3fc69f2-27aa-4aed-842e-9ed544661bfd",
					Name:           "Health",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          1,
					CoverURL:       "",
				},
				{
					ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					ID:             "ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a",
					Name:           "TODO Books",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          0,
					CoverURL:       "",
				},
				{
					ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					ID:             "606ecfd6-2a49-4cc2-911c-a0113ebcf0e6",
					Name:           "Guitar stuff",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          0,
					CoverURL:       "",
				},
			},
		}

		s.Equal(expected, project)
	})

	s.Run("Get", func() {
		project, err := s.store.Projects.Get("2f146153-ee2f-4968-a241-11a4f00bf212")
		s.Require().NoError(err)

		expected := &store.Project{
			ID:   "2f146153-ee2f-4968-a241-11a4f00bf212",
			Name: "Business",
			Boards: []*store.Board{
				{
					ProjectID:      "2f146153-ee2f-4968-a241-11a4f00bf212",
					ID:             "311fef19-eb2a-4c04-98d7-3653d2271293",
					Name:           "Ideas",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          0,
					CoverURL:       "",
				},
				{
					ProjectID:      "2f146153-ee2f-4968-a241-11a4f00bf212",
					ID:             "250bf7a7-ad51-4d62-85cd-c554e4d5f686",
					Name:           "Householding",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          2,
					CoverURL:       "",
				},
			},
		}

		s.Equal(expected, project)
	})

	s.Run("Get - not found", func() {
		actual, err := s.store.Projects.Get("0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		project := &store.Project{Name: "New test"}
		err := s.store.Projects.Add(project)

		s.Require().NoError(err)
		s.Equal("New test", project.Name)
		s.NotEmpty(project.ID)
		s.Nil(project.Boards)

		{
			all, err := s.store.Projects.GetAll()
			s.Require().NoError(err)
			s.Len(all, 5)
		}

		project.Name = "Updated"
		err = s.store.Projects.Update(project)
		s.Require().NoError(err)

		project, err = s.store.Projects.Get(project.ID)
		s.Require().NoError(err)
		s.Equal("Updated", project.Name)

		err = s.store.Projects.Delete(project.ID)
		s.Require().NoError(err)

		project, err = s.store.Projects.Get(project.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(project)
	})
}

func (s *storeSuite) TestBoardsStore() {
	s.Run("Get", func() {
		actual, err := s.store.Boards.Get("29e247c3-69f1-4397-8bab-b1dd10ae28b2")
		s.Require().NoError(err)

		expected := &store.Board{
			ID:             "29e247c3-69f1-4397-8bab-b1dd10ae28b2",
			ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
			Name:           "Pet projects",
			Archived:       false,
			DateCreated:    testDate,
			DateLastViewed: testDate,
			Color:          0,
			CoverURL:       "",
			TaskLists: []*store.TaskList{
				{
					ID:          "2fcff999-fa20-419b-84b9-023d81a7688e",
					BoardID:     "29e247c3-69f1-4397-8bab-b1dd10ae28b2",
					Name:        "In progress",
					Position:    100,
					Archived:    false,
					DateCreated: testDate,
					Color:       0,
					Tasks:       nil,
				},
				{
					ID:          "32f0de22-cc36-4604-9187-f115b45662bd",
					BoardID:     "29e247c3-69f1-4397-8bab-b1dd10ae28b2",
					Name:        "Ideas",
					Position:    300,
					Archived:    false,
					DateCreated: testDate,
					Color:       0,
					Tasks:       nil,
				},
				{
					ID:          "3b8fcd44-4b59-4fa8-ae12-6ca22ddabd01",
					BoardID:     "29e247c3-69f1-4397-8bab-b1dd10ae28b2",
					Name:        "Done",
					Position:    200,
					Archived:    false,
					DateCreated: testDate,
					Color:       0,
					Tasks:       nil,
				},
			},
		}

		s.Equal(expected, actual)
	})

	s.Run("Get - not found", func() {
		actual, err := s.store.Boards.Get("0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		board := &store.Board{Name: "New test", ProjectID: "fd5f451d-fac6-4bc7-a677-34adb39a6701"}
		err := s.store.Boards.Add(board)

		s.Require().NoError(err)
		s.Equal("New test", board.Name)
		s.Equal(0, board.Color)
		s.Equal(false, board.Archived)
		s.Equal("", board.CoverURL)
		s.NotEmpty(board.ID)
		s.Nil(board.TaskLists)

		{
			project, err := s.store.Projects.Get("fd5f451d-fac6-4bc7-a677-34adb39a6701")
			s.Require().NoError(err)
			s.Len(project.Boards, 5)
		}

		board.Name = "Updated"
		err = s.store.Boards.Update(board)
		s.Require().NoError(err)

		board, err = s.store.Boards.Get(board.ID)
		s.Require().NoError(err)
		s.Equal("Updated", board.Name)

		err = s.store.Boards.Delete(board.ID)
		s.Require().NoError(err)

		board, err = s.store.Boards.Get(board.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(board)
	})
}

func (s *storeSuite) TestTaskListsStore() {
	s.Run("Get", func() {
		actual, err := s.store.TaskLists.Get("93892ed8-bd3d-4f8e-b820-bf9a5043bc1d")
		s.Require().NoError(err)

		expected := &store.TaskList{
			ID:          "93892ed8-bd3d-4f8e-b820-bf9a5043bc1d",
			BoardID:     "f3fc69f2-27aa-4aed-842e-9ed544661bfd",
			Name:        "Food",
			Position:    200,
			Archived:    false,
			DateCreated: testDate,
			Color:       0,
			Tasks: []*store.Task{
				{
					ID:          "522a2569-caf5-4c59-8d95-5670ed8378d3",
					TaskListID:  "93892ed8-bd3d-4f8e-b820-bf9a5043bc1d",
					Name:        "Chips",
					Text:        "Text",
					Position:    100,
					DateCreated: testDate,
					DueDate:     &testDate,
				},
				{
					ID:          "0f10f18a-bd51-4822-a44f-f5786baf5d07",
					TaskListID:  "93892ed8-bd3d-4f8e-b820-bf9a5043bc1d",
					Name:        "Noodles",
					Position:    200,
					DateCreated: testDate,
				},
			},
		}

		s.Equal(expected, actual)
	})

	s.Run("Get - not found", func() {
		actual, err := s.store.TaskLists.Get("0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		list := &store.TaskList{Name: "New test", BoardID: "29e247c3-69f1-4397-8bab-b1dd10ae28b2"}
		err := s.store.TaskLists.Add(list)

		s.Require().NoError(err)
		s.Equal("New test", list.Name)
		s.Equal(0, list.Color)
		s.Equal(false, list.Archived)
		s.NotEmpty(list.ID)
		s.Nil(list.Tasks)

		{
			project, err := s.store.Boards.Get("29e247c3-69f1-4397-8bab-b1dd10ae28b2")
			s.Require().NoError(err)
			s.Len(project.TaskLists, 4)
		}

		list.Name = "Updated"
		err = s.store.TaskLists.Update(list)
		s.Require().NoError(err)

		list, err = s.store.TaskLists.Get(list.ID)
		s.Require().NoError(err)
		s.Equal("Updated", list.Name)

		err = s.store.TaskLists.Delete(list.ID)
		s.Require().NoError(err)

		list, err = s.store.TaskLists.Get(list.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(list)
	})
}

func (s *storeSuite) TestTasksStore() {
	s.Run("Get", func() {
		actual, err := s.store.Tasks.Get("3a0c9a3b-bbec-4047-9822-1c4806c2a258")
		s.Require().NoError(err)

		expected := &store.Task{
			ID:          "3a0c9a3b-bbec-4047-9822-1c4806c2a258",
			TaskListID:  "2fcff999-fa20-419b-84b9-023d81a7688e",
			Name:        "Refactor geometry",
			Position:    100,
			DateCreated: testDate,
			Comments: []*store.Comment{
				{
					ID:          "56b3893f-1887-4b7b-bae9-31f959ecca68",
					TaskID:      "3a0c9a3b-bbec-4047-9822-1c4806c2a258",
					Text:        "Read a book about computational geometry",
					DateCreated: testDate,
					Author:      "User",
				},
				{
					ID:          "fbc762dc-2895-45bc-ae56-a73a7fa6e8b5",
					TaskID:      "3a0c9a3b-bbec-4047-9822-1c4806c2a258",
					Text:        "Checkout Coursera",
					DateCreated: testDate,
					Author:      "User",
				},
			},
		}

		s.Equal(expected, actual)
	})

	s.Run("Get - not found", func() {
		actual, err := s.store.Tasks.Get("0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		task := &store.Task{
			Name:       "New test",
			Text:       "xxx",
			TaskListID: "2fcff999-fa20-419b-84b9-023d81a7688e",
			Position:   777,
		}
		err := s.store.Tasks.Add(task)

		s.Require().NoError(err)
		s.Equal("New test", task.Name)
		s.Equal("xxx", task.Text)
		s.Equal(int64(777), task.Position)
		s.NotEmpty(task.ID)
		s.Nil(task.DueDate)
		s.Nil(task.Comments)

		{
			list, err := s.store.TaskLists.Get("2fcff999-fa20-419b-84b9-023d81a7688e")
			s.Require().NoError(err)
			s.Len(list.Tasks, 2)
		}

		task.Name = "Updated"
		task.Text = "yyy"
		err = s.store.Tasks.Update(task)
		s.Require().NoError(err)

		task, err = s.store.Tasks.Get(task.ID)
		s.Require().NoError(err)
		s.Equal("Updated", task.Name)
		s.Equal("yyy", task.Text)

		err = s.store.Tasks.Delete(task.ID)
		s.Require().NoError(err)

		task, err = s.store.Tasks.Get(task.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(task)
	})
}

func (s *storeSuite) TestCommentsStore() {
	s.Run("Get", func() {
		actual, err := s.store.Comments.Get("4d715efa-8e2f-4488-99d2-5b69e7a43aec")
		s.Require().NoError(err)

		expected := &store.Comment{
			ID:          "4d715efa-8e2f-4488-99d2-5b69e7a43aec",
			TaskID:      "522a2569-caf5-4c59-8d95-5670ed8378d3",
			Author:      "User",
			Text:        "https://www.jamieoliver.com/recipes/vegetables-recipes/the-perfect-chips/",
			DateCreated: testDate,
		}

		s.Equal(expected, actual)
	})

	s.Run("Get - not found", func() {
		actual, err := s.store.Comments.Get("0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		comment := &store.Comment{
			Text:   "xxx",
			TaskID: "3a0c9a3b-bbec-4047-9822-1c4806c2a258",
		}
		err := s.store.Comments.Add(comment)

		s.Require().NoError(err)
		s.Equal("xxx", comment.Text)
		s.NotEmpty(comment.ID)

		{
			list, err := s.store.Tasks.Get("3a0c9a3b-bbec-4047-9822-1c4806c2a258")
			s.Require().NoError(err)
			s.Len(list.Comments, 3)
		}

		comment.Text = "yyy"
		err = s.store.Comments.Update(comment)
		s.Require().NoError(err)

		comment, err = s.store.Comments.Get(comment.ID)
		s.Require().NoError(err)
		s.Equal("yyy", comment.Text)

		err = s.store.Comments.Delete(comment.ID)
		s.Require().NoError(err)

		comment, err = s.store.Comments.Get(comment.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(comment)
	})
}
