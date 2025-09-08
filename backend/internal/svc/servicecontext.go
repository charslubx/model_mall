package svc

import (
	"context"
	"log"

	"model_mall_backend/backend/internal/config"
	"model_mall_backend/backend/internal/repository"
)

type ServiceContext struct {
	Config      config.Config
	LogHelper   *LogHelper
	PGHelper    *PGHelper
	MySqlHelper *MySqlHelper
	RedisHelper *RedisHelper
	OrmHelper   *OrmHelper
	
	// 仓库层
	UserRepo           *repository.UserRepository
	RoleRepo           *repository.RoleRepository
	PermissionRepo     *repository.PermissionRepository
	RolePermissionRepo *repository.RolePermissionRepository
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
		Config:      c,
		LogHelper:   logHelper,
		PGHelper:    pgHelper,
		MySqlHelper: mysqlHelper,
		RedisHelper: redisHelper,
		OrmHelper:   ormHelper,
	}

	// 初始化仓库层
	svcCtx.UserRepo = repository.NewUserRepository(svcCtx)
	svcCtx.RoleRepo = repository.NewRoleRepository(svcCtx)
	svcCtx.PermissionRepo = repository.NewPermissionRepository(svcCtx)
	svcCtx.RolePermissionRepo = repository.NewRolePermissionRepository(svcCtx)

	return svcCtx
}
