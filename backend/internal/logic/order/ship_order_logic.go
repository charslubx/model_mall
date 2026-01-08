// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 商户发货
func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipOrderLogic {
	return &ShipOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShipOrderLogic) ShipOrder(req *types.ShipOrderRequest, orderId string) (resp *types.ShipOrderResponse, err error) {
	// 获取商户ID
	merchantId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 转换订单ID
	id, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID")
	}

	// 查询订单
	order, err := l.svcCtx.Repos.OrderRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// TODO: 验证商户权限 - 检查订单中的商品是否属于该商户
	// 这里简化处理，实际应该检查订单商品的merchant_id

	// 检查订单状态
	if order.Status != models.OrderStatusPaid {
		return nil, fmt.Errorf("订单状态不允许发货")
	}

	// 更新物流信息
	now := time.Now()
	order.Status = models.OrderStatusShipped
	order.TrackingNumber = req.TrackingNumber
	// order.ShippingCompany = req.ShippingCompany // Order模型中没有ShippingCompany字段
	order.ShippedAt = &now

	err = l.svcCtx.Repos.OrderRepo.Update(l.ctx, order)
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %v", err)
	}

	logx.Infof("商户 %d 为订单 %s 发货，物流单号: %s", merchantId, order.OrderNo, req.TrackingNumber)

	resp = &types.ShipOrderResponse{
		Success:   true,
		Status:    "shipped",
		ShippedAt: now.Format("2006-01-02T15:04:05Z07:00"),
	}

	return resp, nil
}
