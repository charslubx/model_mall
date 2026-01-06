// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付订单
func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayOrderLogic) PayOrder(req *types.PayOrderRequest, orderId string) (resp *types.PayOrderResponse, err error) {
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
	if order.Status != "pending" {
		return nil, fmt.Errorf("订单状态不允许支付")
	}

	// 更新支付方式
	order.PaymentMethod = req.PaymentMethod

	// 模拟支付成功，更新订单状态
	now := time.Now()
	order.Status = "paid"
	order.PaymentStatus = "paid"
	order.PaidAt = &now

	err = l.svcCtx.Repos.OrderRepo.Update(l.ctx, order)
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 生成支付URL和二维码(模拟)
	paymentUrl := fmt.Sprintf("https://payment.example.com/pay?orderNo=%s", order.OrderNo)
	qrCode := fmt.Sprintf("https://payment.example.com/qr/%s.png", order.OrderNo)

	logx.Infof("订单 %s 支付成功", order.OrderNo)

	resp = &types.PayOrderResponse{
		PaymentUrl: paymentUrl,
		OrderNo:    order.OrderNo,
		QrCode:     qrCode,
	}

	return resp, nil
}
