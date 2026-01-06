// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户详情
func NewGetAdminUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminUserDetailLogic {
	return &GetAdminUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminUserDetailLogic) GetAdminUserDetail(userId string) (resp *types.AdminUserDetail, err error) {
	// 获取管理员ID
	adminId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 转换用户ID
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的用户ID")
	}

	// 查询用户信息
	user, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// TODO: 查询用户的订单统计
	recentOrders := []types.RecentOrder{
		{
			OrderNo: "ORD20250105123456",
			Date:    "2025-01-05T10:30:00Z",
			Total:   208.00,
			Status:  "completed",
		},
	}

	resp = &types.AdminUserDetail{
		Id:           userId,
		Name:         user.Name,
		Email:        user.Email,
		Phone:        user.Phone,
		Type:         user.UserType,
		Status:       user.Status,
		RegisterDate: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastLogin:    user.LastLoginAt.Format("2006-01-02T15:04:05Z07:00"),
		Avatar:       user.Avatar,
		OrderCount:   15,
		TotalSpent:   2345.67,
		RecentOrders: recentOrders,
	}

	logx.Infof("管理员 %d 查询用户详情: %s", adminId, userId)

	return resp, nil
}
