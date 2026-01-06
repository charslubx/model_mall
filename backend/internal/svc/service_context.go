// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"model_mall_backend/backend/internal/config"
	"model_mall_backend/backend/internal/middleware"
)

type ServiceContext struct {
	Config        config.Config
	JWTMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		JWTMiddleware: middleware.NewJWTMiddleware().Handle,
	}
}
