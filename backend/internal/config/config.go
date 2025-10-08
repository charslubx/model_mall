package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	PostgreSQL struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
		SSLMode  string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	MySQL struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	}
	ModelService struct {
		BaseURL string
		APIKey  string
		Timeout int // 超时时间（秒）
	}
	Upload struct {
		MaxSize      int64  // 最大文件大小（字节）
		AllowedTypes string // 允许的文件类型，逗号分隔
		StoragePath  string // 存储路径
		BaseURL      string // 访问基础URL
	}
}
