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

type DeleteCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除购物车商品
func NewDeleteCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCartItemLogic {
	return &DeleteCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCartItemLogic) DeleteCartItem(itemId string) (resp *types.BaseResponse, err error) {
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

	// 删除购物车项
	err = l.svcCtx.Repos.CartRepo.Delete(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("删除购物车项失败: %v", err)
	}

	resp = &types.BaseResponse{
		Code:    200,
		Message: "删除成功",
		Data:    map[string]bool{"success": true},
	}

	return resp, nil
}
