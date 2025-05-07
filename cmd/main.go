// Package main 是课程管理系统的入口包
package main

import (
	"github.com/gin-gonic/gin" // Web框架
	"github.com/joho/godotenv" // 环境变量管理
	"log"

	v1 "github.com/fynn404/go-vue-admin/api/v1"
	"github.com/fynn404/go-vue-admin/internal/config"
	"github.com/fynn404/go-vue-admin/internal/handler"
	"github.com/fynn404/go-vue-admin/internal/middleware"
	"github.com/fynn404/go-vue-admin/internal/utils"
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
	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 创建处理器
	h := handler.NewHandler()

	// 设置路由
	v1.SetupRoutes(r, h)

	// 启动服务器
	// 默认端口为8080
	port := utils.GetEnv("SERVER_PORT", "10000")
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
