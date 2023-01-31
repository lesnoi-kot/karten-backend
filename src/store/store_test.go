package store_test

import (
	"context"
	"errors"
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
		DSN:         os.Getenv("STORE_DSN"),
		Logger:      zap.NewNop().Sugar(),
		Debug:       false,
		FileStorage: nil, // TODO
	})

	s.Require().NoError(err)
}

func (s *storeSuite) TestProjectsStore() {
	s.Run("GetAll", func() {
		all, err := s.store.Projects.GetAll(context.Background())

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
		project, err := s.store.Projects.Get(context.Background(), "fd5f451d-fac6-4bc7-a677-34adb39a6701")
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
				},
				{
					ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					ID:             "f3fc69f2-27aa-4aed-842e-9ed544661bfd",
					Name:           "Health",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          1,
				},
				{
					ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					ID:             "ea716cd0-0d2b-4aa9-9a00-e5fce1f6670a",
					Name:           "TODO Books",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          0,
				},
				{
					ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
					ID:             "606ecfd6-2a49-4cc2-911c-a0113ebcf0e6",
					Name:           "Guitar stuff",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          0,
				},
			},
		}

		s.Equal(expected, project)
	})

	s.Run("Get", func() {
		project, err := s.store.Projects.Get(context.Background(), "2f146153-ee2f-4968-a241-11a4f00bf212")
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
				},
				{
					ProjectID:      "2f146153-ee2f-4968-a241-11a4f00bf212",
					ID:             "250bf7a7-ad51-4d62-85cd-c554e4d5f686",
					Name:           "Householding",
					Archived:       false,
					DateCreated:    testDate,
					DateLastViewed: testDate,
					Color:          2,
				},
			},
		}

		s.Equal(expected, project)
	})

	s.Run("Get - not found", func() {
		actual, err := s.store.Projects.Get(context.Background(), "0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		project := &store.Project{Name: "New test"}
		err := s.store.Projects.Add(context.Background(), project)

		s.Require().NoError(err)
		s.Equal("New test", project.Name)
		s.NotEmpty(project.ID)
		s.Nil(project.Boards)

		{
			all, err := s.store.Projects.GetAll(context.Background())
			s.Require().NoError(err)
			s.Len(all, 5)
		}

		project.Name = "Updated"
		err = s.store.Projects.Update(context.Background(), project)
		s.Require().NoError(err)

		project, err = s.store.Projects.Get(context.Background(), project.ID)
		s.Require().NoError(err)
		s.Equal("Updated", project.Name)

		err = s.store.Projects.Delete(context.Background(), project.ID)
		s.Require().NoError(err)

		project, err = s.store.Projects.Get(context.Background(), project.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(project)
	})
}

func (s *storeSuite) TestBoardsStore() {
	s.Run("Get", func() {
		actual, err := s.store.Boards.Get(context.Background(), "29e247c3-69f1-4397-8bab-b1dd10ae28b2")
		s.Require().NoError(err)

		expected := &store.Board{
			ID:             "29e247c3-69f1-4397-8bab-b1dd10ae28b2",
			ProjectID:      "fd5f451d-fac6-4bc7-a677-34adb39a6701",
			Name:           "Pet projects",
			Archived:       false,
			DateCreated:    testDate,
			DateLastViewed: testDate,
			Color:          0,
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
		actual, err := s.store.Boards.Get(context.Background(), "0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		board := &store.Board{Name: "New test", ProjectID: "fd5f451d-fac6-4bc7-a677-34adb39a6701"}
		err := s.store.Boards.Add(context.Background(), board)

		s.Require().NoError(err)
		s.Equal("New test", board.Name)
		s.Equal(0, board.Color)
		s.Equal(false, board.Archived)
		s.Equal(nil, board.CoverID)
		s.NotEmpty(board.ID)
		s.Nil(board.TaskLists)

		{
			project, err := s.store.Projects.Get(context.Background(), "fd5f451d-fac6-4bc7-a677-34adb39a6701")
			s.Require().NoError(err)
			s.Len(project.Boards, 5)
		}

		board.Name = "Updated"
		err = s.store.Boards.Update(context.Background(), board)
		s.Require().NoError(err)

		board, err = s.store.Boards.Get(context.Background(), board.ID)
		s.Require().NoError(err)
		s.Equal("Updated", board.Name)

		err = s.store.Boards.Delete(context.Background(), board.ID)
		s.Require().NoError(err)

		board, err = s.store.Boards.Get(context.Background(), board.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(board)
	})
}

func (s *storeSuite) TestTaskListsStore() {
	s.Run("Get", func() {
		actual, err := s.store.TaskLists.Get(context.Background(), "93892ed8-bd3d-4f8e-b820-bf9a5043bc1d")
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
		actual, err := s.store.TaskLists.Get(context.Background(), "0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		list := &store.TaskList{Name: "New test", BoardID: "29e247c3-69f1-4397-8bab-b1dd10ae28b2"}
		err := s.store.TaskLists.Add(context.Background(), list)

		s.Require().NoError(err)
		s.Equal("New test", list.Name)
		s.Equal(0, list.Color)
		s.Equal(false, list.Archived)
		s.NotEmpty(list.ID)
		s.Nil(list.Tasks)

		{
			project, err := s.store.Boards.Get(context.Background(), "29e247c3-69f1-4397-8bab-b1dd10ae28b2")
			s.Require().NoError(err)
			s.Len(project.TaskLists, 4)
		}

		list.Name = "Updated"
		err = s.store.TaskLists.Update(context.Background(), list)
		s.Require().NoError(err)

		list, err = s.store.TaskLists.Get(context.Background(), list.ID)
		s.Require().NoError(err)
		s.Equal("Updated", list.Name)

		err = s.store.TaskLists.Delete(context.Background(), list.ID)
		s.Require().NoError(err)

		list, err = s.store.TaskLists.Get(context.Background(), list.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(list)
	})
}

func (s *storeSuite) TestTasksStore() {
	s.Run("Get", func() {
		actual, err := s.store.Tasks.Get(context.Background(), "3a0c9a3b-bbec-4047-9822-1c4806c2a258")
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
		actual, err := s.store.Tasks.Get(context.Background(), "0ee97545-1d21-4b10-bd58-65e1682b9ec1")
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
		err := s.store.Tasks.Add(context.Background(), task)

		s.Require().NoError(err)
		s.Equal("New test", task.Name)
		s.Equal("xxx", task.Text)
		s.Equal(int64(777), task.Position)
		s.NotEmpty(task.ID)
		s.Nil(task.DueDate)
		s.Nil(task.Comments)

		{
			list, err := s.store.TaskLists.Get(context.Background(), "2fcff999-fa20-419b-84b9-023d81a7688e")
			s.Require().NoError(err)
			s.Len(list.Tasks, 2)
		}

		task.Name = "Updated"
		task.Text = "yyy"
		err = s.store.Tasks.Update(context.Background(), task)
		s.Require().NoError(err)

		task, err = s.store.Tasks.Get(context.Background(), task.ID)
		s.Require().NoError(err)
		s.Equal("Updated", task.Name)
		s.Equal("yyy", task.Text)

		err = s.store.Tasks.Delete(context.Background(), task.ID)
		s.Require().NoError(err)

		task, err = s.store.Tasks.Get(context.Background(), task.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(task)
	})
}

func (s *storeSuite) TestCommentsStore() {
	s.Run("Get", func() {
		actual, err := s.store.Comments.Get(context.Background(), "4d715efa-8e2f-4488-99d2-5b69e7a43aec")
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
		actual, err := s.store.Comments.Get(context.Background(), "0ee97545-1d21-4b10-bd58-65e1682b9ec1")
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(actual)
	})

	s.Run("Add -> Update -> Delete", func() {
		comment := &store.Comment{
			Text:   "xxx",
			TaskID: "3a0c9a3b-bbec-4047-9822-1c4806c2a258",
		}
		err := s.store.Comments.Add(context.Background(), comment)

		s.Require().NoError(err)
		s.Equal("xxx", comment.Text)
		s.NotEmpty(comment.ID)

		{
			list, err := s.store.Tasks.Get(context.Background(), "3a0c9a3b-bbec-4047-9822-1c4806c2a258")
			s.Require().NoError(err)
			s.Len(list.Comments, 3)
		}

		comment.Text = "yyy"
		err = s.store.Comments.Update(context.Background(), comment)
		s.Require().NoError(err)

		comment, err = s.store.Comments.Get(context.Background(), comment.ID)
		s.Require().NoError(err)
		s.Equal("yyy", comment.Text)

		err = s.store.Comments.Delete(context.Background(), comment.ID)
		s.Require().NoError(err)

		comment, err = s.store.Comments.Get(context.Background(), comment.ID)
		s.ErrorIs(err, store.ErrNotFound)
		s.Nil(comment)
	})
}

func (s *storeSuite) TestTxStore() {
	s.Run("Delete all projects and rollback", func() {
		txStore, err := s.store.BeginTx(context.Background())
		s.Require().NoError(err)

		all, err := txStore.Projects.GetAll(context.Background())
		s.Require().NoError(err)

		beforeLen := len(all)
		for _, p := range all {
			err = txStore.Projects.Delete(context.Background(), p.ID)
			s.Require().NoError(err)
		}

		all, err = txStore.Projects.GetAll(context.Background())
		s.Require().NoError(err)
		s.Require().Empty(all)

		err = txStore.Rollback()
		s.Require().NoError(err)

		all, err = s.store.Projects.GetAll(context.Background())
		s.Require().NoError(err)
		s.Require().Len(all, beforeLen)
	})

	s.Run("RunInTx - delete all projects and rollback", func() {
		beforeLen := 43645654

		err := s.store.RunInTx(context.Background(), func(ctx context.Context, store *store.TxStore) error {
			all, err := store.Projects.GetAll(ctx)
			s.Require().NoError(err)

			beforeLen = len(all)

			for _, p := range all {
				err = store.Projects.Delete(ctx, p.ID)
				s.Require().NoError(err)
			}

			all, err = store.Projects.GetAll(ctx)
			s.Require().NoError(err)
			s.Require().Empty(all)

			return errors.New("You know I want to rollback now")
		})
		s.Require().Error(err)

		all, err := s.store.Projects.GetAll(context.Background())
		s.Require().NoError(err)
		s.Require().Len(all, beforeLen)
	})

	s.Run("RunInTx - commit empty", func() {
		ctx := context.Background()
		err := s.store.RunInTx(ctx, func(ctx context.Context, s *store.TxStore) error {
			return nil
		})
		s.Require().NoError(err)
	})

	s.Run("RunInTx - double rollback", func() {
		ctx := context.Background()
		err := s.store.RunInTx(ctx, func(ctx context.Context, s *store.TxStore) error {
			s.Rollback()
			return errors.New("xxx")
		})
		s.Require().ErrorContains(err, "transaction has already been committed or rolled back")
	})

	s.Run("RunInTx - panic", func() {
		ctx := context.Background()

		s.Panics(func() {
			s.store.RunInTx(ctx, func(ctx context.Context, txStore *store.TxStore) error {
				err := txStore.Projects.Delete(ctx, "fd5f451d-fac6-4bc7-a677-34adb39a6701")
				s.Require().NoError(err)

				_, err = txStore.Projects.Get(ctx, "fd5f451d-fac6-4bc7-a677-34adb39a6701")
				s.Require().ErrorIs(err, store.ErrNotFound)

				panic("Test panic")
			})
		})

		_, err := s.store.Projects.Get(ctx, "fd5f451d-fac6-4bc7-a677-34adb39a6701")
		s.Require().NoError(err)
	})
}
