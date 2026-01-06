// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cart

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddToCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加商品到购物车
func NewAddToCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddToCartLogic {
	return &AddToCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddToCartLogic) AddToCart(req *types.AddToCartRequest) (resp *types.AddToCartResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 转换商品ID
	productId, err := strconv.ParseInt(req.ProductId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的商品ID")
	}

	// 检查商品是否存在
	product, err := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("商品不存在")
	}

	// 检查商品状态
	if product.Status == 0 {
		return nil, fmt.Errorf("商品已下架")
	}

	// 检查库存
	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("库存不足")
	}

	// 检查是否已在购物车
	existing, _ := l.svcCtx.Repos.CartRepo.FindByUserAndProduct(l.ctx, userId, productId, req.Color, req.Size)

	var cartId int64
	if existing != nil {
		// 更新数量
		existing.Quantity += req.Quantity
		// 检查总数量是否超过库存
		if existing.Quantity > product.Stock {
			return nil, fmt.Errorf("库存不足")
		}
		err = l.svcCtx.Repos.CartRepo.Update(l.ctx, existing)
		cartId = existing.ID
	} else {
		// 新增购物车项
		cart := &models.Cart{
			UserID:    userId,
			ProductID: productId,
			Quantity:  req.Quantity,
			Color:     req.Color,
			Size:      req.Size,
			Selected:  true,
		}
		err = l.svcCtx.Repos.CartRepo.Create(l.ctx, cart)
		cartId = cart.ID
	}

	if err != nil {
		return nil, fmt.Errorf("添加购物车失败: %v", err)
	}

	resp = &types.AddToCartResponse{
		Id:        fmt.Sprintf("%d", cartId),
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	}

	return resp, nil
}
