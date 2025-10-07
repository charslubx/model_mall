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
	Model struct {
		Name     string
		Version  string
		Type     string // local 或 remote
		Path     string // 本地模型路径或远程服务端点
	}
}
