package api

import (
	"net/http"

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
	// 从上下文中获取用户 ID
	userIDInterface, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID not found in context"})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 调用 Service 层获取权限树
	permissions, err := r.roleService.GetPermissionsTree(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回权限树
	Resp.Succ(c, permissions)
}
func (r *RoleHandler) List(c *gin.Context) {
	data, err := r.roleService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
