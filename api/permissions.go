package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/service"
)

type PermissionsHandler struct {
	permissionsService *service.PermissionsService
}

func NewPermissionsHandler() *PermissionsHandler {
	return &PermissionsHandler{
		permissionsService: &service.PermissionsService{},
	}
}

func (p *PermissionsHandler) List(c *gin.Context) {
	data, err := p.permissionsService.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取权限列表失败"})
		return
	}
	Resp.Succ(c, data)
}

func (p *PermissionsHandler) ListPage(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}
func (p *PermissionsHandler) Add(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (p *PermissionsHandler) Delete(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}

func (p *PermissionsHandler) PatchPermission(c *gin.Context) {
	var data = &inout.RoleListRes{}

	Resp.Succ(c, data)
}
