package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/tiamxu/cactus/logic/api"
	"github.com/tiamxu/cactus/logic/middleware"

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
	// r.Static("/static", "./static")
	// r.LoadHTMLGlob("static/templates/*")
	// r.GET("/links", linkHandler.RenderIndexPage)

	// ================== 开放路由（无需鉴权） ==================
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)    // 登录
		auth.GET("/captcha", authHandler.Captcha) // 验证码
	}

	// ================== 受保护路由（需 JWT 鉴权） ==================
	api := r.Group("")
	api.Use(middleware.JWTAuthMiddleware()) // 应用 JWT 中间件
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/logout", authHandler.Logout)     // 退出登录
			auth.POST("/password", authHandler.Password) // 修改密码
		}

		// 用户管理
		user := api.Group("/user")
		{
			user.GET("", userHandler.List)                        // 用户列表
			user.POST("", userHandler.Add)                        // 新增用户
			user.DELETE("/:id", userHandler.Delete)               // 删除用户
			user.PATCH("/password/reset/:id", userHandler.Update) // 重置密码
			user.PATCH("/:id", userHandler.Update)                // 更新用户信息
			user.PATCH("/profile/:id", userHandler.Profile)       // 更新个人资料
			user.GET("/detail", userHandler.Detail)               // 用户详情
		}

		// 角色管理
		role := api.Group("/role")
		{
			role.GET("", roleHandler.List)                             // 角色列表
			role.POST("", roleHandler.Add)                             // 新增角色
			role.PATCH("/:id", roleHandler.Update)                     // 更新角色
			role.DELETE("/:id", roleHandler.Delete)                    // 删除角色
			role.PATCH("/users/add/:id", roleHandler.AddUser)          // 角色添加用户
			role.PATCH("/users/remove/:id", roleHandler.RemoveUser)    // 角色移除用户
			role.GET("/page", roleHandler.ListPage)                    // 角色分页列表
			role.GET("/permissions/tree", roleHandler.PermissionsTree) // 角色权限树
		}

		// 权限管理
		permission := api.Group("/permission")
		{
			permission.POST("", permissionHandler.Add)                  // 新增权限
			permission.PATCH("/:id", permissionHandler.PatchPermission) // 更新权限
			permission.DELETE("/:id", permissionHandler.Delete)         // 删除权限
			permission.GET("/tree", permissionHandler.List)             // 权限树
			permission.GET("/menu/tree", permissionHandler.List)        // 菜单权限树
		}

		// 项目管理
		project := api.Group("/project")
		{
			project.GET("", projectHandler.List)            // 项目列表
			project.GET("/:id", projectHandler.List)        // 项目详情
			project.POST("", projectHandler.List)           // 新增项目（注：这里似乎有误，应该是 Add）
			project.PATCH("/:id", projectHandler.List)      // 更新项目（注：这里似乎有误，应该是 Update）
			project.DELETE("/:id", projectHandler.List)     // 删除项目（注：这里似乎有误，应该是 Delete）
			project.PUT("/:id/status", projectHandler.List) // 更新项目状态（注：这里似乎有误，应该是 UpdateStatus）
		}

		// 链接管理
		link := api.Group("/links")
		{
			link.GET("", linkHandler.List)          // 链接列表
			link.POST("", linkHandler.Add)          // 新增链接
			link.PUT("/:id", linkHandler.Update)    // 更新链接
			link.DELETE("/:id", linkHandler.Delete) // 删除链接
		}
	}
}
