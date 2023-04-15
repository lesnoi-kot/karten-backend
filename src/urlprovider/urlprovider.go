package urlprovider

import (
	"net/url"

	"github.com/lesnoi-kot/karten-backend/src/settings"
	"github.com/lesnoi-kot/karten-backend/src/store"
)

func GetFileURL(file *store.File) string {
	if file == nil {
		return ""
	}

	url, err := url.JoinPath(settings.AppConfig.MediaURL, file.StorageObjectID)
	if err != nil {
		return ""
	}
	return url
}
