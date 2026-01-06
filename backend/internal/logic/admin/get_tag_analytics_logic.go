// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTagAnalyticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取标签分析数据
func NewGetTagAnalyticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTagAnalyticsLogic {
	return &GetTagAnalyticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTagAnalyticsLogic) GetTagAnalytics() (resp *types.TagAnalytics, err error) {
	// 获取管理员ID
	adminId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// TODO: 实现标签统计分析
	// 这里使用模拟数据
	resp = &types.TagAnalytics{
		TopTags: []types.TagStat{
			{Name: "舒适", Count: 1245, Percentage: 18},
			{Name: "时尚", Count: 1120, Percentage: 16},
			{Name: "百搭", Count: 980, Percentage: 14},
			{Name: "简约", Count: 856, Percentage: 12},
			{Name: "休闲", Count: 745, Percentage: 11},
		},
		TagsByCategory: []types.CategoryTags{
			{
				Category: "上衣",
				TopTags:  []string{"舒适", "透气", "百搭", "简约"},
			},
			{
				Category: "裤装",
				TopTags:  []string{"修身", "显瘦", "百搭", "时尚"},
			},
			{
				Category: "裙装",
				TopTags:  []string{"优雅", "时尚", "显瘦", "甜美"},
			},
		},
		TagGrowth: []types.TagGrowth{
			{Name: "新品", Growth: 25},
			{Name: "时尚", Growth: 15},
			{Name: "热卖", Growth: 12},
		},
	}

	logx.Infof("管理员 %d 查询标签分析数据", adminId)

	return resp, nil
}
