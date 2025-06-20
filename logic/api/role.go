package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/logic/service"
	"github.com/tiamxu/cactus/types"
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
		c.JSON(http.StatusBadRequest, RespError(c, nil, "user ID not found in context"))

		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, RespError(c, nil, "invalid user ID"))

		return
	}

	// 调用 Service 层获取权限树
	permissions, err := r.roleService.GetPermissionsTree(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}

	// 返回权限树
	c.JSON(http.StatusOK, RespSuccess(c, permissions))

}
func (r *RoleHandler) List(c *gin.Context) {
	data, err := r.roleService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, data))

}

func (r *RoleHandler) ListPage(c *gin.Context) {
	var name = c.DefaultQuery("name", "")
	var enable = c.DefaultQuery("enable", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	data, err := r.roleService.ListPage(enable, name, pageNo, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, data))
}

func (r *RoleHandler) Update(c *gin.Context) {
	var params types.PatchRoleReq
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数解析失败"))

		return
	}
	err = r.roleService.Update(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))

}

func (r *RoleHandler) Add(c *gin.Context) {
	var params types.AddRoleReq
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数解析失败"))

		return
	}
	err = r.roleService.Add(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))

}

func (r *RoleHandler) Delete(c *gin.Context) {
	roleID := c.Param("id")
	err := r.roleService.Delete(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, ""))

}

func (r *RoleHandler) AddUser(c *gin.Context) {
	var params types.PatchRoleOpeateUserReq
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	uid, _ := strconv.Atoi(c.Param("id"))
	params.Id = uid
	r.roleService.AddUser(params)
	c.JSON(http.StatusOK, RespSuccess(c, ""))
}

func (r *RoleHandler) RemoveUser(c *gin.Context) {
	var params types.PatchRoleOpeateUserReq
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	uid, _ := strconv.Atoi(c.Param("id"))
	params.Id = uid
	r.roleService.RemoveUser(params)
	c.JSON(http.StatusOK, RespSuccess(c, ""))
}
