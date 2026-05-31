package config

import (
	"encoding/json"
	"os"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	ITick    ITickConfig    `json:"itick"`
	Database DatabaseConfig `json:"database"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// ITickConfig iTick API 配置
type ITickConfig struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

// DatabaseConfig 数据库配置（预留）
type DatabaseConfig struct {
	DSN string `json:"dsn"`
}

// LoadConfig 从JSON文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
		ITick: ITickConfig{
			BaseURL: "https://api-free.itick.org",
			Token:   "",
		},
	}
}
