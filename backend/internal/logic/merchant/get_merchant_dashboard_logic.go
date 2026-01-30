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

	// 获取仓储层实例
	productRepo := l.svcCtx.Repos.ProductRepo
	orderRepo := l.svcCtx.Repos.OrderRepo

	// 统计商品总数
	productsCount, err := productRepo.CountByMerchant(l.ctx, merchantId)
	if err != nil {
		logx.Errorf("统计商品总数失败: %v", err)
		return nil, fmt.Errorf("统计数据失败")
	}

	// 统计总销量
	salesCount, err := productRepo.GetTotalSalesByMerchant(l.ctx, merchantId)
	if err != nil {
		logx.Errorf("统计商品销量失败: %v", err)
		return nil, fmt.Errorf("统计数据失败")
	}

	// 统计总收入
	revenue, err := orderRepo.GetRevenueSumByMerchant(l.ctx, merchantId)
	if err != nil {
		logx.Errorf("统计总收入失败: %v", err)
		return nil, fmt.Errorf("统计数据失败")
	}

	// 统计订单总数
	ordersCount, err := orderRepo.CountByMerchant(l.ctx, merchantId)
	if err != nil {
		logx.Errorf("统计订单总数失败: %v", err)
		return nil, fmt.Errorf("统计数据失败")
	}

	// 统计待处理订单数（待支付）
	pendingOrders, err := orderRepo.CountByMerchantAndStatus(l.ctx, merchantId, 0)
	if err != nil {
		logx.Errorf("统计待处理订单失败: %v", err)
		return nil, fmt.Errorf("统计数据失败")
	}

	// 统计今日销售额
	todaySales, err := orderRepo.GetTodaySalesByMerchant(l.ctx, merchantId)
	if err != nil {
		logx.Errorf("统计今日销售额失败: %v", err)
		return nil, fmt.Errorf("统计数据失败")
	}

	// 统计今日订单数
	todayOrders, err := orderRepo.GetTodayOrdersByMerchant(l.ctx, merchantId)
	if err != nil {
		logx.Errorf("统计今日订单数失败: %v", err)
		return nil, fmt.Errorf("统计数据失败")
	}

	resp = &types.MerchantDashboard{
		ProductsCount: productsCount,
		SalesCount:    salesCount,
		Revenue:       revenue,
		OrdersCount:   ordersCount,
		PendingOrders: pendingOrders,
		TodaySales:    todaySales,
		TodayOrders:   todayOrders,
	}

	logx.Infof("商户 %d 查询数据统计成功", merchantId)

	return resp, nil
}
