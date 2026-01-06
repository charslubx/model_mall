// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除商品（商户）
func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductLogic) DeleteProduct(productId string) (resp *types.BaseResponse, err error) {
	// 获取当前商户ID
	merchantID, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 转换商品ID
	id, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的商品ID")
	}

	// 查询商品
	product, err := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("商品不存在")
	}

	// 验证权限
	if product.MerchantID != merchantID {
		return nil, fmt.Errorf("无权限操作此商品")
	}

	// 删除商品
	err = l.svcCtx.Repos.ProductRepo.Delete(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("删除商品失败: %v", err)
	}

	resp = &types.BaseResponse{
		Code:    200,
		Message: "删除成功",
		Data:    map[string]bool{"success": true},
	}

	return resp, nil
}
