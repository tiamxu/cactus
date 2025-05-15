package repo

import (
	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/logic/model"

	"github.com/tiamxu/kit/sql"
)

type NavigationDB struct {
	*sql.DB
}

func NewNavigationDB() *NavigationDB {
	return &NavigationDB{NewDBClient()}
}
func (db NavigationDB) GetAllLinks() ([]model.NavigationLink, error) {
	var links []model.NavigationLink
	err := db.Select(&links, "SELECT * FROM navigation_links ORDER BY category, title")
	return links, err
}

func (db NavigationDB) GetLinkByID(id int) (model.NavigationLink, error) {
	var link model.NavigationLink
	err := db.Get(&link, "SELECT * FROM navigation_links WHERE id = ?", id)
	return link, err
}

func (db NavigationDB) CreateLink(link inout.CreateLinkRequest) (int, error) {
	result, err := db.Exec(
		"INSERT INTO navigation_links (title, url, icon, category, description) VALUES (?, ?, ?, ?, ?)",
		link.Title, link.URL, link.Icon, link.Category, link.Description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (db NavigationDB) UpdateLink(id int, link inout.UpdateLinkRequest) error {
	_, err := db.Exec(
		"UPDATE navigation_links SET title = ?, url = ?, icon = ?, category = ?, description = ? WHERE id = ?",
		link.Title, link.URL, link.Icon, link.Category, link.Description, id)
	return err
}

func (db NavigationDB) DeleteLink(id int) error {
	_, err := db.Exec("DELETE FROM navigation_links WHERE id = ?", id)
	return err
}
