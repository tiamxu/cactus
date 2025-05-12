package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/tiamxu/cactus/api"
	"github.com/tiamxu/cactus/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte("captch"))))

	authHandler := api.NewAuthHandler()
	userHandler := api.NewUserHandler()
	roleHandler := api.NewRoleHandler()
	permissionHandler := api.NewPermissionsHandler()
	projectHandler := api.NewProjectHandler()
	linkHandler := api.NewNavigationHandler()

	r.GET("/links", linkHandler.RenderIndexPage)

	r.POST("/auth/login", authHandler.Login)
	r.GET("/auth/captcha", authHandler.Captcha)

	r.Use(middleware.Jwt())
	r.POST("/auth/logout", authHandler.Logout)
	r.POST("/auth/password", authHandler.Password)

	r.GET("/user", userHandler.List)
	r.POST("/user", userHandler.Add)
	r.DELETE("/user/:id", userHandler.Delete)
	r.PATCH("/user/password/reset/:id", userHandler.Update)
	r.PATCH("/user/:id", userHandler.Update)
	r.PATCH("/user/profile/:id", userHandler.Profile)
	r.GET("/user/detail", userHandler.Detail)

	r.GET("/role", roleHandler.List)
	r.POST("/role", roleHandler.Add)
	r.PATCH("/role/:id", roleHandler.Update)
	r.DELETE("/role/:id", roleHandler.Delete)
	r.PATCH("/role/users/add/:id", roleHandler.AddUser)
	r.PATCH("/role/users/remove/:id", roleHandler.RemoveUser)
	r.GET("/role/page", roleHandler.ListPage)
	r.GET("/role/permissions/tree", roleHandler.PermissionsTree)

	r.POST("/permission", permissionHandler.Add)
	r.PATCH("/permission/:id", permissionHandler.PatchPermission)
	r.DELETE("/permission/:id", permissionHandler.Delete)
	r.GET("/permission/tree", permissionHandler.List)
	r.GET("/permission/menu/tree", permissionHandler.List)

	r.GET("/project", projectHandler.List)
	r.GET("/project/:id", projectHandler.List)
	r.POST("/project", projectHandler.List)
	r.PATCH("/project/:id", projectHandler.List)
	r.DELETE("/project/:id", projectHandler.List)
	r.PUT("/project/:id/status", projectHandler.List)

}
