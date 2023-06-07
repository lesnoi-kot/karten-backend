package settings

var AppConfig appConfig

type appConfig struct {
	StoreDSN        string   `env:"STORE_DSN,notEmpty,unset"`
	APIBindAddress  string   `env:"API_HOST,notEmpty"`
	APIPrefix       string   `env:"API_PREFIX"`
	FrontendURL     string   `env:"FRONTEND_URL,notEmpty"`
	BackendURL      string   `env:"BACKEND_URL,notEmpty"`
	MediaURL        string   `env:"MEDIA_URL,notEmpty"`
	CookieDomain    string   `env:"COOKIE_DOMAIN"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH,notEmpty"`
	AllowOrigins    []string `env:"ALLOW_ORIGINS,notEmpty" envSeparator:","`

	Debug       bool `env:"DEBUG"`
	EnableGuest bool `env:"ENABLE_GUEST"`

	GithubClientID     string `env:"GITHUB_CLIENT_ID,notEmpty"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET,notEmpty,unset"`

	SessionsSecretKey string `env:"SESSIONS_SECRET_KEY,notEmpty,unset"`
	SessionsStorePath string `env:"SESSIONS_STORE_PATH,notEmpty"`
}

type projectsConfig struct {
	AvatarThumbnailSize uint // px
}

var Projects = projectsConfig{
	AvatarThumbnailSize: 80,
}
