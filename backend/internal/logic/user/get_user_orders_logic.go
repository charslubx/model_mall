// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户订单历史
func NewGetUserOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOrdersLogic {
	return &GetUserOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserOrdersLogic) GetUserOrders(req *types.GetOrdersRequest) (resp *types.GetOrdersResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询订单列表
	orders, total, err := l.svcCtx.Repos.OrderRepo.GetByUserID(l.ctx, userId, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %v", err)
	}

	// 构造响应
	orderList := make([]types.OrderListItem, 0)
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

		// 状态文本映射
		statusTextMap := map[string]string{
			"pending":   "待付款",
			"paid":      "已付款",
			"shipped":   "已发货",
			"completed": "已完成",
			"cancelled": "已取消",
		}

		orderList = append(orderList, types.OrderListItem{
			Id:         fmt.Sprintf("%d", order.ID),
			OrderNo:    order.OrderNo,
			Date:       order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Status:     order.Status,
			StatusText: statusTextMap[order.Status],
			Total:      order.Total,
			ItemCount:  len(items),
			Items:      orderItems,
		})
	}

	resp = &types.GetOrdersResponse{
		Orders:   orderList,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	return resp, nil
}
