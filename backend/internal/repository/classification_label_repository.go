package repository

import (
	"model_mall_backend/backend/internal/models"

	"gorm.io/gorm"
)

type ClassificationLabelRepository struct {
	db *gorm.DB
}

func NewClassificationLabelRepository(db *gorm.DB) *ClassificationLabelRepository {
	return &ClassificationLabelRepository{db: db}
}

// Create 创建分类标签
func (r *ClassificationLabelRepository) Create(label *models.ClassificationLabel) error {
	return r.db.Create(label).Error
}

// BatchCreate 批量创建分类标签
func (r *ClassificationLabelRepository) BatchCreate(labels []models.ClassificationLabel) error {
	return r.db.Create(&labels).Error
}

// GetByTaskID 根据任务ID获取标签列表
func (r *ClassificationLabelRepository) GetByTaskID(taskID int64) ([]models.ClassificationLabel, error) {
	var labels []models.ClassificationLabel
	err := r.db.Where("task_id = ?", taskID).Order("confidence DESC").Find(&labels).Error
	return labels, err
}

// GetByImageID 根据图片ID获取标签列表
func (r *ClassificationLabelRepository) GetByImageID(imageID int64) ([]models.ClassificationLabel, error) {
	var labels []models.ClassificationLabel
	err := r.db.Where("image_id = ?", imageID).Order("confidence DESC").Find(&labels).Error
	return labels, err
}

// GetByImageIDWithPagination 根据图片ID获取标签列表（带分页）
func (r *ClassificationLabelRepository) GetByImageIDWithPagination(imageID int64, page, pageSize int) ([]models.ClassificationLabel, int64, error) {
	var labels []models.ClassificationLabel
	var total int64

	offset := (page - 1) * pageSize
	
	if err := r.db.Model(&models.ClassificationLabel{}).Where("image_id = ?", imageID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("image_id = ?", imageID).
		Order("confidence DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&labels).Error

	return labels, total, err
}

// GetByLabelName 根据标签名称搜索
func (r *ClassificationLabelRepository) GetByLabelName(labelName string, page, pageSize int) ([]models.ClassificationLabel, int64, error) {
	var labels []models.ClassificationLabel
	var total int64

	offset := (page - 1) * pageSize
	
	if err := r.db.Model(&models.ClassificationLabel{}).Where("label_name LIKE ?", "%"+labelName+"%").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("label_name LIKE ?", "%"+labelName+"%").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&labels).Error

	return labels, total, err
}

// DeleteByTaskID 删除任务相关的所有标签
func (r *ClassificationLabelRepository) DeleteByTaskID(taskID int64) error {
	return r.db.Where("task_id = ?", taskID).Delete(&models.ClassificationLabel{}).Error
}

// DeleteByImageID 删除图片相关的所有标签
func (r *ClassificationLabelRepository) DeleteByImageID(imageID int64) error {
	return r.db.Where("image_id = ?", imageID).Delete(&models.ClassificationLabel{}).Error
}
