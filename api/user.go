package api

import (
	"net/http"
	"strconv"

	"github.com/tiamxu/cactus/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: &service.UserService{},
	}
}
func (a *UserHandler) Detail(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid not found in context"})
		return
	}

	userID := uid.(int)

	userDetail, err := a.userService.GetUserDetail(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(userDetail.Roles) > 0 {
		userDetail.CurrentRole = userDetail.Roles[0]
	}
	Resp.Succ(c, userDetail)
}
func (a *UserHandler) List(c *gin.Context) {
	gender := c.DefaultQuery("gender", "")
	enable := c.DefaultQuery("enable", "")
	username := c.DefaultQuery("username", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)

	userListRes, err := a.userService.GetUserList(gender, enable, username, pageNo, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Resp.Succ(c, userListRes)

}
func (a *UserHandler) Profile(c *gin.Context) {
}
func (a *UserHandler) Update(c *gin.Context) {
}
func (a *UserHandler) Add(c *gin.Context) {
}
func (a *UserHandler) Delete(c *gin.Context) {
}

// CreateUser 创建用户
// func (h *UserHandler) CreateUser(c *gin.Context) {
// 	var req service.CreateUserRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		response.Error(c, 400, "无效的请求参数")
// 		return
// 	}

// 	if err := h.userService.Create(&req); err != nil {
// 		response.Error(c, 400, err.Error())
// 		return
// 	}

// 	response.Success(c, nil)
// }

// GetUser 获取用户详情
// func (h *UserHandler) GetUser(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		response.Error(c, 400, "无效的用户ID")
// 		return
// 	}

// 	user, err := h.userService.GetByID(uint(id))
// 	if err != nil {
// 		response.Error(c, 404, "用户不存在")
// 		return
// 	}

// 	response.Success(c, user)
// }

// ListUsers 用户列表
// func (h *UserHandler) ListUsers(c *gin.Context) {
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

// 	users, err := h.userService.List(page, pageSize)
// 	if err != nil {
// 		fmt.Errorf("获取用户列表失败:", "error", err)
// 		response.Error(c, 500, "服务器内部错误")
// 		return
// 	}

// 	response.Success(c, gin.H{
// 		"items": users,
// 		"total": len(users),
// 	})
// }

// UpdateUser 更新用户
// func (h *UserHandler) UpdateUser(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		response.Error(c, 400, "无效的用户ID")
// 		return
// 	}

// 	var req service.UpdateUserRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		response.Error(c, 400, "无效的请求参数")
// 		return
// 	}

// 	user, err := h.userService.Update(uint(id), &req)
// 	if err != nil {
// 		response.Error(c, 400, err.Error())
// 		return
// 	}

// 	response.Success(c, user)
// }

// Delete删除用户
// func (h *UserHandler) Delete(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		response.Error(c, 400, "无效的用户ID")
// 		return
// 	}

// 	if err := h.userService.Delete(uint(id)); err != nil {
// 		fmt.Errorf("删除用户失败",
// 			"error", err.Error(),
// 			"user_id", id)
// 		response.Error(c, 500, "删除用户失败")
// 		return
// 	}

// 	response.Success(c, nil)
// }
