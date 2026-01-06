// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package cart

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清空购物车
func NewClearCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCartLogic {
	return &ClearCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearCartLogic) ClearCart() (resp *types.BaseResponse, err error) {
	// 获取用户ID
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 清空用户购物车
	err = l.svcCtx.Repos.CartRepo.DeleteByUserID(l.ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("清空购物车失败: %v", err)
	}

	resp = &types.BaseResponse{
		Code:    200,
		Message: "清空成功",
		Data:    map[string]bool{"success": true},
	}

	return resp, nil
}
