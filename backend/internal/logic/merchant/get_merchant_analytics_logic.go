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

type GetMerchantAnalyticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商户数据分析
func NewGetMerchantAnalyticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantAnalyticsLogic {
	return &GetMerchantAnalyticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantAnalyticsLogic) GetMerchantAnalytics(req *types.GetMerchantAnalyticsRequest) (resp *types.MerchantAnalytics, err error) {
	// 获取商户ID
	merchantId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// TODO: 根据timeRange查询实际数据
	// 这里使用模拟数据
	resp = &types.MerchantAnalytics{
		SalesData: []types.SalesDataPoint{
			{Date: "2025-01-01", Sales: 1200.00, Orders: 12},
			{Date: "2025-01-02", Sales: 1800.00, Orders: 18},
			{Date: "2025-01-03", Sales: 1500.00, Orders: 15},
			{Date: "2025-01-04", Sales: 2100.00, Orders: 21},
			{Date: "2025-01-05", Sales: 1900.00, Orders: 19},
		},
		CategorySales: []types.CategorySales{
			{Category: "上衣", Sales: 42, Percentage: 35, Revenue: 4158.00},
			{Category: "裤装", Sales: 28, Percentage: 23, Revenue: 4452.00},
			{Category: "裙装", Sales: 18, Percentage: 15, Revenue: 2790.00},
			{Category: "外套", Sales: 32, Percentage: 27, Revenue: 6400.00},
		},
		TopProducts: []types.TopProduct{
			{
				Id:      "p001",
				Name:    "简约纯棉T恤",
				Sales:   120,
				Revenue: 11880.00,
				Image:   "https://cdn.example.com/products/p001.jpg",
			},
			{
				Id:      "p002",
				Name:    "休闲牛仔裤",
				Sales:   95,
				Revenue: 14250.00,
				Image:   "https://cdn.example.com/products/p002.jpg",
			},
		},
	}

	logx.Infof("商户 %d 查询数据分析，时间范围: %s", merchantId, req.TimeRange)

	return resp, nil
}
