// Package main 是课程管理系统的入口包
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin" // Web框架
	"github.com/joho/godotenv" // 环境变量管理

	"go-vue-admin/config" // 项目配置
	"go-vue-admin/routes" // 路由定义
)

// init 函数在main函数之前执行，用于初始化系统配置
func init() {
	// 加载环境变量配置文件
	// 如果.env文件不存在，记录警告但继续执行
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 未找到.env文件，将使用默认配置")
	}

	// 初始化数据库连接
	// 会自动创建数据库表结构（如果不存在）
	config.InitDB()

	// 初始化Redis连接
	// 用于缓存和会话管理
	config.InitRedis()
}

// main 函数是程序的入口点
func main() {
	// 设置Gin的运行模式
	// 可选值: debug/release/test
	gin.SetMode(getEnv("GIN_MODE", "debug"))

	// 创建默认的Gin路由引擎
	// 包含了Logger和Recovery中间件
	r := gin.Default()

	// 配置CORS（跨域资源共享）中间件
	r.Use(func(c *gin.Context) {
		// 允许的源（域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", getEnv("ALLOWED_ORIGINS", "*"))
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
	})

	// 设置路由
	// 包括API路由和静态文件服务
	routes.SetupRoutes(r)

	// 启动服务器
	// 默认端口为8080
	port := getEnv("SERVER_PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// getEnv 获取环境变量值，如果环境变量不存在则返回默认值
// key: 环境变量名
// defaultValue: 默认值
// 返回: 环境变量值或默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
