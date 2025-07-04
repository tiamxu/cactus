package api

//导航
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/types"

	"github.com/tiamxu/cactus/logic/service"
)

type NavigationHandler struct {
	service *service.NavigationService
}

func NewNavigationHandler() *NavigationHandler {
	return &NavigationHandler{
		service: service.NewNavigationService(),
	}
}

func (h *NavigationHandler) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req types.NavigationCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))

		return
	}

	err := h.service.Add(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, ""))

}

func (h *NavigationHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "invalid id"))
		return
	}

	link, err := h.service.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, RespError(c, err, "link not found"))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, link))
}

func (h *NavigationHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var req types.NavigationListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数错误"))
		return
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	data, err := h.service.List(ctx, &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, "error"))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, data))

}

func (h *NavigationHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "invalid id"))

		return
	}

	var req types.NavigationUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, ""))

		return
	}

	if err := h.service.Update(ctx, id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, ""))
}

func (h *NavigationHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "invalid id"))

		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, ""))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, ""))

}

func (h *NavigationHandler) PublicList(c *gin.Context) {
	var req types.NavigationListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "参数错误"))
		return
	}
	if req.PageNo < 1 {
		req.PageNo = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	data, total, err := h.service.ListPublicNavigationList(req.PageNo, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, "查询导航链接信息失败"))
		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, map[string]interface{}{
		"total":    total,
		"pageData": data,
	}))
}
