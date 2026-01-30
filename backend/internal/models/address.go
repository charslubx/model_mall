package models

import "time"

// Address 收货地址模型
type Address struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	Name      string    `db:"name" json:"name"`            // 收货人姓名
	Phone     string    `db:"phone" json:"phone"`          // 联系电话
	Province  string    `db:"province" json:"province"`    // 省份
	City      string    `db:"city" json:"city"`            // 城市
	District  string    `db:"district" json:"district"`    // 区/县
	Address   string    `db:"address" json:"address"`      // 详细地址
	IsDefault bool      `db:"is_default" json:"isDefault"` // 是否默认地址
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
