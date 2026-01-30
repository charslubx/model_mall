// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品分类列表
func NewGetCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoriesLogic {
	return &GetCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoriesLogic) GetCategories() (resp *types.GetCategoriesResponse, err error) {
	// 查询所有分类
	categories, err := l.svcCtx.Repos.ProductRepo.GetCategories(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	// 转换为响应格式
	categoryList := make([]types.Category, 0, len(categories))
	for _, c := range categories {
		// 统计该分类下的商品数量
		var count int64
		l.svcCtx.OrmHelper.GetDB().Model(&struct {
			ID int64 `gorm:"column:id"`
		}{}).Table("products").Where("category_id = ? AND status = 1", c.ID).Count(&count)

		categoryList = append(categoryList, types.Category{
			Id:    fmt.Sprintf("%d", c.ID),
			Name:  c.Name,
			Slug:  "", // 数据库中没有slug字段，返回空字符串
			Count: int(count),
			Icon:  c.Icon,
		})
	}

	resp = &types.GetCategoriesResponse{
		Categories: categoryList,
	}

	return resp, nil
}
