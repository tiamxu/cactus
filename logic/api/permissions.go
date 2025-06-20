package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/logic/service"
	"github.com/tiamxu/cactus/types"
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
		c.JSON(http.StatusBadRequest, RespError(c, err, "获取权限列表失败"))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, data))
}

func (p *PermissionsHandler) ListPage(c *gin.Context) {
	var data = &types.RoleListPageRes{}
	var name = c.DefaultQuery("name", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)

	data, err := p.permissionsService.ListPage(name, pageNo, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, data))

}
func (p *PermissionsHandler) Add(c *gin.Context) {
	var params types.AddPermissionReq
	err := c.Bind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))

		return
	}
	err = p.permissionsService.Add(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))

}

func (p *PermissionsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := p.permissionsService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))
}

func (p *PermissionsHandler) PatchPermission(c *gin.Context) {
	var params types.PatchPermissionReq
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	err = p.permissionsService.PatchPermission(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))
}
