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
	permissionHandler := api.NewPermissionsHandler()

	r.POST("/auth/login", authHandler.Login)
	r.GET("/auth/captcha", authHandler.Captcha)

	r.Use(middleware.Jwt())
	r.POST("/auth/logout", authHandler.Logout)
	r.POST("/auth/password", authHandler.Logout)

	// r.GET("/user", userHandler.ListUsers)
	// r.POST("/user", userHandler.Add)
	// r.DELETE("/user/:id", userHandler.Delete)
	// r.PATCH("/user/password/reset/:id", userHandlerUpdate)
	// r.PATCH("/user/:id", userHandler.Update)
	// r.PATCH("/user/profile/:id", api.User.Profile)
	r.GET("/user/detail", userHandler.Detail)

	// r.GET("/role", api.Role.List)
	// r.POST("/role", api.Role.Add)
	// r.PATCH("/role/:id", api.Role.Update)
	// r.DELETE("/role/:id", api.Role.Delete)
	// r.PATCH("/role/users/add/:id", api.Role.AddUser)
	// r.PATCH("/role/users/remove/:id", api.Role.RemoveUser)
	// r.GET("/role/page", api.Role.ListPage)
	r.GET("/role/permissions/tree", permissionHandler.PermissionsTree)

	// r.POST("/permission", api.Permissions.Add)
	// r.PATCH("/permission/:id", api.Permissions.PatchPermission)
	// r.DELETE("/permission/:id", api.Permissions.Delete)
	// r.GET("/permission/tree", api.Permissions.List)
	// r.GET("/permission/menu/tree", api.Permissions.List)

}
