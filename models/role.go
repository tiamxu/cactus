package models

import "github.com/jmoiron/sqlx"

type Role struct {
	ID     int    `db:"id" json:"id"`
	Code   string `db:"code" json:"code"`
	Name   string `db:"name" json:"name"`
	Enable bool   `db:"enable" json:"enable"`
}

func (Role) TableName() string {
	return "role"
}

func GetRolesByUserID(userId int) ([]Role, error) {
	query := `
		SELECT r.* 
		FROM role r
		JOIN user_roles_role urr ON r.id = urr.roleId
		WHERE urr.userId = ?`

	var roles []Role
	err := DB.Select(&roles, query, userId)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func GetRolesByUserIDs(userIds []int) (map[int][]Role, error) {
	if len(userIds) == 0 {
		return nil, nil
	}

	// 使用JOIN一次性查询所有用户角色
	query := `
		SELECT urr.userId, r.* 
		FROM role r
		JOIN user_roles_role urr ON r.id = urr.roleId
		WHERE urr.userId IN (?)`

	query, args, err := sqlx.In(query, userIds)
	if err != nil {
		return nil, err
	}

	type userRole struct {
		UserID int `db:"userId"`
		Role
	}

	var userRoles []userRole
	err = DB.Select(&userRoles, query, args...)
	if err != nil {
		return nil, err
	}

	// 构建用户ID到角色列表的映射
	result := make(map[int][]Role)
	for _, ur := range userRoles {
		result[ur.UserID] = append(result[ur.UserID], ur.Role)
	}

	return result, nil
}
