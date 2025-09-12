package logic

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/repository"
	"model_mall_backend/backend/internal/service"
	"model_mall_backend/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// ImageClassificationLogic 图片分类业务逻辑
type ImageClassificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	repo   *repository.ImageClassificationRepository
	modelService *service.ModelService
}

// NewImageClassificationLogic 创建图片分类业务逻辑
func NewImageClassificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageClassificationLogic {
	return &ImageClassificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		repo:   repository.NewImageClassificationRepository(svcCtx.DB),
		modelService: service.NewModelService(
			svcCtx.Config.ModelService.Endpoint,
			svcCtx.Config.Upload.Path,
		),
	}
}

// ClassifyImage 分类图片
func (l *ImageClassificationLogic) ClassifyImage(req *models.ImageClassificationReq) (*models.ImageClassificationResp, error) {
	// 验证图片数据
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if err := l.modelService.ValidateImage(req.ImageData, maxSize); err != nil {
		return nil, fmt.Errorf("图片验证失败: %v", err)
	}
	
	// 创建分类记录
	classification := &models.ImageClassification{
		ImageName:   req.ImageName,
		ImageSize:   int64(len(req.ImageData)),
		ImageFormat: l.modelService.GetImageFormat(req.ImageData),
		ModelName:   req.ModelName,
		Status:      2, // 处理中
		UserID:      req.UserID,
	}
	
	// 如果需要保存图片
	if req.SaveImage {
		imagePath, err := l.modelService.SaveImage(req.ImageData, req.ImageName)
		if err != nil {
			l.Errorf("保存图片失败: %v", err)
		} else {
			classification.ImagePath = imagePath
		}
	}
	
	// 先创建记录
	if err := l.repo.CreateClassification(l.ctx, classification); err != nil {
		return nil, fmt.Errorf("创建分类记录失败: %v", err)
	}
	
	// 调用模型服务
	modelResp, err := l.modelService.ProcessImage(l.ctx, req)
	if err != nil {
		// 更新状态为失败
		l.repo.UpdateClassificationStatus(l.ctx, classification.ID, 0)
		return nil, fmt.Errorf("模型处理失败: %v", err)
	}
	
	if !modelResp.Success {
		// 更新状态为失败
		l.repo.UpdateClassificationStatus(l.ctx, classification.ID, 0)
		return nil, fmt.Errorf("模型处理失败: %s", modelResp.Message)
	}
	
	// 更新分类记录
	classification.ProcessTime = modelResp.ProcessTime
	classification.Status = 1 // 成功
	
	// 计算总体置信度（取最高置信度）
	if len(modelResp.Predictions) > 0 {
		maxConfidence := 0.0
		for _, pred := range modelResp.Predictions {
			if pred.Confidence > maxConfidence {
				maxConfidence = pred.Confidence
			}
		}
		classification.Confidence = maxConfidence
	}
	
	// 保存分类结果和标签
	if err := l.repo.CreateClassificationWithLabels(l.ctx, classification, modelResp.Predictions); err != nil {
		return nil, fmt.Errorf("保存分类结果失败: %v", err)
	}
	
	// 重新获取完整的分类记录
	result, err := l.repo.GetClassificationByID(l.ctx, classification.ID)
	if err != nil {
		return nil, fmt.Errorf("获取分类结果失败: %v", err)
	}
	
	// 转换为响应格式
	resp := &models.ImageClassificationResp{
		ID:          result.ID,
		ImagePath:   result.ImagePath,
		ImageName:   result.ImageName,
		ModelName:   result.ModelName,
		ProcessTime: result.ProcessTime,
		Confidence:  result.Confidence,
		Status:      result.Status,
		Labels:      result.Labels,
		CreatedAt:   result.CreatedAt,
	}
	
	return resp, nil
}

// GetClassification 获取分类记录
func (l *ImageClassificationLogic) GetClassification(id int64) (*models.ImageClassificationResp, error) {
	classification, err := l.repo.GetClassificationByID(l.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取分类记录失败: %v", err)
	}
	
	resp := &models.ImageClassificationResp{
		ID:          classification.ID,
		ImagePath:   classification.ImagePath,
		ImageName:   classification.ImageName,
		ModelName:   classification.ModelName,
		ProcessTime: classification.ProcessTime,
		Confidence:  classification.Confidence,
		Status:      classification.Status,
		Labels:      classification.Labels,
		CreatedAt:   classification.CreatedAt,
	}
	
	return resp, nil
}

// GetUserClassifications 获取用户分类记录列表
func (l *ImageClassificationLogic) GetUserClassifications(userID int64, page, pageSize int) (*ClassificationListResp, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	
	classifications, total, err := l.repo.GetClassificationsByUserID(l.ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("获取用户分类记录失败: %v", err)
	}
	
	// 转换为响应格式
	items := make([]models.ImageClassificationResp, 0, len(classifications))
	for _, classification := range classifications {
		items = append(items, models.ImageClassificationResp{
			ID:          classification.ID,
			ImagePath:   classification.ImagePath,
			ImageName:   classification.ImageName,
			ModelName:   classification.ModelName,
			ProcessTime: classification.ProcessTime,
			Confidence:  classification.Confidence,
			Status:      classification.Status,
			Labels:      classification.Labels,
			CreatedAt:   classification.CreatedAt,
		})
	}
	
	return &ClassificationListResp{
		Items: items,
		Total: total,
		Page:  page,
		PageSize: pageSize,
	}, nil
}

// GetStatistics 获取统计信息
func (l *ImageClassificationLogic) GetStatistics(userID *int64) (map[string]interface{}, error) {
	return l.repo.GetStatistics(l.ctx, userID)
}

// DeleteClassification 删除分类记录
func (l *ImageClassificationLogic) DeleteClassification(id int64) error {
	return l.repo.DeleteClassification(l.ctx, id)
}

// SearchClassifications 搜索分类记录
func (l *ImageClassificationLogic) SearchClassifications(req *repository.SearchClassificationReq) (*ClassificationListResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	
	classifications, total, err := l.repo.SearchClassifications(l.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("搜索分类记录失败: %v", err)
	}
	
	// 转换为响应格式
	items := make([]models.ImageClassificationResp, 0, len(classifications))
	for _, classification := range classifications {
		items = append(items, models.ImageClassificationResp{
			ID:          classification.ID,
			ImagePath:   classification.ImagePath,
			ImageName:   classification.ImageName,
			ModelName:   classification.ModelName,
			ProcessTime: classification.ProcessTime,
			Confidence:  classification.Confidence,
			Status:      classification.Status,
			Labels:      classification.Labels,
			CreatedAt:   classification.CreatedAt,
		})
	}
	
	return &ClassificationListResp{
		Items: items,
		Total: total,
		Page:  req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ClassificationListResp 分类列表响应
type ClassificationListResp struct {
	Items    []models.ImageClassificationResp `json:"items"`
	Total    int64                            `json:"total"`
	Page     int                              `json:"page"`
	PageSize int                              `json:"page_size"`
}