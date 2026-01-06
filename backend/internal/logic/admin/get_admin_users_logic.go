// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户列表
func NewGetAdminUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminUsersLogic {
	return &GetAdminUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminUsersLogic) GetAdminUsers(req *types.GetAdminUsersRequest) (resp *types.GetAdminUsersResponse, err error) {
	// 获取管理员ID
	adminId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// TODO: 根据type、status、keyword查询用户列表
	// 这里使用模拟数据
	users := []types.AdminUserInfo{
		{
			Id:           "u123456",
			Name:         "张三",
			Email:        "zhangsan@example.com",
			Type:         "customer",
			Status:       "active",
			RegisterDate: "2023-01-15",
			LastLogin:    "2025-01-05T10:30:00Z",
			Avatar:       "https://cdn.example.com/avatars/u123456.jpg",
			OrderCount:   15,
			TotalSpent:   2345.67,
		},
		{
			Id:           "m001",
			Name:         "时尚优选",
			Email:        "merchant@example.com",
			Type:         "merchant",
			Status:       "active",
			RegisterDate: "2022-12-05",
			LastLogin:    "2025-01-06T08:20:00Z",
			Avatar:       "https://cdn.example.com/merchants/m001.jpg",
			OrderCount:   856,
			TotalSpent:   98765.43,
		},
	}

	logx.Infof("管理员 %d 查询用户列表", adminId)

	resp = &types.GetAdminUsersResponse{
		Users:    users,
		Total:    2568,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	return resp, nil
}
