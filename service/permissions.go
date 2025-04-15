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

func (s *PermissionsService) GetPermissionsTree(userID int) ([]models.Permission, error) {
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
