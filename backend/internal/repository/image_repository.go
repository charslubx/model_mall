package repository

import (
	"model_mall_backend/backend/internal/models"

	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// Create 创建图片记录
func (r *ImageRepository) Create(image *models.Image) error {
	return r.db.Create(image).Error
}

// GetByID 根据ID获取图片
func (r *ImageRepository) GetByID(id int64) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("id = ? AND status = 1", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByUserID 获取用户的图片列表
func (r *ImageRepository) GetByUserID(userID int64, page, pageSize int) ([]models.Image, int64, error) {
	var images []models.Image
	var total int64

	offset := (page - 1) * pageSize
	
	if err := r.db.Model(&models.Image{}).Where("user_id = ? AND status = 1", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ? AND status = 1", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&images).Error

	return images, total, err
}

// GetByMD5 根据MD5获取图片
func (r *ImageRepository) GetByMD5(md5 string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("md5 = ? AND status = 1", md5).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// Update 更新图片信息
func (r *ImageRepository) Update(image *models.Image) error {
	return r.db.Save(image).Error
}

// Delete 软删除图片
func (r *ImageRepository) Delete(id int64) error {
	return r.db.Model(&models.Image{}).Where("id = ?", id).Update("status", 0).Error
}
