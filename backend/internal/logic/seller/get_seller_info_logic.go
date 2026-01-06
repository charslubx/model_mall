// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package seller

import (
	"context"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSellerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取卖家信息
func NewGetSellerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSellerInfoLogic {
	return &GetSellerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSellerInfoLogic) GetSellerInfo(sellerId string) (resp *types.SellerInfo, err error) {
	// 转换卖家ID
	id, err := strconv.ParseInt(sellerId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的卖家ID")
	}

	// 查询卖家信息
	seller, err := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("卖家不存在")
	}

	// 验证是否为商户
	if seller.UserType != "merchant" {
		return nil, fmt.Errorf("该用户不是商户")
	}

	// TODO: 查询卖家的统计数据
	resp = &types.SellerInfo{
		Id:            sellerId,
		Name:          seller.MerchantName,
		Avatar:        seller.Avatar,
		Description:   seller.Description,
		TotalSales:    1024,
		Rating:        4.8,
		ReviewsCount:  456,
		JoinDate:      seller.CreatedAt.Format("2006-01-02"),
		ProductsCount: 128,
	}

	return resp, nil
}
