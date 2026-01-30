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

		// 解析图片
		var images []string
		if p.Images != "" {
			_ = json.Unmarshal([]byte(p.Images), &images)
		}
		image := ""
		if len(images) > 0 {
			image = images[0]
		}
		// 如果images为空，使用主图片
		if image == "" && p.Image != "" {
			image = p.Image
		}

		// 解析颜色选项
		var colors []types.ColorOption
		if p.Colors != "" {
			_ = json.Unmarshal([]byte(p.Colors), &colors)
		}
		// 如果JSON字段为空，从product_colors表查询
		if len(colors) == 0 {
			var colorRows []struct {
				Name  string
				Value string
				Hex   string
			}
			l.svcCtx.OrmHelper.GetDB().Table("product_colors").
				Select("name, value, hex").
				Where("product_id = ?", p.ID).
				Find(&colorRows)
			for _, row := range colorRows {
				colors = append(colors, types.ColorOption{
					Name:  row.Name,
					Value: row.Value,
					Hex:   row.Hex,
				})
			}
		}

		// 解析尺码
		var sizes []string
		if p.Sizes != "" {
			_ = json.Unmarshal([]byte(p.Sizes), &sizes)
		}
		// 如果JSON字段为空，从product_sizes表查询
		if len(sizes) == 0 {
			var sizeRows []struct {
				Size string
			}
			l.svcCtx.OrmHelper.GetDB().Table("product_sizes").
				Select("size").
				Where("product_id = ?", p.ID).
				Find(&sizeRows)
			for _, row := range sizeRows {
				sizes = append(sizes, row.Size)
			}
		}

		// 获取分类名称
		categoryName := ""
		if p.CategoryID > 0 {
			var category struct{ Name string }
			l.svcCtx.OrmHelper.GetDB().Table("categories").Select("name").Where("id = ?", p.CategoryID).First(&category)
			categoryName = category.Name
		}

		// 获取商户信息
		seller := types.SellerInfo{
			Id:     fmt.Sprintf("%d", p.SellerID),
			Name:   "商户",
			Rating: 5.0,
		}
		var merchantProfile struct {
			ShopName   string
			ShopAvatar string
			Rating     float64
		}
		err := l.svcCtx.OrmHelper.GetDB().Table("merchant_profiles").
			Select("shop_name, shop_avatar, rating").
			Where("user_id = ?", p.SellerID).
			First(&merchantProfile).Error
		if err == nil {
			seller.Name = merchantProfile.ShopName
			seller.Avatar = merchantProfile.ShopAvatar
			seller.Rating = merchantProfile.Rating
		}

		productList = append(productList, types.ProductListItem{
			Id:          fmt.Sprintf("%d", p.ID),
			Name:        p.Name,
			Category:    categoryName,
			Price:       p.Price,
			Image:       image,
			Images:      images,
			Rating:      p.Rating,
			Sales:       p.Sales,
			Stock:       p.Stock,
			Tags:        tags,
			Colors:      colors,
			Sizes:       sizes,
			Description: p.Description,
			Seller:      seller,
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
