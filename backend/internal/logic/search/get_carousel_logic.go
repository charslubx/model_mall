// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package search

import (
	"context"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCarouselLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取轮播图数据
func NewGetCarouselLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCarouselLogic {
	return &GetCarouselLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCarouselLogic) GetCarousel() (resp *types.GetCarouselResponse, err error) {
	// TODO: 从数据库查询轮播图配置
	// 这里使用模拟数据
	banners := []types.Banner{
		{
			Id:    "banner001",
			Image: "https://cdn.example.com/banners/banner1.jpg",
			Title: "新春大促",
			Link:  "/products?category=上衣",
			Order: 1,
		},
		{
			Id:    "banner002",
			Image: "https://cdn.example.com/banners/banner2.jpg",
			Title: "冬季新品",
			Link:  "/products?tag=新品",
			Order: 2,
		},
		{
			Id:    "banner003",
			Image: "https://cdn.example.com/banners/banner3.jpg",
			Title: "限时特惠",
			Link:  "/products?sortBy=price-asc",
			Order: 3,
		},
	}

	logx.Info("获取轮播图数据成功")

	resp = &types.GetCarouselResponse{
		Banners: banners,
	}

	return resp, nil
}
