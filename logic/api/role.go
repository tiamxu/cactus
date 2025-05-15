package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/logic/service"
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
	var name = c.DefaultQuery("name", "")
	var enable = c.DefaultQuery("enable", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	data, err := r.roleService.ListPage(enable, name, pageNo, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Resp.Succ(c, data)
}

func (r *RoleHandler) Update(c *gin.Context) {
	var params inout.PatchRoleReq
	err := c.ShouldBindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, "参数解析失败: "+err.Error())
		return
	}
	err = r.roleService.Update(params)
	if err != nil {
		return
	}
	Resp.Succ(c, err)
}

func (r *RoleHandler) Add(c *gin.Context) {
	var params inout.AddRoleReq
	err := c.ShouldBind(&params)
	if err != nil {
		Resp.Err(c, 20001, "参数解析失败: "+err.Error())
		return
	}
	err = r.roleService.Add(params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")
}

func (r *RoleHandler) Delete(c *gin.Context) {
	roleID := c.Param("id")
	err := r.roleService.Delete(roleID)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	Resp.Succ(c, "")
}

func (r *RoleHandler) AddUser(c *gin.Context) {
	var params inout.PatchRoleOpeateUserReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	uid, _ := strconv.Atoi(c.Param("id"))
	params.Id = uid
	r.roleService.AddUser(params)
	Resp.Succ(c, "")
}

func (r *RoleHandler) RemoveUser(c *gin.Context) {
	var params inout.PatchRoleOpeateUserReq
	err := c.BindJSON(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	uid, _ := strconv.Atoi(c.Param("id"))
	params.Id = uid
	r.roleService.RemoveUser(params)
	Resp.Succ(c, "")
}
