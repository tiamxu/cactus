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

func (h *NavigationHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	data, err := h.service.List(ctx, pageNo, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespError(c, err, "error"))

		return
	}
	c.JSON(http.StatusOK, RespSuccess(c, data))

}

func (h *NavigationHandler) GetLinkByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, RespError(c, err, "invalid id"))
		return
	}

	link, err := h.service.GetLinkByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, RespError(c, err, "link not found"))

		return
	}

	c.JSON(http.StatusOK, RespSuccess(c, link))
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
