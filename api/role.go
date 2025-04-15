package api

import "github.com/tiamxu/cactus/service"

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: &service.RoleService{},
	}
}
