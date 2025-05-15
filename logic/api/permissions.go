package api

import (
	"net/http"
	"strconv"

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
	var data = &inout.RoleListPageRes{}
	var name = c.DefaultQuery("name", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)

	data, err := p.permissionsService.ListPage(name, pageNo, pageSize)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	Resp.Succ(c, data)
}
func (p *PermissionsHandler) Add(c *gin.Context) {
	var params inout.AddPermissionReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = p.permissionsService.Add(params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}

func (p *PermissionsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := p.permissionsService.Delete(id)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}

func (p *PermissionsHandler) PatchPermission(c *gin.Context) {
	var params inout.PatchPermissionReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = p.permissionsService.PatchPermission(params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}
