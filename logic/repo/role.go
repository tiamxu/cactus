package repo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/tiamxu/cactus/logic/model"
)

func GetPermissionsTree(userID int) ([]model.Permission, error) {
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
			return []model.Permission{}, nil
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
			return []model.Permission{}, nil
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
	var onePermissList []model.Permission
	err = DB.Select(&onePermissList, baseQuery, args...)
	if err != nil {
		fmt.Println("error:", err)
		return nil, errors.New("查询一级权限失败")
	}

	// 构建权限树
	for i, perm := range onePermissList {
		// 查询二级权限
		var twoPerissList []model.Permission
		err = DB.Select(&twoPerissList,
			"SELECT * FROM permission WHERE parentId = ? ORDER BY `order` ASC",
			perm.ID)
		if err != nil {
			return nil, errors.New("查询二级权限失败")
		}

		for i2, perm2 := range twoPerissList {
			// 查询三级权限
			var twoPerissList2 []model.Permission
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
func GetRolesByID(roleIDs []int) ([]*model.Role, error) {
	query, args, err := sqlx.In(`SELECT id, code, name, enable FROM role WHERE id IN (?)`, roleIDs)
	if err != nil {
		return nil, err
	}
	query = DB.Rebind(query) // 重新绑定查询语句
	var roles []*model.Role
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

func GetRolesByUserId(userId int) ([]*model.Role, error) {
	query := `
		SELECT r.* 
		FROM role r
		JOIN user_roles_role urr ON r.id = urr.roleId
		WHERE urr.userId = ?`

	var roles []*model.Role
	err := DB.Select(&roles, query, userId)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func GetRolesByUserIds(userIds []int) (map[int][]*model.Role, error) {
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

	var userRoles []model.UserRolesRole
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
	var roles []model.Role
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
	roleMap := make(map[int]model.Role)
	for _, role := range roles {
		roleMap[role.ID] = role
	}

	// 构建最终结果：用户ID到角色列表的映射
	result := make(map[int][]*model.Role)
	for userId, rIds := range userRoleMap {
		for _, rId := range rIds {
			if role, ok := roleMap[rId]; ok {
				result[userId] = append(result[userId], &role)
			}
		}
	}

	return result, nil
}

// func GetRolesCountByWhereName(name string) (int64, error) {
// 	var total int64
// 	var args []interface{}

// 	query := "SELECT COUNT(*) FROM roles"
// 	if name != "" {
// 		whereClause := " WHERE name LIKE ?"
// 		query += whereClause
// 		args = append(args, "%"+name+"%")

// 	}
// 	// 执行计数查询
// 	err := DB.Get(&total, query, args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return total, nil

// }
func GetRolesCountWhereByNameEnable(name string, enable string, pageNo, pageSize int) ([]*model.Role, int64, error) {
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
	if enable != "" {
		ena := enable == "1"
		baseQuery += " AND enable = ?"
		countQuery += " AND enable = ?"
		args = append(args, ena)
	}
	// countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") AS t"

	// 执行计数查询
	err := DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	// 添加分页条件
	pageQuery := baseQuery + " ORDER BY id LIMIT ? OFFSET ?"
	pageArgs := append(args, pageSize, (pageNo-1)*pageSize)
	var roleList []*Role
	// // 执行分页查询
	err = DB.Select(&roleList, pageQuery, pageArgs...)
	if err != nil {
		return nil, 0, errors.New("查询角色列表失败")
	}
	return roleList, total, nil
}

func GetRolesCountWhereByName(name string, pageNo, pageSize int) ([]*model.Role, int64, error) {
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
	// countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") AS t"

	// 执行计数查询
	err := DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	// 添加分页条件
	pageQuery := baseQuery + " LIMIT ? OFFSET ?"
	pageArgs := append(args, pageSize, (pageNo-1)*pageSize)
	var roleList []*Role
	// // 执行分页查询
	err = DB.Select(&roleList, pageQuery, pageArgs...)
	if err != nil {
		return nil, 0, errors.New("查询角色列表失败")
	}
	return roleList, total, nil
}
func GetRolesList() ([]*Role, error) {
	var roles []*Role
	err := DB.Select(&roles, "SELECT * FROM role")

	if err != nil {
		return nil, errors.New("查询角色失败")
	}
	return roles, nil
}

func UpdateRoleWhereByCondition(roleId *int, name, code *string, enable *bool, permissionIds []int) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	if name != nil || enable != nil || code != nil {
		updateQuery := "UPDATE role SET "
		var setClauses []string
		var args []interface{}

		if name != nil {
			setClauses = append(setClauses, "name = ?")
			args = append(args, *name)
		}
		if enable != nil {
			setClauses = append(setClauses, "enable = ?")
			args = append(args, *enable)
		}
		if code != nil {
			setClauses = append(setClauses, "code = ?")
			args = append(args, *code)
		}

		updateQuery += strings.Join(setClauses, ", ") + " WHERE id = ?"
		args = append(args, roleId)

		_, err = tx.Exec(updateQuery, args...)
		if err != nil {
			tx.Rollback()
			return errors.New("更新角色信息失败")
		}
	}

	// 更新角色权限关联
	if permissionIds != nil {
		// 删除旧权限关联
		_, err = tx.Exec("DELETE FROM role_permissions_permission WHERE roleId = ?", roleId)
		if err != nil {
			tx.Rollback()
			return errors.New("删除旧权限关联失败")
		}

		// 添加新权限关联
		if len(permissionIds) > 0 {
			// 构建批量插入语句
			query := "INSERT INTO role_permissions_permission (permissionId, roleId) VALUES "
			var placeholders []string
			var insertArgs []interface{}

			for _, permId := range permissionIds {
				placeholders = append(placeholders, "(?, ?)")
				insertArgs = append(insertArgs, permId, roleId)
			}

			query += strings.Join(placeholders, ", ")
			_, err = tx.Exec(query, insertArgs...)
			if err != nil {
				tx.Rollback()
				return errors.New("添加新权限关联失败")
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return errors.New("事务提交失败")
	}
	return nil
}
func AddRoleWhereByCondition(name string, code string, enable bool, permissionIds []int) error {
	tx, err := DB.Beginx()
	if err != nil {
		return errors.New("事务开启失败")

	}

	// 插入角色记录
	res, err := tx.Exec(
		"INSERT INTO role (code, name, enable) VALUES (?, ?, ?)",
		code, name, enable)
	if err != nil {
		tx.Rollback()
		return errors.New("创建角色失败")

	}

	// 获取自增ID
	roleID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return errors.New("获取角色ID失败")

	}

	// 添加权限关联
	if len(permissionIds) > 0 {
		// 构建批量插入语句
		query := "INSERT INTO role_permissions_permission (roleId, permissionId) VALUES "
		var placeholders []string
		var args []interface{}

		for _, permID := range permissionIds {
			placeholders = append(placeholders, "(?, ?)")
			args = append(args, roleID, permID)
		}

		query += strings.Join(placeholders, ", ")
		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return errors.New("添加权限关联失败")

		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return errors.New("事务提交失败")

	}
	return nil
}
func DeleteRolesWhereById(roleId string) error {
	tx, err := DB.Beginx()
	if err != nil {
		return errors.New("事务开启失败")
	}

	//  删除用户角色关联
	_, err = tx.Exec("DELETE FROM user_roles_role WHERE roleId = ?", roleId)
	if err != nil {
		tx.Rollback()
		return errors.New("删除用户角色关联失败")

	}

	//  删除角色权限关联
	_, err = tx.Exec("DELETE FROM role_permissions_permission WHERE roleId = ?", roleId)
	if err != nil {
		tx.Rollback()
		return errors.New("删除角色权限关联失败")

	}

	// 删除角色记录
	_, err = tx.Exec("DELETE FROM role WHERE id = ?", roleId)
	if err != nil {
		tx.Rollback()
		return errors.New("删除角色记录失败")

	}

	//  提交事务
	if err := tx.Commit(); err != nil {
		return errors.New("事务提交失败")

	}
	return nil
}

func AddUserRolesByWhereId(userIds []int, roleId int) error {
	tx, err := DB.Beginx()
	if err != nil {
		return errors.New("事务开启失败")
	}

	// 5. 删除现有关联
	if len(userIds) > 0 {
		// 构建 IN 条件
		query, args, err := sqlx.In("DELETE FROM user_roles_role WHERE userId IN (?) AND roleId = ?", userIds, roleId)
		if err != nil {
			tx.Rollback()
			return errors.New("删除旧关联失败")

		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return errors.New("删除旧关联失败")

		}
	}

	// 6. 添加新关联
	if len(userIds) > 0 {
		// 构建批量插入语句
		query := "INSERT INTO user_roles_role (userId, roleId) VALUES "
		var placeholders []string
		var insertArgs []interface{}

		for _, userID := range userIds {
			placeholders = append(placeholders, "(?, ?)")
			insertArgs = append(insertArgs, userID, roleId)
		}

		query += strings.Join(placeholders, ", ")
		_, err = tx.Exec(query, insertArgs...)
		if err != nil {
			tx.Rollback()
			return errors.New("添加新关联失败")

		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return errors.New("事务提交失败")

	}
	return nil
}

func RemoveUserRolesByWhereId(userIds []int, roleId int) error {
	if len(userIds) == 0 {
		return errors.New("用户id为空")

	}
	query, args, err := sqlx.In("DELETE FROM user_roles_role WHERE userId IN (?) AND roleId = ?", userIds, roleId)
	if err != nil {
		return errors.New("构建删除语句失败")
	}

	_, err = DB.Exec(query, args...)
	if err != nil {
		return errors.New("删除关联失败")
	}
	return nil
}
