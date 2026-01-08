// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package seller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSellerProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取卖家商品列表
func NewGetSellerProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSellerProductsLogic {
	return &GetSellerProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSellerProductsLogic) GetSellerProducts(req *types.GetSellerProductsRequest, sellerId string) (resp *types.GetSellerProductsResponse, err error) {
	// 转换卖家ID
	id, err := strconv.ParseInt(sellerId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的卖家ID")
	}

	// 查询卖家商品列表
	products, total, err := l.svcCtx.Repos.ProductRepo.GetBySellerID(l.ctx, id, req.Category, req.SortBy, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取商品列表失败: %v", err)
	}

	// 构造响应
	productList := make([]types.ProductListItem, 0)
	for _, product := range products {
		// 解析图片
		var images []string
		if product.Images != "" {
			_ = json.Unmarshal([]byte(product.Images), &images)
		}
		image := ""
		if len(images) > 0 {
			image = images[0]
		}

		// 解析标签
		var tags []string
		if product.Tags != "" {
			_ = json.Unmarshal([]byte(product.Tags), &tags)
		}

		// 获取卖家信息
		seller, _ := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, product.SellerID)
		sellerInfo := types.SellerInfo{
			Id:   fmt.Sprintf("%d", product.SellerID),
			Name: "未知商户",
		}
		if seller != nil {
			sellerInfo.Name = seller.MerchantName
			sellerInfo.Avatar = seller.Avatar
			sellerInfo.Rating = 4.5
		}

		productList = append(productList, types.ProductListItem{
			Id:       fmt.Sprintf("%d", product.ID),
			Name:     product.Name,
			Category: product.Category,
			Price:    product.Price,
			Image:    image,
			Rating:   product.Rating,
			Sales:    int(product.Sales),
			Stock:    int(product.Stock),
			Tags:     tags,
			Seller:   sellerInfo,
		})
	}

	resp = &types.GetSellerProductsResponse{
		Products: productList,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	return resp, nil
}
