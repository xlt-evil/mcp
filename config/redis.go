package config

import (
	"fmt"
	"time"
)

// RedisConfig Redis配置结构
type RedisConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`

	// 连接池配置
	Pool struct {
		MaxIdle     int           `yaml:"max_idle" json:"max_idle"`
		MaxActive   int           `yaml:"max_active" json:"max_active"`
		IdleTimeout time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	} `yaml:"pool" json:"pool"`

	// 超时配置
	Timeout struct {
		Connect time.Duration `yaml:"connect" json:"connect"`
		Read    time.Duration `yaml:"read" json:"read"`
		Write   time.Duration `yaml:"write" json:"write"`
	} `yaml:"timeout" json:"timeout"`

	// 日志配置
	Logging struct {
		Enabled bool   `yaml:"enabled" json:"enabled"`
		Level   string `yaml:"level" json:"level"`
		File    string `yaml:"file" json:"file"`
	} `yaml:"logging" json:"logging"`
}

// GetAddr 获取Redis地址
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsValid 验证配置是否有效
func (c *RedisConfig) IsValid() bool {
	return c.Enabled && c.Host != "" && c.Port > 0
}

// GetPassword 获取密码（如果为空则返回空字符串）
func (c *RedisConfig) GetPassword() string {
	return c.Password
}

// GetDB 获取数据库编号
func (c *RedisConfig) GetDB() int {
	return c.DB
}
