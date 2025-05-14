package models

import (
	"time"

	"github.com/tiamxu/kit/sql"
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

type CreateLinkRequest struct {
	Title       string `json:"title" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type UpdateLinkRequest struct {
	Title       string `json:"title"`
	URL         string `json:"url" binding:"omitempty,url"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Description string `json:"description"`
}
type NavigationDB struct {
	*sql.DB
}

func NewNavigationDB() *NavigationDB {
	return &NavigationDB{NewDBClient()}
}
func (db NavigationDB) GetAllLinks() ([]NavigationLink, error) {
	var links []NavigationLink
	err := db.Select(&links, "SELECT * FROM navigation_links ORDER BY category, title")
	return links, err
}

func (db NavigationDB) GetLinkByID(id int) (NavigationLink, error) {
	var link NavigationLink
	err := db.Get(&link, "SELECT * FROM navigation_links WHERE id = ?", id)
	return link, err
}

func (db NavigationDB) CreateLink(link CreateLinkRequest) (int, error) {
	result, err := db.Exec(
		"INSERT INTO navigation_links (title, url, icon, category, description) VALUES (?, ?, ?, ?, ?)",
		link.Title, link.URL, link.Icon, link.Category, link.Description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (db NavigationDB) UpdateLink(id int, link UpdateLinkRequest) error {
	_, err := db.Exec(
		"UPDATE navigation_links SET title = ?, url = ?, icon = ?, category = ?, description = ? WHERE id = ?",
		link.Title, link.URL, link.Icon, link.Category, link.Description, id)
	return err
}

func (db NavigationDB) DeleteLink(id int) error {
	_, err := db.Exec("DELETE FROM navigation_links WHERE id = ?", id)
	return err
}
