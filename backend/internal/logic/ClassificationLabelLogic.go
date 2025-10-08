package logic

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClassificationLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClassificationLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClassificationLabelLogic {
	return &ClassificationLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetImageLabels 获取图片标签
func (l *ClassificationLabelLogic) GetImageLabels(req *types.GetImageLabelsReq) (*types.GetImageLabelsResp, error) {
	// 获取标签列表
	labels, err := l.svcCtx.Repos.ClassificationLabel.GetByImageID(req.ImageID)
	if err != nil {
		return nil, fmt.Errorf("获取标签列表失败: %w", err)
	}

	labelInfos := make([]types.LabelInfo, 0, len(labels))
	for _, label := range labels {
		info := types.LabelInfo{
			ID:         label.ID,
			LabelName:  label.LabelName,
			LabelCode:  label.LabelCode,
			Confidence: label.Confidence,
			CreatedAt:  label.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// 添加边界框信息（如果存在）
		if label.BboxWidth > 0 && label.BboxHeight > 0 {
			info.BBox = &types.BoundingBox{
				X:      label.BboxX,
				Y:      label.BboxY,
				Width:  label.BboxWidth,
				Height: label.BboxHeight,
			}
		}

		// 添加额外数据
		if label.ExtraData != nil {
			info.ExtraData = label.ExtraData
		}

		labelInfos = append(labelInfos, info)
	}

	return &types.GetImageLabelsResp{
		ImageID: req.ImageID,
		Labels:  labelInfos,
	}, nil
}
