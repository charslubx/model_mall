// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商户信息
func NewGetMerchantProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantProfileLogic {
	return &GetMerchantProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantProfileLogic) GetMerchantProfile() (resp *types.MerchantProfile, err error) {
	// 获取商户ID
	merchantId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 查询商户信息
	merchant, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, merchantId)
	if err != nil {
		return nil, fmt.Errorf("商户不存在")
	}

	// TODO: 查询商户的商品数、销量等统计数据
	resp = &types.MerchantProfile{
		Id:            fmt.Sprintf("%d", merchant.ID),
		Name:          merchant.MerchantName,
		Avatar:        merchant.Avatar,
		Description:   merchant.Description,
		ProductsCount: 128,
		SalesCount:    1024,
		Rating:        4.8,
		ReviewsCount:  456,
		CreatedAt:     merchant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return resp, nil
}
