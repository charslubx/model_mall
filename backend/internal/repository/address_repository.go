package repository

import (
	"context"

	"model_mall_backend/backend/internal/models"

	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{
		db: db,
	}
}

// Create 创建地址
func (r *AddressRepository) Create(ctx context.Context, address *models.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

// GetByID 根据ID获取地址
func (r *AddressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	var address models.Address
	err := r.db.WithContext(ctx).First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByUserID 根据用户ID获取地址列表
func (r *AddressRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	var addresses []*models.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetByIDAndUserID 根据ID和用户ID获取地址（确保用户只能访问自己的地址）
func (r *AddressRepository) GetByIDAndUserID(ctx context.Context, id, userID int64) (*models.Address, error) {
	var address models.Address
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Update 更新地址
func (r *AddressRepository) Update(ctx context.Context, address *models.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}

// UpdateByID 根据ID更新地址
func (r *AddressRepository) UpdateByID(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Address{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除地址
func (r *AddressRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Address{}, id).Error
}

// DeleteByIDAndUserID 根据ID和用户ID删除地址（确保用户只能删除自己的地址）
func (r *AddressRepository) DeleteByIDAndUserID(ctx context.Context, id, userID int64) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Address{}).Error
}

// ClearDefaultByUserID 清除用户的所有默认地址标记
func (r *AddressRepository) ClearDefaultByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

// SetDefaultByIDAndUserID 设置指定地址为默认地址
func (r *AddressRepository) SetDefaultByIDAndUserID(ctx context.Context, id, userID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先清除该用户的所有默认地址
		if err := tx.Model(&models.Address{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置指定地址为默认
		return tx.Model(&models.Address{}).
			Where("id = ? AND user_id = ?", id, userID).
			Update("is_default", true).Error
	})
}

// GetDefaultByUserID 获取用户的默认地址
func (r *AddressRepository) GetDefaultByUserID(ctx context.Context, userID int64) (*models.Address, error) {
	var address models.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// CountByUserID 统计用户的地址数量
func (r *AddressRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Address{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
