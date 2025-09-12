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
		Endpoint string
		Timeout  int
	}
	Upload struct {
		Path         string
		MaxSize      int64
		AllowedTypes []string
	}
}
