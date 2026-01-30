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

	// 根据分类名称查找分类ID
	var category models.Category
	db := l.svcCtx.OrmHelper.GetDB()
	err = db.Where("name = ?", req.Category).First(&category).Error
	if err != nil {
		// 如果分类不存在，创建新分类
		category = models.Category{
			Name:     req.Category,
			ParentID: 0,
			Level:    1,
			Sort:     0,
			Status:   1,
		}
		if err := db.Create(&category).Error; err != nil {
			return nil, fmt.Errorf("创建分类失败: %v", err)
		}
	}

	// 序列化JSON字段
	imagesJSON, _ := json.Marshal(req.Images)

	// 设置主图片（使用第一张图片）
	var mainImage string
	if len(req.Images) > 0 {
		mainImage = req.Images[0]
	}

	// 创建商品
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  category.ID,
		Price:       req.Price,
		Stock:       req.Stock,
		Image:       mainImage,
		Images:      string(imagesJSON),
		Status:      1, // 默认在售
		SellerID:    merchantID,
		Sales:       0,
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
