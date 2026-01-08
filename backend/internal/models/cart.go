package models

import "time"

// Cart 购物车模型
type Cart struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null;index" json:"user_id"`
	ProductID int64     `gorm:"column:product_id;not null;index" json:"product_id"`
	Quantity  int       `gorm:"column:quantity;default:1" json:"quantity"`
	Color     string    `gorm:"column:color;size:50" json:"color"`
	Size      string    `gorm:"column:size;size:50" json:"size"`
	Selected  bool      `gorm:"column:selected;default:true" json:"selected"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Cart) TableName() string {
	return "cart_items"
}
