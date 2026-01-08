package models

import "time"

// 订单状态常量
const (
	OrderStatusPending   = 0 // 待支付
	OrderStatusPaid      = 1 // 已支付
	OrderStatusShipped   = 2 // 已发货
	OrderStatusCompleted = 3 // 已完成
	OrderStatusCancelled = 4 // 已取消
)

// Order 订单模型
type Order struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo         string     `gorm:"column:order_no;size:50;not null;uniqueIndex" json:"order_no"`
	UserID          int64      `gorm:"column:user_id;not null;index" json:"user_id"`
	SellerID        int64      `gorm:"column:seller_id;not null;index" json:"seller_id"`
	Total           float64    `gorm:"column:total_amount;type:decimal(10,2);not null" json:"total"`
	PaymentMethod   string     `gorm:"column:payment_method;size:50" json:"payment_method"`
	PaymentStatus   int        `gorm:"column:payment_status;default:0" json:"payment_status"`
	Status          int        `gorm:"column:order_status;default:0;index" json:"status"`
	ShippingAddress string     `gorm:"column:shipping_address;type:text" json:"shipping_address"`
	ShippingName    string     `gorm:"column:shipping_name;size:100" json:"shipping_name"`
	ShippingPhone   string     `gorm:"column:shipping_phone;size:20" json:"shipping_phone"`
	TrackingNumber  string     `gorm:"column:tracking_number;size:100" json:"tracking_number"`
	Remark          string     `gorm:"column:remark;type:text" json:"remark"`
	PaidAt          *time.Time `gorm:"column:paid_at" json:"paid_at"`
	ShippedAt       *time.Time `gorm:"column:shipped_at" json:"shipped_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at" json:"completed_at"`
	CancelledAt     *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID   int64     `gorm:"column:order_id;not null;index" json:"order_id"`
	ProductID int64     `gorm:"column:product_id;not null;index" json:"product_id"`
	Name      string    `gorm:"column:name;size:255;not null" json:"name"`
	Image     string    `gorm:"column:image;size:255" json:"image"`
	Price     float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Quantity  int       `gorm:"column:quantity;not null" json:"quantity"`
	Subtotal  float64   `gorm:"column:subtotal;type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}
