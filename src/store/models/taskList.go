package models

import (
	"time"

	"github.com/uptrace/bun"
)

type TaskList struct {
	bun.BaseModel `bun:"table:task_lists"`

	ID          string    `bun:"id,pk" json:"id"`
	BoardID     string    `bun:"board_id" json:"board_id"`
	Name        string    `json:"name"`
	Archived    bool      `json:"archived"`
	Position    int64     `json:"position"`
	DateCreated time.Time `json:"date_created"`
	Color       Color     `json:"color"`

	Tasks []*Task `bun:"rel:has-many,join:id=task_list_id" json:"tasks,omitempty"`
}
