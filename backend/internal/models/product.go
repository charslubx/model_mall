package models

import (
	"time"
)

// Product 商品表
type Product struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:商品ID"`
	Name           string    `json:"name" gorm:"type:varchar(200);not null;comment:商品名称"`
	Description    string    `json:"description" gorm:"type:text;comment:商品描述"`
	Category       string    `json:"category" gorm:"type:varchar(50);not null;index;comment:商品分类"`
	Price          float64   `json:"price" gorm:"type:decimal(10,2);not null;comment:商品价格"`
	Stock          int       `json:"stock" gorm:"not null;default:0;comment:库存数量"`
	Rating         float64   `json:"rating" gorm:"type:decimal(3,2);default:0;comment:商品评分"`
	Reviews        int       `json:"reviews" gorm:"default:0;comment:评价数量"`
	Sales          int       `json:"sales" gorm:"default:0;comment:销量"`
	Images         string    `json:"images" gorm:"type:text;comment:商品图片URL数组(JSON)"`
	Colors         string    `json:"colors" gorm:"type:text;comment:可选颜色(JSON)"`
	Sizes          string    `json:"sizes" gorm:"type:text;comment:可选尺码(JSON)"`
	Tags           string    `json:"tags" gorm:"type:text;comment:商品标签(JSON)"`
	Features       string    `json:"features" gorm:"type:text;comment:商品特点(JSON)"`
	Specifications string    `json:"specifications" gorm:"type:text;comment:规格参数(JSON)"`
	Status         int8      `json:"status" gorm:"type:smallint;default:1;comment:状态 0-下架 1-在售"`
	MerchantID     int64     `json:"merchant_id" gorm:"not null;index;comment:商户ID"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// Category 商品分类表
type Category struct {
	ID    int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"type:varchar(50);not null;uniqueIndex"`
	Slug  string `json:"slug" gorm:"type:varchar(50);not null;uniqueIndex"`
	Count int    `json:"count" gorm:"default:0;comment:商品数量"`
	Icon  string `json:"icon" gorm:"type:varchar(255);comment:图标URL"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Cart 购物车表
type Cart struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	ProductID int64     `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Color     string    `json:"color" gorm:"type:varchar(50)"`
	Size      string    `json:"size" gorm:"type:varchar(20)"`
	Selected  bool      `json:"selected" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Cart) TableName() string {
	return "carts"
}

// Order 订单表
type Order struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo         string     `json:"order_no" gorm:"type:varchar(50);not null;uniqueIndex;comment:订单号"`
	UserID          int64      `json:"user_id" gorm:"not null;index;comment:用户ID"`
	Status          string     `json:"status" gorm:"type:varchar(20);not null;index;comment:订单状态"`
	PaymentMethod   string     `json:"payment_method" gorm:"type:varchar(20);comment:支付方式"`
	PaymentStatus   string     `json:"payment_status" gorm:"type:varchar(20);comment:支付状态"`
	PaidAt          *time.Time `json:"paid_at" gorm:"comment:支付时间"`
	ShippingName    string     `json:"shipping_name" gorm:"type:varchar(50);comment:收货人"`
	ShippingPhone   string     `json:"shipping_phone" gorm:"type:varchar(20);comment:收货电话"`
	ShippingAddress string     `json:"shipping_address" gorm:"type:varchar(500);comment:收货地址"`
	TrackingNumber  string     `json:"tracking_number" gorm:"type:varchar(100);comment:物流单号"`
	ShippingCompany string     `json:"shipping_company" gorm:"type:varchar(50);comment:物流公司"`
	ShippedAt       *time.Time `json:"shipped_at" gorm:"comment:发货时间"`
	Subtotal        float64    `json:"subtotal" gorm:"type:decimal(10,2);comment:小计"`
	ShippingFee     float64    `json:"shipping_fee" gorm:"type:decimal(10,2);default:0;comment:运费"`
	Tax             float64    `json:"tax" gorm:"type:decimal(10,2);default:0;comment:税费"`
	Total           float64    `json:"total" gorm:"type:decimal(10,2);not null;comment:总金额"`
	Note            string     `json:"note" gorm:"type:text;comment:订单备注"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单商品表
type OrderItem struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   int64     `json:"order_id" gorm:"not null;index"`
	ProductID int64     `json:"product_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"type:varchar(200);not null"`
	Image     string    `json:"image" gorm:"type:varchar(500)"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Color     string    `json:"color" gorm:"type:varchar(50)"`
	Size      string    `json:"size" gorm:"type:varchar(20)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}
