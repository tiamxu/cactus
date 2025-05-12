package api

//导航
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/models"
	"github.com/tiamxu/cactus/service"
)

type NavigationHandler struct {
	service *service.NavigationService
}

func NewNavigationHandler() *NavigationHandler {
	return &NavigationHandler{
		service: service.NewNavigationService(),
	}
}

func (h *NavigationHandler) GetAllLinks(c *gin.Context) {
	links, err := h.service.GetAllLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Resp.Succ(c, links)
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

func (h *NavigationHandler) CreateLink(c *gin.Context) {
	var req models.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateLink(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *NavigationHandler) UpdateLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateLink(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NavigationHandler) DeleteLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteLink(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NavigationHandler) RenderIndexPage(c *gin.Context) {
	grouped, err := h.service.RenderIndexPage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Resp.Succ(c, grouped)

}
