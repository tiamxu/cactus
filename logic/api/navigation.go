package api

//导航
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/inout"
	"github.com/tiamxu/cactus/logic/repo"

	"github.com/tiamxu/cactus/logic/service"
)

type NavigationHandler struct {
	service *service.NavigationService
}

func NewNavigationHandler() *NavigationHandler {
	db := repo.NewNavigationDB()
	return &NavigationHandler{
		service: service.NewNavigationService(db),
	}
}

func (h *NavigationHandler) List(c *gin.Context) {
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	data, err := h.service.List(pageNo, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Resp.Succ(c, data)
}

func (h *NavigationHandler) GetLinkByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	link, err := h.service.GetLinkByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	Resp.Succ(c, link)
}

func (h *NavigationHandler) Add(c *gin.Context) {
	var req inout.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Add(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Resp.Succ(c, "")
}

func (h *NavigationHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req inout.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NavigationHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// func (h *NavigationHandler) RenderIndexPage(c *gin.Context) {
// 	grouped, err := h.service.RenderIndexPage()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	Resp.Succ(c, grouped)
// 	// c.HTML(http.StatusOK, "index.html", gin.H{
// 	// 	"groupedLinks": grouped,
// 	// })

// }
