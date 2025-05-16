package repo

import (
	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/logic/model"

	"github.com/tiamxu/kit/sql"
)

var (
	NavigationTableName = " navigation_links "
)

type NavigationDB struct {
	*sql.DB
}

func NewNavigationDB() *NavigationDB {
	return &NavigationDB{NewDBClient()}
}

func (db NavigationDB) GetAllLinks(pageNo, pageSize int) ([]model.NavigationLink, int64, error) {
	var links []model.NavigationLink
	var args []interface{}
	var total int64
	query := "SELECT * FROM " + NavigationTableName + " ORDER BY category, title"
	countQuery := "SELECT COUNT(*) FROM" + NavigationTableName + "WHERE 1=1"
	err := db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}
	pageQuery := query + " LIMIT ? OFFSET ?"
	pageArgs := append(args, pageSize, (pageNo-1)*pageSize)
	err = db.Select(&links, pageQuery, pageArgs...)
	return links, total, err
}

func (db NavigationDB) GetLinkByID(id int) (model.NavigationLink, error) {
	var link model.NavigationLink
	query := "SELECT * FROM " + NavigationTableName + " WHERE id = ?"
	err := db.Get(&link, query, id)
	return link, err
}

func (db NavigationDB) Create(link inout.CreateLinkRequest) error {
	query := "INSERT INTO " + NavigationTableName + " (title, url, icon, category, description) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(
		query,
		link.Title, link.URL, link.Icon, link.Category, link.Description)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	return err
}

func (db NavigationDB) UpdateNavigationWithId(id int, link inout.UpdateLinkRequest) error {
	query := "UPDATE " + NavigationTableName + " SET title = ?, url = ?, icon = ?, category = ?, description = ? WHERE id = ?"
	_, err := db.Exec(
		query,
		link.Title, link.URL, link.Icon, link.Category, link.Description, id)
	return err
}

func (db NavigationDB) DeleteNavigationWithId(id int) error {
	query := "DELETE FROM " + NavigationTableName + "WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}
