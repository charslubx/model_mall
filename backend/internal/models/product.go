package models

import (
	"time"
)

// Product 商品表
type Product struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:商品ID"`
	SellerID      int64     `json:"seller_id" gorm:"column:seller_id;not null;index;comment:卖家ID"`
	CategoryID    int64     `json:"category_id" gorm:"column:category_id;not null;index;comment:商品分类ID"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null;comment:商品名称"`
	Description   string    `json:"description" gorm:"type:text;comment:商品描述"`
	Price         float64   `json:"price" gorm:"type:decimal(10,2);not null;comment:商品价格"`
	OriginalPrice float64   `json:"original_price" gorm:"column:original_price;type:decimal(10,2);comment:原价"`
	Stock         int       `json:"stock" gorm:"not null;default:0;comment:库存数量"`
	Sales         int       `json:"sales" gorm:"default:0;comment:销量"`
	Image         string    `json:"image" gorm:"type:varchar(255);comment:主图片URL"`
	Images        string    `json:"images" gorm:"type:text;comment:商品图片URL数组(JSON)"`
	Status        int8      `json:"status" gorm:"type:smallint;default:1;comment:状态 0-下架 1-在售"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`

	// 以下字段不在数据库表中，用于业务逻辑
	Rating         float64 `json:"rating" gorm:"-"`
	Reviews        int     `json:"reviews" gorm:"-"`
	Colors         string  `json:"colors" gorm:"-"`
	Sizes          string  `json:"sizes" gorm:"-"`
	Tags           string  `json:"tags" gorm:"-"`
	Features       string  `json:"features" gorm:"-"`
	Specifications string  `json:"specifications" gorm:"-"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// Category 商品分类表
type Category struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	ParentID int64  `json:"parent_id" gorm:"column:parent_id;default:0"`
	Level    int    `json:"level" gorm:"default:1"`
	Sort     int    `json:"sort" gorm:"default:0"`
	Icon     string `json:"icon" gorm:"type:varchar(255)"`
	Image    string `json:"image" gorm:"type:varchar(255)"`
	Status   int8   `json:"status" gorm:"type:smallint;default:1"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
