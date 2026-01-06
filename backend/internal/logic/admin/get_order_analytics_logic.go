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

type GetOrderAnalyticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取订单分析数据
func NewGetOrderAnalyticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderAnalyticsLogic {
	return &GetOrderAnalyticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderAnalyticsLogic) GetOrderAnalytics(req *types.GetOrderAnalyticsRequest) (resp *types.OrderAnalytics, err error) {
	// 获取管理员ID
	adminId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// TODO: 根据timeRange查询实际订单分析数据
	// 这里使用模拟数据
	resp = &types.OrderAnalytics{
		OrdersByStatus: []types.OrderStatusStat{
			{Status: "已完成", Count: 9845, Percentage: 62},
			{Status: "已发货", Count: 2568, Percentage: 16},
			{Status: "待发货", Count: 1256, Percentage: 8},
			{Status: "待付款", Count: 1020, Percentage: 6},
			{Status: "已取消", Count: 1000, Percentage: 8},
		},
		OrdersByPayment: []types.PaymentMethodStat{
			{Method: "支付宝", Count: 7845, Percentage: 50},
			{Method: "微信支付", Count: 6258, Percentage: 40},
			{Method: "银联支付", Count: 1586, Percentage: 10},
		},
		RecentSales: []types.SalesDataPoint{
			{Date: "2025-01-01", Sales: 125689.00, Orders: 456},
			{Date: "2025-01-02", Sales: 138456.00, Orders: 512},
			{Date: "2025-01-03", Sales: 142356.00, Orders: 534},
			{Date: "2025-01-04", Sales: 135678.00, Orders: 498},
			{Date: "2025-01-05", Sales: 148923.00, Orders: 556},
		},
		TopCategories: []types.CategorySalesStat{
			{Category: "上衣", Sales: 456789.00, Percentage: 36},
			{Category: "裤装", Sales: 325678.00, Percentage: 26},
			{Category: "裙装", Sales: 198456.00, Percentage: 16},
			{Category: "外套", Sales: 275890.00, Percentage: 22},
		},
	}

	logx.Infof("管理员 %d 查询订单分析数据，时间范围: %s", adminId, req.TimeRange)

	return resp, nil
}
