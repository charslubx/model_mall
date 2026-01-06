package repository

import (
	"context"

	"model_mall_backend/backend/internal/models"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

// Create 添加到购物车
func (r *CartRepository) Create(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

// GetByID 根据ID获取购物车项
func (r *CartRepository) GetByID(ctx context.Context, id int64) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).First(&cart, id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetByUserID 获取用户的购物车列表
func (r *CartRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Cart, error) {
	var carts []*models.Cart
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&carts).Error
	return carts, err
}

// Update 更新购物车项
func (r *CartRepository) Update(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

// UpdateQuantity 更新数量
func (r *CartRepository) UpdateQuantity(ctx context.Context, id int64, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.Cart{}).
		Where("id = ?", id).
		Update("quantity", quantity).Error
}

// Delete 删除购物车项
func (r *CartRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Cart{}, id).Error
}

// DeleteByUserID 清空用户购物车
func (r *CartRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.Cart{}).Error
}

// FindByUserAndProduct 查找用户的特定商品购物车项
func (r *CartRepository) FindByUserAndProduct(ctx context.Context, userID, productID int64, color, size string) (*models.Cart, error) {
	var cart models.Cart
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID)

	if color != "" {
		query = query.Where("color = ?", color)
	}
	if size != "" {
		query = query.Where("size = ?", size)
	}

	err := query.First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}
