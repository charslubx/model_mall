// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cart

import (
	"context"
	"encoding/json"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取购物车列表
func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCartLogic) GetCart() (resp *types.GetCartResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询购物车列表
	carts, err := l.svcCtx.Repos.CartRepo.GetByUserID(l.ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败: %v", err)
	}

	// 构造响应
	items := make([]types.CartItem, 0)
	var subtotal, total float64

	for _, cart := range carts {
		// 查询商品信息
		product, err := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, cart.ProductID)
		if err != nil {
			logx.Errorf("商品不存在: %d", cart.ProductID)
			continue
		}

		// 解析图片
		var images []string
		if product.Images != "" {
			_ = json.Unmarshal([]byte(product.Images), &images)
		}
		image := ""
		if len(images) > 0 {
			image = images[0]
		}

		item := types.CartItem{
			Id:        fmt.Sprintf("%d", cart.ID),
			ProductId: fmt.Sprintf("%d", cart.ProductID),
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  cart.Quantity,
			Color:     cart.Color,
			Size:      cart.Size,
			Image:     image,
			Stock:     product.Stock,
			Selected:  cart.Selected,
		}

		items = append(items, item)
		if cart.Selected {
			subtotal += product.Price * float64(cart.Quantity)
		}
	}

	total = subtotal

	resp = &types.GetCartResponse{
		Items:    items,
		Subtotal: subtotal,
		Total:    total,
	}

	return resp, nil
}
