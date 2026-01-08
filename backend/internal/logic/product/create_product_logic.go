// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建商品（商户）
func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductRequest) (resp *types.CreateProductResponse, err error) {
	// 获取当前商户ID
	merchantID, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 序列化JSON字段
	imagesJSON, _ := json.Marshal(req.Images)
	tagsJSON, _ := json.Marshal(req.Tags)
	colorsJSON, _ := json.Marshal(req.Colors)
	sizesJSON, _ := json.Marshal(req.Sizes)
	featuresJSON, _ := json.Marshal(req.Features)
	specificationsJSON, _ := json.Marshal(req.Specifications)

	// 创建商品
	product := &models.Product{
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		Price:          req.Price,
		Stock:          req.Stock,
		Images:         string(imagesJSON),
		Tags:           string(tagsJSON),
		Colors:         string(colorsJSON),
		Sizes:          string(sizesJSON),
		Features:       string(featuresJSON),
		Specifications: string(specificationsJSON),
		Status:         1, // 默认在售
		SellerID:       merchantID,
		Rating:         0,
		Reviews:        0,
		Sales:          0,
	}

	err = l.svcCtx.Repos.ProductRepo.Create(l.ctx, product)
	if err != nil {
		return nil, fmt.Errorf("创建商品失败: %v", err)
	}

	// 构造响应
	resp = &types.CreateProductResponse{
		Id:        fmt.Sprintf("%d", product.ID),
		Name:      product.Name,
		Status:    "active",
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
	}

	return resp, nil
}
