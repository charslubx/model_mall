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
	SellerID       int64     `json:"seller_id" gorm:"column:seller_id;not null;index;comment:卖家ID"`
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
