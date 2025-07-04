package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/logic/service"
	"github.com/tiamxu/cactus/types"
)

type EnvironmentHandler struct {
	service *service.EnvironmentService
}

func NewEnvironmentHandler() *EnvironmentHandler {
	return &EnvironmentHandler{
		service: service.NewEnvironmentService(),
	}
}

func (h *EnvironmentHandler) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req types.EnvironmentCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数错误"))
		return
	}
	env, err := h.service.Add(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, env))

}

func (h *EnvironmentHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "无效的环境ID"))
		return
	}

	env, err := h.service.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, RespError(c, err, ""))
		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, env))
}

func (h *EnvironmentHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var req types.EnvironmentListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数错误"))
		return
	}
	// 设置默认值
	if req.PageNo <= 0 {
		req.PageNo = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	resp, err := h.service.List(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, resp))

}

func (h *EnvironmentHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "环境ID必须是有效的数字"))

		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, RespError(c, err, "环境ID必须大于0"))

		return
	}
	var req types.EnvironmentUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数错误"))

		return
	}

	if err := h.service.Update(ctx, id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, nil))

}

func (h *EnvironmentHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "无效的环境ID"))
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, nil))

}
