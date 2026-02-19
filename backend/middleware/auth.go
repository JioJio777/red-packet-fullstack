package middleware

import (
	"net/http"
	"strings"

	"red-packet/pkg/response"
	"red-packet/service"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := service.ParseToken(tokenStr)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, 401, "invalid token")
			c.Abort()
			return
		}

		// 把 user_id 存入 context，后续 handler 直接取
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
