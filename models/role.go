package models

import "strings"

type Role struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

func (Role) TableName() string {
	return "role"
}

func GetRolesByUserID(userId int) ([]Role, error) {
	// 先查询用户角色关联
	roleIdsQuery := `
		SELECT roleId 
		FROM user_roles_role 
		WHERE userId = ?`

	rows, err := DB.Query(roleIdsQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roleIds []int
	for rows.Next() {
		var roleId int
		if err := rows.Scan(&roleId); err != nil {
			continue
		}
		roleIds = append(roleIds, roleId)
	}

	if len(roleIds) == 0 {
		return nil, nil
	}

	// 查询角色详细信息
	rolesQuery := `
	SELECT id, code, name, enable 
	FROM role 
	WHERE id IN (` + strings.Repeat("?,", len(roleIds)-1) + "?)"

	// 将roleIds转换为interface{}切片
	args := make([]interface{}, len(roleIds))
	for i, id := range roleIds {
		args[i] = id
	}

	rows, err = DB.Query(rolesQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Enable,
		); err == nil {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

// 辅助函数，生成占位符字符串
func preparePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}
