package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/logic/service"
	"github.com/tiamxu/cactus/pkg/utils"
	"github.com/tiamxu/cactus/types"
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
	c.Header("Content-Type", "image/svg+xml; charset=utf-8")
	c.Data(http.StatusOK, "image/svg+xml", svg)
}

// 登陆
func (h *AuthHandler) Login(c *gin.Context) {
	var params types.LoginReq
	if err := c.Bind(&params); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "请求参数错误"))
		return
	}
	session := sessions.Default(c)
	if params.Captcha != session.Get("captch") {
		c.JSON(http.StatusBadRequest, RespError(c, nil, "验证码不正确"))
		return
	}

	resp, err := h.authService.Authenticate(params.Username, params.Password)
	if err != nil {
		c.JSON(401, RespError(c, err, "认证失败"))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, resp))

}

func (h *AuthHandler) Password(c *gin.Context) {
	var req types.AuthPwReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数错误"))
		return
	}
	uid, _ := c.Get("uid")
	if err := h.authService.ChangePassword(uid.(int), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "更新密码错误"))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, true))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, RespSuccess(c, true))

}
