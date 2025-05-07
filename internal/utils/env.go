package utils

import "os"

// GetEnv 获取环境变量值，如果环境变量不存在则返回默认值
// key: 环境变量名
// defaultValue: 默认值
// 返回: 环境变量值或默认值
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
