package api

import (
	"net/http"
	"strconv"

	"github.com/tiamxu/cactus/logic/service"
	"github.com/tiamxu/cactus/types"

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
		c.JSON(http.StatusBadRequest, RespError(c, nil, "uid not found in context"))

		return
	}

	userID := uid.(int)

	userDetail, err := a.userService.GetUserDetail(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))

		return
	}
	if len(userDetail.Roles) > 0 {
		userDetail.CurrentRole = userDetail.Roles[0]
	}
	c.JSON(http.StatusOK, RespSuccess(c, userDetail))

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
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, userListRes))

}
func (a *UserHandler) Profile(c *gin.Context) {
	var params types.PatchProfileUserReq
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	err := a.userService.UpdateProfile(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))

}

func (a *UserHandler) Update(c *gin.Context) {
	var params types.PatchUserReq
	err := c.BindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	err = a.userService.Update(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))

}
func (a *UserHandler) Add(c *gin.Context) {
	var params types.AddUserReq
	err := c.Bind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	err = a.userService.Add(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))

}
func (a *UserHandler) Delete(c *gin.Context) {
	uid := c.Param("id")
	s, err := strconv.Atoi(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	err = a.userService.Delete(s)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, ""))

}
