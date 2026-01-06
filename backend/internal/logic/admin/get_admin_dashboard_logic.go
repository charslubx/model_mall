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

type GetAdminDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取系统数据统计
func NewGetAdminDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminDashboardLogic {
	return &GetAdminDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminDashboardLogic) GetAdminDashboard() (resp *types.AdminDashboard, err error) {
	// 获取管理员ID
	adminId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// TODO: 实现系统统计数据查询
	// 这里使用模拟数据
	resp = &types.AdminDashboard{
		TotalUsers:     2568,
		TotalMerchants: 156,
		TotalCustomers: 2412,
		TotalOrders:    15689,
		TotalRevenue:   1256789.45,
		TotalProducts:  8754,
		TodayOrders:    234,
		TodayRevenue:   45678.90,
		ActiveUsers:    1856,
	}

	logx.Infof("管理员 %d 查询系统统计数据", adminId)

	return resp, nil
}
