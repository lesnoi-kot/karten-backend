package urlprovider

import (
	"net/url"

	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func GetFileURL(file *store.File) string {
	url, _ := url.JoinPath(settings.AppConfig.MediaURL, file.StorageObjectID)
	return url
}
