package repository

import (
	"context"

	"model_mall_backend/backend/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Create 创建商品
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID 根据ID获取商品
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update 更新商品
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// UpdateByID 根据ID更新商品
func (r *ProductRepository) UpdateByID(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除商品
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

// List 获取商品列表
func (r *ProductRepository) List(ctx context.Context, page, pageSize int, category, keyword, sortBy string, minPrice, maxPrice float64) ([]*models.Product, int64, error) {
	db := r.db.WithContext(ctx).Model(&models.Product{})

	// 只查询在售商品
	db = db.Where("status = ?", 1)

	// 构建查询条件
	if category != "" {
		db = db.Where("category = ?", category)
	}
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	if minPrice > 0 {
		db = db.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		db = db.Where("price <= ?", maxPrice)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch sortBy {
	case "price-asc":
		db = db.Order("price ASC")
	case "price-desc":
		db = db.Order("price DESC")
	case "name-asc":
		db = db.Order("name ASC")
	case "name-desc":
		db = db.Order("name DESC")
	case "sales":
		db = db.Order("sales DESC")
	case "newest":
		db = db.Order("created_at DESC")
	default:
		db = db.Order("created_at DESC")
	}

	// 分页查询
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	var products []*models.Product
	err := db.Offset(offset).Limit(pageSize).Find(&products).Error

	return products, total, err
}

// ListByMerchant 获取商户的商品列表
func (r *ProductRepository) ListByMerchant(ctx context.Context, merchantID int64, page, pageSize int, category, status, keyword string) ([]*models.Product, int64, error) {
	db := r.db.WithContext(ctx).Model(&models.Product{})

	// 商户ID条件
	db = db.Where("merchant_id = ?", merchantID)

	// 构建查询条件
	if category != "" {
		db = db.Where("category = ?", category)
	}
	if status != "" {
		if status == "active" {
			db = db.Where("status = ?", 1)
		} else if status == "inactive" {
			db = db.Where("status = ?", 0)
		}
	}
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	var products []*models.Product
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error

	return products, total, err
}

// GetCategories 获取所有分类
func (r *ProductRepository) GetCategories(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

// IncrementSales 增加销量
func (r *ProductRepository) IncrementSales(ctx context.Context, productID int64, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", productID).
		UpdateColumn("sales", gorm.Expr("sales + ?", quantity)).Error
}

// DecrementStock 减少库存
func (r *ProductRepository) DecrementStock(ctx context.Context, productID int64, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error
}

// CheckStock 检查库存
func (r *ProductRepository) CheckStock(ctx context.Context, productID int64, quantity int) (bool, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Select("stock").First(&product, productID).Error
	if err != nil {
		return false, err
	}
	return product.Stock >= quantity, nil
}

// IncrementStock 增加库存(用于取消订单等场景)
func (r *ProductRepository) IncrementStock(ctx context.Context, productID int64, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}

// Search 搜索商品(别名方法,方便调用)
func (r *ProductRepository) Search(ctx context.Context, keyword, category string, minPrice, maxPrice float64, sortBy string, page, pageSize int) ([]*models.Product, int64, error) {
	return r.List(ctx, page, pageSize, category, keyword, sortBy, minPrice, maxPrice)
}

// GetByMerchantID 获取商户商品列表(别名方法)
func (r *ProductRepository) GetByMerchantID(ctx context.Context, merchantID int64, category, sortBy string, page, pageSize int) ([]*models.Product, int64, error) {
	return r.ListByMerchant(ctx, merchantID, page, pageSize, category, "", "")
}

// GetByIDs 批量获取商品(用于优化N+1查询)
func (r *ProductRepository) GetByIDs(ctx context.Context, ids []int64) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error
	return products, err
}
