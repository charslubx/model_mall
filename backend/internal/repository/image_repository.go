package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"modelmall/backend/internal/models"
)


type ImageRepository struct {
	db *gorm.DB
}

// NewImageRepository 创建图片仓储实例
func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		db: db,
	}
}

// CreateImage 创建图片记录
func (r *ImageRepository) CreateImage(ctx context.Context, image *models.Image) (int64, error) {
	now := time.Now()
	image.CreatedAt = now
	image.UpdatedAt = now
	
	result := r.db.WithContext(ctx).Create(image)
	if result.Error != nil {
		return 0, fmt.Errorf("创建图片记录失败: %w", result.Error)
	}
	
	return image.ID, nil
}

// GetImageByID 根据ID获取图片
func (r *ImageRepository) GetImageByID(ctx context.Context, id int64) (*models.Image, error) {
	var image models.Image
	result := r.db.WithContext(ctx).First(&image, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("图片不存在")
		}
		return nil, fmt.Errorf("查询图片失败: %w", result.Error)
	}
	
	return &image, nil
}

// UpdateImageStatus 更新图片状态
func (r *ImageRepository) UpdateImageStatus(ctx context.Context, id int64, status int) error {
	result := r.db.WithContext(ctx).Model(&models.Image{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})
	
	if result.Error != nil {
		return fmt.Errorf("更新图片状态失败: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("图片不存在")
	}
	
	return nil
}

// ListImages 列出图片（带分页）
func (r *ImageRepository) ListImages(ctx context.Context, offset, limit int) ([]*models.Image, error) {
	var images []*models.Image
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&images)
	
	if result.Error != nil {
		return nil, fmt.Errorf("查询图片列表失败: %w", result.Error)
	}
	
	return images, nil
}

// CreateClassification 创建分类记录
func (r *ImageRepository) CreateClassification(ctx context.Context, classification *models.ImageClassification) error {
	classification.CreatedAt = time.Now()
	
	result := r.db.WithContext(ctx).Create(classification)
	if result.Error != nil {
		return fmt.Errorf("创建分类记录失败: %w", result.Error)
	}
	
	return nil
}

// GetClassificationsByImageID 获取图片的所有分类标签
func (r *ImageRepository) GetClassificationsByImageID(ctx context.Context, imageID int64) ([]*models.ImageClassification, error) {
	var classifications []*models.ImageClassification
	result := r.db.WithContext(ctx).
		Where("image_id = ?", imageID).
		Order("confidence DESC").
		Find(&classifications)
	
	if result.Error != nil {
		return nil, fmt.Errorf("查询分类记录失败: %w", result.Error)
	}
	
	return classifications, nil
}

// DeleteImage 删除图片记录
func (r *ImageRepository) DeleteImage(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.Image{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除图片失败: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("图片不存在")
	}
	
	return nil
}