package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/inout"
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

	user, _, err := h.authService.Authenticate(params.Username, params.Password)
	if err != nil {
		Resp.Err(c, 20001, "账号或密码不正确")
		return
	}
	Resp.Succ(c, inout.LoginRes{
		AccessToken: utils.GenerateToken(user.ID),
	})

}

func (h *AuthHandler) Password(c *gin.Context) {
	var req inout.AuthPwReq
	err := c.Bind(&req)
	if err != nil {
		Resp.Err(c, 20001, err.Error())
		return
	}
	uid, _ := c.Get("uid")
	if err := h.authService.ChangePassword(uid.(int), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Resp.Succ(c, true)
}
func (h *AuthHandler) Logout(c *gin.Context) {
	Resp.Succ(c, true)
}
