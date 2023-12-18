package store

import (
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type (
	UserID     = int
	FileID     = string
	EntityID   = string
	LabelID    = int
	DateString = string
	Color      = int
)

const GuestUserID = 1

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

	Avatar *File `bun:"rel:has-one,join:avatar_id=id"`
}

func (user User) IsGuest() bool {
	return user.ID == GuestUserID
}

type File struct {
	bun.BaseModel `bun:"table:files"`

	ID              FileID `bun:",pk"`
	StorageObjectID string
	Name            string
	MimeType        string
	Size            int
}

func (file *File) IsImage() bool {
	if file == nil {
		return false
	}

	return strings.HasPrefix(file.MimeType, "image/")
}

type ImageFile struct {
	bun.BaseModel `bun:"table:files"`
	File

	Thumbnails []*File `bun:"m2m:image_thumbnails,join:Original=Thumbnail"`
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
	Labels    []*Label    `bun:"rel:has-many,join:id=board_id"`
	Project   *Project    `bun:"rel:belongs-to,join:project_id=id"`
	Cover     *File       `bun:"rel:has-one,join:cover_id=id"`
}

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	ID      EntityID `bun:",pk,autoincrement"`
	UserID  UserID
	ShortID string `bun:",nullzero"`
	Name    string

	AvatarID FileID     `bun:",nullzero"`
	Avatar   *ImageFile `bun:"rel:has-one,join:avatar_id=id"`

	Boards []*Board `bun:"rel:has-many,join:id=project_id"`
}

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	ID                  EntityID `bun:",pk"`
	UserID              UserID
	ShortID             string
	TaskListID          EntityID
	Name                string
	Text                string
	Position            int64
	SpentTime           int64
	Archived            bool
	DateCreated         time.Time
	DateStartedTracking *time.Time `bun:",nullzero"`
	DueDate             *time.Time `bun:",nullzero"`

	Comments    []*Comment `bun:"rel:has-many,join:id=task_id"`
	Attachments []*File    `bun:"m2m:task_files,join:Task=File"`
	Labels      []*Label   `bun:"m2m:task_labels,join:Task=Label"`

	// Rendered Text markdown
	HTML string `bun:"-"`
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

	ID          EntityID `bun:",pk"`
	UserID      UserID
	TaskID      EntityID
	Text        string
	DateCreated time.Time

	// Rendered Text markdown
	HTML string `bun:"-"`

	Author      *User   `bun:"rel:belongs-to,join:user_id=id"`
	Attachments []*File `bun:"m2m:comment_files,join:Comment=File"`
}

type Label struct {
	bun.BaseModel `bun:"table:labels"`

	ID      LabelID `bun:",pk"`
	BoardID EntityID
	UserID  UserID
	Name    string
	Color   int
}

type ImageThumbnailAssoc struct {
	bun.BaseModel `bun:"table:image_thumbnails"`

	ID      FileID `bun:",pk"` // Thumbnail File.ID
	ImageID FileID // Original image File.ID

	// ORM many-to-many magic:
	Original  *File `bun:"rel:belongs-to,join:image_id=id"`
	Thumbnail *File `bun:"rel:belongs-to,join:id=id"`
}

type CoverImageToFileAssoc struct {
	bun.BaseModel `bun:"table:default_cover_images"`

	ID EntityID `bun:",pk"`
}

type AttachmentToTaskAssoc struct {
	bun.BaseModel `bun:"table:task_files"`

	TaskID EntityID `bun:",pk"`
	FileID FileID   `bun:",pk"`

	// ORM many-to-many magic:
	Task *Task `bun:"rel:belongs-to,join:task_id=id"`
	File *File `bun:"rel:belongs-to,join:file_id=id"`
}

type AttachmentToCommentAssoc struct {
	bun.BaseModel `bun:"table:comment_files"`

	CommentID EntityID `bun:",pk"`
	FileID    FileID   `bun:",pk"`

	// ORM many-to-many magic:
	Comment *Comment `bun:"rel:belongs-to,join:comment_id=id"`
	File    *File    `bun:"rel:belongs-to,join:file_id=id"`
}

type LabelToTaskAssoc struct {
	bun.BaseModel `bun:"table:task_labels"`

	LabelID LabelID  `bun:",pk"`
	TaskID  EntityID `bun:",pk"`

	// ORM many-to-many magic:
	Label *Label `bun:"rel:belongs-to,join:label_id=id"`
	Task  *Task  `bun:"rel:belongs-to,join:task_id=id"`
}
