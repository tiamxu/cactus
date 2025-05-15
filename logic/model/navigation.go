package model

import (
	"time"
)

type NavigationLink struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	URL         string    `db:"url" json:"url"`
	Icon        string    `db:"icon" json:"icon"`
	Category    string    `db:"category" json:"category"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (NavigationLink) TableName() string {
	return "navigations"
}
