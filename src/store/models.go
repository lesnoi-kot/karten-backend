package store

import (
	"time"

	"github.com/lesnoi-kot/karten-backend/src/filestorage"
	"github.com/uptrace/bun"
)

type (
	Color      = int
	DateString = string
)

const (
	NoColor Color = 0
	Red     Color = 1
	Yellow  Color = 2
	Green   Color = 3
)

type File struct {
	bun.BaseModel `bun:"table:files"`

	ID              string             `bun:",pk" json:"id"`
	StorageObjectID filestorage.FileID `json:"storage_object_id"`
	Name            string             `json:"name"`
	MimeType        string             `json:"mime_type"`
	Size            int                `json:"size"`
}

type ImageThumbnailAssoc struct {
	bun.BaseModel `bun:"table:image_thumbnails"`

	ID      string `bun:",pk" json:"id"` // Thumbnail File.ID
	ImageID string `json:"image_id"`     // Original image File.ID

	// ORM many-to-many magic:
	Original  *File `bun:"rel:belongs-to,join:image_id=id" json:"-"`
	Thumbnail *File `bun:"rel:belongs-to,join:id=id" json:"-"`
}

type ImageFile struct {
	bun.BaseModel `bun:"table:files"`
	File

	Thumbnails []File `bun:"m2m:image_thumbnails,join:Original=Thumbnail" json:"thumbnails,omitempty"`
}

type Board struct {
	bun.BaseModel `bun:"table:boards"`

	ID             string    `bun:",pk" json:"id"`
	Name           string    `json:"name"`
	ProjectID      string    `json:"project_id"`
	Archived       bool      `json:"archived"`
	DateCreated    time.Time `json:"date_created"`
	DateLastViewed time.Time `json:"date_last_viewed"`
	Color          Color     `json:"color"`
	CoverURL       string    `bun:"cover_url,nullzero" json:"cover_url"`

	TaskLists []*TaskList `bun:"rel:has-many,join:id=board_id" json:"task_lists,omitempty"`
}

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	ID   string `bun:",pk,autoincrement" json:"id"`
	Name string `json:"name"`

	AvatarID string     `bun:",nullzero" json:"-"`
	Avatar   *ImageFile `bun:"rel:has-one,join:avatar_id=id" json:"avatar,omitempty"`

	Boards []*Board `bun:"rel:has-many,join:id=project_id" json:"boards,omitempty"`
}

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	ID          string     `bun:",pk,autoincrement" json:"id"`
	TaskListID  string     `bun:"task_list_id" json:"task_list_id"`
	Name        string     `json:"name"`
	Text        string     `json:"text"`
	Position    int64      `json:"position"`
	Archived    bool       `json:"archived"`
	DateCreated time.Time  `json:"date_created"`
	DueDate     *time.Time `bun:"due_date,nullzero" json:"due_date"`

	Comments []*Comment `bun:"rel:has-many,join:id=task_id" json:"comments,omitempty"`
}

type TaskList struct {
	bun.BaseModel `bun:"table:task_lists"`

	ID          string    `bun:",pk" json:"id"`
	BoardID     string    `bun:"board_id" json:"board_id"`
	Name        string    `json:"name"`
	Archived    bool      `json:"archived"`
	Position    int64     `json:"position"`
	DateCreated time.Time `json:"date_created"`
	Color       Color     `json:"color"`

	Tasks []*Task `bun:"rel:has-many,join:id=task_list_id" json:"tasks,omitempty"`
}

type Comment struct {
	bun.BaseModel `bun:"table:comments"`

	ID          string    `bun:",pk" json:"id"`
	TaskID      string    `json:"task_id"`
	Author      string    `json:"author"`
	Text        string    `json:"text"`
	DateCreated time.Time `json:"date_created"`
}
