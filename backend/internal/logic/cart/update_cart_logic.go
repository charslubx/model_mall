// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cart

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新购物车商品数量
func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCartLogic) UpdateCart(req *types.UpdateCartRequest, itemId string) (resp *types.UpdateCartResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 转换购物车项ID
	id, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的购物车项ID")
	}

	// 查询购物车项
	cart, err := l.svcCtx.Repos.CartRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("购物车项不存在")
	}

	// 验证权限
	if cart.UserID != userId {
		return nil, fmt.Errorf("无权限操作此购物车项")
	}

	// 检查商品库存
	product, err := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, cart.ProductID)
	if err != nil {
		return nil, fmt.Errorf("商品不存在")
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("库存不足")
	}

	// 更新数量
	err = l.svcCtx.Repos.CartRepo.UpdateQuantity(l.ctx, id, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("更新购物车失败: %v", err)
	}

	resp = &types.UpdateCartResponse{
		Id:       itemId,
		Quantity: req.Quantity,
	}

	return resp, nil
}
