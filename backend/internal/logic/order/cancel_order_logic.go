// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消订单
func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderRequest, orderId string) (resp *types.CancelOrderResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
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

	// 验证权限
	if order.UserID != userId {
		return nil, fmt.Errorf("无权限操作此订单")
	}

	// 检查订单状态
	if order.Status == models.OrderStatusShipped || order.Status == models.OrderStatusCompleted {
		return nil, fmt.Errorf("订单已发货或已完成，无法取消")
	}

	if order.Status == models.OrderStatusCancelled {
		return nil, fmt.Errorf("订单已取消")
	}

	// 获取订单项，恢复库存
	items, _ := l.svcCtx.Repos.OrderRepo.GetOrderItems(l.ctx, order.ID)
	for _, item := range items {
		_ = l.svcCtx.Repos.ProductRepo.IncrementStock(l.ctx, item.ProductID, item.Quantity)
	}

	// 更新订单状态
	err = l.svcCtx.Repos.OrderRepo.UpdateStatus(l.ctx, id, models.OrderStatusCancelled)
	if err != nil {
		return nil, fmt.Errorf("取消订单失败: %v", err)
	}

	// 如果已支付，需要退款
	refundStatus := "none"
	refundAmount := 0.0
	if order.PaymentStatus == 1 { // 1表示已支付
		refundStatus = "processing"
		refundAmount = order.Total
		logx.Infof("订单 %s 退款处理中，金额: %.2f", order.OrderNo, refundAmount)
	}

	logx.Infof("用户 %d 取消订单: %s, 原因: %s", userId, order.OrderNo, req.Reason)

	resp = &types.CancelOrderResponse{
		Success:      true,
		RefundStatus: refundStatus,
		RefundAmount: refundAmount,
	}

	return resp, nil
}
