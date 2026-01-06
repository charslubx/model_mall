// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product

import (
	"context"
	"encoding/json"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商户商品列表
func NewGetMerchantProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantProductsLogic {
	return &GetMerchantProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantProductsLogic) GetMerchantProducts(req *types.GetMerchantProductsRequest) (resp *types.GetMerchantProductsResponse, err error) {
	// 获取当前商户ID
	merchantID, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 调用repository获取商户商品列表
	products, total, err := l.svcCtx.Repos.ProductRepo.ListByMerchant(
		l.ctx,
		merchantID,
		req.Page,
		req.PageSize,
		req.Category,
		req.Status,
		req.Keyword,
	)
	if err != nil {
		return nil, fmt.Errorf("获取商品列表失败: %v", err)
	}

	// 转换为响应格式
	productList := make([]types.MerchantProductItem, 0, len(products))
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

		// 状态转换
		status := "active"
		if p.Status == 0 {
			status = "inactive"
		}

		productList = append(productList, types.MerchantProductItem{
			Id:        fmt.Sprintf("%d", p.ID),
			Name:      p.Name,
			Category:  p.Category,
			Price:     p.Price,
			Stock:     p.Stock,
			Status:    status,
			Image:     image,
			Tags:      tags,
			Sales:     p.Sales,
			CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp = &types.GetMerchantProductsResponse{
		Products: productList,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	return resp, nil
}
