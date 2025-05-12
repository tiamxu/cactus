package models

import (
	"time"
)

type NavigationLink struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	URL         string    `db:"url"`
	Icon        string    `db:"icon"`
	Category    string    `db:"category"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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

func GetAllLinks() ([]NavigationLink, error) {
	var links []NavigationLink
	err := DB.Select(&links, "SELECT * FROM navigation_links ORDER BY category, title")
	return links, err
}

func GetLinkByID(id int) (NavigationLink, error) {
	var link NavigationLink
	err := DB.Get(&link, "SELECT * FROM navigation_links WHERE id = ?", id)
	return link, err
}

func CreateLink(link CreateLinkRequest) (int, error) {
	result, err := DB.Exec(
		"INSERT INTO navigation_links (title, url, icon, category, description) VALUES (?, ?, ?, ?, ?)",
		link.Title, link.URL, link.Icon, link.Category, link.Description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateLink(id int, link UpdateLinkRequest) error {
	_, err := DB.Exec(
		"UPDATE navigation_links SET title = ?, url = ?, icon = ?, category = ?, description = ? WHERE id = ?",
		link.Title, link.URL, link.Icon, link.Category, link.Description, id)
	return err
}

func DeleteLink(id int) error {
	_, err := DB.Exec("DELETE FROM navigation_links WHERE id = ?", id)
	return err
}
