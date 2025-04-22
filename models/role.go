package models

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Role struct {
	ID     int    `db:"id" json:"id"`
	Code   string `db:"code" json:"code"`
	Name   string `db:"name" json:"name"`
	Enable bool   `db:"enable" json:"enable"`
}

func (Role) TableName() string {
	return "role"
}
func GetPermissionsTree(userID int) ([]Permission, error) {
	// 检查是否是管理员
	var adminRole int64
	err := DB.Get(&adminRole,
		"SELECT COUNT(*) FROM user_roles_role WHERE userId = ? AND roleId = 1",
		userID)
	if err != nil {
		return nil, errors.New("查询管理员状态失败")
	}
	// 构建基础查询
	baseQuery := "SELECT * FROM permission WHERE parentId IS NULL ORDER BY `order` ASC"
	var args []interface{}
	// 非管理员权限过滤
	if adminRole == 0 {
		// 查询用户拥有的角色ID列表

		var uroleIds []int64
		err := DB.Select(&uroleIds, "SELECT roleId FROM user_roles_role WHERE userId = ?", userID)
		if err != nil {
			return nil, errors.New("查询用户角色失败")
		}
		// 如果用户没有角色，返回空权限树
		if len(uroleIds) == 0 {
			return []Permission{}, nil
		}
		// 查询用户的权限 ID 列表
		queryPermissionIDs := `SELECT permissionId FROM role_permissions_permission WHERE roleId IN (?)`
		query, inArgs, err := sqlx.In(queryPermissionIDs, uroleIds)
		query = DB.Rebind(query) // 重新绑定查询语句
		if err != nil {
			return nil, errors.New("构建权限查询失败")
		}

		var rpermisId []int64
		err = DB.Select(&rpermisId, query, inArgs...)
		if err != nil {
			return nil, errors.New("查询角色权限失败")
		}
		// 如果用户没有权限，返回空权限树

		if len(rpermisId) == 0 {
			return []Permission{}, nil
		}

		// 添加权限过滤条件
		permQuery, permArgs, err := sqlx.In("id IN (?)", rpermisId)
		permQuery = DB.Rebind(permQuery) // 重新绑定查询语句
		if err != nil {
			return nil, errors.New("构建权限过滤条件失败")

		}

		baseQuery += " AND " + permQuery
		args = append(args, permArgs...)

	}
	// 查询一级权限
	var onePermissList []Permission
	err = DB.Select(&onePermissList, baseQuery, args...)
	if err != nil {
		return nil, errors.New("查询一级权限失败")
	}

	// 构建权限树
	for i, perm := range onePermissList {
		// 查询二级权限
		var twoPerissList []Permission
		err = DB.Select(&twoPerissList,
			"SELECT * FROM permission WHERE parentId = ? ORDER BY `order` ASC",
			perm.ID)
		if err != nil {
			return nil, errors.New("查询二级权限失败")
		}

		for i2, perm2 := range twoPerissList {
			// 查询三级权限
			var twoPerissList2 []Permission
			err = DB.Select(&twoPerissList2,
				"SELECT * FROM permission WHERE parentId = ? ORDER BY `order` ASC",
				perm2.ID)
			if err != nil {
				return nil, errors.New("查询三级权限失败")

			}
			twoPerissList[i2].Children = twoPerissList2
		}
		onePermissList[i].Children = twoPerissList
	}

	return onePermissList, nil
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

func GetRolesCountByWhereName(name string) (int64, error) {
	var total int64
	var args []interface{}

	query := "SELECT COUNT(*) FROM roles"
	if name != "" {
		whereClause := " WHERE name LIKE ?"
		query += whereClause
		args = append(args, "%"+name+"%")

	}
	// 执行计数查询
	err := DB.Get(&total, query, args...)
	if err != nil {
		return 0, err
	}
	return total, nil

}
func GetRolesCountWhereByName(name string, enable int, pageNo, pageSize int) ([]*Role, int64, error) {
	baseQuery := "SELECT * FROM role WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM role WHERE 1=1"
	var args []interface{}
	var total int64
	if name != "" {
		whereClause := " AND name LIKE ?"
		baseQuery += whereClause
		countQuery += whereClause
		args = append(args, "%"+name+"%")
	}
	// 执行计数查询
	err := DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	// 添加分页条件
	pageQuery := baseQuery + " LIMIT ? OFFSET ?"
	pageArgs := append(args, pageSize, (pageNo-1)*pageSize)
	var userList []*Role
	// // 执行分页查询
	err = DB.Select(&userList, pageQuery, pageArgs...)
	if err != nil {
		return nil, 0, errors.New("查询角色列表失败")
	}
	return userList, total, nil
}

func GetRolesList() ([]*Role, error) {
	var roles []*Role
	err := DB.Select(&roles, "SELECT * FROM role")

	if err != nil {
		return nil, errors.New("查询角色失败")
	}
	return roles, nil
}
