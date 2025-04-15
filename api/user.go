package api

import (
	"fmt"
	"naive-admin-go/db"
	"naive-admin-go/model"
	"strconv"

	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/models/response"
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
func (h *UserHandler) Detail(c *gin.Context) {
	var data = &inout.UserDetailRes{}
	var uid, _ = c.Get("uid")
	db.Dao.Model(model.User{}).Where("id=?", uid).Find(&data)
	db.Dao.Model(model.Profile{}).Where("userId=?", uid).Find(&data.Profile)
	urolIdList := db.Dao.Model(model.UserRolesRole{}).Where("userId=?", uid).Select("roleId")
	db.Dao.Model(model.Role{}).Where("id IN (?)", urolIdList).Find(&data.Roles)
	if len(data.Roles) > 0 {
		data.CurrentRole = data.Roles[0]
	}
	Resp.Succ(c, data)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "无效的请求参数")
		return
	}

	if err := h.userService.Create(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUser 获取用户详情
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "用户不存在")
		return
	}

	response.Success(c, user)
}

// ListUsers 用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, err := h.userService.List(page, pageSize)
	if err != nil {
		fmt.Errorf("获取用户列表失败:", "error", err)
		response.Error(c, 500, "服务器内部错误")
		return
	}

	response.Success(c, gin.H{
		"items": users,
		"total": len(users),
	})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "无效的请求参数")
		return
	}

	user, err := h.userService.Update(uint(id), &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "无效的用户ID")
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		fmt.Errorf("删除用户失败",
			"error", err.Error(),
			"user_id", id)
		response.Error(c, 500, "删除用户失败")
		return
	}

	response.Success(c, nil)
}
