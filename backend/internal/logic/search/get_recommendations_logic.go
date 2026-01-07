// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取推荐商品
func NewGetRecommendationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendationsLogic {
	return &GetRecommendationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecommendationsLogic) GetRecommendations(req *types.GetRecommendationsRequest) (resp *types.GetRecommendationsResponse, err error) {
	var products []types.RecommendationProduct
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}

	// 基于商品推荐
	if req.ProductId != "" {
		productId, err := strconv.ParseInt(req.ProductId, 10, 64)
		if err == nil {
			product, _ := l.svcCtx.Repos.ProductRepo.GetByID(l.ctx, productId)
			if product != nil {
				// 查询同类商品
				similarProducts, _, _ := l.svcCtx.Repos.ProductRepo.Search(
					l.ctx,
					"",
					product.Category,
					0,
					0,
					"sales",
					1,
					limit,
				)

				for _, p := range similarProducts {
					if p.ID == productId {
						continue // 跳过自己
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

					products = append(products, types.RecommendationProduct{
						Id:     fmt.Sprintf("%d", p.ID),
						Name:   p.Name,
						Price:  p.Price,
						Image:  image,
						Rating: p.Rating,
						Reason: "相似商品",
					})

					if len(products) >= limit {
						break
					}
				}
			}
		}
	} else {
		// 基于用户推荐或热门商品
		hotProducts, _, _ := l.svcCtx.Repos.ProductRepo.Search(
			l.ctx,
			"",
			"",
			0,
			0,
			"sales",
			1,
			limit,
		)

		for _, p := range hotProducts {
			// 解析图片
			var images []string
			if p.Images != "" {
				_ = json.Unmarshal([]byte(p.Images), &images)
			}
			image := ""
			if len(images) > 0 {
				image = images[0]
			}

			products = append(products, types.RecommendationProduct{
				Id:     fmt.Sprintf("%d", p.ID),
				Name:   p.Name,
				Price:  p.Price,
				Image:  image,
				Rating: p.Rating,
				Reason: "热门商品",
			})
		}
	}

	logx.Infof("生成推荐商品 %d 个", len(products))

	resp = &types.GetRecommendationsResponse{
		Products: products,
	}

	return resp, nil
}
