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

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新商品（商户）
func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductRequest, productId string) (resp *types.UpdateProductResponse, err error) {
	// 获取当前商户ID
	merchantID, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 转换商品ID
	id, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的商品ID")
	}

	// 查询商品
	product, err := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("商品不存在")
	}

	// 验证权限
	if product.MerchantID != merchantID {
		return nil, fmt.Errorf("无权限操作此商品")
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Stock >= 0 {
		updates["stock"] = req.Stock
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if len(req.Images) > 0 {
		imagesJSON, _ := json.Marshal(req.Images)
		updates["images"] = string(imagesJSON)
	}
	if len(req.Tags) > 0 {
		tagsJSON, _ := json.Marshal(req.Tags)
		updates["tags"] = string(tagsJSON)
	}
	if len(req.Colors) > 0 {
		colorsJSON, _ := json.Marshal(req.Colors)
		updates["colors"] = string(colorsJSON)
	}
	if len(req.Sizes) > 0 {
		sizesJSON, _ := json.Marshal(req.Sizes)
		updates["sizes"] = string(sizesJSON)
	}
	if len(req.Features) > 0 {
		featuresJSON, _ := json.Marshal(req.Features)
		updates["features"] = string(featuresJSON)
	}
	if req.Specifications.Material != "" || req.Specifications.Care != "" {
		specificationsJSON, _ := json.Marshal(req.Specifications)
		updates["specifications"] = string(specificationsJSON)
	}

	// 执行更新
	err = l.svcCtx.Repos.ProductRepo.UpdateByID(l.ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("更新商品失败: %v", err)
	}

	resp = &types.UpdateProductResponse{
		Id:        productId,
		Updated:   true,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return resp, nil
}
