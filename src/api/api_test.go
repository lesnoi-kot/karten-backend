package api_test

import (
	"context"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/api"
	"github.com/lesnoi-kot/karten-backend/src/store"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var testDate = time.Unix(0, 0).UTC()

type (
	mockBoardsStore    struct{ mock.Mock }
	mockProjectsStore  struct{ mock.Mock }
	mockTaskListsStore struct{ mock.Mock }
	mockTasksStore     struct{ mock.Mock }
	mockCommentsStore  struct{ mock.Mock }
)

type baseAPITestSuite struct {
	suite.Suite

	store *store.Store
	api   *api.APIService

	projectsMock  *mockProjectsStore
	boardsMock    *mockBoardsStore
	taskListsMock *mockTaskListsStore
	tasksMock     *mockTasksStore
	commentsMock  *mockCommentsStore
}

func (suite *baseAPITestSuite) init() {
	suite.projectsMock = new(mockProjectsStore)
	suite.boardsMock = new(mockBoardsStore)
	suite.taskListsMock = new(mockTaskListsStore)
	suite.tasksMock = new(mockTasksStore)
	suite.commentsMock = new(mockCommentsStore)

	suite.store = &store.Store{
		Entities: store.Entities{
			Projects:  suite.projectsMock,
			Boards:    suite.boardsMock,
			TaskLists: suite.taskListsMock,
			Tasks:     suite.tasksMock,
			Comments:  suite.commentsMock,
		},
	}
	suite.api = api.NewAPI(api.APIConfig{
		Store:     suite.store,
		Logger:    zap.NewNop().Sugar(),
		APIPrefix: "",
		Debug:     true,
	})
}

func (suite *baseAPITestSuite) SetupTest() {
	suite.init()
}

func (m mockBoardsStore) Get(ctx context.Context, id string) (*store.Board, error) {
	args := m.Called(id)
	return args.Get(0).(*store.Board), args.Error(1)
}

func (m mockBoardsStore) Add(ctx context.Context, a *store.Board) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m mockBoardsStore) Update(ctx context.Context, a *store.Board) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m mockBoardsStore) Delete(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m mockBoardsStore) UpdateColumns(ctx context.Context, item *store.Board, columns ...string) error {
	args := m.Called(item, columns)
	return args.Error(0)
}

func (m mockProjectsStore) Get(ctx context.Context, id string) (*store.Project, error) {
	args := m.Called(id)
	return args.Get(0).(*store.Project), args.Error(1)
}

func (m mockProjectsStore) GetAll(ctx context.Context) ([]*store.Project, error) {
	args := m.Called()
	return args.Get(0).([]*store.Project), args.Error(1)
}

func (m mockProjectsStore) Add(ctx context.Context, item *store.Project) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockProjectsStore) Update(ctx context.Context, item *store.Project) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockProjectsStore) Delete(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m mockTaskListsStore) Get(ctx context.Context, id string) (*store.TaskList, error) {
	args := m.Called(id)
	return args.Get(0).(*store.TaskList), args.Error(1)
}

func (m mockTaskListsStore) Add(ctx context.Context, item *store.TaskList) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockTaskListsStore) Update(ctx context.Context, item *store.TaskList) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockTaskListsStore) Delete(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m mockTasksStore) Get(ctx context.Context, id string) (*store.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*store.Task), args.Error(1)
}

func (m mockTasksStore) Add(ctx context.Context, item *store.Task) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockTasksStore) Update(ctx context.Context, item *store.Task) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockTasksStore) Delete(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m mockCommentsStore) Get(ctx context.Context, id string) (*store.Comment, error) {
	args := m.Called(id)
	return args.Get(0).(*store.Comment), args.Error(1)
}

func (m mockCommentsStore) Add(ctx context.Context, item *store.Comment) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockCommentsStore) Update(ctx context.Context, item *store.Comment) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m mockCommentsStore) Delete(ctx context.Context, id string) error {
	args := m.Called(id)
	return args.Error(0)
}
