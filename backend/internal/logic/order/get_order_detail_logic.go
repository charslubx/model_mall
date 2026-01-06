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

type GetOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取订单详情
func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderDetailLogic) GetOrderDetail(orderId string) (resp *types.OrderDetail, err error) {
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
		return nil, fmt.Errorf("无权限查看此订单")
	}

	// 获取订单项
	items, err := l.svcCtx.Repos.OrderRepo.GetOrderItems(l.ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("获取订单项失败: %v", err)
	}

	// 构造订单项
	orderItems := make([]types.OrderItemDetail, 0)
	for _, item := range items {
		orderItems = append(orderItems, types.OrderItemDetail{
			Id:        fmt.Sprintf("%d", item.ID),
			ProductId: fmt.Sprintf("%d", item.ProductID),
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Color:     item.Color,
			Size:      item.Size,
			Image:     item.Image,
		})
	}

	// 状态文本映射
	statusTextMap := map[string]string{
		"pending":   "待付款",
		"paid":      "已付款",
		"shipped":   "已发货",
		"completed": "已完成",
		"cancelled": "已取消",
	}

	// 支付方式文本
	paymentMethodMap := map[string]string{
		"alipay": "支付宝",
		"wechat": "微信支付",
		"union":  "银联支付",
	}

	// 构造时间线
	timeline := []types.OrderTimeline{
		{
			Date:        order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Status:      "订单已创建",
			Description: "您的订单已提交",
		},
	}

	if order.PaidAt != nil {
		timeline = append(timeline, types.OrderTimeline{
			Date:        order.PaidAt.Format("2006-01-02T15:04:05Z07:00"),
			Status:      "支付成功",
			Description: "订单已支付",
		})
	}

	if order.ShippedAt != nil {
		timeline = append(timeline, types.OrderTimeline{
			Date:        order.ShippedAt.Format("2006-01-02T15:04:05Z07:00"),
			Status:      "已发货",
			Description: fmt.Sprintf("快递已揽收，运单号：%s", order.TrackingNumber),
		})
	}

	if order.CompletedAt != nil {
		timeline = append(timeline, types.OrderTimeline{
			Date:        order.CompletedAt.Format("2006-01-02T15:04:05Z07:00"),
			Status:      "已完成",
			Description: "订单已完成",
		})
	}

	resp = &types.OrderDetail{
		Id:         orderId,
		OrderNo:    order.OrderNo,
		Date:       order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Status:     order.Status,
		StatusText: statusTextMap[order.Status],
		Items:      orderItems,
		Shipping: types.ShippingInfo{
			Method:          "快递配送",
			Address:         order.ShippingAddress,
			Recipient:       order.ShippingName,
			Phone:           order.ShippingPhone,
			TrackingNumber:  order.TrackingNumber,
			ShippingCompany: order.ShippingCompany,
		},
		Payment: types.PaymentInfo{
			Method:     order.PaymentMethod,
			MethodText: paymentMethodMap[order.PaymentMethod],
			Subtotal:   order.Subtotal,
			Shipping:   order.ShippingFee,
			Tax:        order.Tax,
			Total:      order.Total,
			PaidAt: func() string {
				if order.PaidAt != nil {
					return order.PaidAt.Format("2006-01-02T15:04:05Z07:00")
				}
				return ""
			}(),
		},
		Timeline: timeline,
	}

	return resp, nil
}
