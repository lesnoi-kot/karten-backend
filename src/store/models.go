package store

import (
	"strings"
	"time"

	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/uptrace/bun"
)

type (
	UserID     = int
	FileID     = string
	EntityID   = string
	DateString = string
	Color      = int
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID          UserID `bun:",pk,autoincrement"`
	SocialID    string
	AvatarID    FileID `bun:",nullzero"`
	Name        string
	Login       string
	Email       string
	URL         string
	DateCreated time.Time
}

type File struct {
	bun.BaseModel `bun:"table:files"`

	ID              FileID             `bun:",pk" json:"id"`
	StorageObjectID filestorage.FileID `json:"storage_object_id"`
	Name            string             `json:"name"`
	MimeType        string             `json:"mime_type"`
	Size            int                `json:"size"`
}

func (file *File) IsImage() bool {
	if file == nil {
		return false
	}

	return strings.HasPrefix(file.MimeType, "image/")
}

type ImageThumbnailAssoc struct {
	bun.BaseModel `bun:"table:image_thumbnails"`

	ID      FileID `bun:",pk"` // Thumbnail File.ID
	ImageID FileID // Original image File.ID

	// ORM many-to-many magic:
	Original  *File `bun:"rel:belongs-to,join:image_id=id" json:"-"`
	Thumbnail *File `bun:"rel:belongs-to,join:id=id" json:"-"`
}

type ImageFile struct {
	bun.BaseModel `bun:"table:files"`
	File

	Thumbnails []*File `bun:"m2m:image_thumbnails,join:Original=Thumbnail" json:"thumbnails,omitempty"`
}

type CoverImageToFileAssoc struct {
	bun.BaseModel `bun:"table:default_cover_images"`

	ID EntityID `bun:",pk" json:"id"`
}

type Board struct {
	bun.BaseModel `bun:"table:boards"`

	ID             EntityID `bun:",pk"`
	UserID         UserID
	ShortID        string
	Name           string
	ProjectID      string
	Archived       bool
	Favorite       bool
	DateCreated    time.Time
	DateLastViewed time.Time
	Color          Color
	CoverID        *FileID `bun:"cover_id,nullzero"`

	TaskLists []*TaskList `bun:"rel:has-many,join:id=board_id"`
	Cover     *File       `bun:"rel:has-one,join:cover_id=id"`
}

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	ID      EntityID `bun:",pk,autoincrement"`
	UserID  UserID
	ShortID string
	Name    string

	AvatarID FileID     `bun:",nullzero"`
	Avatar   *ImageFile `bun:"rel:has-one,join:avatar_id=id"`

	Boards []*Board `bun:"rel:has-many,join:id=project_id"`
}

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	ID          EntityID `bun:",pk,autoincrement" json:"id"`
	UserID      UserID
	ShortID     string     `json:"short_id"`
	TaskListID  EntityID   `bun:"task_list_id" json:"task_list_id"`
	Name        string     `json:"name"`
	Text        string     `json:"text"`
	Position    int64      `json:"position"`
	Archived    bool       `json:"archived"`
	DateCreated time.Time  `json:"date_created"`
	DueDate     *time.Time `bun:"due_date,nullzero" json:"due_date,omitempty"`

	Comments []*Comment `bun:"rel:has-many,join:id=task_id" json:"comments,omitempty"`
}

type TaskList struct {
	bun.BaseModel `bun:"table:task_lists"`

	ID          EntityID  `bun:",pk" json:"-"`
	UserID      UserID    `json:"-"`
	BoardID     EntityID  `bun:"board_id" json:"-"`
	Name        string    `json:"-"`
	Archived    bool      `json:"-"`
	Position    int64     `json:"-"`
	DateCreated time.Time `json:"-"`
	Color       Color     `json:"-"`

	Tasks []*Task `bun:"rel:has-many,join:id=task_list_id" json:"tasks,omitempty"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comments"`

	ID          EntityID `bun:",pk" json:"id"`
	UserID      UserID
	TaskID      EntityID
	Text        string
	DateCreated time.Time
}
