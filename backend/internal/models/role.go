package models

import (
	"time"
)

// Role 角色表
type Role struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:角色ID"`
	Name        string    `json:"name" gorm:"type:varchar(50);uniqueIndex;not null;comment:角色名称"`
	Code        string    `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:角色代码"`
	Description string    `json:"description" gorm:"type:varchar(255);comment:角色描述"`
	Status      int8      `json:"status" gorm:"type:smallint;default:1;comment:状态 0-禁用 1-正常"`
	Sort        int       `json:"sort" gorm:"default:0;comment:排序"`
	IsSystem    bool      `json:"is_system" gorm:"default:false;comment:是否系统角色"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	
	// 关联关系
	RolePermissions []RolePermission `json:"role_permissions" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permissions     []Permission     `json:"permissions" gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	Name        string  `json:"name" validate:"required,min=2,max=50"`
	Code        string  `json:"code" validate:"required,min=2,max=50,alphanum"`
	Description string  `json:"description" validate:"omitempty,max=255"`
	Status      int8    `json:"status" validate:"omitempty,oneof=0 1"`
	Sort        int     `json:"sort" validate:"omitempty,min=0"`
	PermissionIDs []int64 `json:"permission_ids" validate:"omitempty"`
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	ID          int64   `json:"id" validate:"required"`
	Name        string  `json:"name" validate:"omitempty,min=2,max=50"`
	Code        string  `json:"code" validate:"omitempty,min=2,max=50,alphanum"`
	Description string  `json:"description" validate:"omitempty,max=255"`
	Status      int8    `json:"status" validate:"omitempty,oneof=0 1"`
	Sort        int     `json:"sort" validate:"omitempty,min=0"`
	PermissionIDs []int64 `json:"permission_ids" validate:"omitempty"`
}

// RoleResp 角色响应
type RoleResp struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Description string             `json:"description"`
	Status      int8               `json:"status"`
	Sort        int                `json:"sort"`
	IsSystem    bool               `json:"is_system"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Permissions []PermissionResp   `json:"permissions,omitempty"`
}

// RoleListReq 角色列表请求
type RoleListReq struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `json:"keyword" validate:"omitempty,max=50"`
	Status   *int8  `json:"status" validate:"omitempty,oneof=0 1"`
}

// RoleOption 角色选项（用于下拉选择）
type RoleOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}