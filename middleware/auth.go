package middleware

import (
	"github.com/tiamxu/cactus/api"
	"github.com/tiamxu/cactus/utils"

	"github.com/gin-gonic/gin"
)

// func JWTAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := c.GetHeader("Authorization")
// 		if token == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
// 			return
// 		}
// 		claims, err := utils.ParseToken(token)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
// 			return
// 		}

// 		c.Set("userID", claims.UserID)
// 		c.Next()
// 	}
// }

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			api.Resp.Err(c, 10002, "请求未携带token，无权限访问")
			c.Abort()
			return
		}
		j := utils.NewJWT()
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.TokenExpired {
				api.Resp.Err(c, 10002, "授权已过期")
				c.Abort()
				return
			}
			api.Resp.Err(c, 10002, err.Error())
			c.Abort()
			return
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("uid", claims.UID)
		c.Next()
	}
}
