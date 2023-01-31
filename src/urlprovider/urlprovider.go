package urlprovider

import (
	"fmt"

	"github.com/lesnoi-kot/karten-backend/src/store"
)

func GetFileURL(file *store.File) string {
	return fmt.Sprintf("media/%s", file.StorageObjectID)
}
