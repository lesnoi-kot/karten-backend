package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Board struct {
	bun.BaseModel `bun:"table:boards"`

	ID             string    `bun:"id,pk" json:"id"`
	Name           string    `json:"name"`
	ProjectID      string    `json:"project_id"`
	Archived       bool      `json:"archived"`
	DateCreated    time.Time `json:"date_created"`
	DateLastViewed time.Time `json:"date_last_viewed"`
	Color          Color     `json:"color"`
	CoverURL       string    `bun:"cover_url,nullzero" json:"cover_url"`

	TaskLists []*TaskList `bun:"rel:has-many,join:id=board_id" json:"task_lists,omitempty"`
}
