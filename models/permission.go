package models

import (
	"errors"
)

type Permission struct {
	ID          int          `db:"id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Code        string       `db:"code" json:"code"`
	Type        string       `db:"type" json:"type"`
	ParentId    *int         `db:"parentId" json:"parentId"`
	Path        string       `db:"path" json:"path"`
	Redirect    string       `db:"redirect" json:"redirect"`
	Icon        string       `db:"icon" json:"icon"`
	Component   string       `db:"component" json:"component"`
	Layout      string       `db:"layout" json:"layout"`
	KeepAlive   int          `db:"keepAlive" json:"keepAlive"`
	Method      string       `db:"method" json:"method"`
	Description string       `db:"description" json:"description"`
	Show        int          `db:"show" json:"show"`
	Enable      int          `db:"enable" json:"enable"`
	Order       int          `db:"order" json:"order"`
	Children    []Permission `db:"children" json:"children"`
}

func (Permission) TableName() string {
	return "permission"
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

func InsertPermission(p Permission) error {
	query := `
        INSERT INTO permission (
            name, code, type, parentId, path, 
            icon, component, layout, keepAlive,
            ` + "`show`" + `, enable, ` + "`order`" + `
        ) VALUES (
            :name, :code, :type, :parentId, :path, 
            :icon, :component, :layout, :keepAlive, 
            :show, :enable, :order
        )`

	// args := []interface{}{
	// 	p.Name,
	// 	p.Code,
	// 	p.Type,
	// 	p.ParentId, // 默认parent_id为NULL
	// 	p.Path,
	// 	p.Icon,
	// 	p.Component,
	// 	p.Layout,
	// 	p.KeepAlive,
	// 	p.Show,
	// 	p.Enable,
	// 	p.Order,
	// }
	args := map[string]interface{}{
		"name":      p.Name,
		"code":      p.Code,
		"type":      p.Type,
		"parentId":  p.ParentId,
		"path":      p.Path,
		"icon":      p.Icon,
		"component": p.Component,
		"layout":    p.Layout,
		"keepAlive": p.KeepAlive,
		"method":    p.Method, // 修正了原来的错误
		"show":      p.Show,
		"enable":    p.Enable,
		"order":     p.Order,
	}
	// 处理ParentId，如果不为0则设置值
	// if p.ParentId != nil {
	// 	args[3] = p.ParentId
	// }

	// 执行插入操作
	_, err := DB.NamedExec(query, args)
	if err != nil {
		return err
	}
	return nil
}

func DeletePermissionByWhereId(permId string) error {
	tx, err := DB.Beginx()
	if err != nil {
		return errors.New("事务开启失败")
	}
	// 删除角色权限关联表中的记录
	_, err = tx.Exec("DELETE FROM role_permissions_permission WHERE permissionId = ?", permId)
	if err != nil {
		tx.Rollback()
		return errors.New("删除角色权限关联失败")

	}

	// 删除权限表中的记录
	_, err = tx.Exec("DELETE FROM permission WHERE id = ?", permId)
	if err != nil {
		tx.Rollback()
		return errors.New("删除权限记录失败")

	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return errors.New("事务提交失败")
	}
	return nil
}

func UpdatePermissionByWhere(p Permission) error {
	query := `
        UPDATE permission
        SET name = :name,
            code = :code,
            type = :type,
            parent_id = :parent_id,
            path = :path,
            icon = :icon,
            component = :component,
            layout = :layout,
            keep_alive = :keep_alive,
            method = :method,
            show = :show,
            enable = :enable,
            "order" = :order
        WHERE id = :id
    `

	args := map[string]interface{}{
		"id":        p.ID,
		"name":      p.Name,
		"code":      p.Code,
		"type":      p.Type,
		"parentId":  p.ParentId,
		"path":      p.Path,
		"icon":      p.Icon,
		"component": p.Component,
		"layout":    p.Layout,
		"keepAlive": p.KeepAlive,
		"method":    p.Method, // 修正了原来的错误
		"show":      p.Show,
		"enable":    p.Enable,
		"order":     p.Order,
	}

	_, err := DB.NamedExec(query, args)
	if err != nil {
		return err
	}
	return nil
}
