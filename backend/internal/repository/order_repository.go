package repository

import (
	"context"
	"strconv"

	"model_mall_backend/backend/internal/models"

	"gorm.io/gorm"
)

// stringStatusToInt 将字符串状态转换为整数
func stringStatusToInt(status string) int {
	statusMap := map[string]int{
		"pending":   models.OrderStatusPending,
		"paid":      models.OrderStatusPaid,
		"shipped":   models.OrderStatusShipped,
		"completed": models.OrderStatusCompleted,
		"cancelled": models.OrderStatusCancelled,
	}
	if s, ok := statusMap[status]; ok {
		return s
	}
	// 尝试直接解析为整数
	if i, err := strconv.Atoi(status); err == nil {
		return i
	}
	return -1 // 无效状态
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// Create 创建订单
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// CreateWithItems 创建订单及订单项
func (r *OrderRepository) CreateWithItems(ctx context.Context, order *models.Order, items []*models.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 创建订单项
		for _, item := range items {
			item.OrderID = order.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByID 根据ID获取订单
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取订单
func (r *OrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderItems 获取订单商品列表
func (r *OrderRepository) GetOrderItems(ctx context.Context, orderID int64) ([]*models.OrderItem, error) {
	var items []*models.OrderItem
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Find(&items).Error
	return items, err
}

// Update 更新订单
func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus 更新订单状态
func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status int) error {
	return r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("order_status", status).Error
}

// ListByUserID 获取用户的订单列表
func (r *OrderRepository) ListByUserID(ctx context.Context, userID int64, page, pageSize int, status string) ([]*models.Order, int64, error) {
	db := r.db.WithContext(ctx).Model(&models.Order{})

	// 用户ID条件
	db = db.Where("user_id = ?", userID)

	// 状态条件
	if status != "" {
		statusInt := stringStatusToInt(status)
		if statusInt >= 0 {
			db = db.Where("order_status = ?", statusInt)
		}
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
	var orders []*models.Order
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error

	return orders, total, err
}

// ListByMerchant 获取商户的订单列表
func (r *OrderRepository) ListByMerchant(ctx context.Context, merchantID int64, page, pageSize int, status string) ([]*models.Order, int64, error) {
	db := r.db.WithContext(ctx).Model(&models.Order{})

	// 通过订单项关联查询商户的订单
	db = db.Joins("INNER JOIN order_items ON orders.id = order_items.order_id").
		Joins("INNER JOIN products ON order_items.product_id = products.id").
		Where("products.seller_id = ?", merchantID).
		Group("orders.id")

	// 状态条件
	if status != "" {
		statusInt := stringStatusToInt(status)
		if statusInt >= 0 {
			db = db.Where("orders.order_status = ?", statusInt)
		}
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
	var orders []*models.Order
	err := db.Offset(offset).Limit(pageSize).Order("orders.created_at DESC").Find(&orders).Error

	return orders, total, err
}

// GetByUserID 获取用户订单列表(别名方法)
func (r *OrderRepository) GetByUserID(ctx context.Context, userID int64, status string, page, pageSize int) ([]*models.Order, int64, error) {
	return r.ListByUserID(ctx, userID, page, pageSize, status)
}

// GetBySellerID 获取商户订单列表(别名方法)
func (r *OrderRepository) GetBySellerID(ctx context.Context, merchantID int64, status string, page, pageSize int) ([]*models.Order, int64, error) {
	return r.ListByMerchant(ctx, merchantID, page, pageSize, status)
}
