package repo

import (
	"errors"

	"github.com/tiamxu/cactus/logic/model"
)

func GetPermissionsList() ([]model.Permission, error) {
	var onePermissList []model.Permission
	err := DB.Select(&onePermissList, "SELECT * FROM permission WHERE parentId IS NULL ORDER BY `order` ASC")
	if err != nil {
		return nil, errors.New("查询一级权限失败")
	}
	// 遍历一级权限，查询二级权限
	for i, perm := range onePermissList {
		var twoPerissList []model.Permission
		err := DB.Select(&twoPerissList, "SELECT * FROM permission WHERE parentId = ? ORDER BY `order` ASC", perm.ID)
		if err != nil {
			return nil, errors.New("询二级权限失败")
		}
		// 遍历二级权限，查询三级权限
		for i2, perm2 := range twoPerissList {
			var twoPerissList2 []model.Permission
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

func InsertPermission(p model.Permission) error {
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

func UpdatePermissionByWhere(p model.Permission) error {
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
