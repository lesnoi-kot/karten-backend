package api

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"

	"go.uber.org/zap"
	"go.uber.org/zap/zapio"
)

const (
	SHUTDOWN_TIMEOUT = 10_000
)

type APIService struct {
	handler     *echo.Echo
	store       *store.Store
	logger      *zap.SugaredLogger
	fileStorage filestorage.FileStorage
	apiPrefix   string
	frontendURL string
}

type APIConfig struct {
	Store        *store.Store
	Logger       *zap.SugaredLogger
	FileStorage  filestorage.FileStorage
	FrontendURL  string
	APIPrefix    string
	AllowOrigins []string
	CookieDomain string
	Debug        bool
}

func NewAPI(cfg APIConfig) *APIService {
	api := &APIService{
		handler:     echo.New(),
		store:       cfg.Store,
		logger:      cfg.Logger,
		fileStorage: cfg.FileStorage,
		apiPrefix:   cfg.APIPrefix,
		frontendURL: cfg.FrontendURL,
	}

	api.handler.Debug = cfg.Debug
	api.handler.Logger.SetOutput(
		&zapio.Writer{Log: cfg.Logger.Desugar(), Level: zap.DebugLevel},
	)
	api.handler.Validator = newEchoValidator()
	api.handler.HTTPErrorHandler = api.errorHandler

	securityConfig := middleware.SecureConfig{
		ContentSecurityPolicy: "default-src 'none';",
		ReferrerPolicy:        "same-origin",
	}

	corsConfig := middleware.CORSConfig{
		AllowOrigins:     cfg.AllowOrigins,
		AllowCredentials: true, // Allow cookies in cross origin requests.
	}

	csrfConfig := middleware.CSRFConfig{
		CookieSameSite: http.SameSiteStrictMode,
		CookieDomain:   cfg.CookieDomain,
		CookieSecure:   !cfg.Debug,
	}

	api.handler.Pre(middleware.RemoveTrailingSlash())
	api.handler.Use(
		middleware.Logger(),
		middleware.SecureWithConfig(securityConfig),
		middleware.CORSWithConfig(corsConfig),
		middleware.CSRFWithConfig(csrfConfig),
		middleware.BodyLimit("10M"),

		session.Middleware(
			sessions.NewFilesystemStore("", []byte(settings.AppConfig.SessionsSecretKey)),
		),

		parseError,
	)

	if cfg.Debug {
		// Emulate delay to debug frontend ux.
		api.handler.Use(emulateDelay)

		// Serve user uploaded media.
		api.handler.Group("/media", middleware.Static("./media"))

		// Proxy the frontend application in DEBUG mode.
		spaURL, _ := url.Parse("http://localhost:3000/")
		balancer := middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: spaURL},
		})
		api.handler.Group("/client", middleware.Proxy(balancer))
	} else {
		api.handler.Use(middleware.Recover())
	}

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

func (a APIService) Server() *http.Server {
	return a.handler.Server
}

func (a APIService) Prefix() string {
	return a.apiPrefix
}

func initRoutes(api *APIService) {
	root := api.handler.Group(api.apiPrefix)

	root.GET("/ping", api.ping)
	root.GET("/cover-images", api.getCoverImages)

	initProjectsRoute(root, api)
	initBoardsRoute(root, api)
	initTaskListsRoute(root, api)
	initTasksRoute(root, api)
	initCommentsRoute(root, api)
	initUserRoutes(root, api)
}

func initUserRoutes(root *echo.Group, api *APIService) {
	root.GET("/oauth-callback", api.oauthCallback)

	subroute := root.Group("/users")
	subroute.GET("/self", api.getCurrentUser)
}

func initProjectsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/projects")

	subroute.GET("", api.getProjects)
	subroute.POST("", api.addProject)
	subroute.DELETE("", api.deleteProjects)

	subroute.GET("/:id", api.getProject)
	subroute.PATCH("/:id", api.editProject)
	subroute.DELETE("/:id", api.deleteProject)

	subroute.POST("/:id/boards", api.addBoard)
	subroute.DELETE("/:id/boards", api.clearProject)
}

func initBoardsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/boards")

	subroute.GET("/:id", api.getBoard)
	subroute.PATCH("/:id", api.editBoard)
	subroute.DELETE("/:id", api.deleteBoard)

	subroute.PUT("/:id/favorite", api.favoriteBoard)
	subroute.DELETE("/:id/favorite", api.unfavoriteBoard)

	subroute.POST("/:id/task-lists", api.addTaskList)
}

func initTaskListsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/task-lists")

	subroute.GET("/:id", api.getTaskList)
	subroute.PATCH("/:id", api.editTaskList)
	subroute.DELETE("/:id", api.deleteTaskList)

	subroute.POST("/:id/tasks", api.addTask)
}

func initTasksRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/tasks")

	subroute.GET("/:id", api.getTask)
	subroute.PATCH("/:id", api.editTask)
	subroute.DELETE("/:id", api.deleteTask)

	subroute.POST("/:id/comments", api.addComment)
}

func initCommentsRoute(root *echo.Group, api *APIService) {
	subroute := root.Group("/comments")

	subroute.PATCH("/:id", api.editComment)
	subroute.DELETE("/:id", api.deleteComment)
}

func (api *APIService) ping(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (api *APIService) errorHandler(err error, c echo.Context) {
	api.handler.DefaultHTTPErrorHandler(err, c)
}
