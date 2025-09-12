package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"model_mall_backend/backend/internal/models"
	"gorm.io/gorm"
)

// ImageClassificationRepository 图片分类存储库
type ImageClassificationRepository struct {
	db *gorm.DB
}

// NewImageClassificationRepository 创建图片分类存储库
func NewImageClassificationRepository(db *gorm.DB) *ImageClassificationRepository {
	return &ImageClassificationRepository{
		db: db,
	}
}

// CreateClassification 创建分类记录
func (r *ImageClassificationRepository) CreateClassification(ctx context.Context, classification *models.ImageClassification) error {
	return r.db.WithContext(ctx).Create(classification).Error
}

// CreateClassificationWithLabels 创建分类记录和标签（事务）
func (r *ImageClassificationRepository) CreateClassificationWithLabels(ctx context.Context, classification *models.ImageClassification, predictions []models.ModelPrediction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建分类记录
		if err := tx.Create(classification).Error; err != nil {
			return fmt.Errorf("创建分类记录失败: %v", err)
		}
		
		// 创建标签记录
		if len(predictions) > 0 {
			labels := make([]models.ImageClassificationLabel, 0, len(predictions))
			for _, pred := range predictions {
				label := models.ImageClassificationLabel{
					ClassificationID: classification.ID,
					LabelName:        pred.Label,
					LabelCode:        pred.Code,
					Confidence:       pred.Confidence,
				}
				
				// 如果有边界框信息，序列化为JSON
				if pred.BoundingBox != nil {
					boundingBoxJSON, err := json.Marshal(pred.BoundingBox)
					if err != nil {
						return fmt.Errorf("序列化边界框失败: %v", err)
					}
					label.BoundingBox = string(boundingBoxJSON)
				}
				
				labels = append(labels, label)
			}
			
			if err := tx.CreateInBatches(labels, 100).Error; err != nil {
				return fmt.Errorf("创建标签记录失败: %v", err)
			}
		}
		
		return nil
	})
}

// GetClassificationByID 根据ID获取分类记录
func (r *ImageClassificationRepository) GetClassificationByID(ctx context.Context, id int64) (*models.ImageClassification, error) {
	var classification models.ImageClassification
	err := r.db.WithContext(ctx).
		Preload("Labels").
		Preload("User").
		First(&classification, id).Error
	if err != nil {
		return nil, err
	}
	return &classification, nil
}

// GetClassificationsByUserID 根据用户ID获取分类记录列表
func (r *ImageClassificationRepository) GetClassificationsByUserID(ctx context.Context, userID int64, page, pageSize int) ([]models.ImageClassification, int64, error) {
	var classifications []models.ImageClassification
	var total int64
	
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	// 获取总数
	if err := query.Model(&models.ImageClassification{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.
		Preload("Labels").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&classifications).Error
	
	return classifications, total, err
}

// GetClassificationsByModelName 根据模型名称获取分类记录
func (r *ImageClassificationRepository) GetClassificationsByModelName(ctx context.Context, modelName string, page, pageSize int) ([]models.ImageClassification, int64, error) {
	var classifications []models.ImageClassification
	var total int64
	
	query := r.db.WithContext(ctx).Where("model_name = ?", modelName)
	
	// 获取总数
	if err := query.Model(&models.ImageClassification{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	err := query.
		Preload("Labels").
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&classifications).Error
	
	return classifications, total, err
}

// UpdateClassificationStatus 更新分类状态
func (r *ImageClassificationRepository) UpdateClassificationStatus(ctx context.Context, id int64, status int8) error {
	return r.db.WithContext(ctx).
		Model(&models.ImageClassification{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// DeleteClassification 删除分类记录（级联删除标签）
func (r *ImageClassificationRepository) DeleteClassification(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除标签
		if err := tx.Where("classification_id = ?", id).Delete(&models.ImageClassificationLabel{}).Error; err != nil {
			return fmt.Errorf("删除标签失败: %v", err)
		}
		
		// 删除分类记录
		if err := tx.Delete(&models.ImageClassification{}, id).Error; err != nil {
			return fmt.Errorf("删除分类记录失败: %v", err)
		}
		
		return nil
	})
}

// GetStatistics 获取分类统计信息
func (r *ImageClassificationRepository) GetStatistics(ctx context.Context, userID *int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	query := r.db.WithContext(ctx).Model(&models.ImageClassification{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	
	// 总分类数
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_classifications"] = totalCount
	
	// 成功分类数
	var successCount int64
	if err := query.Where("status = ?", 1).Count(&successCount).Error; err != nil {
		return nil, err
	}
	stats["success_classifications"] = successCount
	
	// 失败分类数
	var failedCount int64
	if err := query.Where("status = ?", 0).Count(&failedCount).Error; err != nil {
		return nil, err
	}
	stats["failed_classifications"] = failedCount
	
	// 按模型统计
	var modelStats []struct {
		ModelName string `json:"model_name"`
		Count     int64  `json:"count"`
	}
	if err := query.Select("model_name, count(*) as count").
		Group("model_name").
		Find(&modelStats).Error; err != nil {
		return nil, err
	}
	stats["model_statistics"] = modelStats
	
	// 按日期统计（最近7天）
	var dateStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := query.Select("DATE(created_at) as date, count(*) as count").
		Where("created_at >= CURRENT_DATE - INTERVAL '7 days'").
		Group("DATE(created_at)").
		Order("date DESC").
		Find(&dateStats).Error; err != nil {
		return nil, err
	}
	stats["date_statistics"] = dateStats
	
	return stats, nil
}

// SearchClassifications 搜索分类记录
func (r *ImageClassificationRepository) SearchClassifications(ctx context.Context, req *SearchClassificationReq) ([]models.ImageClassification, int64, error) {
	var classifications []models.ImageClassification
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.ImageClassification{})
	
	// 添加搜索条件
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	
	if req.ModelName != "" {
		query = query.Where("model_name LIKE ?", "%"+req.ModelName+"%")
	}
	
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	
	if req.ImageName != "" {
		query = query.Where("image_name LIKE ?", "%"+req.ImageName+"%")
	}
	
	if req.StartTime != nil {
		query = query.Where("created_at >= ?", *req.StartTime)
	}
	
	if req.EndTime != nil {
		query = query.Where("created_at <= ?", *req.EndTime)
	}
	
	if req.MinConfidence != nil {
		query = query.Where("confidence >= ?", *req.MinConfidence)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Preload("Labels").
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&classifications).Error
	
	return classifications, total, err
}

// SearchClassificationReq 搜索分类请求
type SearchClassificationReq struct {
	UserID        *int64     `json:"user_id"`
	ModelName     string     `json:"model_name"`
	Status        *int8      `json:"status"`
	ImageName     string     `json:"image_name"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	MinConfidence *float64   `json:"min_confidence"`
	Page          int        `json:"page"`
	PageSize      int        `json:"page_size"`
}