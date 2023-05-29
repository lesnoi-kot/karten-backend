package settings

var AppConfig appConfig

type appConfig struct {
	StoreDSN        string   `env:"STORE_DSN,notEmpty,unset"`
	APIBindAddress  string   `env:"API_HOST,notEmpty"`
	APIPrefix       string   `env:"API_PREFIX"`
	FrontendURL     string   `env:"FRONTEND_URL"`
	MediaURL        string   `env:"MEDIA_URL"`
	CookieDomain    string   `env:"COOKIE_DOMAIN,notEmpty"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH,notEmpty"`
	AllowOrigins    []string `env:"ALLOW_ORIGINS,notEmpty" envSeparator:","`
	Debug           bool     `env:"DEBUG"`

	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`

	SessionsSecretKey string `env:"SESSIONS_SECRET_KEY,notEmpty"`
	SessionsStorePath string `env:"SESSIONS_STORE_PATH"`
}

type projectsConfig struct {
	AvatarThumbnailSize uint // px
}

var Projects = projectsConfig{
	AvatarThumbnailSize: 80,
}
