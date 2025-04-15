package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (a *PermissionsHandler) PermissionsTree(c *gin.Context) {
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
	permissions, err := a.permissionsService.GetPermissionsTree(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回权限树
	Resp.Succ(c, permissions)
}
