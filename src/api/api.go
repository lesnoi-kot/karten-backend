package api

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

const (
	SHUTDOWN_TIMEOUT = 10_000
)

type APIService struct {
	handler *echo.Echo
	store   *bun.DB
	logger  *zap.Logger
}

type APIConfig struct {
	DB     *bun.DB
	Logger *zap.Logger
}

func NewAPI(c APIConfig) *APIService {
	api := &APIService{
		handler: echo.New(),
		store:   c.DB,
		logger:  c.Logger,
	}

	api.handler.Validator = &Validator{validator.New()}

	api.handler.Pre(middleware.RemoveTrailingSlash())
	api.handler.Use(middleware.Logger(), middleware.Recover())

	initRoutes(api)

	return api
}

func (a APIService) Start(address string) error {
	return a.handler.Start(address)
}

func (a APIService) Shutdown() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		SHUTDOWN_TIMEOUT*time.Millisecond,
	)

	defer cancel()
	return a.handler.Shutdown(ctx)
}

func initRoutes(api *APIService) {
	root := api.handler.Group("/api")

	initProjectsRoute(root, api)
	initBoardsRoute(root, api)
	initTaskListsRoute(root, api)
	initTasksRoute(root, api)
	initCommentsRoute(root, api)
}

func initProjectsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/projects")

	subroute.GET("", api.getProjects)
	subroute.POST("", api.addProject)

	subroute.GET("/:id", api.getProject, requireId)
	subroute.PATCH("/:id", api.editProject, requireId)
	subroute.DELETE("/:id", api.deleteProject, requireId)

	subroute.POST("/:id/boards", api.addBoard, requireId)
}

func initBoardsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/boards")

	subroute.GET("/:id", api.getBoard, requireId)
	subroute.PATCH("/:id", api.editBoard, requireId)
	subroute.DELETE("/:id", api.deleteBoard, requireId)

	subroute.POST("/:id/task-lists", api.addTaskList, requireId)
}

func initTaskListsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/task-lists")

	subroute.GET("/:id", api.getTaskList, requireId)
	subroute.PATCH("/:id", api.editTaskList, requireId)
	subroute.DELETE("/:id", api.deleteTaskList, requireId)

	subroute.POST("/:id/tasks", api.addTask, requireId)
}

func initTasksRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/tasks")

	subroute.GET("/:id", api.getTask, requireId)
	subroute.PATCH("/:id", api.editTask, requireId)
	subroute.DELETE("/:id", api.deleteTask, requireId)

	subroute.POST("/:id/comments", api.addComment, requireId)
}

func initCommentsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/comments")

	subroute.PATCH("/:id", api.editComment, requireId)
	subroute.DELETE("/:id", api.deleteComment, requireId)
}
