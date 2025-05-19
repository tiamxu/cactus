package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/tiamxu/cactus/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供 Token"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Token 格式错误"})
			return
		}
		token := parts[1]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Token 不能为空"})
			return
		}

		claims, err := utils.NewJWT().ParseToken(token)

		if err != nil {
			switch err {
			case utils.ErrTokenExpired:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token 已过期"})
			case utils.ErrTokenMalformed:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Malformed token"})
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "解析token失败"})
			}
			return
		}
		if claims.UID <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
			return
		}

		if time.Until(claims.ExpiresAt.Time) < 5*time.Minute {
			newToken, err := utils.NewJWT().RefreshToken(token)
			if err == nil {
				c.Header("New-Token", newToken) // 返回新 Token 给客户端
			}
		}

		c.Set("uid", claims.UID)
		c.Next()
	}
}
