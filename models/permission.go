package models

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

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
	// 检查是否是管理员
	var adminRole int64
	err := DB.Get(&adminRole,
		"SELECT COUNT(*) FROM user_roles_role WHERE userId = ? AND roleId = 1",
		userID)
	if err != nil {
		return nil, errors.New("查询管理员状态失败")
	}
	// 构建基础查询
	baseQuery := "SELECT * FROM permissions WHERE parentId IS NULL ORDER BY `order` ASC"
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
			"SELECT * FROM permissions WHERE parentId = ? ORDER BY `order` ASC",
			perm.ID)
		if err != nil {
			return nil, errors.New("查询二级权限失败")
		}

		for i2, perm2 := range twoPerissList {
			// 查询三级权限
			var twoPerissList2 []Permission
			err = DB.Select(&twoPerissList2,
				"SELECT * FROM permissions WHERE parentId = ? ORDER BY `order` ASC",
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

func GetPermissionsList() ([]Permission, error) {
	var onePermissList []Permission
	err := DB.Select(&onePermissList, "SELECT * FROM permission WHERE parentId IS NULL ORDER BY `order` ASC")
	if err != nil {
		return nil, errors.New("查询一级权限失败")
	}

	// 遍历一级权限，查询二级权限
	for i, perm := range onePermissList {
		var twoPerissList []Permission
		err := DB.Select(&twoPerissList, "SELECT * FROM permission WHERE parentId = ? ORDER BY `order` ASC", perm.ID)
		if err != nil {
			return nil, errors.New("询二级权限失败")
		}
		// 遍历二级权限，查询三级权限
		for i2, perm2 := range twoPerissList {
			var twoPerissList2 []Permission
			err := DB.Select(&twoPerissList2, "SELECT * FROM permission WHERE parentId = ? ORDER BY `order` ASC", perm2.ID)
			if err != nil {
				return nil, errors.New("查询三级权限失败")
			}
			twoPerissList[i2].Children = twoPerissList2
		}

		onePermissList[i].Children = twoPerissList
	}

	return onePermissList, nil
}
