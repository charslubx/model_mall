package models

import (
	"time"
)

// User 用户表
type User struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID"`
	Username     string     `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名（登录用）"`
	Email        string     `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;comment:邮箱"`
	Phone        string     `json:"phone" gorm:"type:varchar(20);uniqueIndex;comment:手机号"`
	Password     string     `json:"-" gorm:"type:varchar(255);not null;comment:密码哈希"`
	Avatar       string     `json:"avatar" gorm:"type:varchar(255);comment:头像URL"`
	Nickname     string     `json:"nickname" gorm:"type:varchar(50);comment:昵称/显示名称"`
	Gender       int8       `json:"gender" gorm:"type:smallint;default:0;comment:性别 0-未知 1-男 2-女"`
	Birthday     *time.Time `json:"birthday" gorm:"type:date;comment:生日"`
	Status       int8       `json:"status" gorm:"type:smallint;default:1;comment:状态 0-禁用 1-正常"`
	UserType     string     `json:"user_type" gorm:"type:varchar(20);not null;default:'customer';comment:用户类型 customer-顾客 merchant-商户 admin-管理员"`
	MerchantName string     `json:"merchant_name" gorm:"type:varchar(100);comment:商户名称（仅商户使用）"`
	Description  string     `json:"description" gorm:"type:text;comment:商户描述"`
	RoleID       int64      `json:"role_id" gorm:"not null;comment:角色ID"`
	LastLoginAt  *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	LastLoginIP  string     `json:"last_login_ip" gorm:"type:varchar(45);comment:最后登录IP"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`

	// 关联关系
	Role Role `json:"role" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
	Gender   int8   `json:"gender" validate:"omitempty,oneof=0 1 2"`
	RoleID   int64  `json:"role_id" validate:"required"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	ID       int64  `json:"id" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
	Avatar   string `json:"avatar" validate:"omitempty,url"`
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
	Gender   int8   `json:"gender" validate:"omitempty,oneof=0 1 2"`
	Status   int8   `json:"status" validate:"omitempty,oneof=0 1"`
	RoleID   int64  `json:"role_id" validate:"omitempty"`
}

// UserResp 用户响应
type UserResp struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	Nickname    string     `json:"nickname"`
	Gender      int8       `json:"gender"`
	Birthday    *time.Time `json:"birthday"`
	Status      int8       `json:"status"`
	RoleID      int64      `json:"role_id"`
	RoleName    string     `json:"role_name"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserListReq 用户列表请求
type UserListReq struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `json:"keyword" validate:"omitempty,max=50"`
	Status   *int8  `json:"status" validate:"omitempty,oneof=0 1"`
	RoleID   *int64 `json:"role_id" validate:"omitempty"`
}
