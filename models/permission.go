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
