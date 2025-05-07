package middleware

import (
	"github.com/fynn404/go-vue-admin/internal/utils"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 创建CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的源（域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", utils.GetEnv("ALLOWED_ORIGINS", "*"))
		// 允许携带凭证（cookies等）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求（OPTIONS）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
