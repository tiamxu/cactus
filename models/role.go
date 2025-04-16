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

// 查询角色信息
func GetRolesByID(roleIDs []int) ([]*Role, error) {
	query, args, err := sqlx.In(`SELECT id, code, name, enable FROM role WHERE id IN (?)`, roleIDs)
	if err != nil {
		return nil, err
	}
	query = DB.Rebind(query) // 重新绑定查询语句
	var roles []*Role
	err = DB.Select(&roles, query, args...)
	return roles, err
}

func GetRolesIdByUserID(userId int) ([]int, error) {
	query := `
		SELECT r.id 
		FROM role r
		JOIN user_roles_role urr ON r.id = urr.roleId
		WHERE urr.userId = ?`

	var roleIds []int
	err := DB.Select(&roleIds, query, userId)
	if err != nil {
		return nil, err
	}

	return roleIds, nil
}

func GetRolesByUserId(userId int) ([]*Role, error) {
	query := `
		SELECT r.* 
		FROM role r
		JOIN user_roles_role urr ON r.id = urr.roleId
		WHERE urr.userId = ?`

	var roles []*Role
	err := DB.Select(&roles, query, userId)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func GetRolesByUserIds(userIds []int) (map[int][]*Role, error) {
	if len(userIds) == 0 {
		return nil, nil
	}

	// 获取用户角色关系
	query, args, err := sqlx.In(`
		SELECT userId, roleId FROM user_roles_role 
		WHERE userId IN (?)
	`, userIds)
	if err != nil {
		return nil, err
	}

	var userRoles []UserRolesRole
	err = DB.Select(&userRoles, query, args...)
	if err != nil {
		return nil, err
	}

	// 收集所有角色ID
	roleIds := make([]int, 0, len(userRoles))
	userRoleMap := make(map[int][]int)
	for _, ur := range userRoles {
		roleIds = append(roleIds, ur.RoleId)
		userRoleMap[ur.UserId] = append(userRoleMap[ur.UserId], ur.RoleId)
	}

	// 获取所有角色
	var roles []Role
	if len(roleIds) > 0 {
		roleQuery, roleArgs, err := sqlx.In("SELECT * FROM role WHERE id IN (?)", roleIds)
		if err != nil {
			return nil, err
		}
		err = DB.Select(&roles, roleQuery, roleArgs...)
		if err != nil {
			return nil, err
		}
	}

	// 构建角色ID到角色的映射
	roleMap := make(map[int]Role)
	for _, role := range roles {
		roleMap[role.ID] = role
	}

	// 构建最终结果：用户ID到角色列表的映射
	result := make(map[int][]*Role)
	for userId, rIds := range userRoleMap {
		for _, rId := range rIds {
			if role, ok := roleMap[rId]; ok {
				result[userId] = append(result[userId], &role)
			}
		}
	}

	return result, nil
}
