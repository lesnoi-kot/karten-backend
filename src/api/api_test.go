package api_test

import (
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
)

type baseAPITestSuite struct {
	suite.Suite

	store *store.Store
	api   *api.APIService

	projectsMock  *mockProjectsStore
	boardsMock    *mockBoardsStore
	taskListsMock *mockTaskListsStore
}

func (suite *baseAPITestSuite) init() {
	suite.projectsMock = new(mockProjectsStore)
	suite.boardsMock = new(mockBoardsStore)
	suite.taskListsMock = new(mockTaskListsStore)

	suite.store = &store.Store{
		DB:        nil,
		Projects:  suite.projectsMock,
		Boards:    suite.boardsMock,
		TaskLists: suite.taskListsMock,
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

func (m mockBoardsStore) Get(id string) (*store.Board, error) {
	args := m.Called(id)
	return args.Get(0).(*store.Board), args.Error(1)
}

func (m mockBoardsStore) Add(a *store.Board) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m mockBoardsStore) Update(a *store.Board) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m mockBoardsStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m mockProjectsStore) Get(id string) (*store.Project, error) {
	args := m.Called(id)
	return args.Get(0).(*store.Project), args.Error(1)
}

func (m mockProjectsStore) GetAll() ([]*store.Project, error) {
	args := m.Called()
	return args.Get(0).([]*store.Project), args.Error(1)
}

func (m mockProjectsStore) Add(project *store.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m mockProjectsStore) Update(project *store.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m mockProjectsStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m mockTaskListsStore) Get(id string) (*store.TaskList, error) {
	args := m.Called(id)
	return args.Get(0).(*store.TaskList), args.Error(1)
}

func (m mockTaskListsStore) Add(project *store.TaskList) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m mockTaskListsStore) Update(project *store.TaskList) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m mockTaskListsStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
