package logic

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecognitionTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecognitionTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecognitionTaskLogic {
	return &RecognitionTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetTaskStatus 获取任务状态
func (l *RecognitionTaskLogic) GetTaskStatus(req *types.GetTaskStatusReq) (*types.GetTaskStatusResp, error) {
	// 获取任务信息
	task, err := l.svcCtx.Repos.RecognitionTask.GetByTaskID(req.TaskID)
	if err != nil {
		return nil, fmt.Errorf("任务不存在: %w", err)
	}

	resp := &types.GetTaskStatusResp{
		TaskID:      task.TaskID,
		ImageID:     task.ImageID,
		Status:      task.Status,
		Progress:    task.Progress,
		ResultCount: task.ResultCount,
		ErrorMsg:    task.ErrorMessage,
		CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if task.CompletedAt.Valid {
		resp.CompletedAt = task.CompletedAt.Time.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}

// GetTaskList 获取任务列表
func (l *RecognitionTaskLogic) GetTaskList(req *types.GetTaskListReq) (*types.GetTaskListResp, error) {
	userID := l.getUserIDFromContext()
	if userID == 0 {
		return nil, fmt.Errorf("未授权")
	}

	tasks, total, err := l.svcCtx.Repos.RecognitionTask.GetByUserID(userID, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取任务列表失败: %w", err)
	}

	list := make([]types.TaskInfo, 0, len(tasks))
	for _, task := range tasks {
		info := types.TaskInfo{
			TaskID:      task.TaskID,
			ImageID:     task.ImageID,
			ModelName:   task.ModelName,
			Status:      task.Status,
			Progress:    task.Progress,
			ResultCount: task.ResultCount,
			CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		
		if task.CompletedAt.Valid {
			info.CompletedAt = task.CompletedAt.Time.Format("2006-01-02 15:04:05")
		}
		
		list = append(list, info)
	}

	return &types.GetTaskListResp{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

func (l *RecognitionTaskLogic) getUserIDFromContext() int64 {
	userIDVal := l.ctx.Value("user_id")
	if userIDVal == nil {
		return 0
	}
	
	if userID, ok := userIDVal.(int64); ok {
		return userID
	}
	
	return 0
}
