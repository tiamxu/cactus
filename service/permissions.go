package service

import (
	"errors"

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
