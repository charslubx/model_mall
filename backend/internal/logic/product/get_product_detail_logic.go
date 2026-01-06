// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品详情
func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductDetailLogic) GetProductDetail(productId string) (resp *types.ProductDetail, err error) {
	// 转换商品ID
	id, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的商品ID")
	}

	// 查询商品信息
	product, err := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("商品不存在")
	}

	// 检查商品状态
	if product.Status == 0 {
		return nil, fmt.Errorf("商品已下架")
	}

	// 解析JSON字段
	var images []string
	if product.Images != "" {
		_ = json.Unmarshal([]byte(product.Images), &images)
	}

	var colors []types.ColorOption
	if product.Colors != "" {
		_ = json.Unmarshal([]byte(product.Colors), &colors)
	}

	var sizes []string
	if product.Sizes != "" {
		_ = json.Unmarshal([]byte(product.Sizes), &sizes)
	}

	var tags []string
	if product.Tags != "" {
		_ = json.Unmarshal([]byte(product.Tags), &tags)
	}

	var features []string
	if product.Features != "" {
		_ = json.Unmarshal([]byte(product.Features), &features)
	}

	var specifications types.Specifications
	if product.Specifications != "" {
		_ = json.Unmarshal([]byte(product.Specifications), &specifications)
	}

	// 获取卖家信息
	seller := types.SellerInfo{
		Id:     fmt.Sprintf("%d", product.MerchantID),
		Name:   "商户名称", // TODO: 从用户表查询
		Rating: 4.8,
	}

	// 获取相关推荐商品(简化处理)
	relatedProducts := []types.RelatedProduct{}

	// 构造响应
	resp = &types.ProductDetail{
		Id:              fmt.Sprintf("%d", product.ID),
		Name:            product.Name,
		Description:     product.Description,
		Category:        product.Category,
		Price:           product.Price,
		Stock:           product.Stock,
		Rating:          product.Rating,
		Reviews:         product.Reviews,
		Images:          images,
		Colors:          colors,
		Sizes:           sizes,
		Tags:            tags,
		Features:        features,
		Specifications:  specifications,
		Seller:          seller,
		RelatedProducts: relatedProducts,
		CreatedAt:       product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       product.UpdatedAt.Format(time.RFC3339),
	}

	return resp, nil
}
