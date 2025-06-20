package repo

import (
	"context"

	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/logic/model"

	"github.com/tiamxu/kit/sql"
)

var (
	NavigationTableName = "navigation_links"
)

type NavigationDB struct {
	*sql.DB
}

func NewNavigationDB() *NavigationDB {
	return &NavigationDB{NewDBClient()}
}

func GetAllLinks(ctx context.Context, pageNo, pageSize int) ([]model.NavigationLink, int64, error) {
	var links []model.NavigationLink
	var total int64

	countQuery := "SELECT COUNT(*) FROM " + NavigationTableName + " WHERE 1=1"
	err := DB.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT * FROM " + NavigationTableName + " ORDER BY category, title LIMIT ? OFFSET ?"
	offset := (pageNo - 1) * pageSize
	err = DB.Select(&links, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return links, total, nil
}

func GetLinkByID(ctx context.Context, id int) (model.NavigationLink, error) {
	var link model.NavigationLink
	query := "SELECT * FROM " + NavigationTableName + " WHERE id = ?"
	err := DB.GetContext(ctx, &link, query, id)
	return link, err
}

func InsertLink(ctx context.Context, link inout.CreateLinkRequest) error {
	query := "INSERT INTO " + NavigationTableName +
		" (title, url, icon, category, description) VALUES (?, ?, ?, ?, ?)"
	result, err := DB.ExecContext(ctx,
		query,
		link.Title, link.URL, link.Icon, link.Category, link.Description)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	return err
}

func UpdateNavigationWithId(ctx context.Context, id int, link inout.UpdateLinkRequest) error {
	query := "UPDATE " + NavigationTableName +
		" SET title = ?, url = ?, icon = ?, category = ?, description = ? WHERE id = ?"
	_, err := DB.ExecContext(ctx,
		query,
		link.Title, link.URL, link.Icon, link.Category, link.Description, id)
	return err
}

func DeleteNavigationWithId(ctx context.Context, id int) error {
	query := "DELETE FROM " + NavigationTableName + " WHERE id = ?"
	_, err := DB.ExecContext(ctx, query, id)
	return err
}
