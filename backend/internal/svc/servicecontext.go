package svc

import (
	"context"
	"log"

	"model_mall_backend/backend/internal/config"
	"model_mall_backend/backend/internal/middleware"
	"model_mall_backend/backend/internal/repository"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config             config.Config
	JWTMiddleware      rest.Middleware
	LogHelper          *LogHelper
	PGHelper           *PGHelper
	MySqlHelper        *MySqlHelper
	RedisHelper        *RedisHelper
	OrmHelper          *OrmHelper
	Repos              *repository.Repositories
	ModelServiceClient *ModelServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx := context.Background()

	// 初始化日志助手
	logHelper := NewLogHelper(ctx)

	// 初始化PostgreSQL连接
	pgHelper, err := NewPGHelper(ctx, &c, logHelper)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	// 初始化MySQL连接
	mysqlHelper, err := NewMySqlHelper(&c, logHelper)
	if err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}

	// 初始化Redis连接
	redisHelper, err := NewRedisHelper(ctx, &c, logHelper)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 初始化GORM连接
	ormHelper, err := NewOrmHelper(ctx, &c, logHelper)
	if err != nil {
		log.Fatalf("Failed to initialize GORM: %v", err)
	}

	svcCtx := &ServiceContext{
		Config:        c,
		JWTMiddleware: middleware.NewJwtMiddlewareWithRedis(c.Auth.AccessSecret, redisHelper.GetClient()).Handle,
		LogHelper:     logHelper,
		PGHelper:      pgHelper,
		MySqlHelper:   mysqlHelper,
		RedisHelper:   redisHelper,
		OrmHelper:     ormHelper,
	}

	// 初始化仓库层
	svcCtx.Repos = repository.NewRepositories(ormHelper.GetDB())

	// 初始化模型服务客户端
	svcCtx.ModelServiceClient = NewModelServiceClient(
		c.ModelService.BaseURL,
		c.ModelService.APIKey,
	)

	return svcCtx
}
