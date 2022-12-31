package models

import "github.com/uptrace/bun"

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	ID   string `bun:"id,pk,autoincrement" json:"id"`
	Name string `json:"name"`

	Boards []*Board `bun:"rel:has-many,join:id=project_id" json:"boards,omitempty"`
}
