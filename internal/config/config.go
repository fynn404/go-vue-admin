// Package config 提供配置文件的加载和管理功能
package config

import (
	"fmt"
	"time"

	"gopkg.in/ini.v1"
)

// Config 全局配置结构体，包含所有配置项
type Config struct {
	Server   ServerConfig   // 服务器配置
	Database DatabaseConfig // 数据库配置
	Redis    RedisConfig    // Redis配置
	JWT      JWTConfig      // JWT认证配置
	CORS     CORSConfig     // 跨域配置
	Log      LogConfig      // 日志配置
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Port string // 服务器监听端口
	Mode string // 运行模式（debug/release）
}

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host            string        // 数据库服务器地址
	Port            string        // 数据库端口
	Name            string        // 数据库名称
	User            string        // 数据库用户名
	Password        string        // 数据库密码
	MaxIdleConns    int           // 最大空闲连接数
	MaxOpenConns    int           // 最大打开连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
}

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host     string // Redis服务器地址
	Port     string // Redis端口
	Password string // Redis密码
	DB       int    // Redis数据库索引
	PoolSize int    // 连接池大小
}

// JWTConfig JWT认证配置结构体
type JWTConfig struct {
	Secret      string // JWT签名密钥
	ExpireHours int    // Token过期时间（小时）
}

// CORSConfig 跨域配置结构体
type CORSConfig struct {
	AllowedOrigins string // 允许的跨域来源
}

// LogConfig 日志配置结构体
type LogConfig struct {
	Level      string // 日志级别
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小（MB）
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 旧日志文件保留天数
	Compress   bool   // 是否压缩旧日志文件
}

// GlobalConfig 全局配置实例，其他包可以直接使用
var GlobalConfig Config

// Init 初始化配置，从指定的INI文件中加载配置
// file: 配置文件路径
// 返回错误信息，如果加载失败则返回相应的错误
func Init(file string) error {
	cfg, err := ini.Load(file)
	if err != nil {
		return fmt.Errorf("failed to load config file: %v", err)
	}

	// 加载服务器配置
	GlobalConfig.Server.Port = cfg.Section("server").Key("port").String()
	GlobalConfig.Server.Mode = cfg.Section("server").Key("mode").String()

	// 加载数据库配置
	GlobalConfig.Database.Host = cfg.Section("database").Key("host").String()
	GlobalConfig.Database.Port = cfg.Section("database").Key("port").String()
	GlobalConfig.Database.Name = cfg.Section("database").Key("name").String()
	GlobalConfig.Database.User = cfg.Section("database").Key("user").String()
	GlobalConfig.Database.Password = cfg.Section("database").Key("password").String()
	GlobalConfig.Database.MaxIdleConns = cfg.Section("database").Key("max_idle_conns").MustInt(10)
	GlobalConfig.Database.MaxOpenConns = cfg.Section("database").Key("max_open_conns").MustInt(100)
	GlobalConfig.Database.ConnMaxLifetime = time.Duration(cfg.Section("database").Key("conn_max_lifetime").MustInt(3600)) * time.Second

	// 加载Redis配置
	GlobalConfig.Redis.Host = cfg.Section("redis").Key("host").String()
	GlobalConfig.Redis.Port = cfg.Section("redis").Key("port").String()
	GlobalConfig.Redis.Password = cfg.Section("redis").Key("password").String()
	GlobalConfig.Redis.DB = cfg.Section("redis").Key("db").MustInt(0)
	GlobalConfig.Redis.PoolSize = cfg.Section("redis").Key("pool_size").MustInt(10)

	// 加载JWT配置
	GlobalConfig.JWT.Secret = cfg.Section("jwt").Key("secret").String()
	GlobalConfig.JWT.ExpireHours = cfg.Section("jwt").Key("expire_hours").MustInt(24)

	// 加载跨域配置
	GlobalConfig.CORS.AllowedOrigins = cfg.Section("cors").Key("allowed_origins").String()

	// 加载日志配置
	GlobalConfig.Log.Level = cfg.Section("log").Key("level").String()
	GlobalConfig.Log.FilePath = cfg.Section("log").Key("file_path").String()
	GlobalConfig.Log.MaxSize = cfg.Section("log").Key("max_size").MustInt(100)
	GlobalConfig.Log.MaxBackups = cfg.Section("log").Key("max_backups").MustInt(10)
	GlobalConfig.Log.MaxAge = cfg.Section("log").Key("max_age").MustInt(30)
	GlobalConfig.Log.Compress = cfg.Section("log").Key("compress").MustBool(true)

	return nil
}

// GetDSN 获取MySQL数据库连接字符串
// 返回格式化后的DSN字符串，用于数据库连接
func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		GlobalConfig.Database.User,
		GlobalConfig.Database.Password,
		GlobalConfig.Database.Host,
		GlobalConfig.Database.Port,
		GlobalConfig.Database.Name)
}

// GetRedisAddr 获取Redis连接地址
// 返回格式化后的Redis地址字符串，用于Redis连接
func GetRedisAddr() string {
	return fmt.Sprintf("%s:%s",
		GlobalConfig.Redis.Host,
		GlobalConfig.Redis.Port)
}
