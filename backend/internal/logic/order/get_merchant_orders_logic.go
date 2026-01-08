// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商户订单列表
func NewGetMerchantOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantOrdersLogic {
	return &GetMerchantOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantOrdersLogic) GetMerchantOrders(req *types.GetOrdersRequest) (resp *types.GetMerchantOrdersResponse, err error) {
	// 获取商户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询商户订单列表
	orders, total, err := l.svcCtx.Repos.OrderRepo.GetBySellerID(l.ctx, userId, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取商户订单列表失败: %v", err)
	}

	// 构造响应
	orderList := make([]types.MerchantOrderItem, 0)
	for _, order := range orders {
		// 获取订单项
		items, _ := l.svcCtx.Repos.OrderRepo.GetOrderItems(l.ctx, order.ID)

		orderItems := make([]types.OrderItemDetail, 0)
		for _, item := range items {
			orderItems = append(orderItems, types.OrderItemDetail{
				ProductId: fmt.Sprintf("%d", item.ProductID),
				Name:      item.Name,
				Image:     item.Image,
				Quantity:  item.Quantity,
				Price:     item.Price,
			})
		}

		// 获取客户信息
		customer, _ := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, order.UserID)
		customerInfo := types.CustomerInfo{
			Name:   "未知用户",
			Avatar: "",
		}
		if customer != nil {
			customerInfo.Name = customer.Nickname
			customerInfo.Avatar = customer.Avatar
		}

		statusStr := OrderStatusToString(order.Status)
		orderList = append(orderList, types.MerchantOrderItem{
			Id:         fmt.Sprintf("%d", order.ID),
			OrderNo:    order.OrderNo,
			Customer:   customerInfo,
			Date:       order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Status:     statusStr,
			StatusText: OrderStatusTextMap[statusStr],
			Total:      order.Total,
			Items:      orderItems,
			Payment:    PaymentMethodTextMap[order.PaymentMethod],
			Address:    order.ShippingAddress,
		})
	}

	resp = &types.GetMerchantOrdersResponse{
		Orders:   orderList,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	return resp, nil
}
