package response

import (
	"naive-admin-go/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRes struct {
	AccessToken string `json:"accessToken"`
}

type UserDetailRes struct {
	model.User
	Profile     *model.Profile `json:"profile"`
	Roles       []*model.Role  `json:"roles" `
	CurrentRole *model.Role    `json:"currentRole"`
}

type RoleListRes []*model.Role

type UserListItem struct {
	ID         int           `json:"id"`
	Username   string        `json:"username"`
	Enable     bool          `json:"enable"`
	CreateTime time.Time     `json:"createTime"`
	UpdateTime time.Time     `json:"updateTime"`
	Gender     int           `json:"gender"`
	Avatar     string        `json:"avatar"`
	Address    string        `json:"address"`
	Email      string        `json:"email"`
	Roles      []*model.Role `json:"roles"`
}
type UserListRes struct {
	PageData []UserListItem `json:"pageData"`
	Total    int64          `json:"total"`
}
type RoleListPageItem struct {
	model.Role
	PermissionIds []int64 `json:"permissionIds" gorm:"-"`
}
type RoleListPageRes struct {
	PageData []RoleListPageItem `json:"pageData"`
	Total    int64              `json:"total"`
}

// Response 基础序列化器
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// DataList 带有总数的Data结构
type DataList struct {
	Item  interface{} `json:"item"`
	Total uint        `json:"total"`
}

// TokenData 带有token的Data结构
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: 0,
		Data:   data,
		Msg:    "success",
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Status: code,
		Data:   nil,
		Msg:    message,
		Error:  "error",
	})
}
