// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品列表
func NewGetProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductsLogic {
	return &GetProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductsLogic) GetProducts(req *types.GetProductsRequest) (resp *types.GetProductsResponse, err error) {
	// 调用repository获取商品列表
	products, total, err := l.svcCtx.Repos.ProductRepo.List(
		l.ctx,
		req.Page,
		req.PageSize,
		req.Category,
		req.Keyword,
		req.SortBy,
		req.MinPrice,
		req.MaxPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("获取商品列表失败: %v", err)
	}

	// 转换为响应格式
	productList := make([]types.ProductListItem, 0, len(products))
	for _, p := range products {
		// 解析tags
		var tags []string
		if p.Tags != "" {
			_ = json.Unmarshal([]byte(p.Tags), &tags)
		}

		// 解析第一张图片
		var images []string
		if p.Images != "" {
			_ = json.Unmarshal([]byte(p.Images), &images)
		}
		image := ""
		if len(images) > 0 {
			image = images[0]
		}

		// 获取卖家信息(简化处理,实际应该查询用户表)
		seller := types.SellerInfo{
			Id:   fmt.Sprintf("%d", p.SellerID),
			Name: "商户名称", // TODO: 从用户表查询
		}

		productList = append(productList, types.ProductListItem{
			Id:       fmt.Sprintf("%d", p.ID),
			Name:     p.Name,
			Category: p.Category,
			Price:    p.Price,
			Image:    image,
			Rating:   p.Rating,
			Sales:    p.Sales,
			Stock:    p.Stock,
			Tags:     tags,
			Seller:   seller,
		})
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	resp = &types.GetProductsResponse{
		Products:   productList,
		Total:      int(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	return resp, nil
}
