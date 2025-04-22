package service

import (
	"errors"

	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/models"
)

type RoleService struct {
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (r *RoleService) GetPermissionsTree(userID int) ([]models.Permission, error) {
	// 调用 Models 层的方法获取权限树
	permissions, err := models.GetPermissionsTree(userID)
	if err != nil {
		return nil, err
	}

	if len(permissions) == 0 {
		return nil, errors.New("no permissions found for the user")
	}

	return permissions, nil
}

func (r *RoleService) List() ([]*models.Role, error) {
	data, err := models.GetRolesList()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RoleService) ListPage(enable int, username string, pageNo, pageSize int) (*inout.RoleListPageRes, error) {
	var data = inout.RoleListPageRes{
		PageData: make([]inout.RoleListPageItem, 0),
	}
	users, total, err := models.GetRolesCountWhereByName(username, enable, pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询用户资料信息失败")
	}
	data.Total = total
	for i, user := range users {
		var perIdList []int

		perIdList, err = models.GetPermissionsIdsByWhere(user.ID)
		if err != nil {
			return nil, err
		}
		data.PageData[i].PermissionIds = perIdList
	}
	return &data, nil
}

func (r *RoleService) Update() {

}

func (r *RoleService) Add() {

}

func (r *RoleService) Delete() {

}

func (r *RoleService) AddUser() {

}
func (r *RoleService) RemoveUser() {

}
