package service

import (
	"errors"

	"github.com/tiamxu/cactus/logic/model"
	"github.com/tiamxu/cactus/logic/repo"
	"github.com/tiamxu/cactus/types"
)

type RoleService struct {
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (r *RoleService) GetPermissionsTree(userID int) ([]model.Permission, error) {
	// 调用 Models 层的方法获取权限树
	permissions, err := repo.GetPermissionsTree(userID)
	if err != nil {
		return nil, err
	}

	if len(permissions) == 0 {
		return nil, errors.New("no permissions found for the user")
	}

	return permissions, nil
}

func (r *RoleService) List() ([]*model.Role, error) {
	data, err := repo.GetRolesList()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RoleService) ListPage(enable, username string, pageNo, pageSize int) (*types.RoleListPageRes, error) {
	var data = types.RoleListPageRes{
		PageData: make([]types.RoleListPageItem, 0),
	}
	roles, total, err := repo.GetRolesCountWhereByNameEnable(username, enable, pageNo, pageSize)
	if err != nil {
		return nil, errors.New("查询角色信息失败")
	}
	data.Total = total
	if len(roles) == 0 {
		return &data, nil
	}
	//预分配足够容量的切片
	data.PageData = make([]types.RoleListPageItem, len(roles))

	for i, role := range roles {
		perIdList, err := repo.GetPermissionsIdsByWhere(role.ID)
		if err != nil {
			return nil, err
		}
		data.PageData[i] = types.RoleListPageItem{
			Role:          *role,
			PermissionIds: perIdList,
		}

	}

	return &data, nil
}

func (r *RoleService) Update(req types.PatchRoleReq) error {
	err := repo.UpdateRoleWhereByCondition(&req.Id, req.Name, req.Code, req.Enable, req.PermissionIds)
	if err != nil {
		return nil
	}
	return nil
}

func (r *RoleService) Add(params types.AddRoleReq) error {
	err := repo.AddRoleWhereByCondition(params.Name, params.Code, params.Enable, params.PermissionIds)
	if err != nil {
		return nil
	}
	return nil
}

func (r *RoleService) Delete(id string) error {
	err := repo.DeleteRolesWhereById(id)
	if err != nil {
		return nil
	}
	return nil
}

func (r *RoleService) AddUser(params types.PatchRoleOpeateUserReq) error {
	err := repo.AddUserRolesByWhereId(params.UserIds, params.Id)
	if err != nil {
		return nil
	}
	return nil
}
func (r *RoleService) RemoveUser(params types.PatchRoleOpeateUserReq) error {
	err := repo.RemoveUserRolesByWhereId(params.UserIds, params.Id)
	if err != nil {
		return nil
	}
	return nil
}
