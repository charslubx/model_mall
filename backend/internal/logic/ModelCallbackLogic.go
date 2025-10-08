package logic

import (
	"context"
	"fmt"
	"strings"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModelCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModelCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModelCallbackLogic {
	return &ModelCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// HandleCallback 处理模型服务回调
func (l *ModelCallbackLogic) HandleCallback(req *types.ModelCallbackReq) (*types.ModelCallbackResp, error) {
	// 获取任务信息
	task, err := l.svcCtx.Repos.RecognitionTask.GetByTaskID(req.TaskID)
	if err != nil {
		logx.Errorf("任务不存在: %v", err)
		return &types.ModelCallbackResp{
			Success: false,
			Message: "任务不存在",
		}, nil
	}

	// 根据状态处理
	switch strings.ToLower(req.Status) {
	case "pending":
		// 待处理状态
		err = l.svcCtx.Repos.RecognitionTask.UpdateStatus(req.TaskID, models.TaskStatusPending)
		
	case "processing":
		// 处理中状态，更新进度
		err = l.svcCtx.Repos.RecognitionTask.UpdateStatus(req.TaskID, models.TaskStatusProcessing)
		if err == nil && req.Progress > 0 {
			err = l.svcCtx.Repos.RecognitionTask.UpdateProgress(req.TaskID, req.Progress)
		}
		
	case "completed":
		// 完成状态，保存识别结果
		err = l.handleCompletedTask(task, req)
		
	case "failed":
		// 失败状态，保存错误信息
		err = l.svcCtx.Repos.RecognitionTask.UpdateError(req.TaskID, req.Error)
		
	default:
		logx.Errorf("未知的任务状态: %s", req.Status)
		return &types.ModelCallbackResp{
			Success: false,
			Message: "未知的任务状态",
		}, nil
	}

	if err != nil {
		logx.Errorf("更新任务状态失败: %v", err)
		return &types.ModelCallbackResp{
			Success: false,
			Message: fmt.Sprintf("更新任务状态失败: %v", err),
		}, nil
	}

	return &types.ModelCallbackResp{
		Success: true,
		Message: "success",
	}, nil
}

// handleCompletedTask 处理已完成的任务
func (l *ModelCallbackLogic) handleCompletedTask(task *models.RecognitionTask, req *types.ModelCallbackReq) error {
	// 更新任务状态为完成
	if err := l.svcCtx.Repos.RecognitionTask.UpdateStatus(req.TaskID, models.TaskStatusCompleted); err != nil {
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	// 如果没有识别结果，直接返回
	if len(req.Results) == 0 {
		return nil
	}

	// 保存识别结果
	labels := make([]models.ClassificationLabel, 0, len(req.Results))
	for _, result := range req.Results {
		label := models.ClassificationLabel{
			TaskID:     task.ID,
			ImageID:    task.ImageID,
			LabelName:  result.Name,
			LabelCode:  result.Code,
			Confidence: result.Confidence,
		}

		// 保存边界框信息
		if result.BBox != nil {
			label.BboxX = result.BBox.X
			label.BboxY = result.BBox.Y
			label.BboxWidth = result.BBox.Width
			label.BboxHeight = result.BBox.Height
		}

		// 保存额外数据
		if result.Extra != nil {
			label.ExtraData = result.Extra
		}

		labels = append(labels, label)
	}

	// 批量保存标签
	if err := l.svcCtx.Repos.ClassificationLabel.BatchCreate(labels); err != nil {
		return fmt.Errorf("保存识别结果失败: %w", err)
	}

	// 更新任务的结果数量
	task.ResultCount = len(labels)
	if err := l.svcCtx.Repos.RecognitionTask.Update(task); err != nil {
		logx.Errorf("更新任务结果数量失败: %v", err)
	}

	return nil
}
