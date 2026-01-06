// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建订单
func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) (resp *types.CreateOrderResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 生成订单号
	orderNo := fmt.Sprintf("ORD%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)

	// 计算订单金额
	var subtotal float64
	orderItems := make([]*models.OrderItem, 0)

	for _, item := range req.Items {
		productId, _ := strconv.ParseInt(item.ProductId, 10, 64)
		product, err := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, productId)
		if err != nil {
			return nil, fmt.Errorf("商品不存在: %s", item.ProductId)
		}

		// 检查商品状态
		if product.Status == 0 {
			return nil, fmt.Errorf("商品 %s 已下架", product.Name)
		}

		// 检查库存
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("商品 %s 库存不足", product.Name)
		}

		// 解析图片
		var images []string
		if product.Images != "" {
			json.Unmarshal([]byte(product.Images), &images)
		}
		image := ""
		if len(images) > 0 {
			image = images[0]
		}

		orderItem := &models.OrderItem{
			ProductID: productId,
			Name:      product.Name,
			Image:     image,
			Price:     product.Price,
			Quantity:  item.Quantity,
			Color:     item.Color,
			Size:      item.Size,
		}
		orderItems = append(orderItems, orderItem)
		subtotal += product.Price * float64(item.Quantity)
	}

	// 计算运费和总金额
	shippingFee := 10.0
	if subtotal >= 99 {
		shippingFee = 0 // 满99免运费
	}
	total := subtotal + shippingFee

	// 创建订单
	order := &models.Order{
		OrderNo:         orderNo,
		UserID:          userId,
		Status:          "pending",
		PaymentMethod:   req.PaymentMethod,
		PaymentStatus:   "unpaid",
		ShippingName:    req.Address.Name,
		ShippingPhone:   req.Address.Phone,
		ShippingAddress: fmt.Sprintf("%s%s%s%s", req.Address.Province, req.Address.City, req.Address.District, req.Address.Detail),
		Subtotal:        subtotal,
		ShippingFee:     shippingFee,
		Tax:             0,
		Total:           total,
		Note:            req.Note,
	}

	// 使用事务创建订单
	err = l.svcCtx.Repos.OrderRepo.CreateWithItems(l.ctx, order, orderItems)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %v", err)
	}

	// 扣减库存
	for _, item := range orderItems {
		_ = l.svcCtx.Repos.ProductRepo.DecrementStock(l.ctx, item.ProductID, item.Quantity)
	}

	// 生成支付URL(模拟)
	paymentUrl := fmt.Sprintf("https://payment.example.com/pay?orderNo=%s", orderNo)

	logx.Infof("用户 %d 创建订单成功: %s, 总金额: %.2f", userId, orderNo, total)

	resp = &types.CreateOrderResponse{
		OrderId:    fmt.Sprintf("%d", order.ID),
		OrderNo:    orderNo,
		Total:      total,
		PaymentUrl: paymentUrl,
	}

	return resp, nil
}
