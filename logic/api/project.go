package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tiamxu/cactus/service"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		projectService: &service.ProjectService{},
	}
}

func (r *ProjectHandler) List(c *gin.Context) {

}

func (r *ProjectHandler) Add(c *gin.Context) {

}

func (r *ProjectHandler) Delete(c *gin.Context) {

}
