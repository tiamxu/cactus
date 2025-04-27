package service

import (
	"errors"
	"fmt"

	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/models"
)

type PermissionsService struct {
}

func NewPermissionsServiceService() *PermissionsService {
	return &PermissionsService{}
}

func (p *PermissionsService) List() ([]models.Permission, error) {
	data, err := models.GetPermissionsList()
	if err != nil {
		return nil, errors.New("获取权限列表错误")
	}
	return data, nil
}

func (p *PermissionsService) ListPage(username string, pageNo, pageSize int) (*inout.RoleListPageRes, error) {
	var data = inout.RoleListPageRes{
		PageData: make([]inout.RoleListPageItem, 0),
	}
	roles, total, err := models.GetRolesCountWhereByName(username, pageNo, pageSize)
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

func (p *PermissionsService) Add(params inout.AddPermissionReq) error {
	perm := models.Permission{
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		ParentId:  params.ParentId,
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: IsTrue(params.KeepAlive),
		Method:    params.Component,
		Show:      IsTrue(params.Show),
		Enable:    IsTrue(params.Enable),
		Order:     params.Order,
	}
	err := models.InsertPermission(perm)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	return nil
}

func (p *PermissionsService) Delete(id string) error {

	err := models.DeletePermissionByWhereId(id)
	if err != nil {
		return errors.New("删除权限错误")
	}
	return nil
}

func (p *PermissionsService) PatchPermission(params inout.PatchPermissionReq) error {
	perm := models.Permission{
		ID:        params.Id,
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		ParentId:  params.ParentId,
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: params.KeepAlive,
		Method:    params.Component,
		Show:      params.Show,
		Enable:    params.Enable,
		Order:     params.Order,
	}
	err := models.UpdatePermissionByWhere(perm)
	if err != nil {
		return errors.New("更新权限信息失败")
	}
	return nil
}

func IsTrue(v bool) int {
	if v {
		return 1
	}
	return 0
}
