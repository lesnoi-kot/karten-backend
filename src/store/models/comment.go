package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Comment struct {
	bun.BaseModel `bun:"table:comments"`

	ID          string `bun:"id,pk"`
	TaskID      string
	Author      string
	Text        string
	DateCreated time.Time
}
