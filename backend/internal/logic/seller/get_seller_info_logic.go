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

func (l *GetSellerInfoLogic) GetSellerInfo(sellerId string) (resp *types.SellerDetail, err error) {
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

	// 从merchant_profiles表查询商户详细信息
	var merchantProfile struct {
		ShopName        string
		ShopDescription string
		ShopAvatar      string
		Rating          float64
		TotalSales      int
		ProductsCount   int
		ReviewsCount    int
		JoinDate        string
	}

	err = l.svcCtx.OrmHelper.GetDB().Table("merchant_profiles").
		Select("shop_name, shop_description, shop_avatar, rating, total_sales, products_count, reviews_count, join_date").
		Where("user_id = ?", id).
		First(&merchantProfile).Error

	if err != nil {
		logx.Errorf("查询商户详情失败 user_id=%d, error=%v", id, err)
		// 如果查不到merchant_profiles，使用默认值
		resp = &types.SellerDetail{
			Id:            sellerId,
			Name:          seller.MerchantName,
			Avatar:        seller.Avatar,
			Description:   seller.Description,
			TotalSales:    0,
			Rating:        5.0,
			ReviewsCount:  0,
			JoinDate:      seller.CreatedAt.Format("2006-01-02"),
			ProductsCount: 0,
		}
	} else {
		// 使用merchant_profiles中的数据
		resp = &types.SellerDetail{
			Id:            sellerId,
			Name:          merchantProfile.ShopName,
			Avatar:        merchantProfile.ShopAvatar,
			Description:   merchantProfile.ShopDescription,
			TotalSales:    merchantProfile.TotalSales,
			Rating:        merchantProfile.Rating,
			ReviewsCount:  merchantProfile.ReviewsCount,
			JoinDate:      merchantProfile.JoinDate,
			ProductsCount: merchantProfile.ProductsCount,
		}
	}

	return resp, nil
}
