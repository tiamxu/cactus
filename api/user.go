package api

import (
	"net/http"
	"strconv"

	"github.com/tiamxu/cactus/inout"
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
	var params inout.PatchProfileUserReq
	if err := c.BindJSON(&params); err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err := a.userService.UpdateProfile(params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, err)

}

func (a *UserHandler) Update(c *gin.Context) {
	var params inout.PatchUserReq
	err := a.userService.Update(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, err)

}
func (a *UserHandler) Add(c *gin.Context) {
	var params inout.AddUserReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = a.userService.Add(params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")

}
func (a *UserHandler) Delete(c *gin.Context) {
	uid := c.Param("id")
	s, err := strconv.Atoi(uid)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	err = a.userService.Delete(s)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	Resp.Succ(c, "")

}
