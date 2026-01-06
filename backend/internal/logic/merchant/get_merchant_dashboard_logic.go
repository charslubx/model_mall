// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商户数据统计
func NewGetMerchantDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantDashboardLogic {
	return &GetMerchantDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantDashboardLogic) GetMerchantDashboard() (resp *types.MerchantDashboard, err error) {
	// 获取商户ID
	merchantId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// TODO: 实现商户统计数据查询
	// 这里使用模拟数据
	resp = &types.MerchantDashboard{
		ProductsCount: 128,
		SalesCount:    1024,
		Revenue:       98765.43,
		OrdersCount:   856,
		PendingOrders: 23,
		TodaySales:    5432.10,
		TodayOrders:   45,
	}

	logx.Infof("商户 %d 查询数据统计", merchantId)

	return resp, nil
}
