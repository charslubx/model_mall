// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package search

import (
	"context"
	"encoding/json"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 全局搜索商品
func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	// 查询商品列表
	products, total, err := l.svcCtx.Repos.ProductRepo.Search(
		l.ctx,
		req.Keyword,
		req.Category,
		req.MinPrice,
		req.MaxPrice,
		req.SortBy,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("搜索失败: %v", err)
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
		seller, _ := l.svcCtx.Repos.UserRepo.GetByID(l.ctx, product.MerchantID)
		sellerInfo := types.SellerInfo{
			Id:   fmt.Sprintf("%d", product.MerchantID),
			Name: "未知商户",
		}
		if seller != nil {
			sellerInfo.Name = seller.MerchantName
			sellerInfo.Avatar = seller.Avatar
			sellerInfo.Rating = 4.5 // TODO: 从数据库获取实际评分
		}

		productList = append(productList, types.ProductListItem{
			Id:       fmt.Sprintf("%d", product.ID),
			Name:     product.Name,
			Price:    product.Price,
			Image:    image,
			Category: product.Category,
			Rating:   product.Rating,
			Sales:    int(product.Sales),
			Stock:    int(product.Stock),
			Tags:     tags,
			Seller:   sellerInfo,
		})
	}

	logx.Infof("搜索关键词: %s, 找到 %d 个结果", req.Keyword, total)

	resp = &types.SearchResponse{
		Products: productList,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}

	return resp, nil
}
