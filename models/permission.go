package models

import "github.com/jmoiron/sqlx"

type Permission struct {
	ID          int          `db:"id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Code        string       `db:"code" json:"code"`
	Type        string       `db:"type" json:"type"`
	ParentId    *int         `db:"parentId" json:"parentId"`
	Path        *string      `db:"path" json:"path"`
	Redirect    *string      `db:"redirect" json:"redirect"`
	Icon        *string      `db:"icon" json:"icon"`
	Component   *string      `db:"component" json:"component"`
	Layout      *string      `db:"layout" json:"layout"`
	KeepAlive   *int         `db:"keepAlive" json:"keepAlive"`
	Method      *string      `db:"method" json:"method"`
	Description *string      `db:"description" json:"description"`
	Show        int          `db:"show" json:"show"`
	Enable      int          `db:"enable" json:"enable"`
	Order       *int         `db:"order" json:"order"`
	Children    []Permission `db:"children" json:"children"`
}

func (Permission) TableName() string {
	return "permission"
}

func GetPermissionsTree(userID int) ([]Permission, error) {
	// 查询用户的角色 ID 列表
	var roleIDs []int
	queryRoleIDs := `SELECT roleId FROM user_roles_role WHERE userId = ?`
	if err := DB.Select(&roleIDs, queryRoleIDs, userID); err != nil {
		return nil, err
	}

	// 如果用户没有角色，返回空权限树
	if len(roleIDs) == 0 {
		return []Permission{}, nil
	}

	// 查询用户的权限 ID 列表
	var permissionIDs []int
	queryPermissionIDs := `SELECT permissionId FROM role_permissions_permission WHERE roleId IN (?)`
	query, args, err := sqlx.In(queryPermissionIDs, roleIDs)
	if err != nil {
		return nil, err
	}
	query = DB.Rebind(query) // 重新绑定查询语句
	if err := DB.Select(&permissionIDs, query, args...); err != nil {
		return nil, err
	}

	// 如果用户没有权限，返回空权限树
	if len(permissionIDs) == 0 {
		return []Permission{}, nil
	}

	// 查询所有权限，并构建树形结构
	queryPermissions := `
		SELECT 
			p.id, p.name, p.code, p.type, p.parentId, p.path, p.icon, p.redirect, p.component , p.layout,
			p.keepAlive, p.method, p.description, p.show, p.enable, p.order
		FROM permission p
		WHERE p.id IN (?)
		ORDER BY p.order ASC
	`
	queryPermissions, args, err = sqlx.In(queryPermissions, permissionIDs)
	if err != nil {
		return nil, err
	}
	queryPermissions = DB.Rebind(queryPermissions) // 重新绑定查询语句
	var permissions []Permission
	if err := DB.Select(&permissions, queryPermissions, args...); err != nil {
		return nil, err
	}

	// 构建树形结构
	permissionMap := make(map[int]*Permission)
	var roots []Permission

	for i := range permissions {
		permissions[i].Children = []Permission{}
		permissionMap[permissions[i].ID] = &permissions[i]
		if permissions[i].ParentId == nil {
			roots = append(roots, permissions[i])
		}
	}

	for i := range permissions {
		if permissions[i].ParentId != nil {
			parent, exists := permissionMap[*permissions[i].ParentId]
			if exists {
				parent.Children = append(parent.Children, permissions[i])
			}
		}
	}

	return roots, nil
}
