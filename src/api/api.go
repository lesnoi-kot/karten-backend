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
	"github.com/lesnoi-kot/karten-backend/src/userservice"

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
	debug       bool
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
		debug:       cfg.Debug,
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
		CookiePath:     "/",
	}

	sessionStore := sessions.NewFilesystemStore(
		settings.AppConfig.SessionsStorePath,
		[]byte(settings.AppConfig.SessionsSecretKey),
	)
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		Domain:   settings.AppConfig.CookieDomain,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   !cfg.Debug,
	}

	api.handler.Pre(middleware.RemoveTrailingSlash())
	api.handler.Use(
		middleware.Logger(),
		middleware.SecureWithConfig(securityConfig),
		middleware.CORSWithConfig(corsConfig),
		middleware.CSRFWithConfig(csrfConfig),
		middleware.BodyLimit("10M"),

		session.Middleware(sessionStore),

		parseError,
	)

	if cfg.Debug {
		// Emulate delay to debug frontend ux.
		api.handler.Use(emulateDelay)

		// Serve user uploaded media.
		if localStorage, ok := api.fileStorage.(filestorage.FileSystemStorage); ok {
			api.handler.Group("/media", middleware.Static(localStorage.RootPath))
		}

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

func initRoutes(api *APIService) {
	requireAuth := api.makeRequireAuthMiddleware()

	root := api.handler.Group(api.apiPrefix)

	root.GET("/ping", api.ping)
	root.GET("/cover-images", api.getCoverImages)
	root.GET("/oauth-callback", api.oauthCallback)

	if api.debug {
		root.POST("/login", api.guestLogIn)
	}

	users := root.Group("/users", requireAuth)
	users.GET("/self", api.getCurrentUser)
	users.POST("/self/logout", api.logOut)

	projects := root.Group("/projects", requireAuth)
	projects.GET("", api.getProjects)
	projects.POST("", api.addProject)
	projects.DELETE("", api.deleteProjects)
	projects.GET("/:id", api.getProject)
	projects.PATCH("/:id", api.editProject)
	projects.DELETE("/:id", api.deleteProject)
	projects.POST("/:id/boards", api.addBoard)
	projects.DELETE("/:id/boards", api.clearProject)

	boards := root.Group("/boards", requireAuth)
	boards.GET("/:id", api.getBoard)
	boards.PATCH("/:id", api.editBoard)
	boards.DELETE("/:id", api.deleteBoard)
	boards.PUT("/:id/favorite", api.favoriteBoard)
	boards.DELETE("/:id/favorite", api.unfavoriteBoard)
	boards.POST("/:id/task-lists", api.addTaskList)

	taskLists := root.Group("/task-lists", requireAuth)
	taskLists.GET("/:id", api.getTaskList)
	taskLists.PATCH("/:id", api.editTaskList)
	taskLists.DELETE("/:id", api.deleteTaskList)
	taskLists.POST("/:id/tasks", api.addTask)

	tasks := root.Group("/tasks", requireAuth)
	tasks.GET("/:id", api.getTask)
	tasks.PATCH("/:id", api.editTask)
	tasks.DELETE("/:id", api.deleteTask)
	tasks.POST("/:id/comments", api.addComment)

	comments := root.Group("/comments", requireAuth)
	comments.PATCH("/:id", api.editComment)
	comments.DELETE("/:id", api.deleteComment)

	files := root.Group("/files", requireAuth)
	files.POST("/image", api.uploadImage)
}

func (api *APIService) ping(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (api *APIService) errorHandler(err error, c echo.Context) {
	api.handler.DefaultHTTPErrorHandler(err, c)
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

func (api *APIService) getUserService(c echo.Context) *userservice.UserService {
	userID, _ := getUserID(c)

	return &userservice.UserService{
		Context: context.Background(),
		UserID:  userID,
		Store:   api.store,
	}
}
