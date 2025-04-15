package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/models"
	"github.com/tiamxu/cactus/service"
	"github.com/tiamxu/cactus/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// 验证码
func (h *AuthHandler) Captcha(c *gin.Context) {
	svg, code := utils.GenerateCaptcha(80, 40)
	session := sessions.Default(c)
	session.Set("captch", code)
	session.Save()
	// 设置 Content-Type 为 "image/svg+xml"
	c.Header("Content-Type", "image/svg+xml; charset=utf-8")
	// 返回验证码
	c.Data(http.StatusOK, "image/svg+xml", svg)
}

// 登陆
func (h *AuthHandler) Login(c *gin.Context) {
	var params inout.LoginReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	session := sessions.Default(c)
	if params.Captcha != session.Get("captch") {
		Resp.Err(c, 20001, "验证码不正确")
		return
	}

	user, token, err := h.authService.Authenticate(params.Username, params.Password)
	fmt.Println("user:", user)
	if err != nil {
		fmt.Printf("auth failed : %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}
	// Resp.Succ(c, inout.LoginRes{
	// 	AccessToken: utils.GenerateToken(user.ID),
	// })
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  userResponse(user),
	})
}

func (h *AuthHandler) password(c *gin.Context) {
	var params inout.AuthPwReq
	err := c.Bind(&params)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}

	Resp.Succ(c, true)
}
func (h *AuthHandler) Logout(c *gin.Context) {
	Resp.Succ(c, true)
}
func userResponse(u *models.User) gin.H {
	return gin.H{
		"id":       u.ID,
		"username": u.Username,
	}
}
