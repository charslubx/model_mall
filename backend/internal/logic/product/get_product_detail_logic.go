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

	// 从product_features表查询商品特点
	var features []string
	var featureRows []struct {
		Feature string
		Sort    int
	}
	result := l.svcCtx.OrmHelper.GetDB().Table("product_features").
		Select("feature, sort").
		Where("product_id = ?", product.ID).
		Order("sort ASC, id ASC").
		Find(&featureRows)

	logx.Infof("Query product_features for product_id=%d, rows=%d, error=%v", product.ID, len(featureRows), result.Error)

	for _, row := range featureRows {
		features = append(features, row.Feature)
	}

	// 从product_specifications表查询商品规格参数
	var specRows []struct {
		SpecKey   string
		SpecValue string
		Sort      int
	}
	result2 := l.svcCtx.OrmHelper.GetDB().Table("product_specifications").
		Select("spec_key, spec_value, sort").
		Where("product_id = ?", product.ID).
		Order("sort ASC, id ASC").
		Find(&specRows)

	logx.Infof("Query product_specifications for product_id=%d, rows=%d, error=%v", product.ID, len(specRows), result2.Error)

	// 构建规格参数对象 - 将所有规格参数放入map中
	specifications := make(map[string]interface{})
	for _, row := range specRows {
		specifications[row.SpecKey] = row.SpecValue
	}
	logx.Infof("Specifications map: %+v", specifications)

	// 获取分类名称
	categoryName := ""
	if product.CategoryID > 0 {
		var category struct{ Name string }
		l.svcCtx.OrmHelper.GetDB().Table("categories").Select("name").Where("id = ?", product.CategoryID).First(&category)
		categoryName = category.Name
	}

	// 获取卖家信息 - 从merchant_profiles表查询
	seller := types.SellerInfo{
		Id:     fmt.Sprintf("%d", product.SellerID),
		Name:   "商户",
		Rating: 5.0,
	}

	var merchantProfile struct {
		ShopName   string
		ShopAvatar string
		Rating     float64
	}

	err = l.svcCtx.OrmHelper.GetDB().Table("merchant_profiles").
		Select("shop_name, shop_avatar, rating").
		Where("user_id = ?", product.SellerID).
		First(&merchantProfile).Error

	if err == nil {
		seller.Name = merchantProfile.ShopName
		seller.Avatar = merchantProfile.ShopAvatar
		seller.Rating = merchantProfile.Rating
	} else {
		logx.Infof("未找到商户信息 seller_id=%d, error=%v", product.SellerID, err)
	}

	// 获取相关推荐商品（基于同类目的商品）
	relatedProducts := []types.RelatedProduct{}
	var relatedProductsData []struct {
		ID    int64
		Name  string
		Price float64
		Image string
	}

	result3 := l.svcCtx.OrmHelper.GetDB().Table("products").
		Select("id, name, price, image").
		Where("category_id = ? AND id != ? AND status = 1", product.CategoryID, product.ID).
		Order("sales DESC, rating DESC").
		Limit(6).
		Find(&relatedProductsData)

	logx.Infof("Query related products for category_id=%d, rows=%d, error=%v", product.CategoryID, len(relatedProductsData), result3.Error)

	for _, p := range relatedProductsData {
		relatedProducts = append(relatedProducts, types.RelatedProduct{
			Id:    fmt.Sprintf("%d", p.ID),
			Name:  p.Name,
			Price: p.Price,
			Image: p.Image,
		})
	}

	// 构造响应
	resp = &types.ProductDetail{
		Id:              fmt.Sprintf("%d", product.ID),
		Name:            product.Name,
		Description:     product.Description,
		Category:        categoryName,
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
