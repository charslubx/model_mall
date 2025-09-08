package models

import (
	"time"
)

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:关联ID"`
	RoleID       int64     `json:"role_id" gorm:"not null;index;comment:角色ID"`
	PermissionID int64     `json:"permission_id" gorm:"not null;index;comment:权限ID"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	
	// 关联关系
	Role       Role       `json:"role" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// RolePermissionBatchReq 批量设置角色权限请求
type RolePermissionBatchReq struct {
	RoleID        int64   `json:"role_id" validate:"required"`
	PermissionIDs []int64 `json:"permission_ids" validate:"required"`
}

// RolePermissionResp 角色权限响应
type RolePermissionResp struct {
	ID           int64              `json:"id"`
	RoleID       int64              `json:"role_id"`
	PermissionID int64              `json:"permission_id"`
	Role         RoleResp           `json:"role,omitempty"`
	Permission   PermissionResp     `json:"permission,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
}