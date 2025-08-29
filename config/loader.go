package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 主配置结构
type Config struct {
	Database DatabaseConfig `yaml:"database" json:"database"`
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 如果配置文件不存在，使用默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return LoadDefaultConfig(), nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// LoadDefaultConfig 加载默认配置
func LoadDefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Enabled:  true,
			Driver:   "mysql",
			Host:     "mysql.shop-cluster.svc.cluster.local",
			Port:     3306,
			User:     "root",
			Password: "7H4wpXcP6VUF9Z%9",
			Name:     "shop_kingdee",
		},
	}
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, configPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetDatabaseConfig 获取数据库配置
func (c *Config) GetDatabaseConfig() *DatabaseConfig {
	return &c.Database
}
