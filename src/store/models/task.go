package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	ID          string    `bun:"id,pk,autoincrement" json:"id"`
	TaskListID  string    `bun:"task_list_id" json:"task_list_id"`
	Name        string    `json:"name"`
	Text        string    `json:"text"`
	Position    int64     `json:"position"`
	DateCreated time.Time `json:"date_created"`
	DueDate     time.Time `bun:"due_date,nullzero" json:"due_date"`

	Comments []*Comment `bun:"rel:has-many,join:id=task_id" json:"comments,omitempty"`
}
