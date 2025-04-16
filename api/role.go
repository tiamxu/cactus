package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/service"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: &service.RoleService{},
	}
}

func (r *RoleHandler) PermissionsTree(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (r *RoleHandler) List(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (r *RoleHandler) ListPage(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (r *RoleHandler) Update(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (r *RoleHandler) Add(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (r *RoleHandler) Delete(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (r *RoleHandler) AddUser(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (r *RoleHandler) RemoveUser(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}
