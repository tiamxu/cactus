package repo

import (
	"context"
	"fmt"
	"strings"

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

func GetAllLinks(ctx context.Context, pageNo, pageSize int) ([]*model.NavigationLink, int64, error) {
	var links []*model.NavigationLink
	var total int64
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	countQuery := "SELECT COUNT(*) FROM " + NavigationTableName + " WHERE status = 1"
	err := DB.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT * FROM " + NavigationTableName + " WHERE status = 1  ORDER BY category, title LIMIT ? OFFSET ?"
	err = DB.Select(&links, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	return links, total, nil
}

func GetLinkByID(ctx context.Context, id int) (*model.NavigationLink, error) {
	var link model.NavigationLink
	query := "SELECT * FROM " + NavigationTableName + " WHERE id = ?"
	err := DB.GetContext(ctx, &link, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("链接不存在")
		}
		return nil, fmt.Errorf("查询链接失败: %v", err)
	}
	return &link, err
}
func ListNavigationLinks(ctx context.Context, title string, category string, status *int, page, pageSize int) ([]*model.NavigationLink, int64, error) {
	var links []*model.NavigationLink
	var total int64

	// 构建查询条件
	var where []string
	var args []interface{}

	if title != "" {
		where = append(where, "title LIKE ?")
		args = append(args, "%"+title+"%")
	}
	if category != "" {
		where = append(where, "category LIKE ?")
		args = append(args, "%"+category+"%")
	}
	if status != nil {
		where = append(where, "status = ?")
		args = append(args, *status)
	}
	// 计算总数
	countQuery := "SELECT COUNT(*) FROM " + NavigationTableName
	if len(where) > 0 {
		countQuery += " WHERE " + strings.Join(where, " AND ")
	}
	err := DB.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询链接总数失败: %v", err)
	}

	// 查询数据
	query := "SELECT * FROM " + NavigationTableName
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	err = DB.SelectContext(ctx, &links, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询链接列表失败: %v", err)
	}

	return links, total, nil
}
func InsertLink(ctx context.Context, link *model.NavigationLink) error {
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

func UpdateNavigationWithId(ctx context.Context, id int, link *model.NavigationLink) error {
	query := "UPDATE " + NavigationTableName +
		" SET title = ?, url = ?, icon = ?, category = ?, description = ?,status = ? WHERE id = ?"
	_, err := DB.ExecContext(ctx,
		query,
		link.Title, link.URL, link.Icon, link.Category, link.Description, link.Status, id)
	return err
}

func DeleteNavigationWithId(ctx context.Context, id int) error {
	query := "DELETE FROM " + NavigationTableName + " WHERE id = ?"
	_, err := DB.ExecContext(ctx, query, id)
	return err
}
