package service

import (
	"errors"
	"fmt"

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

func (r *RoleService) ListPage(enable, username string, pageNo, pageSize int) (*inout.RoleListPageRes, error) {
	var data = inout.RoleListPageRes{
		PageData: make([]inout.RoleListPageItem, 0),
	}
	roles, total, err := models.GetRolesCountWhereByName(username, enable, pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询角色信息失败")
	}
	data.Total = total
	if len(roles) == 0 {
		return &data, nil
	}
	//预分配足够容量的切片
	data.PageData = make([]inout.RoleListPageItem, len(roles))

	for i, role := range roles {
		perIdList, err := models.GetPermissionsIdsByWhere(role.ID)
		if err != nil {
			return nil, err
		}
		data.PageData[i] = inout.RoleListPageItem{
			Role:          *role,
			PermissionIds: perIdList,
		}

	}

	return &data, nil
}

func (r *RoleService) Update(req inout.PatchRoleReq) error {
	err := models.UpdateRoleWhereByCondition(&req.Id, req.Name, req.Code, req.Enable, req.PermissionIds)
	if err != nil {
		return nil
	}
	return nil
}

func (r *RoleService) Add(params inout.AddRoleReq) error {
	fmt.Println("###params", params)
	err := models.AddRoleWhereByCondition(params.Name, params.Code, params.Enable, params.PermissionIds)
	if err != nil {
		return nil
	}
	return nil
}

func (r *RoleService) Delete(id string) error {
	err := models.DeleteRolesWhereById(id)
	if err != nil {
		return nil
	}
	return nil
}

func (r *RoleService) AddUser(params inout.PatchRoleOpeateUserReq) error {
	err := models.AddUserRolesByWhereId(params.UserIds, params.Id)
	if err != nil {
		return nil
	}
	return nil
}
func (r *RoleService) RemoveUser(params inout.PatchRoleOpeateUserReq) error {
	err := models.RemoveUserRolesByWhereId(params.UserIds, params.Id)
	if err != nil {
		return nil
	}
	return nil
}
